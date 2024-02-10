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

package sets

import (
	"maps"
)

// Empty is public since it is used by some internal API objects for conversions between external
// string arrays and internal sets, and conversion logic requires public types today.
type Empty struct{}

// Set is a set of the same type elements, implemented via map[comparable]struct{} for minimal memory consumption.
type Set[T comparable] map[T]Empty

// New creates a new Set from a list of values.
func New[T comparable](values ...T) Set[T] {
	s := Set[T]{}
	s.Insert(values...)
	return s
}

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

// UnsortedList returns the elements of the set as a list.
func (s Set[T]) UnsortedList() []T {
	list := make([]T, 0, len(s))
	for item := range s {
		list = append(list, item)
	}
	return list
}

// Union returns a new set with all items from both sets.
func Union[T comparable](sets ...Set[T]) Set[T] {
	s := Set[T]{}
	for _, set := range sets {
		for item := range set {
			s.Insert(item)
		}
	}
	return s
}

// Remove returns a new set with all items from the input set
// that are not in the other set.
func Remove[T comparable](a Set[T], sets ...Set[T]) Set[T] {
	compactedItems := maps.Clone(a)
	for _, set := range sets {
		compactedItems.Delete(set.UnsortedList()...)
	}
	return compactedItems
}
