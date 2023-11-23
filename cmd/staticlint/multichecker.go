package main

import (
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/kisielk/errcheck/errcheck"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"strings"
)

var Analyzer = &analysis.Analyzer{
	Name: "osExitInMain",
	Doc:  "reports the usage of os.Exit in main function",
	Run:  osExitInMain,
}

// osExitInMain Поиск вызовов функции os.Exit() в функции main(). Используется AST (Abstract Syntax Tree)
// для обхода всех узлов в файле и проверки каждого выражения, является ли оно вызовом функции os.Exit(). Если это
// так, то выводится сообщение об ошибке.
func osExitInMain(pass *analysis.Pass) (interface{}, error) {
	// функцией ast.Inspect проходим по всем узлам AST
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			// объявление функции
			fn, ok := node.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" {
				return true
			}

			// обход всех выражений в теле функции main()
			for _, stmt := range fn.Body.List {
				exprStmt, okExpr := stmt.(*ast.ExprStmt)
				if !okExpr {
					continue
				}

				// Вызов функции
				call, okCall := exprStmt.X.(*ast.CallExpr)
				if !okCall {
					continue
				}

				selector, okSelector := call.Fun.(*ast.SelectorExpr)
				if !okSelector {
					continue
				}

				// Вызов функции из пакета os
				identX, okIdendX := selector.X.(*ast.Ident)
				if !okIdendX || identX.String() != "os" {
					continue
				}

				// Функция Exit
				if selector.Sel.String() == "Exit" {
					pass.Reportf(stmt.Pos(), "call os.Exit in main function")
				}
			}

			return true
		})
	}

	return nil, nil

}

func main() {
	analyzers := []*analysis.Analyzer{
		Analyzer,
		ineffassign.Analyzer,
		errcheck.Analyzer,
	}

	// все анализаторы класса SA пакета staticcheck
	for i := range staticcheck.Analyzers {
		if strings.HasPrefix(staticcheck.Analyzers[i].Analyzer.Name, "SA") {
			analyzers = append(analyzers, staticcheck.Analyzers[i].Analyzer)
		}
	}

	// CheckErrorStrings
	for i := range stylecheck.Analyzers {
		if stylecheck.Analyzers[i].Analyzer.Name == "ST1005" {
			analyzers = append(analyzers, stylecheck.Analyzers[i].Analyzer)
		}
	}

	// Run all analyzers
	multichecker.Main(analyzers...)
}
