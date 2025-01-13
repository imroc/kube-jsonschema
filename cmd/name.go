package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

func ParseName(name, groupFromUrl string) (group, version, kind string) {
	if name == "meta.apis.pkg.apimachinery.io.k8s" {
		name = "meta.apis.pkg.apimachinery.k8s.io"
	}
	ss := strings.Split(name, ".")
	if len(ss) < 2 {
		panic(fmt.Sprintf("invalid name: %s", name))
	}
	kind = strings.ToLower(ss[len(ss)-1])
	ss = ss[:len(ss)-1] // 去除 "kind"
	version = ss[len(ss)-1]
	ss = ss[:len(ss)-1] // 去除 "version"
	if len(ss) > 0 {    // "io.k8s" --> "k8s.io"
		if len(ss) > 1 && ss[0] == "k8s" && ss[1] == "io" {
			ss[0] = "io"
			ss[1] = "k8s"
		}
		slices.Reverse(ss)
		group = strings.Join(ss, ".")
	} else {
		group = groupFromUrl
	}
	return
}

type FileInfo struct {
	Filename string
	Group    string
	Version  string
	Kind     string
}

func GetFileInfo(name, groupFromUrl string) FileInfo {
	group, version, kind := ParseName(name, groupFromUrl)
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
