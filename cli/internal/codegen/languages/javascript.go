package languages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kvnbanunu/melke-playground/cli/internal/codegen/types"
)

type JavaScriptGenerator struct {
	config *types.Config
}

func NewJavaScriptGenerator(config *types.Config) *JavaScriptGenerator {
	return &JavaScriptGenerator{config: config}
}

func (g *JavaScriptGenerator) Generate() error {
	for _, file := range g.config.Files {
		path := filepath.Join(g.config.ProjectName, "source", "src", file.Name+".js")
		content := g.generateContent(file)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (g *JavaScriptGenerator) generateContent(file types.FileConfig) string {
	var sb strings.Builder

	// Add JSDoc types for better IDE support
	sb.WriteString("/**\n * @typedef {Object} Types\n")
	for _, typ := range g.config.Types {
		sb.WriteString(fmt.Sprintf(" * @typedef {Object} %s\n", typ.Name))
		for _, field := range typ.Fields {
			sb.WriteString(fmt.Sprintf(" * @property {%s} %s\n",
				g.jsDocType(field.Type),
				field.Name))
		}
	}
	sb.WriteString(" */\n\n")

	// Generate class definitions
	for _, typ := range g.config.Types {
		sb.WriteString(fmt.Sprintf("class %s {\n", typ.Name))

		// Constructor
		sb.WriteString("    constructor() {\n")
		for _, field := range typ.Fields {
			sb.WriteString(fmt.Sprintf("        /**\n         * @type {%s}\n         */\n",
				g.jsDocType(field.Type)))
			sb.WriteString(fmt.Sprintf("        this.%s = %s;\n",
				field.Name,
				g.jsDefaultValue(field.Type)))
		}
		sb.WriteString("    }\n\n")

		// Generate methods
		for _, method := range typ.Methods {
			// JSDoc for method
			sb.WriteString("    /**\n")
			for _, param := range method.Parameters {
				sb.WriteString(fmt.Sprintf("     * @param {%s} %s\n",
					g.jsDocType(param.Type),
					param.Name))
			}
			if method.ReturnType != "" && method.ReturnType != "void" {
				sb.WriteString(fmt.Sprintf("     * @returns {%s}\n",
					g.jsDocType(method.ReturnType)))
			}
			sb.WriteString("     */\n")

			params := make([]string, len(method.Parameters))
			for i, param := range method.Parameters {
				params[i] = param.Name
			}

			sb.WriteString(fmt.Sprintf("    %s(%s) {\n",
				method.Name,
				strings.Join(params, ", ")))

			if method.ReturnType != "" && method.ReturnType != "void" {
				sb.WriteString(fmt.Sprintf("        return %s;\n",
					g.jsDefaultValue(method.ReturnType)))
			}

			sb.WriteString("    }\n\n")
		}

		sb.WriteString("}\n\n")
	}

	// Generate standalone functions
	for _, fn := range file.Functions {
		// JSDoc for function
		sb.WriteString("/**\n")
		for _, param := range fn.Parameters {
			sb.WriteString(fmt.Sprintf(" * @param {%s} %s\n",
				g.jsDocType(param.Type),
				param.Name))
		}
		if fn.ReturnType != "" && fn.ReturnType != "void" {
			sb.WriteString(fmt.Sprintf(" * @returns {%s}\n",
				g.jsDocType(fn.ReturnType)))
		}
		sb.WriteString(" */\n")

		params := make([]string, len(fn.Parameters))
		for i, param := range fn.Parameters {
			params[i] = param.Name
		}

		sb.WriteString(fmt.Sprintf("function %s(%s) {\n",
			fn.Name,
			strings.Join(params, ", ")))

		if fn.ReturnType != "" && fn.ReturnType != "void" {
			sb.WriteString(fmt.Sprintf("    return %s;\n",
				g.jsDefaultValue(fn.ReturnType)))
		}

		sb.WriteString("}\n\n")
	}

	// Export all classes and functions
	sb.WriteString("module.exports = {\n")
	for _, typ := range g.config.Types {
		sb.WriteString(fmt.Sprintf("    %s,\n", typ.Name))
	}
	for _, fn := range file.Functions {
		sb.WriteString(fmt.Sprintf("    %s,\n", fn.Name))
	}
	sb.WriteString("};\n")

	return sb.String()
}

func (g *JavaScriptGenerator) jsDocType(typeStr string) string {
	switch typeStr {
	case "int", "float", "double":
		return "number"
	case "char*", "const char*", "string":
		return "string"
	case "bool":
		return "boolean"
	default:
		if strings.Contains(typeStr, "*") {
			return strings.TrimSuffix(typeStr, "*") + "|null"
		}
		return typeStr
	}
}

func (g *JavaScriptGenerator) jsDefaultValue(typeStr string) string {
	switch g.jsDocType(typeStr) {
	case "number":
		return "0"
	case "string":
		return "''"
	case "boolean":
		return "false"
	default:
		return "null"
	}
}
