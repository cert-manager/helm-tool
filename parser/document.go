package parser

import (
	"log"
	"os"
	"strings"

	"github.com/cert-manager/helm-docgen/heuristics"
	"gopkg.in/yaml.v3"
)

const (
	TagSection  = "docs:section"
	TagIgnore   = "docs:ignore"
	TagType     = "docs:type"
	TagDefault  = "docs:default"
	TagProperty = "docs:property"
)

type Document struct {
	Sections []Section
}

type Section struct {
	Name        string
	Description Comment
	Properties  []Property
}

type Property struct {
	Path        Path
	Description Comment
	Type        Type
	Default     string
}

type Type string

const (
	TypeUnknown   Type = "unknown"
	TypeString    Type = "string"
	TypeNumber    Type = "number"
	TypeBool      Type = "bool"
	TypeTimestamp Type = "timestamp"
	TypeArray     Type = "array"
	TypeObject    Type = "object"
)

func (t Type) String() string {
	return string(t)
}

func (t Type) SchemaString() string {
	switch t {
	case TypeString, TypeNumber, TypeArray, TypeObject:
		return string(t)

	case TypeBool:
		return "boolean"

	case TypeTimestamp:
		return "string"

	default:
		return ""
	}
}

type Node struct {
	Path         Path
	HeadComments []Comment
	FootComment  []Comment
	RawNode      *yaml.Node
}

func Load(filename string) (*Document, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var root yaml.Node
	if err := yaml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	document := Document{Sections: make([]Section, 1)}
	node := Node{
		RawNode:      &root,
		HeadComments: parseComments(root.HeadComment),
		FootComment:  parseComments(root.FootComment),
	}
	err = walk(node, func(node Node) (bool, error) {
		comment := pop(&node.HeadComments)

		parseCommentsOntoDocument(node.Path.Parent(), &document, node.HeadComments)
		defer parseCommentsOntoDocument(node.Path.Parent(), &document, node.FootComment)

		// If we have a comment instructing us to skip this node, obey it
		if comment.Tags.GetBool(TagIgnore) {
			return true, nil
		}

		// An end node is a node we find a property at, this is usually a scalar
		// node, but can be a map or sequence if the user uses the
		// +docs:property tag (or if they have no values).
		if !isEndNode(node, comment) {
			parseCommentsOntoDocument(node.Path.Parent(), &document, []Comment{comment})
			return false, nil
		}

		sectionIdx := len(document.Sections) - 1
		document.Sections[sectionIdx].Properties = append(document.Sections[sectionIdx].Properties, Property{
			Path:        node.Path,
			Description: comment,
			Type:        getTypeOf(node, comment),
			Default:     getDefaultValue(node, comment),
		})

		return true, nil
	})

	return &document, err
}

func parseCommentsOntoDocument(path Path, document *Document, comments []Comment) {
	for _, comment := range comments {
		switch {
		case comment.Tags.GetBool(TagSection):
			document.Sections = append(document.Sections, Section{
				Name:        comment.Tags.GetString(TagSection),
				Description: comment,
			})
		case comment.Tags.GetBool(TagProperty):
			// Search for a code block in the comments, we can try and infer
			// information from it
			codeIdx := -1
			for i, segment := range comment.Segments {
				if segment.Type == heuristics.ContentTypeYaml {
					codeIdx = i
				}
			}

			parsedNode := Node{
				HeadComments: []Comment{comment},
			}

			if codeIdx != -1 {
				parsedSuccessfully := false

				codeSegment := comment.Segments[codeIdx]
				var node yaml.Node
				yaml.Unmarshal([]byte(codeSegment.String()), &node)

				// Document node
				if len(node.Content) != 0 {
					// Mapping node
					if node.Content[0].Kind == yaml.MappingNode {
						// Ensure single value
						if len(node.Content[0].Content) == 2 {
							keyNode := node.Content[0].Content[0]
							valueNode := node.Content[0].Content[1]
							parsedNode.Path = path.WithProperty(keyNode.Value)
							parsedNode.RawNode = valueNode
							parsedSuccessfully = true
						}
					}
				}

				// Remove the code block from the comment
				if parsedSuccessfully {
					newComment := Comment{Tags: comment.Tags}
					for i, segment := range comment.Segments {
						if i == codeIdx {
							continue
						}

						newComment.Segments = append(newComment.Segments, segment)
					}
					comment = newComment
				}
			}

			// If we cant calculate the path, we should warn
			name := comment.Tags.GetString(TagProperty)
			if name == "" {
				name = parsedNode.Path.String()
				if name == "" {
					log.Println("could not calculate undefined property name")
					continue

				}
			}

			path, err := ParsePath(name)
			if err != nil {
				log.Printf("could not parse property path %q: %s\n", name, err)
				continue
			}

			sectionIdx := len(document.Sections) - 1
			document.Sections[sectionIdx].Properties = append(document.Sections[sectionIdx].Properties, Property{
				Path:        path,
				Description: comment,
				Type:        getTypeOf(parsedNode, comment),
				Default:     "",
			})
		}

	}
}

func walk(root Node, fn func(node Node) (bool, error)) error {
	// Call the function for every node, we the method can decide to stop
	// walking this branch as part of this call
	stop, err := fn(root)
	if err != nil {
		return err
	}

	if stop {
		return nil
	}

	// For any node type that nests further nodes, recurse the walk function
	switch root.RawNode.Kind {
	case yaml.SequenceNode:
		for i, node := range root.RawNode.Content {
			n := Node{
				Path:         root.Path.WithIndex(i),
				HeadComments: parseComments(root.RawNode.HeadComment),
				FootComment:  parseComments(root.RawNode.FootComment),
				RawNode:      node,
			}

			if err := walk(n, fn); err != nil {
				return err
			}
		}
	case yaml.MappingNode:
		for i := 0; i < len(root.RawNode.Content); i += 2 {
			keyNode := root.RawNode.Content[i]
			valueNode := root.RawNode.Content[i+1]

			n := Node{
				Path:         root.Path.WithProperty(keyNode.Value),
				HeadComments: parseComments(keyNode.HeadComment),
				FootComment:  parseComments(keyNode.FootComment),
				RawNode:      valueNode,
			}

			if err := walk(n, fn); err != nil {
				return err
			}
		}
	case yaml.DocumentNode:
		for _, node := range root.RawNode.Content {
			n := Node{
				Path:         root.Path,
				RawNode:      node,
				HeadComments: parseComments(node.HeadComment),
				FootComment:  parseComments(node.FootComment),
			}

			if err := walk(n, fn); err != nil {
				return err
			}
		}
	case yaml.AliasNode:
		n := Node{
			Path:         root.Path,
			HeadComments: parseComments(root.RawNode.HeadComment),
			FootComment:  parseComments(root.RawNode.FootComment),
			RawNode:      root.RawNode.Alias,
		}

		if err := walk(n, fn); err != nil {
			return err
		}
	}

	return nil
}

// isEndNode returns true if the yaml node is considered one that should
// be documented as a parameter.
//
// This could be because its a node containing a scalar value, an empty map or
// array, or the user may have used the +docs:param tag to specify the node
// as a parameter.
func isEndNode(n Node, c Comment) bool {
	switch {
	case n.RawNode.Kind == yaml.DocumentNode:
		return false
	case n.RawNode.Kind == yaml.ScalarNode:
		return true
	case c.Tags.GetBool(TagProperty):
		return true
	case n.RawNode.Kind == yaml.MappingNode:
		return len(n.RawNode.Content) == 0
	case n.RawNode.Kind == yaml.SequenceNode:
		return len(n.RawNode.Content) == 0
	default:
		return false
	}
}

func getDefaultValue(n Node, c Comment) string {
	if def := c.Tags.GetString(TagDefault); def != "" {
		return def
	}

	// "clean" the object by parsing to an object and back
	var value any
	var clone yaml.Node
	n.RawNode.Decode(&value)
	clone.Encode(&value)

	// Encode into a string
	var sb strings.Builder
	encoder := yaml.NewEncoder(&sb)
	encoder.SetIndent(2)
	encoder.Encode(clone)
	return strings.TrimSpace(sb.String())
}

// Remove the last element from a slice and
// return it
func pop[T any](s *[]T) T {
	var def T
	l := len(*s)
	if l == 0 {
		return def
	}

	v := (*s)[l-1]
	*s = (*s)[:l-1]

	return v
}

func getTypeOf(node Node, comment Comment) Type {
	if typ := comment.Tags.GetString(TagType); typ != "" {
		return Type(typ)
	}

	if node.RawNode == nil {
		return TypeUnknown
	}

	switch node.RawNode.ShortTag() {
	case "!!bool":
		return TypeBool
	case "!!str":
		return TypeString
	case "!!int":
		return TypeNumber
	case "!!float":
		return TypeNumber
	case "!!timestamp":
		return TypeTimestamp
	case "!!seq":
		return TypeArray
	case "!!map":
		return TypeObject
	default:
		return TypeUnknown
	}
}
