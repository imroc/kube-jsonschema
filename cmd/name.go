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

func GetFilename(name, group, version, kind string) string {
	if !strings.Contains(group, ".") {
		group, version, kind = ParseName(name)
	}
	if strings.HasSuffix(group, ".api.k8s.io") {
		group = strings.TrimSuffix(group, ".api.k8s.io")
	}
	return filepath.Join(group, fmt.Sprintf("%s_%s", kind, version))
}
