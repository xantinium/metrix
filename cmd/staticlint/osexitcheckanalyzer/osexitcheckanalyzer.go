// Пакет osexitcheckanalyzer содержит кастомный
// анализатор для проверки прямых вызовов [os.Exit]
// в функции main пакета main.
package osexitcheckanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OSExitCheckAnalyzer анализатор для проверки
// прямых вызовов [os.Exit] в функции main пакета main.
var OSExitCheckAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check direct calls to os.Exit in main function of main package",
	Run: func(pass *analysis.Pass) (any, error) {
		if pass.Pkg.Name() != "main" {
			return nil, nil
		}

		for _, file := range pass.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				f, ok := node.(*ast.FuncDecl)
				if !ok {
					return true
				}

				if f.Name.String() == "main" {
					checkFuncBody(pass, f.Body)
				}

				return true
			})
		}

		return nil, nil
	},
}

func checkFuncBody(pass *analysis.Pass, body *ast.BlockStmt) {
	ast.Inspect(body, func(node ast.Node) bool {
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := selExpr.X.(*ast.Ident)
		if !ok {
			return true
		}

		if ident.Name == "os" && selExpr.Sel.Name == "Exit" {
			pass.Reportf(callExpr.Pos(), "main func calls to os.Exit")
		}

		return true
	})
}
