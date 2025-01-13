package cmd

import "testing"

func TestRegURL(t *testing.T) {
	url := "http://127.0.0.1:59878/openapi/v3/apis/subresources.kubevirt.io/v1?hash=23DD72781956371BA9197C26C0B2F65F428A419890615B21364EE4173E2D065C6105CD3F2E99CB4339DB5DC4F10A22EDE02A30BCA9BC558EF48A7B409A4C7C67"
	ss := regURL.FindStringSubmatch(url)
	if ss[1] != "subresources.kubevirt.io" {
		t.Error("failed to match url")
	}
}
