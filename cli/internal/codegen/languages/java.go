package languages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kvnbanunu/melke-playground/cli/internal/codegen/types"
)

type JavaGenerator struct {
	config *types.Config
}

func NewJavaGenerator(config *types.Config) *JavaGenerator {
	return &JavaGenerator{config: config}
}

func (g *JavaGenerator) Generate() error {
	// Create package directory
	packageDir := filepath.Join(g.config.ProjectName, "source", "src", "main", "java",
		strings.ToLower(g.config.ProjectName))
	if err := os.MkdirAll(packageDir, 0755); err != nil {
		return err
	}

	// Generate a file for each class
	for _, typ := range g.config.Types {
		path := filepath.Join(packageDir, typ.Name+".java")
		content := g.generateClass(typ)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}

	// Generate utility class for standalone functions
	for _, file := range g.config.Files {
		if len(file.Functions) > 0 {
			path := filepath.Join(packageDir, file.Name+"Utils.java")
			content := g.generateUtils(file)
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *JavaGenerator) generateClass(typ types.TypeConfig) string {
	var sb strings.Builder

	// Package declaration
	packageName := strings.ToLower(g.config.ProjectName)
	sb.WriteString(fmt.Sprintf("package %s;\n\n", packageName))

	// Class JavaDoc
	sb.WriteString(fmt.Sprintf("/**\n * %s class\n */\n", typ.Name))

	// Class definition
	sb.WriteString(fmt.Sprintf("public class %s {\n", typ.Name))

	// Fields
	for _, field := range typ.Fields {
		access := field.Access
		if access == "" {
			access = "private"
		}
		sb.WriteString(fmt.Sprintf("    %s %s %s;\n",
			access,
			g.javaType(field.Type),
			field.Name))
	}
	sb.WriteString("\n")

	// Default constructor
	sb.WriteString(fmt.Sprintf("    public %s() {\n", typ.Name))
	for _, field := range typ.Fields {
		sb.WriteString(fmt.Sprintf("        this.%s = %s;\n",
			field.Name,
			g.javaDefaultValue(field.Type)))
	}
	sb.WriteString("    }\n\n")

	// Getters and setters
	for _, field := range typ.Fields {
		// Getter
		capitalizedField := strings.Title(field.Name)
		javaType := g.javaType(field.Type)

		sb.WriteString(fmt.Sprintf("    public %s get%s() {\n",
			javaType,
			capitalizedField))
		sb.WriteString(fmt.Sprintf("        return %s;\n",
			field.Name))
		sb.WriteString("    }\n\n")

		// Setter
		sb.WriteString(fmt.Sprintf("    public void set%s(%s %s) {\n",
			capitalizedField,
			javaType,
			field.Name))
		sb.WriteString(fmt.Sprintf("        this.%s = %s;\n",
			field.Name,
			field.Name))
		sb.WriteString("    }\n\n")
	}

	// Methods
	for _, method := range typ.Methods {
		access := method.Access
		if access == "" {
			access = "public"
		}

		params := make([]string, len(method.Parameters))
		for i, param := range method.Parameters {
			params[i] = fmt.Sprintf("%s %s",
				g.javaType(param.Type),
				param.Name)
		}

		returnType := "void"
		if method.ReturnType != "" && method.ReturnType != "void" {
			returnType = g.javaType(method.ReturnType)
		}

		// Method JavaDoc
		sb.WriteString("    /**\n")
		for _, param := range method.Parameters {
			sb.WriteString(fmt.Sprintf("     * @param %s the %s parameter\n",
				param.Name,
				param.Name))
		}
		if returnType != "void" {
			sb.WriteString("     * @return the result\n")
		}
		sb.WriteString("     */\n")

		sb.WriteString(fmt.Sprintf("    %s %s %s(%s) {\n",
			access,
			returnType,
			method.Name,
			strings.Join(params, ", ")))

		if returnType != "void" {
			sb.WriteString(fmt.Sprintf("        return %s;\n",
				g.javaDefaultValue(method.ReturnType)))
		}

		sb.WriteString("    }\n\n")
	}

	sb.WriteString("}\n")
	return sb.String()
}

func (g *JavaGenerator) generateUtils(file types.FileConfig) string {
	var sb strings.Builder

	// Package declaration
	packageName := strings.ToLower(g.config.ProjectName)
	sb.WriteString(fmt.Sprintf("package %s;\n\n", packageName))

	// Class JavaDoc
	sb.WriteString(fmt.Sprintf("/**\n * Utility functions for %s\n */\n",
		file.Name))

	// Class definition
	sb.WriteString(fmt.Sprintf("public class %sUtils {\n",
		strings.Title(file.Name)))

	// Private constructor to prevent instantiation
	sb.WriteString(fmt.Sprintf("    private %sUtils() {\n",
		strings.Title(file.Name)))
	sb.WriteString("        // Utility class, no instantiation\n")
	sb.WriteString("    }\n\n")

	// Utility functions
	for _, fn := range file.Functions {
		params := make([]string, len(fn.Parameters))
		for i, param := range fn.Parameters {
			params[i] = fmt.Sprintf("%s %s",
				g.javaType(param.Type),
				param.Name)
		}

		returnType := "void"
		if fn.ReturnType != "" && fn.ReturnType != "void" {
			returnType = g.javaType(fn.ReturnType)
		}

		// Method JavaDoc
		sb.WriteString("    /**\n")
		for _, param := range fn.Parameters {
			sb.WriteString(fmt.Sprintf("     * @param %s the %s parameter\n",
				param.Name,
				param.Name))
		}
		if returnType != "void" {
			sb.WriteString("     * @return the result\n")
		}
		sb.WriteString("     */\n")

		sb.WriteString(fmt.Sprintf("    public static %s %s(%s) {\n",
			returnType,
			fn.Name,
			strings.Join(params, ", ")))

		if returnType != "void" {
			sb.WriteString(fmt.Sprintf("        return %s;\n",
				g.javaDefaultValue(fn.ReturnType)))
		}

		sb.WriteString("    }\n\n")
	}

	sb.WriteString("}\n")
	return sb.String()
}

func (g *JavaGenerator) javaType(typeStr string) string {
	switch typeStr {
	case "int":
		return "int"
	case "float":
		return "float"
	case "double":
		return "double"
	case "char*", "const char*", "string":
		return "String"
	case "bool":
		return "boolean"
	default:
		if strings.Contains(typeStr, "*") {
			return strings.TrimSuffix(typeStr, "*")
		}
		return typeStr
	}
}

func (g *JavaGenerator) javaDefaultValue(typeStr string) string {
	switch g.javaType(typeStr) {
	case "int":
		return "0"
	case "float":
		return "0.0f"
	case "double":
		return "0.0"
	case "String":
		return "\"\""
	case "boolean":
		return "false"
	default:
		return "null"
	}
}
