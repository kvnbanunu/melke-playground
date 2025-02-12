package languages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kvnbanunu/melke-playground/cli/internal/codegen/types"
)

type GoGenerator struct {
	config *types.Config
}

func NewGoGenerator(config *types.Config) *GoGenerator {
	return &GoGenerator{config: config}
}

func (g *GoGenerator) Generate() error {
	for _, file := range g.config.Files {
		path := filepath.Join(g.config.ProjectName, "source", "src", file.Name+".go")
		content := g.generateContent(file)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (g *GoGenerator) generateContent(file types.FileConfig) string {
	var sb strings.Builder

	// Package declaration
	sb.WriteString(fmt.Sprintf("package %s\n\n", strings.ToLower(g.config.ProjectName)))

	// Generate struct definitions
	for _, typ := range g.config.Types {
		// Struct comment
		sb.WriteString(fmt.Sprintf("// %s represents %s\n", typ.Name, typ.Name))
		sb.WriteString(fmt.Sprintf("type %s struct {\n", typ.Name))

		for _, field := range typ.Fields {
			fieldName := field.Name
			if field.Access == "public" {
				fieldName = strings.Title(fieldName)
			}
			sb.WriteString(fmt.Sprintf("\t%s %s\n",
				fieldName,
				g.goType(field.Type)))
		}
		sb.WriteString("}\n\n")

		// Generate methods
		for _, method := range typ.Methods {
			methodName := method.Name
			if method.Access == "public" {
				methodName = strings.Title(methodName)
			}

			params := make([]string, len(method.Parameters))
			for i, param := range method.Parameters {
				params[i] = fmt.Sprintf("%s %s",
					param.Name,
					g.goType(param.Type))
			}

			returnType := ""
			if method.ReturnType != "" && method.ReturnType != "void" {
				returnType = " " + g.goType(method.ReturnType)
			}

			sb.WriteString(fmt.Sprintf("func (t *%s) %s(%s)%s {\n",
				typ.Name,
				methodName,
				strings.Join(params, ", "),
				returnType))

			// Add default return statement if needed
			if returnType != "" {
				sb.WriteString(fmt.Sprintf("\treturn %s\n",
					g.goDefaultValue(method.ReturnType)))
			}

			sb.WriteString("}\n\n")
		}
	}

	// Generate standalone functions
	for _, fn := range file.Functions {
		fnName := fn.Name
		if fn.Access == "public" {
			fnName = strings.Title(fnName)
		}

		params := make([]string, len(fn.Parameters))
		for i, param := range fn.Parameters {
			params[i] = fmt.Sprintf("%s %s",
				param.Name,
				g.goType(param.Type))
		}

		returnType := ""
		if fn.ReturnType != "" && fn.ReturnType != "void" {
			returnType = " " + g.goType(fn.ReturnType)
		}

		sb.WriteString(fmt.Sprintf("func %s(%s)%s {\n",
			fnName,
			strings.Join(params, ", "),
			returnType))

		if returnType != "" {
			sb.WriteString(fmt.Sprintf("\treturn %s\n",
				g.goDefaultValue(fn.ReturnType)))
		}

		sb.WriteString("}\n\n")
	}

	return sb.String()
}

func (g *GoGenerator) goType(typeStr string) string {
	switch typeStr {
	case "int":
		return "int"
	case "float", "double":
		return "float64"
	case "char*", "const char*", "string":
		return "string"
	case "bool":
		return "bool"
	default:
		if strings.Contains(typeStr, "*") {
			return "*" + strings.TrimSuffix(typeStr, "*")
		}
		return typeStr
	}
}

func (g *GoGenerator) goDefaultValue(typeStr string) string {
	switch g.goType(typeStr) {
	case "int":
		return "0"
	case "float64":
		return "0.0"
	case "string":
		return "\"\""
	case "bool":
		return "false"
	default:
		return "nil"
	}
}
