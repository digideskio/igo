// Copyright 2010 Aaron Jacobs. All rights reserved.
// See the LICENSE file for licensing details.

// The parse package offers utility functions for extracting dependency
// information from Go source files.
package parse

import (
	"go/ast"
	"go/parser"
	"igo/set"
	"regexp"
)

// GetPackageName returns the package name from the supplied .go file source
// code, or the empty string if it could not be properly parsed.
func GetPackageName(source string) string {
	fileNode, err := parser.ParseFile("", source, nil, parser.ImportsOnly)
	if err != nil {
		return ""
	}

	return fileNode.Name.Name()
}

// ExtractImports parses the supplied source code for a .go file and
// returns a set of package names that the file depends upon.
//
// For example, if source looks like the following:
//
//     import (
//       "./bar/baz"
//       "fmt"
//       "os"
//     )
//
//     func DoSomething() {
//       ...
//     }
//
// then the result will be { "./bar/baz", "fmt", "os" }.
//
// An attempt is made to return the imports for the file even if there is a
// syntax error elsewhere in the file.
func ExtractImports(source string) *set.StringSet {
	node, err := parser.ParseFile("", source, nil, parser.ImportsOnly)
	if err != nil {
		return &set.StringSet{}
	}

	var visitor importVisitor
	ast.Walk(&visitor, node)

	return &visitor.imports
}

type importVisitor struct {
	imports set.StringSet
}

var importRegexp *regexp.Regexp = regexp.MustCompile(`"(.+)"`)

func (v *importVisitor) Visit(node interface{}) ast.Visitor {
	switch t := node.(type) {
	case *ast.ImportSpec:
		for _, component := range node.(*ast.ImportSpec).Path {
			matches := importRegexp.MatchStrings(string(component.Value))
			if len(matches) < 2 {
				continue
			}

			v.imports.Insert(matches[1])
		}
	}

	return v
}