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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiffPaths(t *testing.T) {
	type testcase struct {
		a            []string
		b            []string
		wantMissingA []string
		wantMissingB []string
	}

	testcases := []testcase{
		{
			a:            []string{".$", ".$.a", ".$.b", ".$.c"},
			b:            []string{".$", ".$.a", ".$.b", ".$.c"},
			wantMissingA: []string{},
			wantMissingB: []string{},
		},
		{
			a:            []string{".$", ".$.a", ".$.b", ".$.c"},
			b:            []string{".$", ".$.a", ".$.b"},
			wantMissingA: []string{},
			wantMissingB: []string{".$.c"},
		},
		{
			a:            []string{".$", ".$.a", ".$.b"},
			b:            []string{".$", ".$.a", ".$.b", ".$.c"},
			wantMissingA: []string{".$.c"},
			wantMissingB: []string{},
		},
		{
			a:            []string{".$", ".$.a", ".$.a.b"},
			b:            []string{".$", ".$.a"},
			wantMissingA: []string{},
			wantMissingB: []string{},
		},
		{
			a:            []string{".$", ".$.a", ".$.a.b"},
			b:            []string{".$", ".$.a", ".$.a.c"},
			wantMissingA: []string{".$.a.c"},
			wantMissingB: []string{".$.a.b"},
		},
		{
			a:            []string{".$", ".$.a", ".$.a.b.d"},
			b:            []string{".$", ".$.a", ".$.a.c"},
			wantMissingA: []string{".$.a.c"},
			wantMissingB: []string{".$.a.b.d"},
		},
		{
			a:            []string{".$", ".$.a", ".$.a.b.d"},
			b:            []string{".$", ".$.b", ".$.c"},
			wantMissingA: []string{".$.b", ".$.c"},
			wantMissingB: []string{".$.a", ".$.a.b.d"},
		},
		{
			a:            []string{".$", ".$.b", ".$.d.e.f"},
			b:            []string{".$", ".$.a", ".$.b", ".$.b.c", ".$.d.g"},
			wantMissingA: []string{".$.a", ".$.d.g"},
			wantMissingB: []string{".$.d.e.f"},
		},
		{
			a:            []string{".$", ".$.a", ".$.a.b", ".$.a.b.c3"},
			b:            []string{".$", ".$.a", ".$.a.b", ".$.a.b.c1", ".$.a.b.c2", ".$.a.b.c3", ".$.a.b.c4"},
			wantMissingA: []string{".$.a.b.c1", ".$.a.b.c2", ".$.a.b.c4"},
			wantMissingB: []string{},
		},
	}

	for _, tc := range testcases {
		diffA, diffB := DiffPaths(tc.a, tc.b)

		require.EqualValues(t, tc.wantMissingA, diffA)
		require.EqualValues(t, tc.wantMissingB, diffB)
	}
}
