package cmd

import (
	"strings"
	"testing"
)

func TestRegRef(t *testing.T) {
	body := regRef.ReplaceAllStringFunc(`"#/components/schemas/io.k8s.apimachinery.pkg.apis.meta.v1.objectmeta"`, func(s string) string {
		s = strings.TrimPrefix(s, `"#/components/schemas/`)
		s = strings.TrimSuffix(s, `"`)
		filename := GetFilename(s, "", "", "")
		return filename + ".json"
	})
	t.Log((body))
}
