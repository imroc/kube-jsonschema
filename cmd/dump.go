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
	var address, outDir, extraDir, group string
	var pretty, force, index bool

	cmd := &cobra.Command{
		Use:               "dump",
		Short:             "dump json schema from openapi v3 of kubernetes apiserver",
		DisableAutoGenTag: true,
		Args:              cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if address == "" {
				cmd, port, err := runKubectlProxy()
				if err != nil {
					return err
				}
				address = fmt.Sprintf("%s:%s", "127.0.0.1", port)
				defer cmd.Process.Kill()
			}

			url := fmt.Sprintf("http://%s/openapi/v3", address)
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
				err = parseApisEndpoint(outDir, fmt.Sprintf("http://%s%s", address, urlPath), group, pretty, force)
				if err != nil {
					fmt.Println("ERROR:", err.Error())
				}
			}
			if index {
				return runIndex(outDir, extraDir)
			}
			return nil
		},
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd.SetArgs(args)
	flags := cmd.Flags()
	flags.StringVar(&group, "group", "", "API group to dump")
	flags.StringVar(&address, "address", "", "The IP address on which kubectl proxy is serving on.")
	flags.StringVar(&outDir, "out-dir", cwd, "json schema output directory")
	flags.StringVar(&extraDir, "extra-dir", "", "extra json schema directory")
	flags.BoolVar(&pretty, "pretty", true, "whether write json in pretty format")
	flags.BoolVar(&index, "index", false, "whether to index all json schema after dump")
	flags.BoolVar(&force, "force", false, "whether to override the existed json schema file")
	return cmd
}

var regRef = regexp.MustCompile(`"#/components/schemas/([^"]+)"`)

func parseApisEndpoint(outDir, url, groupOnly string, pretty, force bool) error {
	resp, err := req.R().Get(url)
	if err != nil {
		return err
	}
	if resp.GetStatusCode() != http.StatusOK {
		return errors.New("unexpected status: " + resp.GetStatus())
	}
	body := resp.String()
	body = regRef.ReplaceAllStringFunc(body, func(s string) string {
		s = strings.TrimPrefix(s, `"#/components/schemas/`)
		s = strings.TrimSuffix(s, `"`)
		fi := GetFileInfo(s)
		return `"../` + fi.Filename + `.json"`
	})
	// body = regRef.ReplaceAllStringFunc(body, strings.ToLower)
	// body = regRef.ReplaceAllString(body, `"$1.json"`)
	scms := gjson.Get(body, "components.schemas").Map()
	for name, schema := range scms {
		m, ok := schema.Value().(map[string]any)
		if !ok {
			fmt.Printf("WARN: invalid schema: %s\n", schema.Raw)
			continue
		}
		gvk := schema.Get(XGVK_NAME + ".0")
		var group, version, kind string
		if gvk.Exists() {
			group = gvk.Get("group").String()
			version = gvk.Get("version").String()
			kind = gvk.Get("kind").String()
			if version == "" || kind == "" {
				fmt.Printf("WARN: skip empty version or kind: %s\n", name)
				continue
			}
			var apiVersion string
			if group != "" {
				apiVersion = group + "/" + version
			} else {
				apiVersion = version
			}
			setMap(m, []string{apiVersion}, "properties", "apiVersion", "enum")
			setMap(m, []string{kind}, "properties", "kind", "enum")
			m["required"] = []string{"apiVersion", "kind"}
		}
		fi := GetFileInfo(name)
		if groupOnly != "" && fi.Group != groupOnly {
			continue
		}
		filename := strings.ToLower(fi.Filename)
		if groupOnly == "" && !force && schemas.Exists(outDir, filename) {
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
	filePath := filepath.Join(outDir, name+".json")
	fmt.Printf("write %s\n", filePath)
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	indent := ""
	if pretty {
		indent = "  "
	}
	data, err := json.MarshalIndent(v, "", indent)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
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
