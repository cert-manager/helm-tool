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
