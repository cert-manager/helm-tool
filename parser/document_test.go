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

// Scaled-down billion-laughs must not OOM, and — since a shared node's
// properties are now memoized and replayed rather than just walked once —
// must still produce the exact (bounded, polynomial) number of properties
// implied by the fan-out, not an exponential blow-up.
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
	doc, err := Load(path, false)
	require.NoError(t, err)

	total := 0
	for _, s := range doc.Sections {
		total += len(s.Properties)
	}
	// a:10 + b:10*10 + c:10*100 + d:10*1000 + e:10*10000 = 111110
	assert.Equal(t, 111110, total)
}

// A mapping merged into multiple mappings via "<<" must surface its
// properties under each merge site, not just the first.
func TestLoad_MergeKeySharedAcrossMappings(t *testing.T) {
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

	var paths []string
	for _, s := range doc.Sections {
		for _, p := range s.Properties {
			paths = append(paths, p.Path.String())
		}
	}
	assert.Contains(t, paths, "serviceA.replicaCount")
	assert.Contains(t, paths, "serviceA.image")
	assert.Contains(t, paths, "serviceA.port")
	assert.Contains(t, paths, "serviceB.replicaCount")
	assert.Contains(t, paths, "serviceB.image")
	assert.Contains(t, paths, "serviceB.port")
	assert.NotContains(t, paths, "serviceA.<<")
	assert.NotContains(t, paths, "serviceB.<<")
}

// An explicit key must override a merged-in key of the same name regardless
// of whether it's declared before or after the "<<" merge key.
func TestLoad_MergeKeyExplicitOverridesMerged(t *testing.T) {
	yaml := `
defaults: &defaults
  # -- from defaults
  replicaCount: 1

service:
  <<: *defaults
  # -- overridden
  replicaCount: 3
`
	path := writeTemp(t, yaml)
	doc, err := Load(path, false)
	require.NoError(t, err)

	var found *Property
	for _, s := range doc.Sections {
		for i, p := range s.Properties {
			if p.Path.String() == "service.replicaCount" {
				found = &s.Properties[i]
			}
		}
	}
	require.NotNil(t, found, "expected service.replicaCount to be documented")
	assert.Equal(t, "3", found.Default)
}

// Merging multiple mappings via a sequence must resolve conflicting keys in
// favour of the earliest mapping in the sequence.
func TestLoad_MergeKeySequencePrecedence(t *testing.T) {
	yaml := `
a: &a
  value: from-a
b: &b
  value: from-b

service:
  <<: [*a, *b]
`
	path := writeTemp(t, yaml)
	doc, err := Load(path, false)
	require.NoError(t, err)

	var found *Property
	for _, s := range doc.Sections {
		for i, p := range s.Properties {
			if p.Path.String() == "service.value" {
				found = &s.Properties[i]
			}
		}
	}
	require.NotNil(t, found, "expected service.value to be documented")
	assert.Equal(t, "from-a", found.Default)
}

// A mapping that merges itself (directly or transitively) must not cause a
// stack overflow.
func TestLoad_MergeKeySelfReferential(t *testing.T) {
	yaml := "a: &a\n  <<: *a\n  b: 1\n"
	path := writeTemp(t, yaml)
	_, err := Load(path, false)
	// The call must return (possibly with an error) — a stack overflow is
	// fatal and would kill the test process before we reach this line.
	_ = err
}

// A non-scalar value nested inside a mapping merged at more than one site
// must be documented at every site, not just the first. Before memoization
// was added, only the first site to reach the shared "nested" node got its
// properties — later sites silently got nothing, because the cycle/fan-out
// guard treated every repeat encounter as already handled.
func TestLoad_MergeKeyNestedSharedSubtreeDocumentedAtEverySite(t *testing.T) {
	yaml := `
template: &template
  nested:
    # -- from the shared template
    value: shared-default

svcA:
  <<: *template

svcB:
  <<: *template
`
	path := writeTemp(t, yaml)
	doc, err := Load(path, false)
	require.NoError(t, err)

	var paths []string
	for _, s := range doc.Sections {
		for _, p := range s.Properties {
			paths = append(paths, p.Path.String())
		}
	}
	assert.Contains(t, paths, "svcA.nested.value")
	assert.Contains(t, paths, "svcB.nested.value")
}

// The same nested-sharing fix must also apply to plain alias reuse, not just
// "<<" merge keys — a nested non-scalar value reachable from two different
// alias sites must be documented at both.
func TestLoad_AliasNestedSharedSubtreeDocumentedAtEverySite(t *testing.T) {
	yaml := `
template: &template
  nested:
    # -- from the shared template
    value: shared-default

svcA: *template
svcB: *template
`
	path := writeTemp(t, yaml)
	doc, err := Load(path, false)
	require.NoError(t, err)

	var paths []string
	for _, s := range doc.Sections {
		for _, p := range s.Properties {
			paths = append(paths, p.Path.String())
		}
	}
	assert.Contains(t, paths, "svcA.nested.value")
	assert.Contains(t, paths, "svcB.nested.value")
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
