/*
Copyright 2026 The cert-manager Authors.

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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "values-*.yaml")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return f.Name()
}

// Self-referential anchor must not cause a stack overflow.
func TestLoad_SelfReferentialAlias(t *testing.T) {
	path := writeTemp(t, "a: &a\n  b: *a\n")
	_, err := Load(path, false)
	// The call must return (possibly with an error) — a stack overflow is fatal
	// and would kill the test process before we reach this line.
	_ = err
}

// Scaled-down billion-laughs must not OOM.
// Uses 5 levels (10^5 = 100 000 virtual nodes) instead of 9 to keep the test
// fast while still exercising the fan-out path.
func TestLoad_BillionLaughs(t *testing.T) {
	yaml := `
a: &a [1,1,1,1,1,1,1,1,1,1]
b: &b [*a,*a,*a,*a,*a,*a,*a,*a,*a,*a]
c: &c [*b,*b,*b,*b,*b,*b,*b,*b,*b,*b]
d: &d [*c,*c,*c,*c,*c,*c,*c,*c,*c,*c]
e:    [*d,*d,*d,*d,*d,*d,*d,*d,*d,*d]
`
	path := writeTemp(t, yaml)
	_, err := Load(path, false)
	_ = err
}

// Shared anchors must not panic. Note: the visited-set means a non-scalar
// anchor's subtree is walked only at the first alias reference, so serviceB's
// properties are not documented — that trade-off is intentional.
func TestLoad_SharedAnchorDoesNotPanic(t *testing.T) {
	yaml := `
defaults: &defaults
  replicaCount: 1
  image: nginx

serviceA:
  <<: *defaults
  port: 80

serviceB:
  <<: *defaults
  port: 443
`
	path := writeTemp(t, yaml)
	doc, err := Load(path, false)
	require.NoError(t, err)
	require.NotNil(t, doc)
}

// A plain acyclic values file must parse correctly and surface its properties.
func TestLoad_SimpleValues(t *testing.T) {
	yaml := `
# -- Number of replicas
replicaCount: 1
# -- Container image
image: nginx
`
	path := writeTemp(t, yaml)
	doc, err := Load(path, false)
	require.NoError(t, err)
	require.NotEmpty(t, doc.Sections)

	var paths []string
	for _, s := range doc.Sections {
		for _, p := range s.Properties {
			paths = append(paths, p.Path.String())
		}
	}
	assert.Contains(t, paths, "replicaCount")
	assert.Contains(t, paths, "image")
}

// A scalar anchor aliased from multiple keys must surface a property at each path.
func TestLoad_ScalarAnchorAliasedTwice(t *testing.T) {
	path := writeTemp(t, "x: &s 1\ny: *s\n")
	doc, err := Load(path, false)
	require.NoError(t, err)
	require.NotEmpty(t, doc.Sections)
	var paths []string
	for _, s := range doc.Sections {
		for _, p := range s.Properties {
			paths = append(paths, p.Path.String())
		}
	}
	assert.Contains(t, paths, "x")
	assert.Contains(t, paths, "y")
}

// Missing file must return an error, not panic.
func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "does-not-exist.yaml"), false)
	require.Error(t, err)
}
