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
	"strings"
)

func joinPath(path string, segments ...string) string {
	joint := strings.Join(segments, ".")
	joint = strings.TrimLeft(joint, ".")
	path = fmt.Sprintf("%s.%s", path, joint)
	path = strings.TrimRight(path, ".")
	return path
}

func MakeUniform(paths map[string]struct{}) map[string]struct{} {
	results := map[string]struct{}{}

	for path := range paths {
		sections := strings.Split(path, ".")
		buildPath := ""
	SectionLoop:
		for _, section := range sections {
			if section == "" {
				continue SectionLoop
			}

			if strings.HasSuffix(section, "[*]") {
				extraBuildPath := joinPath(buildPath, strings.TrimSuffix(section, "[*]"))
				results[extraBuildPath] = struct{}{}
			}

			buildPath = joinPath(buildPath, section)
			results[buildPath] = struct{}{}
		}
	}

	return results
}

// Empty is public since it is used by some internal API objects for conversions between external
// string arrays and internal sets, and conversion logic requires public types today.
type Empty struct{}

// Set is a set of the same type elements, implemented via map[comparable]struct{} for minimal memory consumption.
type Set[T comparable] map[T]Empty

// Insert adds items to the set.
func (s Set[T]) Insert(items ...T) {
	for _, item := range items {
		s[item] = Empty{}
	}
}

// Has returns true if the set contains the items.
func (s Set[T]) Has(items ...T) bool {
	for _, item := range items {
		if _, ok := s[item]; !ok {
			return false
		}
	}
	return true
}

// Delete removes items from the set.
func (s Set[T]) Delete(items ...T) {
	for _, item := range items {
		delete(s, item)
	}
}

func getSet[T comparable, V comparable](m map[T]Set[V], k T) Set[V] {
	if v, ok := m[k]; ok {
		return v
	}
	zero := Set[V]{}
	m[k] = zero
	return zero
}
