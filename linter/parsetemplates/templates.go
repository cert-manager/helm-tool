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

package parsetemplates

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/cert-manager/helm-tool/linter/parsetemplates/funcs_serdes"
)

func ListTemplatePaths(templatesPath string) ([]string, error) {
	tmpl := template.New("ROOT")

	tmpl.Funcs(funcs_serdes.FuncMap())

	templates := map[*template.Template]struct{}{}

	// parse all templates
	err := filepath.Walk(
		templatesPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			contents, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			t, err := tmpl.New(path).Parse(string(contents))
			if err != nil {
				return err
			}
			templates[t] = struct{}{}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	templatePathStrings, err := ListTemplatePathsFromTemplates(tmpl, templates)
	if err != nil {
		return nil, err
	}

	// remove prefixes
	templatePaths := make([]string, 0, len(templatePathStrings))
	for _, pathString := range templatePathStrings {
		if !strings.HasPrefix(pathString, ".$") {
			continue
		}

		templatePaths = append(templatePaths, strings.TrimPrefix(pathString, ".$"))
	}

	return templatePaths, nil
}

func ListTemplatePathsFromTemplates(
	tmpl *template.Template,
	templates map[*template.Template]struct{},
) ([]string, error) {
	// walk all templates
	templateResults := map[string]map[string]struct{}{}
	templateUsage := map[string][][]string{}
	for _, t := range tmpl.Templates() {
		root := "$"
		tmplPath := "$"
		if _, ok := templates[t]; !ok {
			// we found a directly parsed template
			root = ""
			tmplPath = t.Name()
		}

		result, ok := templateResults[tmplPath]
		if !ok {
			result = map[string]struct{}{}
			templateResults[tmplPath] = result
		}

		walk(t.Root, root, func(path string, node parse.Node, templateName string) {
			if node != nil {
				result[path] = struct{}{}
			} else {
				templateUsage[tmplPath] = append(templateUsage[tmplPath], []string{templateName, path})
			}
		}, func(varname string, assignName, assign, access *string) {
			if assign != nil {
				varPath := fmt.Sprintf("%s%s", tmplPath, varname)
				assignPath := fmt.Sprintf("%s%s", tmplPath, *assignName)
				templateUsage[assignPath] = append(templateUsage[assignPath], []string{varPath, *assign})
			} else {
				varPath := fmt.Sprintf("%s%s", tmplPath, varname)
				result, ok := templateResults[varPath]
				if !ok {
					result = map[string]struct{}{}
					templateResults[varPath] = result
				}

				result[*access] = struct{}{}
			}
		})
	}

	followPath("$", templateUsage, func(path, tmplPath string) {
		for key := range templateResults[tmplPath] {
			templateResults["$"][fmt.Sprintf("%s%s", path, key)] = struct{}{}
		}
	})

	paths := map[string]struct{}{}
	for key := range templateResults["$"] {
		if !strings.HasPrefix(key, "$.Values") {
			continue
		}
		key = strings.TrimPrefix(key, "$.Values")
		key = "$" + key

		paths[key] = struct{}{}
	}

	paths = MakeUniform(paths)

	// sort paths
	values := []string{}
	for key := range paths {
		values = append(values, key)
	}

	return values, nil
}

func followPath(
	key string,
	templateUsage map[string][][]string,
	run func(path string, tpath string),
) {
	for _, usage := range templateUsage[key] {
		followPath(usage[0], templateUsage, func(path, tpath string) {
			run(fmt.Sprintf("%s%s", usage[1], path), tpath)
		})
	}
	run("", key)
}

func walk(
	node parse.Node,
	parentPath string,
	foundPathFn func(path string, node parse.Node, templateName string),
	foundLocalVar func(varname string, assignPath, assign, access *string),
) {
	if node == nil || reflect.ValueOf(node).IsNil() {
		return
	}

	path := parentPath
	switch tn := node.(type) {
	case *parse.FieldNode:
		if len(tn.Ident) == 0 {
		} else if tn.Ident[0] == "$" {
			path = strings.Join(tn.Ident, ".")
		} else {
			path = fmt.Sprintf("%s.%s", parentPath, strings.Join(tn.Ident, "."))
		}
	case *parse.VariableNode:
		if len(tn.Ident) == 0 {
		} else if tn.Ident[0] == "$" {
			path = strings.Join(tn.Ident, ".")
		} else if len(tn.Ident[0]) >= 1 && tn.Ident[0][0] == '$' {
			path := "." + strings.Join(tn.Ident[1:], ".")
			foundLocalVar(tn.Ident[0], nil, nil, &path)
			return
		} else {
			path = fmt.Sprintf("%s.%s", parentPath, strings.Join(tn.Ident, "."))
		}
	}

	foundPathFn(path, node, "")

	switch tn := node.(type) {
	case *parse.ActionNode:
		walk(tn.Pipe, path, foundPathFn, foundLocalVar)
	case *parse.ChainNode:
		walk(tn.Node, path, foundPathFn, foundLocalVar)
	case *parse.CommandNode:
		// handle 'include "test.labels" .' separately
		if len(tn.Args) >= 3 && tn.Args[0].String() == "include" && tn.Args[1].Type() == parse.NodeString {
			foundPathFn(getPath(tn.Args[2], path), nil, tn.Args[1].(*parse.StringNode).Text)
		}
		for _, snode := range tn.Args {
			walk(snode, path, foundPathFn, foundLocalVar)
		}
	case *parse.BranchNode:
		walk(tn.Pipe, path, foundPathFn, foundLocalVar)
		walk(tn.List, path, foundPathFn, foundLocalVar)
		walk(tn.ElseList, path, foundPathFn, foundLocalVar)
	case *parse.ListNode:
		for _, snode := range tn.Nodes {
			walk(snode, path, foundPathFn, foundLocalVar)
		}
	case *parse.PipeNode:
		for _, cmd := range tn.Cmds {
			walk(cmd, path, foundPathFn, foundLocalVar)
		}

		for _, decl := range tn.Decl {
			for _, cmd := range tn.Cmds {
				walk(cmd, path, func(path string, _ parse.Node, _ string) {
					empty := ""
					foundLocalVar(decl.String(), &empty, &path, nil)
				}, func(varname string, _, _, access *string) {
					if access == nil {
						return
					}
					foundLocalVar(decl.String(), &varname, access, nil)
				})

			}
		}
	case *parse.TemplateNode:
		path := getPath(tn.Pipe, path)
		foundPathFn(path, nil, tn.Name)
	case *parse.IfNode:
		walk(&tn.BranchNode, path, foundPathFn, foundLocalVar)
	case *parse.RangeNode:
		walk(tn.Pipe, path,
			func(path string, node parse.Node, templateName string) {},
			func(varname string, assignPath, assign, access *string) {
				if assign == nil {
					return
				}
				newAssign := *assign + ".[*]"
				foundLocalVar(varname, assignPath, &newAssign, access)
			},
		)
		path := getPath(tn.Pipe, path) + ".[*]"
		walk(tn.List, path, foundPathFn, foundLocalVar)
		walk(tn.ElseList, path, foundPathFn, foundLocalVar)
	case *parse.WithNode:
		path := getPath(tn.Pipe, path)
		walk(tn.List, path, foundPathFn, foundLocalVar)
		walk(tn.ElseList, path, foundPathFn, foundLocalVar)
	}
}

func getPath(node parse.Node, parentPath string) string {
	longestPath := ""
	walk(node, parentPath, func(path string, _ parse.Node, _ string) {
		if len(path) > len(longestPath) {
			longestPath = path
		}
	}, func(varname string, assignPath, assign, access *string) {})
	return longestPath
}
