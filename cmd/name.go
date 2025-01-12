package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

func ParseName(name string) (group, version, kind string) {
	ss := strings.Split(name, ".")
	if len(ss) < 3 {
		panic(fmt.Sprintf("invalid name: %s", name))
	}
	kind = strings.ToLower(ss[len(ss)-1])
	ss = ss[:len(ss)-1] // 去除 "kind"
	version = ss[len(ss)-1]
	ss = ss[:len(ss)-1] // 去除 "version"
	slices.Reverse(ss)
	group = strings.Join(ss, ".")
	return
}

type FileInfo struct {
	Filename string
	Group    string
	Version  string
	Kind     string
}

func GetFileInfo(name string) FileInfo {
	group, version, kind := ParseName(name)
	if !strings.Contains(group, ".") {
		group = group + ".api.k8s.io"
	}
	filename := filepath.Join(group, fmt.Sprintf("%s_%s", kind, version))
	return FileInfo{
		Filename: filename,
		Group:    group,
		Version:  version,
		Kind:     kind,
	}
}
