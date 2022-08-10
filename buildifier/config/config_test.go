/*
Copyright 2022 Google LLC

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

package config

import (
	"flag"
	"fmt"
	"strings"
	"testing"
)

func ExampleNew() {
	c := New()
	fmt.Print(c.String())
	// Output:
	// {
	//   "type": "auto"
	// }
}

func ExampleFlagSet() {
	c := New()
	flags := c.FlagSet("buildifier", flag.ExitOnError)
	flags.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s: %s (%q)\n", f.Name, f.Usage, f.DefValue)
	})
	// Output:
	// add_tables: path to JSON file with custom table definitions which will be merged with the built-in tables ("")
	// allowsort: additional sort contexts to treat as safe ("")
	// buildifier_disable: list of buildifier rewrites to disable ("")
	// config: path to .buildifier.json config file ("")
	// d: alias for -mode=diff ("false")
	// diff_command: command to run when the formatting mode is diff (default uses the BUILDIFIER_DIFF, BUILDIFIER_MULTIDIFF, and DISPLAY environment variables to create the diff command) ("")
	// format: diagnostics format: text or json (default text) ("")
	// help: print usage information ("false")
	// lint: lint mode: off, warn, or fix (default off) ("")
	// mode: formatting mode: check, diff, or fix (default fix) ("")
	// multi_diff: the command specified by the -diff_command flag can diff multiple files in the style of tkdiff (default false) ("false")
	// path: assume BUILD file has this path relative to the workspace directory ("")
	// r: find starlark files recursively ("false")
	// tables: path to JSON file with custom table definitions which will replace the built-in tables ("")
	// type: Input file type: build (for BUILD files), bzl (for .bzl files), workspace (for WORKSPACE files), default (for generic Starlark files) or auto (default, based on the filename) ("auto")
	// v: print verbose information to standard error ("false")
	// version: print the version of buildifier ("false")
	// warnings: comma-separated warnings used in the lint mode or "all" ("")
}

func ExampleFlagSet_parse() {
	c := New()
	flags := c.FlagSet("buildifier", flag.ExitOnError)
	flags.Parse([]string{
		"--add_tables=/path/to/add_tables.json",
		"--allowsort=proto_library.deps",
		"--allowsort=proto_library.srcs",
		"--buildifier_disable=unsafesort",
		"--config=/path/to/.buildifier.json",
		"-d",
		"--diff_command=diff",
		"--format=json",
		"--help",
		"--lint=fix",
		"--mode=fix",
		"--multi_diff=true",
		"--path=pkg/foo",
		"-r",
		"--tables=/path/to/tables.json",
		"--type=default",
		"-v",
		"--version",
		"--warnings=+print,-no-effect",
	})
	fmt.Println("help:", c.Help)
	fmt.Println("version:", c.Version)
	fmt.Println("configPath:", c.ConfigPath)
	fmt.Print(c.String())
	// Output:
	// help: true
	// version: true
	// configPath: /path/to/.buildifier.json
	// {
	//   "type": "default",
	//   "format": "json",
	//   "formattingMode": "fix",
	//   "diffMode": true,
	//   "lintMode": "fix",
	//   "warnings": "+print,-no-effect",
	//   "recursive": true,
	//   "verbose": true,
	//   "diffCommand": "diff",
	//   "multiDiff": true,
	//   "tables": "/path/to/tables.json",
	//   "addTables": "/path/to/add_tables.json",
	//   "path": "pkg/foo",
	//   "buildifier_disable": [
	//     "unsafesort"
	//   ],
	//   "allowsort": [
	//     "proto_library.deps",
	//     "proto_library.srcs"
	//   ]
	// }
}

func ExampleFlagSet_validateInputType() {
	c := New()
	flags := c.FlagSet("buildifier", flag.ExitOnError)
	flags.Parse([]string{
		"--type=foo",
	})
	fmt.Print(c.Validate(nil))
	// Output:
	// unrecognized input type foo; valid types are build, bzl, workspace, default, module, auto
}

func ExampleFlagSet_validateFormat() {
	c := New()
	flags := c.FlagSet("buildifier", flag.ExitOnError)
	flags.Parse([]string{
		"--format=foo",
	})
	fmt.Print(c.Validate(nil))
	// Output:
	// unrecognized format foo; valid types are text, json
}

func TestValidate(t *testing.T) {
	for name, tc := range map[string]struct {
		options string
		args    string
		want    error
	}{
		"mode fix ok": {options: "--mode=fix"},
	} {
		t.Run(name, func(t *testing.T) {
			c := New()
			flags := c.FlagSet("buildifier", flag.ExitOnError)
			flags.Parse(strings.Fields(tc.options))
			got := c.Validate(strings.Fields(tc.args))
			if tc.want == nil && got == nil {
				return
			}
			if tc.want == nil && got != nil {
				t.Fatalf("unexpected error: %v", got)
			}
			if tc.want != nil && got == nil {
				t.Fatalf("expected error did not occur: %v", tc.want)
			}
			if tc.want.Error() != got.Error() {
				t.Fatalf("error mismatch: want %v, got %v", tc.want.Error(), got.Error())
			}
		})
	}
}
