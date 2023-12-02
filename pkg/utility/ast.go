package utility

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

func ImportPackage(node *ast.File, projectModule, packageName string) {
	path := filepath.Join(projectModule, "domain", "features", packageName)
	importSpec := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf(`"%v"`, path),
		},
	}

	importDecl := &ast.GenDecl{
		Tok:    token.IMPORT,
		Lparen: token.Pos(1), // for grouping
		Specs:  []ast.Spec{importSpec},
	}

	// Check if there are existing imports, and if so, add to them
	found := false
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok && genDecl.Tok == token.IMPORT {
			genDecl.Specs = append(genDecl.Specs, importSpec)
			found = true
			break
		}
	}

	// If no import declaration exists, add the new one to Decls
	if !found {
		node.Decls = append([]ast.Decl{importDecl}, node.Decls...)
	}
}

func AddAnotherFxOptionsInModule(path, module, projectModule string) string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
	}
	ImportPackage(node, projectModule, module)

	// Traverse the AST and find the fx.Options call
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Module" {
					x.Args = append(x.Args, &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("fx"),
							Sel: ast.NewIdent("Options"),
						},
						Args: []ast.Expr{
							ast.NewIdent(module + ".Module"),
						},
					})
				}
			}
		}
		return true
	})

	// Add the source code in buffer
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, node); err != nil {
		fmt.Println(err)
	}
	formattedCode := buf.String()
	providerToInsert := fmt.Sprintf("fx.Options(%v.Module),", module)
	formattedCode = strings.Replace(formattedCode, providerToInsert, "\n\t"+providerToInsert, 1)
	return formattedCode
}
