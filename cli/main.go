package main

import (
	"github.com/nnathan/testparam/testparam"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(testparam.Analyzer)
}
