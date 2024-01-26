/*
Copyright 2021 The cert-manager Authors.

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

package paths

import (
	"reflect"
	"testing"
)

func TestParsePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected Path
		wantErr  bool
	}{
		{
			name:     "Empty path",
			path:     "",
			expected: Path{},
			wantErr:  false,
		},
		{
			name:     "Single string component",
			path:     "foo",
			expected: Path{mapPathComponent("foo")},
			wantErr:  false,
		},
		{
			name:     "Single index component",
			path:     "foo[0]",
			expected: Path{mapPathComponent("foo"), arrayPathComponent(0)},
			wantErr:  false,
		},
		{
			name:     "Multiple components",
			path:     "foo.bar[0].baz",
			expected: Path{mapPathComponent("foo"), mapPathComponent("bar"), arrayPathComponent(0), mapPathComponent("baz")},
			wantErr:  false,
		},
		{
			name:     "Invalid path 2",
			path:     "foo[0]aa",
			expected: Path{mapPathComponent("foo"), arrayPathComponent(0)},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Parse() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestWithProperty(t *testing.T) {
	rootPath := Path{}.WithProperty("foo").WithProperty("bar").WithProperty("aaaa")

	path1 := rootPath.WithProperty("baz1")
	path2 := rootPath.WithProperty("baz2")

	if path1.String() != "foo.bar.aaaa.baz1" {
		t.Errorf("path1.String() = %v, expected %v", path1.String(), "foo.bar.aaaa.baz1")
	}

	if path2.String() != "foo.bar.aaaa.baz2" {
		t.Errorf("path2.String() = %v, expected %v", path2.String(), "foo.bar.aaaa.baz2")
	}
}

func TestWithIndex(t *testing.T) {
	rootPath := Path{}.WithProperty("foo").WithProperty("bar").WithProperty("aaaa")

	path1 := rootPath.WithIndex(0)
	path2 := rootPath.WithIndex(1)

	if path1.String() != "foo.bar.aaaa[0]" {
		t.Errorf("path1.String() = %v, expected %v", path1.String(), "foo.bar.aaaa[0]")
	}

	if path2.String() != "foo.bar.aaaa[1]" {
		t.Errorf("path2.String() = %v, expected %v", path2.String(), "foo.bar.aaaa[1]")
	}
}
