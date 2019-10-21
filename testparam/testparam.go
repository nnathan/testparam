package testparam

import (
	"fmt"
	"go/ast"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

// Analyzer can be used with the analysis package and frontends.
var Analyzer = &analysis.Analyzer{
	Name:     "testparam",
	Doc:      "check for using correct testing parameter when calling testing.T.Run(...)",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// inImports checks if path* is contained in a slice of import specs which
// is usually provided by the AST parser.
//
// * path should be an import path literal, e.g. `"strings"` as opposed to "strings".
func inImports(imports []*ast.ImportSpec, path string) bool {
	for _, spec := range imports {
		if spec.Path.Value == path {
			return true
		}
	}

	return false
}

func run(pass *analysis.Pass) (interface{}, error) {
	lock.Lock()
	defer func() {
		fmt.Println("finished")
		lock.Unlock()
	}()

	testing := false
	for _, file := range pass.Files {
		if inImports(file.Imports, `"testing"`) {
			testing = true
			break
		}
	}

	if !testing {
		fmt.Println("not testing pass")
		return nil, nil
	}

	fmt.Println("file start")
	ast.Print(pass.Fset, pass.Files)
	fmt.Println("file end")

	fmt.Println("testing pass")
	return filter(pass)
}

func filter(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// We filter only function calls.
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		callee := typeutil.Callee(pass.TypesInfo, call)

		if callee == nil {
			return
		}

		fmt.Printf("callee = %s\nid=%s\ntype=%s\n", callee, callee.Id(), callee.Type())
		//fmt.Printf("callee = %s\n", callee)
		// We only consider function calls with two parameters:
		// the first is the id parameter of type string
		// the second is the fn parameter of either type func(*testing.T) or func(*testing.B)
		//
		// One could call a method expression: (*testing.T).Run(t, id, testfunc)
		// Let's ignore that case for now.
		if len(call.Args) != 2 {
			return
		}

		// There are two ways the call expression can be made, as simple a function call,
		// e.g. f(id, testfunc) or a call to t.Run(id, testfunc).
		//
		// In both cases, f must b
		// selector expressions let us filter on call expressions
		// in the form: x.f where:
		//   f is the "field selector",
		//   x is an identifier (could be a struct/method/interface/variable/package identifer)
		//
		// In any case, we know at some point we only care for expressions such as:
		//   x.Run(string, func(*testing.[TB]))
		// and resolve x to either *testing.T or *testing.B.
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		if sel.Sel.Name != "Run" {
			return
		}

		fmt.Printf("sel.Sel.Name = %s\n", sel.Sel.Name)
		fmt.Println("callexpr start")
		ast.Print(nil, call)
		fmt.Println("callexpr end")
		//fmt.Printf("%s %s %s\n", sel.Sel.Name, nTyp.Obj().Name(), nTyp.Pkg().Path())
		//if sel.Sel.Name != "Exec" &&
		//	nTyp.Obj().Name() != "DB" &&
		//	nTyp.Obj().Pkg().Path() != "database/sql" {
		//	return
		//}

		//arg0 := call.Args[0]
		//typ, ok = pass.TypesInfo.Types[arg0]
		//if !ok || typ.Value == nil {
		//	return
		//}

		//query := constant.StringVal(typ.Value)
		//_, err := pg_query.Parse(query)
		//if err != nil {
		//	pass.Reportf(call.Lparen, "Invalid query: %v", err)
		//	return
		//}
	})

	return nil, nil
}

var lock sync.Mutex

// https://agniva.me/vet/2019/01/21/vet-analyzer.html
// https://arslan.io/2019/06/13/using-go-analysis-to-write-a-custom-linter/
