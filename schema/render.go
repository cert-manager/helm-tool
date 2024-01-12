package schema

import (
	"encoding/json"
	"fmt"

	"github.com/cert-manager/helm-docgen/parser"
	"gopkg.in/yaml.v3"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

type treeLevel struct {
	Path     parser.Path
	Property *parser.Property
	Children []treeLevel
}

func (t *treeLevel) Type() parser.Type {
	if len(t.Children) == 0 && t.Property != nil {
		return t.Property.Type
	}

	if len(t.Children) > 0 {
		firstChild := t.Children[0]
		if parser.IsArrayPathComponent(firstChild.Path.Property()) {
			return parser.TypeArray
		}

		return parser.TypeObject
	}

	return parser.TypeUnknown
}

func (t *treeLevel) add(path parser.Path, property parser.Property) error {
	if path.Equal(t.Path) {
		t.Property = &property
		return nil
	}

	if !t.Path.IsSubPathOf(path) {
		return fmt.Errorf("path %q is not a subpath of %q", t.Path, path)
	}

	for i, child := range t.Children {
		if child.Path.IsSubPathOf(path) {
			child.add(path, property)
			t.Children[i] = child
			return nil
		}
	}

	t.Children = append(t.Children, treeLevel{Path: t.Path.Expand(path, 1)})
	t.Children[len(t.Children)-1].add(path, property)
	return nil
}

func (t *treeLevel) walk(f func(level treeLevel)) {
	f(*t)
	for _, child := range t.Children {
		child.walk(f)
	}
}

func buildTree(document *parser.Document) (treeLevel, error) {
	allProperties := []parser.Property{}
	for _, section := range document.Sections {
		allProperties = append(allProperties, section.Properties...)
	}

	root := treeLevel{}
	for _, property := range allProperties {
		if err := root.add(property.Path, property); err != nil {
			return treeLevel{}, err
		}
	}

	return root, nil
}

func Render(document *parser.Document) (string, error) {
	tree, err := buildTree(document)
	if err != nil {
		return "", err
	}

	definitions := spec.Definitions{}

	tree.walk(func(level treeLevel) {
		levelType := level.Type()

		newSchema := spec.Schema{SchemaProps: spec.SchemaProps{}}

		schemaType := levelType.SchemaString()
		if len(schemaType) > 0 {
			newSchema.SchemaProps.Type = []string{schemaType}
		}

		if level.Property != nil {
			newSchema.SchemaProps.Description = level.Property.Description.String()

			if level.Property.Default != "" {
				var defaultValue interface{}
				if err := yaml.Unmarshal([]byte(level.Property.Default), &defaultValue); err != nil {
					panic(err)
				}
				newSchema.SchemaProps.Default = defaultValue
			}
		}

		switch levelType {
		case parser.TypeArray:
			itemSchema := spec.Schema{SchemaProps: spec.SchemaProps{}}

			if len(level.Children) > 0 {
				firstChild := level.Children[0]
				itemSchema.SchemaProps.Ref = spec.MustCreateRef(fmt.Sprintf("#/$defs/%s", prefixName(firstChild.Path.String())))
			}

			newSchema.SchemaProps.Items = &spec.SchemaOrArray{Schema: &itemSchema}

		case parser.TypeObject:
			properties := map[string]spec.Schema{}

			for _, child := range level.Children {
				properties[parser.SegmentString(child.Path.Property())] = spec.Schema{SchemaProps: spec.SchemaProps{
					Ref: spec.MustCreateRef(fmt.Sprintf("#/$defs/%s", prefixName(child.Path.String()))),
				}}
			}

			newSchema.SchemaProps.Properties = properties
			if len(level.Children) > 0 {
				newSchema.SchemaProps.AdditionalProperties = &spec.SchemaOrBool{Allows: false}
			}
		}

		definitions[prefixName(level.Path.String())] = newSchema
	})

	type JsonSchema struct {
		Schema string           `json:"$schema,omitempty"`
		Ref    string           `json:"$ref,omitempty"`
		Defs   spec.Definitions `json:"$defs,omitempty"`
	}

	data, err := json.Marshal(JsonSchema{
		Schema: "http://json-schema.org/draft-07/schema#",
		Defs:   definitions,
		Ref:    "#/$defs/helm-values",
	})
	if err != nil {
		return "", fmt.Errorf("error serializing api definitions: %w", err)
	}

	return string(data), nil
}

func prefixName(name string) string {
	if name == "" {
		return "helm-values"
	}

	return "helm-values." + name
}
