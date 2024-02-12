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
	"slices"
	"strings"
)

// RemovePrefixes returns a new set with all items from
// the input set that are not a prefix of another item in
// the set or any of the additional sets. We consider a
// string a prefix of another string if the other string
// starts with the first string followed by a period or
// an opening square bracket.
func RemovePrefixes(items Set[string], sets ...Set[string]) Set[string] {
	nonPrefixes := maps.Clone(items)

	values := Union(append(sets, items)...).UnsortedList()
	slices.Sort(values)
	for i := 0; i < len(values); i++ {
		// If the next value is an extension of the current value, skip
		// the current value.
		if i+1 < len(values) && (strings.HasPrefix(values[i+1], values[i]+".") || strings.HasPrefix(values[i+1], values[i]+"[")) {
			nonPrefixes.Delete(values[i])
		}
	}

	return nonPrefixes
}

// RemoveExtensions returns a new set with all items from
// the input set that are not an extension of another item
// in the set or any of the additional sets. We consider a
// string a prefix of another string if the other string
// starts with the first string followed by a period or
// an opening square bracket.
func RemoveExtensions(items Set[string], sets ...Set[string]) Set[string] {
	nonExtensions := maps.Clone(items)

	values := Union(append(sets, items)...).UnsortedList()
	slices.Sort(values)
	for i := 0; i < len(values); i++ {
		// If the next value is an extension of the current value, skip
		// the current value.
		if i+1 < len(values) && (strings.HasPrefix(values[i+1], values[i]+".") || strings.HasPrefix(values[i+1], values[i]+"[")) {
			nonExtensions.Delete(values[i+1])
		}
	}

	return nonExtensions
}
