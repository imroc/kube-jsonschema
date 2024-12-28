package cmd

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// https://raw.githubusercontent.com/imroc/kube-jsonschema/master/schemas/advancedcronjob.json
func NewIndexCmd(args []string) *cobra.Command {
	var outDir, extraDir string
	cmd := &cobra.Command{
		Use:               "index",
		Short:             "generate index for json schema files",
		DisableAutoGenTag: true,
		Args:              cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runIndex(outDir, extraDir)
		},
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cmd.SetArgs(args)
	flags := cmd.Flags()
	flags.StringVar(&outDir, "out-dir", cwd, "json schema output directory")
	flags.StringVar(&extraDir, "extra-dir", "", "extra json schema directory")
	return cmd
}

func runIndex(outDir, extraDir string) error {
	outDir, err := filepath.EvalSymlinks(outDir)
	if err != nil {
		return err
	}
	if extraDir != "" {
		extraDir, err = filepath.EvalSymlinks(extraDir)
		if err != nil {
			return err
		}
	}
	refs := []string{}
	err = walkDir(outDir, &refs, false)
	if err != nil {
		return err
	}
	err = walkDir(extraDir, &refs, true)
	if err != nil {
		return err
	}
	var allJson AllJson
	for _, file := range refs {
		allJson.OneOf = append(allJson.OneOf, Ref{file})
	}
	return writePrettyJson(&allJson, filepath.Join(outDir, "kubernetes.json"))
}

var keys = make(map[string]bool)

func walkDir(outDir string, refs *[]string, isExtra bool) error {
	outDir = strings.TrimSuffix(outDir, "/")
	filepath.WalkDir(outDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.Contains(path, "meta.apis.pkg.apimachinery.k8s.io") {
			return nil
		}
		if filepath.Ext(d.Name()) != ".json" || d.Name() == "kubernetes.json" {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		bs, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		schema := gjson.ParseBytes(bs)
		gvk := schema.Get(XGVK_NAME)
		if !gvk.Exists() {
			return nil
		}
		// if len(gvk.Array()) != 1 {
		// 	return nil
		// }
		if len(schema.Get("properties.apiVersion.enum").Array()) == 0 {
			return nil
		}
		if len(schema.Get("properties.kind.enum").Array()) == 0 {
			return nil
		}
		dir := filepath.Base(filepath.Dir(path))
		key := dir + "/" + d.Name()
		if keys[key] {
			return nil
		}
		keys[key] = true
		if isExtra {
			key = outDir + "/" + key
		}
		*refs = append(*refs, key)
		return nil
	})
	return nil
}

type Ref struct {
	Ref string `json:"$ref"`
}

type AllJson struct {
	OneOf []Ref `json:"oneOf"`
}
