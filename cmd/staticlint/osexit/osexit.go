// Package osexit contains an analyzer that checks for calls to os.Exit in Go code.
//
// The analyzer walks the AST looking for function declarations and checks if the
// function calls os.Exit. If a call to os.Exit is found, the analyzer prints an
// error message with the location of the call.
//
// The analyzer skips functions that are not named "main".
package osexit

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is the analyzer that checks for calls to os.Exit.
var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for calls to os.Exit in Go code",
	Run:  run,
}

// run is the entry point for the analyzer. It walks the AST looking for function
// declarations and checks if the function calls os.Exit. If a call to os.Exit is
// found, the analyzer prints an error message with the location of the call.
func run(pass *analysis.Pass) (interface{}, error) {
	// Iterate over the files in the package.
	for _, file := range pass.Files {
		// Print the name of the file being analyzed.
		//fmt.Println(file.Name)

		// Get the file set for the current package.
		fset := pass.Fset

		// Inspect the AST looking for function declarations.
		ast.Inspect(file, func(node ast.Node) bool {
			// Check if the node is a function declaration.
			if fun, ok := node.(*ast.FuncDecl); ok {
				// Skip functions that are not named "main".
				if fun.Name.Name != "main" {
					//fmt.Println(fun.Name.Name, "--->", file.Name)
					return true
				}

				// Inspect the function body looking for calls to os.Exit.
				for _, node_ := range fun.Body.List {
					// Check if the node is an expression statement.
					if expr, ok := node_.(*ast.ExprStmt); ok {
						// Check if the expression is a call expression.
						if c, ok := expr.X.(*ast.CallExpr); ok {
							// Check if the function being called is os.Exit.
							if s, ok := c.Fun.(*ast.SelectorExpr); ok {
								// Check if the selector expression is of the correct type.
								_, ok := s.X.(*ast.Ident)
								if !ok {
									//fmt.Println("Not *ast.Ident")
									return true
								}
								// Check if the function being called is os.Exit.
								if s.X.(*ast.Ident).Name == "os" && s.Sel.Name == "Exit" {
									// Print an error message with the location of the call.
									pos := fset.Position(c.Pos())
									fmt.Printf("error: call to os.Exit found at %s:%d:%d\n", pos.Filename, pos.Line, pos.Column)
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
