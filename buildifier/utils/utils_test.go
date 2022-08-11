/*
Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"path/filepath"
	"testing"
)

func TestIsStarlarkFile(t *testing.T) {
	tests := []struct {
		filename string
		ok       bool
	}{
		{
			filename: "BUILD",
			ok:       true,
		},
		{
			filename: "BUILD.bazel",
			ok:       true,
		},
		{
			filename: "BUILD.oss",
			ok:       true,
		},
		{
			filename: "BUILD.foo.bazel",
			ok:       true,
		},
		{
			filename: "BUILD.foo.oss",
			ok:       true,
		},
		{
			filename: "build.foo.bazel",
			ok:       false,
		},
		{
			filename: "build.foo.oss",
			ok:       false,
		},
		{
			filename: "build.oss",
			ok:       false,
		},
		{
			filename: "WORKSPACE",
			ok:       true,
		},
		{
			filename: "WORKSPACE.bazel",
			ok:       true,
		},
		{
			filename: "WORKSPACE.oss",
			ok:       true,
		},
		{
			filename: "WORKSPACE.foo.bazel",
			ok:       true,
		},
		{
			filename: "WORKSPACE.foo.oss",
			ok:       true,
		},
		{
			filename: "workspace.foo.bazel",
			ok:       false,
		},
		{
			filename: "workspace.foo.oss",
			ok:       false,
		},
		{
			filename: "workspace.oss",
			ok:       false,
		},
		{
			filename: "build.gradle",
			ok:       false,
		},
		{
			filename: "workspace.xml",
			ok:       false,
		},
		{
			filename: "foo.bzl",
			ok:       true,
		},
		{
			filename: "foo.BZL",
			ok:       false,
		},
		{
			filename: "build.bzl",
			ok:       true,
		},
		{
			filename: "workspace.sky",
			ok:       true,
		},
		{
			filename: "foo.star",
			ok:       true,
		},
		{
			filename: "foo.bar",
			ok:       false,
		},
		{
			filename: "foo.build",
			ok:       false,
		},
		{
			filename: "foo.workspace",
			ok:       false,
		},
	}

	for _, tc := range tests {
		if isStarlarkFile(tc.filename) != tc.ok {
			t.Errorf("Wrong result for %q, want %t", tc.filename, tc.ok)
		}
	}
}

func TestCommonPrefix(t *testing.T) {
	tests := map[string]struct {
		paths []string
		want  string
	}{
		"degenerate": {},
		"single": {
			paths: []string{
				"/src/main/java",
			},
			want: "/src/main/java",
		},
		"common dir": {
			paths: []string{
				"/src/main/java",
				"/src/main/test",
			},
			want: "/src/main",
		},
		"common name": {
			paths: []string{
				"/src/main/java",
				"/src/main/java2",
			},
			want: "/src/main",
		},
		"disjoint": {
			paths: []string{
				"/src/main/java",
				"/test/main/java",
			},
			want: "",
		},
		"windows": {
			paths: []string{
				`c:\src\main\java`,
				`c:\src\main\test`,
			},
			want: `c:\src\main\`,
		},
	}

	for _, tc := range tests {
		got := CommonPrefix(filepath.Separator, tc.paths...)
		if tc.want != got {
			t.Errorf("CommonPrefix: want %q, got %q", tc.want, got)
		}
	}
}
