package main

import (
	"github.com/Krchnk/Selectel_go/internal/loglint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(loglint.Analyzer)
}
