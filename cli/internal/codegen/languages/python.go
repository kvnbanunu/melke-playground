package languages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kvnbanunu/melke-playground/cli/internal/codegen/types"
)

type PythonGenerator struct {
	config *types.Config
}

func NewPythonGenerator(config *types.Config) *PythonGenerator {
	return &PythonGenerator{config: config}
}

func (g *PythonGenerator) Generate() error {
	for _, file := range g.config.Files {
		path := filepath.Join(g.config.ProjectName, "source", "src", file.Name+".py")
		content := g.generateContent(file)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (g *PythonGenerator) generateContent(file types.FileConfig) string {
	var sb strings.Builder

	sb.WriteString("#!/usr/bin/env python3\n")
	sb.WriteString("from typing import List, Optional, Dict, Any\n\n")

	// Generate class definitions
	for _, typ := range g.config.Types {
		sb.WriteString(fmt.Sprintf("class %s:\n", typ.Name))

		// Generate constructor with type hints
		sb.WriteString("    def __init__(self):\n")
		if len(typ.Fields) == 0 {
			sb.WriteString("        pass\n")
		}
		for _, field := range typ.Fields {
			sb.WriteString(fmt.Sprintf("        self.%s: %s = %s\n",
				field.Name,
				g.pythonType(field.Type),
				g.pythonDefaultValue(field.Type)))
		}
		sb.WriteString("\n")

		// Generate methods
		for _, method := range typ.Methods {
			params := []string{"self"}
			for _, param := range method.Parameters {
				params = append(params, fmt.Sprintf("%s: %s",
					param.Name,
					g.pythonType(param.Type)))
			}

			returnHint := ""
			if method.ReturnType != "" && method.ReturnType != "void" {
				returnHint = " -> " + g.pythonType(method.ReturnType)
			}

			sb.WriteString(fmt.Sprintf("    def %s(%s)%s:\n",
				method.Name,
				strings.Join(params, ", "),
				returnHint))

			// Add default return statement if needed
			if method.ReturnType != "" && method.ReturnType != "void" {
				sb.WriteString(fmt.Sprintf("        return %s\n\n",
					g.pythonDefaultValue(method.ReturnType)))
			} else {
				sb.WriteString("        pass\n\n")
			}
		}
	}

	// Generate standalone functions
	for _, fn := range file.Functions {
		params := make([]string, len(fn.Parameters))
		for i, param := range fn.Parameters {
			params[i] = fmt.Sprintf("%s: %s",
				param.Name,
				g.pythonType(param.Type))
		}

		returnHint := ""
		if fn.ReturnType != "" && fn.ReturnType != "void" {
			returnHint = " -> " + g.pythonType(fn.ReturnType)
		}

		sb.WriteString(fmt.Sprintf("def %s(%s)%s:\n",
			fn.Name,
			strings.Join(params, ", "),
			returnHint))

		if fn.ReturnType != "" && fn.ReturnType != "void" {
			sb.WriteString(fmt.Sprintf("    return %s\n\n",
				g.pythonDefaultValue(fn.ReturnType)))
		} else {
			sb.WriteString("    pass\n\n")
		}
	}

	return sb.String()
}

func (g *PythonGenerator) pythonType(cType string) string {
	switch cType {
	case "int":
		return "int"
	case "float", "double":
		return "float"
	case "char*", "const char*", "string":
		return "str"
	case "bool":
		return "bool"
	default:
		if strings.Contains(cType, "*") {
			return "Optional[" + strings.TrimSuffix(cType, "*") + "]"
		}
		return "Any"
	}
}

func (g *PythonGenerator) pythonDefaultValue(cType string) string {
	switch g.pythonType(cType) {
	case "int":
		return "0"
	case "float":
		return "0.0"
	case "str":
		return "\"\""
	case "bool":
		return "False"
	default:
		return "None"
	}
}
