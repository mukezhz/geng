package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/model"
	"github.com/mukezhz/geng/pkg/utility"
	"golang.org/x/tools/go/ast/astutil"
)

func AddController(bruModel model.BruModel) {
	// filePath, methodName string
	filePath := bruModel.GetModulePath("controller.go")
	if !utility.FileExists(filePath) {
		color.Yellowln("File doesn't exists:", filePath)
		return
	}
	methodName := bruModel.Handler
	// Parse the Go file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Identify the struct type by finding a method receiver
	structName := findStruct(node)

	if structName == "" {
		fmt.Println("Struct not found")
		return
	}

	// Check if the method already exists
	methodExists := checkMethodExistance(node, methodName, structName)

	// If the method exists, exit
	if methodExists {
		return
	}

	// Create a new method declaration
	newMethod := createMethodImplementation(bruModel, structName)
	// Add the new method to the file
	node.Decls = append(node.Decls,
		newMethod,
	)

	// Write the modified AST back to the file
	out, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	defer out.Close()
	printer.Fprint(out, fset, node)
	fmt.Printf("Method %s added to struct %s\n", methodName, structName)
}

func findStruct(node *ast.File) string {
	structName := ""
	astutil.Apply(node, func(c *astutil.Cursor) bool {
		if funcDecl, ok := c.Node().(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok {
						structName = ident.Name
						return false
					}
				}
			}
		}
		return true
	}, nil)
	return structName
}

func checkMethodExistance(node *ast.File, methodName, structName string) bool {
	methodExists := false
	astutil.Apply(node, func(c *astutil.Cursor) bool {
		if funcDecl, ok := c.Node().(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok {
						if ident.Name == structName && funcDecl.Name.Name == methodName {
							methodExists = true
							return false
						}
					}
				}
			}
		}
		return true
	}, nil)
	return methodExists
}

func createMethodImplementation(bruModel model.BruModel, structName string) *ast.FuncDecl {
	return &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: fmt.Sprintf("\n// @Title: %s", bruModel.Name),
				},
				{
					Text: fmt.Sprintf("// @Description: %s", bruModel.Description),
				},
				{
					Text: fmt.Sprintf("// @Route %s [%s]", bruModel.Route, bruModel.Method),
				},
			},
		},
		Name: ast.NewIdent(bruModel.Handler),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("c")},
					Type:  &ast.StarExpr{X: ast.NewIdent(structName)},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type:  &ast.SelectorExpr{X: ast.NewIdent("*gin"), Sel: ast.NewIdent("Context")},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("ctx"),
							Sel: ast.NewIdent("JSON"),
						},
						Args: []ast.Expr{
							ast.NewIdent("http.StatusOK"),
							&ast.CompositeLit{
								Type: &ast.MapType{
									Key:   ast.NewIdent("string"),
									Value: ast.NewIdent("interface{}"),
								},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("\"data\""),
										Value: ast.NewIdent("\"ok\""),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
