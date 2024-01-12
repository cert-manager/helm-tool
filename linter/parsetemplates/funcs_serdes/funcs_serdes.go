/*
Copyright 2024 The cert-manager Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
