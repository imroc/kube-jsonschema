package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewCrdCmd(args []string) *cobra.Command {
	var rootDir, schemasDir string
	cmd := &cobra.Command{
		Use:               "crd [ROOT DIR]",
		Short:             "generate json schema from root dir of crd files",
		DisableAutoGenTag: true,
		Args:              cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return generateFromCRDDir(rootDir, schemasDir)
		},
	}
	cmd.SetArgs(args)
	flags := cmd.Flags()
	flags.StringVar(&rootDir, "root-dir", "crds", "root dir of crd files")
	flags.StringVar(&schemasDir, "schemas-dir", "crdSchemas", "root output dir of crd json schema files")
	return cmd
}

type CRDVersion struct {
	Schema struct {
		OpenAPIV3Schema map[string]any `json:"openAPIV3Schema" yaml:"openAPIV3Schema"`
	} `json:"schema" yaml:"schema"`
	Name string `json:"name" yaml:"name"`
}

type CRD struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Metadata   struct {
		Name string `json:"name" yaml:"name"`
	} `json:"metadata" yaml:"metadata"`
	Spec struct {
		Group string `json:"group" yaml:"group"`
		Names struct {
			Kind     string `json:"kind" yaml:"kind"`
			Singular string `json:"singular" yaml:"singular"`
		} `json:"names" yaml:"names"`
		Versions []CRDVersion `json:"versions" yaml:"versions"`
	} `json:"spec" yaml:"spec"`
}

func generateFromCRDDir(rootDir, schemasDir string) error {
	return filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		switch filepath.Ext(path) {
		case ".yaml", ".yml":
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			decoder := yaml.NewDecoder(file)
			for {
				crd := new(CRD)
				err = decoder.Decode(crd)
				{
					err := writeCRDSchemas(schemasDir, crd)
					if err != nil {
						return err
					}
				}
				if errors.Is(err, io.EOF) {
					return nil
				}
			}
		case ".json":
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var crd CRD
			err = json.Unmarshal([]byte(data), &crd)
			if err != nil {
				return err
			}
			return writeCRDSchemas(schemasDir, &crd)
		}
		return nil
	})
}

func writeCRDSchemas(schemaDir string, crd *CRD) error {
	if crd.Kind == "" {
		return nil
	}
	if crd.Kind != "CustomResourceDefinition" {
		fmt.Println("skip", crd.Kind)
		return nil
	}
	group := crd.Spec.Group
	kind := crd.Spec.Names.Singular
	for _, v := range crd.Spec.Versions {
		version := v.Name
		filename := fmt.Sprintf("%s-%s.json", kind, version)
		path := filepath.Join(schemaDir, group)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		err = writeCRDSchema(filepath.Join(path, filename), v.Schema.OpenAPIV3Schema)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeCRDSchema(name string, schema map[string]any) error {
	// 如果同名文件已经存在，不再覆盖
	info, err := os.Stat(name)
	if err == nil && info.Size() > 0 {
		return nil
	}
	// 转成 json 格式并写入文件
	return writePrettyJson(schema, name)
}
