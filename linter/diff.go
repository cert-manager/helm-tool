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

package linter

import (
	"github.com/cert-manager/helm-tool/linter/sets"
)

// DiffPaths returns the paths that are missing from each set.
// We consider a path to be missing if it is not present in the
// other set, and if it is not a prefix or an extension of another
// path in the other set. We consider a string a prefix of another
// string if the other string starts with the first string followed
// by a period or an opening square bracket.
func DiffPaths(pathsA sets.Set[string], pathsB sets.Set[string]) (sets.Set[string], sets.Set[string]) {
	pathsA, pathsB = sets.RemovePrefixes(pathsA), sets.RemovePrefixes(pathsB)

	missingA := sets.Remove(pathsB, pathsA)
	missingA = sets.RemovePrefixes(missingA, pathsA)
	missingA = sets.RemoveExtensions(missingA, pathsA)

	missingB := sets.Remove(pathsA, pathsB)
	missingB = sets.RemovePrefixes(missingB, pathsB)
	missingB = sets.RemoveExtensions(missingB, pathsB)

	return missingA, missingB
}
