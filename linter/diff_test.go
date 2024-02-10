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

	"github.com/cert-manager/helm-tool/linter/sets"
	"github.com/stretchr/testify/require"
)

func TestDiffPaths(t *testing.T) {
	type testcase struct {
		a            sets.Set[string]
		b            sets.Set[string]
		wantMissingA sets.Set[string]
		wantMissingB sets.Set[string]
	}

	testcases := []testcase{
		{
			a:            sets.New("a", "b", "c"),
			b:            sets.New("a", "b", "c"),
			wantMissingA: sets.New[string](),
			wantMissingB: sets.New[string](),
		},
		{
			a:            sets.New("a", "b", "c"),
			b:            sets.New("a", "b"),
			wantMissingA: sets.New[string](),
			wantMissingB: sets.New("c"),
		},
		{
			a:            sets.New("a", "b"),
			b:            sets.New("a", "b", "c"),
			wantMissingA: sets.New("c"),
			wantMissingB: sets.New[string](),
		},
		{
			a:            sets.New("a.b"),
			b:            sets.New("a"),
			wantMissingA: sets.New[string](),
			wantMissingB: sets.New[string](),
		},
		{
			a:            sets.New("a.b"),
			b:            sets.New("a.c"),
			wantMissingA: sets.New("a.c"),
			wantMissingB: sets.New("a.b"),
		},
		{
			a:            sets.New("a.b.d"),
			b:            sets.New("a.c"),
			wantMissingA: sets.New("a.c"),
			wantMissingB: sets.New("a.b.d"),
		},
		{
			a:            sets.New("a.b.d"),
			b:            sets.New("b", "c"),
			wantMissingA: sets.New("b", "c"),
			wantMissingB: sets.New("a.b.d"),
		},
		{
			a:            sets.New("b", "d.e.f"),
			b:            sets.New("a", "b.c", "d.g"),
			wantMissingA: sets.New("a", "d.g"),
			wantMissingB: sets.New("d.e.f"),
		},
		{
			a:            sets.New("a.b.c3"),
			b:            sets.New("a.b.c1", "a.b.c2", "a.b.c3", "a.b.c4"),
			wantMissingA: sets.New("a.b.c1", "a.b.c2", "a.b.c4"),
			wantMissingB: sets.New[string](),
		},
		{
			a:            sets.New("app.logLevel", "app.test", "app"),
			b:            sets.New("app.logLevel", "app.logLevela", "app.name"),
			wantMissingA: sets.New("app.logLevela", "app.name"),
			wantMissingB: sets.New("app.test"),
		},
	}

	for _, tc := range testcases {
		diffA, diffB := DiffPaths(tc.a, tc.b)

		require.ElementsMatch(t, tc.wantMissingA.UnsortedList(), diffA.UnsortedList())
		require.ElementsMatch(t, tc.wantMissingB.UnsortedList(), diffB.UnsortedList())
	}
}
