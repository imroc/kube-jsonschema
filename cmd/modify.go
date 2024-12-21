package cmd

import "fmt"

// 修改 schema 以兼容 yaml-language-server
func modifySchema(m map[string]any) {
	for key, value := range m {
		switch act := handle(key, value); act {
		case CONTINUE:
			if mm, ok := value.(map[string]any); ok {
				modifySchema(mm)
			}
		case STOP:
			continue
		case DELETE:
			delete(m, key)
		default:
			panic(fmt.Sprintf("unknown action %v", act))
		}
	}
}

type ACTION int

const (
	CONTINUE ACTION = iota
	STOP
	DELETE
)

func handle(key string, value any) ACTION {
	switch v := value.(type) {
	case map[string]any:
		// OpenAPI v3 中 exclusiveMinimum/exclusiveMaximum 的值为 bool，转换为 yamlls 中的 number
		if exclusiveMinimum, ok := v["exclusiveMinimum"]; ok {
			if vv, ok := exclusiveMinimum.(bool); ok && vv {
				v["exclusiveMinimum"] = v["minimum"]
				delete(v, "minimum")
			}
		}
		if exclusiveMaximum, ok := v["exclusiveMaximum"]; ok {
			if vv, ok := exclusiveMaximum.(bool); ok && vv {
				v["exclusiveMaximum"] = v["maximum"]
				delete(v, "maximum")
			}
		}
	}
	return CONTINUE
}
