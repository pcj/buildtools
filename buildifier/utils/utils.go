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

// Package utils contains shared methods that can be used by different implementations of
// buildifier binary
package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/buildtools/build"
	"github.com/bazelbuild/buildtools/warn"
)

func isStarlarkFile(name string) bool {
	ext := filepath.Ext(name)
	switch ext {
	case ".bzl", ".sky", ".star":
		return true
	}

	switch ext {
	case ".bazel", ".oss":
		// BUILD.bazel or BUILD.foo.bazel should be treated as Starlark files, same for WORSKSPACE and MODULE
		return strings.HasPrefix(name, "BUILD.") || strings.HasPrefix(name, "WORKSPACE.") || strings.HasPrefix(name, "MODULE.")
	}

	return name == "BUILD" || name == "WORKSPACE"
}

func skip(info os.FileInfo) bool {
	return info.IsDir() && info.Name() == ".git"
}

func CommonPrefix(sep byte, paths ...string) string {
	// Handle special cases.
	switch len(paths) {
	case 0:
		return ""
	case 1:
		return filepath.Clean(paths[0])
	}

	// Note, we treat string as []byte, not []rune as is often
	// done in Go. (And sep as byte, not rune). This is because
	// most/all supported OS' treat paths as string of non-zero
	// bytes. A filename may be displayed as a sequence of Unicode
	// runes (typically encoded as UTF-8) but paths are
	// not required to be valid UTF-8 or in any normalized form
	// (e.g. "é" (U+00C9) and "é" (U+0065,U+0301) are different
	// file names.
	c := []byte(filepath.Clean(paths[0]))

	// We add a trailing sep to handle the case where the
	// common prefix directory is included in the path list
	// (e.g. /home/user1, /home/user1/foo, /home/user1/bar).
	// path.Clean will have cleaned off trailing / separators with
	// the exception of the root directory, "/" (in which case we
	// make it "//", but this will get fixed up to "/" bellow).
	c = append(c, sep)

	// Ignore the first path since it's already in c
	for _, v := range paths[1:] {
		// Clean up each path before testing it
		v = filepath.Clean(v) + string(sep)

		// Find the first non-common byte and truncate c
		if len(v) < len(c) {
			c = c[:len(v)]
		}
		for i := 0; i < len(c); i++ {
			if v[i] != c[i] {
				c = c[:i]
				break
			}
		}
	}

	// Remove trailing non-separator characters and the final separator
	for i := len(c) - 1; i >= 0; i-- {
		if c[i] == sep {
			c = c[:i]
			break
		}
	}

	return string(c)
}

// ExpandDirectories takes a list of file/directory names and returns a list with file names
// by traversing each directory recursively and searching for relevant Starlark files.
func ExpandDirectories(args *[]string) ([]string, error) {
	files := []string{}
	for _, arg := range *args {
		info, err := os.Stat(arg)
		if err != nil {
			return []string{}, err
		}
		if !info.IsDir() {
			files = append(files, arg)
			continue
		}
		err = filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if skip(info) {
				return filepath.SkipDir
			}
			if !info.IsDir() && isStarlarkFile(info.Name()) {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return []string{}, err
		}
	}
	return files, nil
}

// GetParser returns a parser for a given file type
func GetParser(inputType string) func(filename string, data []byte) (*build.File, error) {
	switch inputType {
	case "build":
		return build.ParseBuild
	case "bzl":
		return build.ParseBzl
	case "auto":
		return build.Parse
	case "workspace":
		return build.ParseWorkspace
	case "module":
		return build.ParseModule
	default:
		return build.ParseDefault
	}
}

// getFileReader returns a *FileReader object that reads files from the local
// filesystem if the workspace root is known.
func getFileReader(workspaceRoot string) *warn.FileReader {
	if workspaceRoot == "" {
		return nil
	}

	readFile := func(filename string) ([]byte, error) {
		// Use OS-specific path separators
		filename = strings.ReplaceAll(filename, "/", string(os.PathSeparator))
		path := filepath.Join(workspaceRoot, filename)

		return ioutil.ReadFile(path)
	}

	return warn.NewFileReader(readFile)
}

// Lint calls the linter and returns a list of unresolved findings
func Lint(f *build.File, lint string, warningsList *[]string, verbose bool) []*warn.Finding {
	fileReader := getFileReader(f.WorkspaceRoot)

	switch lint {
	case "warn":
		return warn.FileWarnings(f, *warningsList, nil, warn.ModeWarn, fileReader)
	case "fix":
		warn.FixWarnings(f, *warningsList, verbose, fileReader)
	}
	return nil
}
