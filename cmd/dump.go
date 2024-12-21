package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/imroc/kubeschema/pkg/schemas"
	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

type OpenAPIV3Root struct {
	Paths map[string]PathInfo `json:"paths"`
}

type PathInfo struct {
	ServerRelativeURL string `json:"serverRelativeURL"`
}

func NewDumpCmd(args []string) *cobra.Command {
	var address, outDir string
	var port int
	var pretty, force bool

	cmd := &cobra.Command{
		Use:               "dump",
		Short:             "dump json schema from openapi v3 of kubernetes apiserver",
		DisableAutoGenTag: true,
		Args:              cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			url := fmt.Sprintf("http://%s:%d/openapi/v3", address, port)
			var root OpenAPIV3Root
			_, err := req.R().SetSuccessResult(&root).Get(url)
			if err != nil {
				return err
			}
			if root.Paths == nil {
				return errors.New("invalid openapi v3 format, missing key `paths`")
			}
			if outDir != "" {
				_, err = os.Stat(outDir)
				if os.IsNotExist(err) {
					err = os.Mkdir(outDir, 0755)
					if err != nil {
						return err
					}
				}
			}
			for path, info := range root.Paths {
				if !strings.HasPrefix(path, "apis/") && !strings.HasPrefix(path, "api/") {
					continue
				}
				urlPath := info.ServerRelativeURL
				err = parseApisEndpoint(outDir, fmt.Sprintf("http://%s:%d%s", address, port, urlPath), pretty, force)
				if err != nil {
					fmt.Println("ERROR:", err.Error())
				}
			}
			return nil
		},
	}
	cmd.SetArgs(args)
	flags := cmd.Flags()
	flags.StringVar(&address, "address", "127.0.0.1", "The IP address on which kubectl proxy is serving on.")
	flags.StringVar(&outDir, "out-dir", "kubeschemas", "json schema output directory")
	flags.IntVar(&port, "port", 8001, "The port on which kubectl proxy is listening.")
	flags.BoolVar(&pretty, "pretty", true, "whether write json in pretty format")
	flags.BoolVar(&force, "force", false, "whether to override the existed json schema file")
	return cmd
}

var regRef = regexp.MustCompile(`"#/components/schemas/([^"]+)"`)

func parseApisEndpoint(outDir, url string, pretty, force bool) error {
	resp, err := req.R().Get(url)
	if err != nil {
		return err
	}
	if resp.GetStatusCode() != http.StatusOK {
		return errors.New("unexpected status: " + resp.GetStatus())
	}
	body := resp.String()
	body = regRef.ReplaceAllStringFunc(body, strings.ToLower)
	body = regRef.ReplaceAllString(body, `"$1.json"`)
	scms := gjson.Get(body, "components.schemas").Map()
	for name, schema := range scms {
		m, ok := schema.Value().(map[string]any)
		if !ok {
			fmt.Printf("WARN: invalid schema: %s\n", schema.Raw)
			continue
		}
		var filename string
		gvk := schema.Get(XGVK_NAME + ".0")
		if gvk.Exists() {
			group := gvk.Get("group").String()
			version := gvk.Get("version").String()
			kind := gvk.Get("kind").String()
			if version == "" || kind == "" {
				fmt.Printf("WARN: skip empty version or kind: %s\n", name)
				continue
			}

			var apiVersion string
			if group != "" {
				apiVersion = group + "/" + version
				filename = group + "-" + version
			} else {
				apiVersion = version
				filename = version
			}
			filename = kind + "-" + filename
			setMap(m, []string{apiVersion}, "properties", "apiVersion", "enum")
			setMap(m, []string{kind}, "properties", "kind", "enum")
			m["required"] = []string{"apiVersion", "kind"}
		} else {
			filename = name
		}
		filename = strings.ToLower(filename)
		if !force && schemas.Exists(outDir, filename) {
			continue
		}
		modifySchema(m)
		err = writeJson(outDir, filename, pretty, m)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeJson(outDir, name string, pretty bool, v any) error {
	file := filepath.Join(outDir, name+".json")
	fmt.Printf("write %s\n", file)
	indent := ""
	if pretty {
		indent = "  "
	}
	data, err := json.MarshalIndent(v, "", indent)
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

const XGVK_NAME = "x-kubernetes-group-version-kind"

func setMap(m map[string]any, value any, keys ...string) {
	var ok bool
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		if i == len(keys)-1 {
			m[key] = value
			return
		} else {
			m, ok = m[key].(map[string]any)
			if !ok {
				// fmt.Printf("WARN: invalid type: %T(%s) - %v\n", m[key], key, m)
				return
			}
		}
	}
}
