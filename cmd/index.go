package cmd

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// https://raw.githubusercontent.com/imroc/kube-jsonschema/master/schemas/advancedcronjob.json
func NewIndexCmd(args []string) *cobra.Command {
	var outDir string
	cmd := &cobra.Command{
		Use:               "index",
		Short:             "generate index for json schema files",
		DisableAutoGenTag: true,
		Args:              cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runIndex(outDir)
		},
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cmd.SetArgs(args)
	flags := cmd.Flags()
	flags.StringVar(&outDir, "out-dir", cwd, "json schema output directory")
	return cmd
}

func runIndex(outDir string) error {
	refs := []string{}
	err := walkDir(outDir, &refs)
	if err != nil {
		return err
	}
	var allJson AllJson
	for _, file := range refs {
		allJson.OneOf = append(allJson.OneOf, Ref{file})
	}
	return writePrettyJson(&allJson, filepath.Join(outDir, "kubernetes.json"))
}

func walkDir(outDir string, refs *[]string) error {
	filepath.WalkDir(outDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
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
		if !schema.Get("properties.apiVersion").Exists() {
			return nil
		}
		if !schema.Get("properties.kind").Exists() {
			return nil
		}
		dir := filepath.Base(filepath.Dir(path))
		*refs = append(*refs, dir+"/"+d.Name())
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
