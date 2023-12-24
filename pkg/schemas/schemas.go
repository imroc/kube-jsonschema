package schemas

import (
	"path/filepath"

	"github.com/imroc/kubeschema/pkg/files"
)

func Exists(outDir, name string) bool {
	file := filepath.Join(outDir, name+".json")
	return files.Exists(file)
}
