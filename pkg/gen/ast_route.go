package gen

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"regexp"
	"strings"

	"github.com/gookit/color"
	"github.com/mukezhz/geng/pkg/model"
	"github.com/mukezhz/geng/pkg/utility"
	"golang.org/x/tools/go/ast/astutil"
)

// filePath := "route.go"
// newRoutePath := "/api/users/new"
// newRouteMethod := "GET"
// newRouteHandler := "r.controller.NewHandler"
func AddRoute(bruModel model.BruModel) {
	// filePath, newRoutePath, newRouteMethod, newRouteHandler string
	filePath := bruModel.GetModulePath("route.go")
	newRoutePath := bruModel.Route
	newRouteMethod := bruModel.Method
	newRouteHandler := bruModel.Handler

	if !utility.FileExists(filePath) {
		color.Yellowln("File doesn't exists:", filePath)
		return
	}

	handlerPrefix := "r.controller."

	// Parse the Go file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Find the RegisterUserRoute function
	registerUserRouteFunc, err := findRegisterRoute(node)
	if err != nil {
		fmt.Println("Register function not found")
		return
	}
	endpoint := strings.Replace(newRoutePath, "/api", "", 1)
	// Check if the route already exists
	routeExists := checkRouteExistance(registerUserRouteFunc, endpoint, newRouteMethod)
	// If the route exists, exit
	if routeExists {
		return
	}

	// Create the new route statement
	newRouteStmt := attachController(newRouteMethod, newRoutePath, handlerPrefix, newRouteHandler)

	// Insert the new route at the end of the RegisterUserRoute function body
	registerUserRouteFunc.Body.List = append(registerUserRouteFunc.Body.List, newRouteStmt)

	// Write the modified AST back to the file
	out, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	printer.Fprint(out, fset, node)
}

func checkRouteExistance(
	registerUserRouteFunc *ast.FuncDecl,
	endpoint, newRouteMethod string,
) bool {
	routeExists := false
	astutil.Apply(registerUserRouteFunc.Body, func(c *astutil.Cursor) bool {
		if exprStmt, ok := c.Node().(*ast.ExprStmt); ok {
			if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
				if _, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if len(callExpr.Args) > 1 {
						methodName := callExpr.Fun.(*ast.SelectorExpr).Sel.Name
						if routePath, ok := callExpr.Args[0].(*ast.BasicLit); ok && routePath.Kind == token.STRING {
							if routePath.Value == fmt.Sprintf("\"%s\"", endpoint) && methodName == newRouteMethod {
								routeExists = true
								return false
							}
						}
					}
				}
			}
		}
		return true
	}, nil)
	return routeExists
}

func findRegisterRoute(node *ast.File) (registerUserRouteFunc *ast.FuncDecl, err error) {

	astutil.Apply(node, func(c *astutil.Cursor) bool {
		if funcDecl, ok := c.Node().(*ast.FuncDecl); ok {
			re := regexp.MustCompile(`(Register)?(\w+)?(Route)?`)
			if re.MatchString(funcDecl.Name.Name) {
				registerUserRouteFunc = funcDecl
				return false
			}
		}
		return true
	}, nil)

	if registerUserRouteFunc == nil {
		return nil, errors.New("Register function not found")
	}
	return registerUserRouteFunc, nil
}

func attachController(
	newRouteMethod string,
	newRoutePath string,
	handlerPrefix string,
	newRouteHandler string,
) *ast.ExprStmt {
	newRouteStmt := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("api"),
				Sel: ast.NewIdent(newRouteMethod),
			},
			Args: []ast.Expr{
				&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", strings.Replace(newRoutePath, "/api", "", 1))},
				&ast.SelectorExpr{
					X:   ast.NewIdent(strings.TrimSuffix(handlerPrefix, ".")),
					Sel: ast.NewIdent(newRouteHandler),
				},
			},
		},
	}
	return newRouteStmt
}
