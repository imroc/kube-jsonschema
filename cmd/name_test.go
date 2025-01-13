package cmd

import "testing"

func TestParseName(t *testing.T) {
	group, version, kind := ParseName("io.k8s.api.core.v1.pod", "")
	if group != "core.api.k8s.io" {
		t.Errorf("invalid group %s", group)
	}
	if version != "v1" {
		t.Errorf("invalid group %s", version)
	}
	if kind != "pod" {
		t.Errorf("invalid kind %s", kind)
	}
}
