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
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemovePrefixes(t *testing.T) {
	tests := []struct {
		input      Set[string]
		additional Set[string]
		expected   Set[string]
	}{
		{
			input:      New("ab", "abd", "ab.d"),
			additional: New("a", "ab"),
			expected:   New("abd", "ab.d"),
		},
		{
			input:      New("a", "a.b", "a.d"),
			additional: New("a.b.c"),
			expected:   New("a.d"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			result := RemovePrefixes(tt.input, tt.additional)
			require.ElementsMatch(t, tt.expected.UnsortedList(), result.UnsortedList())
		})
	}
}

func TestRemoveExtensions(t *testing.T) {
	tests := []struct {
		input      Set[string]
		additional Set[string]
		expected   Set[string]
	}{
		{
			input:      New("a.b.c"),
			additional: New("a.b"),
			expected:   New("a.b.c"),
		},
		{
			input:      New("a.b.c", "a.b", "a.d"),
			additional: New("a.b"),
			expected:   New("a.b", "a.d"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			result := RemoveExtensions(tt.input)
			require.ElementsMatch(t, tt.expected.UnsortedList(), result.UnsortedList())
		})
	}

}
