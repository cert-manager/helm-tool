/*
Copyright 2021 The cert-manager Authors.

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

package parser

import (
	"log"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/cert-manager/helm-tool/heuristics"
	"github.com/cert-manager/helm-tool/paths"
)

const (
	TagSection  = "docs:section"
	TagIgnore   = "docs:ignore"
	TagHidden   = "docs:hidden"
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
	Path        paths.Path
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
	Path         paths.Path
	HeadComments []Comment
	FootComment  []Comment
	RawNode      *yaml.Node
}

// sharedNodeFrame accumulates the properties documented while a shared
// (non-end) node's subtree is being walked for the first time, relative to
// that node's own path. Once the walk of that subtree finishes, the frame's
// properties are cached so a later encounter of the same *yaml.Node — via a
// different alias or "<<" merge site — can replay them at its own path
// instead of either silently dropping them or re-walking the whole subtree
// again.
type sharedNodeFrame struct {
	root       paths.Path
	properties []Property
}

func Load(filename string, includeHidden bool) (*Document, error) {
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

	// memo caches the properties documented under a shared node's subtree,
	// relative to the path it was first reached at, keyed by the underlying
	// *yaml.Node so later encounters (via aliases or "<<" merge keys) can
	// replay them rebased onto their own path.
	memo := map[*yaml.Node][]Property{}
	// inProgress guards against a node that, directly or transitively,
	// refers to itself (CWE-674 stack overflow) — it is only ever true while
	// that node's subtree is actively being walked for the first time.
	inProgress := map[*yaml.Node]bool{}
	// frames is a stack of shared nodes currently being walked for the first
	// time, so any property found anywhere below them gets recorded
	// (relative to each) for later replay. Diamond-shaped reuse (the same
	// node reached from several sibling sites) is fine and expected; only a
	// true cycle — re-entering a node already in this stack — is refused.
	var frames []*sharedNodeFrame

	// record appends a property to the document's current section and, for
	// every shared node currently being walked for the first time, also
	// records it (relative to that node's own path) for replay at future
	// encounters.
	record := func(prop Property) {
		sectionIdx := len(document.Sections) - 1
		document.Sections[sectionIdx].Properties = append(document.Sections[sectionIdx].Properties, prop)

		for _, f := range frames {
			if len(prop.Path) < len(f.root) {
				continue
			}
			rel := prop
			rel.Path = append(paths.Path{}, prop.Path[len(f.root):]...)
			f.properties = append(f.properties, rel)
		}
	}

	err = walk(node, func(node Node) (bool, func(), error) {
		comment := pop(&node.HeadComments)

		parseCommentsOntoDocument(node.Path.Parent(), &document, node.HeadComments)
		defer parseCommentsOntoDocument(node.Path.Parent(), &document, node.FootComment)

		// If we have a comment instructing us to skip this node, obey it
		if comment.Tags.GetBool(TagIgnore) {
			return true, nil, nil
		}

		// If we have a comment instructing us to hide this node, obey it if we are not including hidden nodes
		if comment.Tags.GetBool(TagHidden) && !includeHidden {
			return true, nil, nil
		}

		// An end node is a node we find a property at, this is usually a scalar
		// node, but can be a map or sequence if the user uses the
		// +docs:property tag (or if they have no values).
		if !isEndNode(node, comment) {
			parseCommentsOntoDocument(node.Path.Parent(), &document, []Comment{comment})

			// We've fully documented this exact node before, via a different
			// alias or "<<" merge site — replay what we found there instead
			// of silently dropping it.
			if cached, ok := memo[node.RawNode]; ok {
				for _, p := range cached {
					p.Path = append(append(paths.Path{}, node.Path...), p.Path...)
					record(p)
				}
				return true, nil, nil
			}

			// This node is already being walked further up the current call
			// stack — a genuine cycle (e.g. a mapping merging itself).
			// Stop here rather than recursing forever.
			if inProgress[node.RawNode] {
				return true, nil, nil
			}

			inProgress[node.RawNode] = true
			frame := &sharedNodeFrame{root: node.Path}
			frames = append(frames, frame)

			after := func() {
				memo[node.RawNode] = frame.properties
				delete(inProgress, node.RawNode)
				frames = frames[:len(frames)-1]
			}

			return false, after, nil
		}

		record(Property{
			Path:        node.Path,
			Description: comment,
			Type:        getTypeOf(node, comment),
			Default:     getDefaultValue(node, comment),
		})

		return true, nil, nil
	})

	return &document, err
}

func parseCommentsOntoDocument(path paths.Path, document *Document, comments []Comment) {
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

			path, err := paths.Parse(name)
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

// walk performs a depth-first traversal of a yaml node tree, calling fn for
// every node encountered. fn reports whether to stop descending into this
// node's children and may optionally return an "after" function, which is
// called once this node and all of its descendants have been fully visited —
// this lets the caller pair up per-node setup (e.g. entering a subtree for
// the first time) with a matching cleanup step.
func walk(root Node, fn func(node Node) (stop bool, after func(), err error)) error {
	// Call the function for every node, we the method can decide to stop
	// walking this branch as part of this call
	stop, after, err := fn(root)
	if err != nil {
		return err
	}

	if after != nil {
		defer after()
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
		for _, entry := range mappingEntries(root.RawNode, map[*yaml.Node]bool{}) {
			n := Node{
				Path:         root.Path.WithProperty(entry.Key.Value),
				HeadComments: parseComments(entry.Key.HeadComment),
				FootComment:  parseComments(entry.Key.FootComment),
				RawNode:      entry.Value,
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

// mappingEntry is a resolved key/value pair from a mapping node, after
// expanding any "<<" merge keys.
type mappingEntry struct {
	Key   *yaml.Node
	Value *yaml.Node
}

// isMergeKey returns true if n is a "<<" merge key, per the YAML merge key
// convention (not part of the core YAML spec, but widely supported).
func isMergeKey(n *yaml.Node) bool {
	return n.Kind == yaml.ScalarNode && n.Value == "<<" && (n.Tag == "" || n.Tag == "!" || n.ShortTag() == "!!merge")
}

// resolveMergeTargets expands a merge key's value into the mapping nodes it
// refers to. The value may be a single mapping, an alias to one, or a
// sequence of either (for merging in more than one mapping at once).
func resolveMergeTargets(n *yaml.Node) []*yaml.Node {
	switch n.Kind {
	case yaml.AliasNode:
		return resolveMergeTargets(n.Alias)
	case yaml.MappingNode:
		return []*yaml.Node{n}
	case yaml.SequenceNode:
		var out []*yaml.Node
		for _, item := range n.Content {
			out = append(out, resolveMergeTargets(item)...)
		}
		return out
	default:
		// Not a valid merge target, e.g. an alias to a scalar. Ignore it
		// rather than failing the whole document.
		return nil
	}
}

// mappingEntries returns a mapping node's effective key/value pairs with any
// "<<" merge keys expanded in place. Explicit keys always take precedence
// over merged-in ones regardless of declaration order, and where more than
// one merge source defines the same key, the earliest one wins — both match
// the YAML merge key convention.
//
// inProgress guards against a merge target that (directly or transitively)
// merges itself, which would otherwise recurse forever.
func mappingEntries(n *yaml.Node, inProgress map[*yaml.Node]bool) []mappingEntry {
	if inProgress[n] {
		return nil
	}
	inProgress[n] = true
	defer delete(inProgress, n)

	explicit := map[string]bool{}
	for i := 0; i < len(n.Content); i += 2 {
		if key := n.Content[i]; !isMergeKey(key) {
			explicit[key.Value] = true
		}
	}

	seen := map[string]bool{}
	var entries []mappingEntry

	for i := 0; i < len(n.Content); i += 2 {
		key, value := n.Content[i], n.Content[i+1]

		if !isMergeKey(key) {
			if !seen[key.Value] {
				seen[key.Value] = true
				entries = append(entries, mappingEntry{Key: key, Value: value})
			}
			continue
		}

		for _, target := range resolveMergeTargets(value) {
			for _, entry := range mappingEntries(target, inProgress) {
				if explicit[entry.Key.Value] || seen[entry.Key.Value] {
					continue
				}
				seen[entry.Key.Value] = true
				entries = append(entries, entry)
			}
		}
	}

	return entries
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
