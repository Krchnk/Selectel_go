package loglint

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages for style and content",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	ast.Inspect(pass.Files[0], func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		logger := ident.Name
		method := sel.Sel.Name
		if !isLogMethod(logger, method) {
			return true
		}
		if len(call.Args) == 0 {
			return true
		}
		msg, ok := call.Args[0].(*ast.BasicLit)
		if !ok || msg.Kind !=  token.STRING {
			return true
		}
		logMsg := strings.Trim(msg.Value, "\"")
		//lowercase first letter
		if len(logMsg) > 0 && !isLowerFirst(logMsg) {
			pass.Reportf(msg.Pos(), "log message should start with a lowercase letter")
		}
		//English only
		if !isEnglish(logMsg) {
			pass.Reportf(msg.Pos(), "log message should be in English only")
		}
		//no special symbols or emoji
		if hasSpecialOrEmoji(logMsg) {
			pass.Reportf(msg.Pos(), "log message should not contain special symbols or emoji")
		}
		//no sensitive data
		if containsSensitive(logMsg) {
			pass.Reportf(msg.Pos(), "log message should not contain sensitive data")
		}
		return true
	})
	return nil, nil
}

func isLogMethod(logger, method string) bool {
	loggers := []string{"log", "slog", "zap", "logger", "sugar"}
	methods := []string{"Info", "Error", "Warn", "Debug", "Fatal", "Print", "Printf", "Println"}
	for _, l := range loggers {
		if logger == l {
			for _, m := range methods {
				if method == m {
					return true
				}
			}
		}
	}
	return false
}

func isLowerFirst(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return unicode.IsLower(r)
		}
	}
	return true
}

func isEnglish(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') {
			return false
		}
	}
	return true
}

func hasSpecialOrEmoji(s string) bool {
	for _, r := range s {
		if unicode.IsSymbol(r) || unicode.IsPunct(r) || (r >= 0x1F600 && r <= 0x1F64F) || (r >= 0x1F300 && r <= 0x1F5FF) {
			return true
		}
	}
	return false
}

func containsSensitive(s string) bool {
	keywords := []string{"password", "token", "api_key", "secret", "key", "pass", "auth", "session"}
	lower := strings.ToLower(s)
	for _, k := range keywords {
		if strings.Contains(lower, k) {
			return true
		}
	}
	return false
}
