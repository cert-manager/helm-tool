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

package parsetemplates_test

import (
	"fmt"
	"sort"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"

	"github.com/cert-manager/helm-tool/linter/parsetemplates"
	"github.com/cert-manager/helm-tool/linter/parsetemplates/funcs_serdes"
	"github.com/cert-manager/helm-tool/linter/sets"
)

func TestListTemplatePathsFromTemplates(t *testing.T) {
	type testcase struct {
		templates     []string
		expectedPaths []string
	}

	testcases := []testcase{
		{
			templates: []string{
				"{{ .Values.foo }}",
			},
			expectedPaths: []string{
				"foo",
			},
		},
		{
			templates: []string{
				"{{ .Values.foo }}",
				"{{ .Values.bar }}",
			},
			expectedPaths: []string{
				"foo",
				"bar",
			},
		},
		{
			templates: []string{
				"{{ .Values.foo }}",
				"{{ .Values.foo }}",
			},
			expectedPaths: []string{
				"foo",
			},
		},
		{
			templates: []string{
				"{{ .Values.foo }}",
				"{{ .Values.foo }}",
				"{{ .Values.bar }}",
			},
			expectedPaths: []string{
				"foo",
				"bar",
			},
		},
		{
			templates: []string{
				"{{ .foo }}",
				"{{ .Values.bar }}",
			},
			expectedPaths: []string{
				"bar",
			},
		},
		{
			templates: []string{
				"{{ range $key, $value := .Values.test }}{{ end }}",
			},
			expectedPaths: []string{
				"test[*]",
			},
		},
		{
			templates: []string{
				"{{ $aa := .Values.test1.test2 }}",
			},
			expectedPaths: []string{
				"test1.test2",
			},
		},
		{
			templates: []string{
				"{{ $aa := .Values.test1 }}{{ $bb := $aa.test2 }}{{ $bb.test3 }}",
			},
			expectedPaths: []string{
				"test1.test2.test3",
			},
		},
		{
			templates: []string{
				"{{ $value := .Values.test }}{{ $value.value }}",
			},
			expectedPaths: []string{
				"test.value",
			},
		},
		{
			templates: []string{
				"{{ $value := .Values.test1 }}{{ $value := .Values.test2 }}{{ $value.value }}",
			},
			expectedPaths: []string{
				"test1.value",
				"test2.value",
			},
		},
		{
			templates: []string{
				"{{ range $key, $value := .Values.test }}{{ $key.key }}{{ $value.value.test1 }}{{ end }}",
			},
			expectedPaths: []string{
				"test[*].key",
				"test[*].value.test1",
			},
		},
		{
			templates: []string{
				"{{ with .Values.test1 }}{{ .test2 }}{{ end }}",
			},
			expectedPaths: []string{
				"test1.test2",
			},
		},
		{
			templates: []string{
				"{{ with .Values.test1 }}{{ . }}{{ end }}",
			},
			expectedPaths: []string{
				"test1",
			},
		},
		{
			templates: []string{
				"{{ if .Values.test1 }}{{ . }}{{ .Values.test2 }}{{ end }}",
			},
			expectedPaths: []string{
				"test1",
				"test2",
			},
		},
		{
			templates: []string{
				"{{ if .Values.test1 }}{{ end }}",
			},
			expectedPaths: []string{
				"test1",
			},
		},
		{
			templates: []string{
				"{{define \"T1\" }}{{ .test2 }}{{end}} {{ .Values.foo }}",
				"{{ template \"T1\" .Values.test1 }}",
				"{{ .Values.bar }}",
			},
			expectedPaths: []string{
				"test1.test2",
				"foo",
				"bar",
			},
		},
		{
			templates: []string{
				"{{define \"T1\" }}{{ .test1 }}{{end}}",
				"{{define \"T2\" }}{{ template \"T1\" .test2 }}{{end}}",
				"{{define \"T3\" }}{{ template \"T1\" .test2 }}{{ template \"T2\" .test3 }}{{end}}",
				"{{ template \"T1\" .Values.test1 }}{{ template \"T2\" .Values.test1 }}{{ template \"T3\" .Values.test1 }}",
			},
			expectedPaths: []string{
				"test1.test1",
				"test1.test2.test1",
				"test1.test3.test2.test1",
			},
		},
		{
			templates: []string{
				"{{ $name := default .Values.test1 .Values.test2 }}{{ $name.test3 }}",
			},
			expectedPaths: []string{
				"test1.test3",
				"test2.test3",
			},
		},
		{
			templates: []string{
				"{{ $name := (tuple .Values.test1 .Values.test2) }}{{ $name.test3 }}",
			},
			expectedPaths: []string{
				"test1.test3",
				"test2.test3",
			},
		},
		{
			templates: []string{
				"{{ $name := (list .Values.test1 .Values.test2) }}{{ $name.test3 }}",
			},
			expectedPaths: []string{
				"test1.test3",
				"test2.test3",
			},
		},
		{
			templates: []string{
				"{{ define \"T1\" }}{{ .test3 }}{{ end }}",
				"{{ template \"T1\" (tuple .Values.test1 .Values.test2) }}",
			},
			expectedPaths: []string{
				"test1.test3",
				"test2.test3",
			},
		},
		{
			templates: []string{
				"{{ .Values.app.logLevela }}",
				"{{ .Values.app.name }}",
			},
			expectedPaths: []string{
				"app.logLevela",
				"app.name",
			},
		},
		{
			templates: []string{
				"{{- with (or .Values.test1 .Values.test2) }}{{- toYaml . | nindent 8 }}{{- end }}",
			},
			expectedPaths: []string{
				"test1",
				"test2",
			},
		},
		{
			templates: []string{
				"{{- with (or .Values.test1 .Values.test2) }}aaa{{- end }}",
			},
			expectedPaths: []string{
				"test1",
				"test2",
			},
		},
	}

	for _, tc := range testcases {
		tmpl := template.New("ROOT")

		tmpl.Funcs(funcs_serdes.FuncMap())

		templates := sets.Set[*template.Template]{}
		for idx, tem := range tc.templates {
			tpl, err := tmpl.New(fmt.Sprintf("input-item-%d", idx)).Parse(tem)
			if err != nil {
				t.Errorf("error parsing template: %s", err)
			}

			templates[tpl] = struct{}{}
		}

		paths, err := parsetemplates.ListTemplatePathsFromTemplates(tmpl, templates)
		if err != nil {
			t.Errorf("error listing template paths: %s", err)
		}

		sort.Strings(tc.expectedPaths)

		pathList := paths.UnsortedList()
		sort.Strings(pathList)

		require.EqualValues(t, tc.expectedPaths, pathList)
	}
}
