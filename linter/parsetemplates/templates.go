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
	"github.com/cert-manager/helm-tool/linter/sets"
)

func ListTemplatePaths(templatesPath string) (sets.Set[string], error) {
	tmpl := template.New("ROOT")

	tmpl.Funcs(funcs_serdes.FuncMap())

	templates := sets.Set[*template.Template]{}

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

	return ListTemplatePathsFromTemplates(tmpl, templates)
}

func joinPath(path string, segments ...string) string {
	joint := strings.Join(segments, ".")
	joint = strings.TrimLeft(joint, ".")
	path = fmt.Sprintf("%s.%s", path, joint)
	path = strings.TrimRight(path, ".")
	return path
}

func getSet[T comparable, V comparable](m map[T]sets.Set[V], k T) sets.Set[V] {
	if v, ok := m[k]; ok {
		return v
	}
	zero := sets.Set[V]{}
	m[k] = zero
	return zero
}

type templateUsage struct {
	node    string
	context string
}

const (
	RootNode = "<root-node>"
	RootPath = "<root-path>"
)

func ListTemplatePathsFromTemplates(
	tmpl *template.Template,
	templates sets.Set[*template.Template],
) (sets.Set[string], error) {
	// templateResults lists all property paths that are used in a template
	templateResults := map[string]sets.Set[string]{}
	// templateUsage lists all templates that are used in a template
	templateUsages := map[string]sets.Set[templateUsage]{}
	for _, t := range tmpl.Templates() {
		var selfNode, selfPath string
		if _, ok := templates[t]; ok {
			// we found a template that is the output of a
			// parsed file in the templates folder
			selfNode = RootNode
			selfPath = RootPath
		} else {
			// we found a template that is NOT the output
			// of a parsed file in the templates folder
			selfNode = t.Name()
			selfPath = ""
		}

		results := getSet(templateResults, selfNode)
		usages := getSet(templateUsages, selfNode)

		walk(t.Root, selfNode, selfPath,
			// Found path to a value
			func(path string) {
				results.Insert(path)
			},
			// Found template call
			func(templateName string, context string) {
				usages.Insert(templateUsage{templateName, context})
			},
			// Found local variable usage
			func(varname string, path string) {
				getSet(templateResults, joinPath(selfNode, varname)).Insert(path)
			},
			// Found local variable definition
			func(varname string, node, path string) {
				getSet(templateUsages, node).Insert(templateUsage{joinPath(selfNode, varname), path})
			},
		)
	}

	followPath(RootNode, sets.Set[string]{}, templateUsages, func(node, path string) {
		for key := range templateResults[node] {
			if strings.HasPrefix(key, RootPath) {
				templateResults[RootNode].Insert(key)
			} else {
				templateResults[RootNode].Insert(joinPath(path, key))
			}
		}
	})

	paths := sets.Set[string]{}
	for key := range templateResults[RootNode] {
		if !strings.HasPrefix(key, joinPath(RootPath, "Values")+".") {
			continue
		}
		paths.Insert(strings.TrimPrefix(key, joinPath(RootPath, "Values")+"."))
	}

	return sets.RemovePrefixes(paths), nil
}

func followPath(
	node string,
	visited sets.Set[string],
	templateUsage map[string]sets.Set[templateUsage],
	run func(node string, path string),
) {
	if visited.Has(node) {
		return
	}
	visited.Insert(node)

	for usage := range templateUsage[node] {
		// Recursively follow the path until we reach the <root-node>
		followPath(usage.node, visited, templateUsage, func(node, path string) {
			run(node, joinPath(usage.context, path))
		})
	}
	run(node, "")

	visited.Delete(node)
}

func walk(
	node parse.Node,
	parentNode string,
	parentPath string,
	// foundPathFn is called when a used path is found
	foundPathFn func(path string),
	// foundTemplateFn is called when a template-like call is found
	foundTemplateFn func(templateName string, context string),
	// foundVarUsageFn is called when a variable is used
	foundVarUsageFn func(varname string, path string),
	// foundVarDefFn is called when a variable is defined
	foundVarDefFn func(varname string, node, path string),
) {
	if node == nil || reflect.ValueOf(node).IsNil() {
		return
	}

	switch tn := node.(type) {
	case *parse.DotNode:
		foundPathFn(parentPath)
	case *parse.FieldNode:
		if len(tn.Ident) == 0 {
		} else if tn.Ident[0] == "$" {
			foundPathFn(joinPath(RootPath, tn.Ident[1:]...))
		} else {
			foundPathFn(joinPath(parentPath, tn.Ident...))
		}
		return
	case *parse.VariableNode:
		if len(tn.Ident) == 0 {
		} else if tn.Ident[0] == "$" {
			foundPathFn(joinPath(RootPath, tn.Ident[1:]...))
		} else if len(tn.Ident[0]) >= 1 && tn.Ident[0][0] == '$' {
			foundVarUsageFn(tn.Ident[0], joinPath("", tn.Ident[1:]...))
		} else {
			foundPathFn(joinPath(parentPath, tn.Ident...))
		}
		return
	}

	switch tn := node.(type) {
	case *parse.ActionNode:
		walk(tn.Pipe, parentNode, parentPath, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
	case *parse.ChainNode:
		walk(tn.Node, parentNode, parentPath, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
	case *parse.CommandNode:
		// handle 'include "test.labels" .' separately
		if len(tn.Args) >= 3 && tn.Args[0].String() == "include" && tn.Args[1].Type() == parse.NodeString {
			foundTemplateFn(
				tn.Args[1].(*parse.StringNode).Text,
				getPath(tn.Args[2], parentNode, parentPath),
			)
		}
		for _, snode := range tn.Args {
			walk(snode, parentNode, parentPath, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
		}
	case *parse.BranchNode:
		walk(tn.Pipe, parentNode, parentPath, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
		walk(tn.List, parentNode, parentPath, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
		walk(tn.ElseList, parentNode, parentPath, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
	case *parse.ListNode:
		for _, snode := range tn.Nodes {
			walk(snode, parentNode, parentPath, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
		}
	case *parse.PipeNode:
		for _, cmd := range tn.Cmds {
			walk(cmd, parentNode, parentPath, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
		}

		for _, decl := range tn.Decl {
			for _, cmd := range tn.Cmds {
				walk(cmd, parentNode, parentPath,
					func(path string) {
						foundVarDefFn(decl.String(), parentNode, path)
					},
					func(templateName, context string) {
						foundVarDefFn(decl.String(), templateName, context)
					},
					func(varname string, path string) {
						foundVarDefFn(decl.String(), joinPath(parentNode, varname), path)
					},
					func(varname string, node, path string) {
						// TODO: maybe do something here?
					},
				)

			}
		}
	case *parse.TemplateNode:
		walk(tn.Pipe, parentNode, parentPath,
			func(path string) {
				foundTemplateFn(tn.Name, path)
			},
			func(templateName, context string) {},
			func(varname string, path string) {},
			func(varname string, node, path string) {},
		)
	case *parse.IfNode:
		walk(&tn.BranchNode, parentNode, parentPath, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
	case *parse.RangeNode:
		walk(tn.Pipe, parentNode, parentPath,
			func(path string) {
				foundPathFn(path + "[*]")
			},
			func(templateName, context string) {
				foundTemplateFn(templateName, context+"[*]")
			},
			func(varname string, path string) {
				foundVarUsageFn(varname, path+"[*]")
			},
			func(varname string, node, path string) {
				foundVarDefFn(varname, node, path+"[*]")
			},
		)
		path := getPath(tn.Pipe, parentNode, parentPath) + "[*]"
		walk(tn.List, parentNode, path, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
		walk(tn.ElseList, parentNode, path, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
	case *parse.WithNode:
		path := getPath(tn.Pipe, parentNode, parentPath)
		walk(tn.List, parentNode, path, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
		walk(tn.ElseList, parentNode, path, foundPathFn, foundTemplateFn, foundVarUsageFn, foundVarDefFn)
	}
}

func getPath(node parse.Node, parentNode string, parentPath string) string {
	longestPath := parentPath
	walk(node, parentNode, parentPath,
		func(path string) {
			if len(path) > len(longestPath) {
				longestPath = path
			}
		},
		func(_, context string) {
			if len(context) > len(longestPath) {
				longestPath = context
			}
		},
		func(varname string, path string) {},
		func(varname string, node, path string) {},
	)
	return longestPath
}
