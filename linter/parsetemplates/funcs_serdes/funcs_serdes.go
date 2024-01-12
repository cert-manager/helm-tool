package funcs_serdes

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func FuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	fmap := template.FuncMap{
		"toToml":        func(v interface{}) string { return "" },
		"toYaml":        func(v interface{}) string { return "" },
		"fromYaml":      func(str string) map[string]interface{} { return map[string]interface{}{} },
		"fromYamlArray": func(str string) []interface{} { return nil },
		"toJson":        func(v interface{}) string { return "" },
		"fromJson":      func(str string) map[string]interface{} { return map[string]interface{}{} },
		"fromJsonArray": func(str string) []interface{} { return nil },
		"lookup": func(string, string, string, string) (map[string]interface{}, error) {
			return map[string]interface{}{}, nil
		},
		"include":  func(name string, data interface{}) (string, error) { return "", nil },
		"required": func(warn string, val interface{}) (interface{}, error) { return nil, nil },
		"fail":     func(msg string) (string, error) { return "", nil },
		"tpl":      func(tpl string, parentContext interface{}) (string, error) { return "", nil },
	}

	for key, value := range fmap {
		f[key] = value
	}

	return f
}
