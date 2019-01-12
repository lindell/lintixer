package fixer

import (
	"go/ast"
	"go/token"
	"strings"
)

// NonCapitalError changes any error that is capitalized non lower case
func NonCapitalError(n ast.Node) bool {
	anyChanged := false
	ast.Inspect(n, func(n ast.Node) bool {
		call, ok := isCallDot(n, "errors", "New")
		if !ok {
			return true
		}

		// We need at least one argument
		if len(call.Args) == 0 {
			return true
		}

		var firstArg *ast.BasicLit
		if firstArg, ok = call.Args[0].(*ast.BasicLit); !ok {
			return true
		}

		if firstArg.Kind != token.STRING {
			return true
		}

		oldStr := firstArg.Value
		newStr := lowerFirst(firstArg.Value)
		if oldStr != newStr {
			firstArg.Value = newStr
			anyChanged = true
		}

		return true
	})
	return anyChanged
}

func isCallDot(node ast.Node, first, second string) (call *ast.CallExpr, valid bool) {
	var funcSelector *ast.SelectorExpr
	var x *ast.Ident
	var ok bool

	// Make sure we have errors.New call
	if call, ok = node.(*ast.CallExpr); !ok {
		return nil, false
	}
	if funcSelector, ok = call.Fun.(*ast.SelectorExpr); !ok {
		return nil, false
	}
	if x, ok = funcSelector.X.(*ast.Ident); !ok {
		return nil, false
	}
	if x.Name != first || funcSelector.Sel.Name != second {
		return nil, false
	}

	return call, true
}

func lowerFirst(str string) string {
	if len(str) == 2 {
		return str
	}
	return str[0:1] + strings.ToLower(str[1:2]) + str[2:]
}
