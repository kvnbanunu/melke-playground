package languages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kvnbanunu/melke-playground/cli/internal/codegen/types"
)

type CPPGenerator struct {
	config *types.Config
}

func NewCPPGenerator(config *types.Config) *CPPGenerator {
	return &CPPGenerator{config: config}
}

func (g *CPPGenerator) Generate() error {
	for _, file := range g.config.Files {
		// Generate header file
		headerPath := filepath.Join(g.config.ProjectName, "source", "include", file.Name+".hpp")
		headerContent := g.generateHeader(file)
		if err := os.WriteFile(headerPath, []byte(headerContent), 0644); err != nil {
			return err
		}

		// Generate source file
		sourcePath := filepath.Join(g.config.ProjectName, "source", "src", file.Name+".cpp")
		sourceContent := g.generateSource(file)
		if err := os.WriteFile(sourcePath, []byte(sourceContent), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (g *CPPGenerator) generateHeader(file types.FileConfig) string {
	var sb strings.Builder

	guardName := strings.ToUpper(file.Name) + "_HPP"
	sb.WriteString(fmt.Sprintf("#ifndef %s\n", guardName))
	sb.WriteString(fmt.Sprintf("#define %s\n\n", guardName))
	sb.WriteString("#include <string>\n\n")

	// Generate class definitions
	for _, typ := range g.config.Types {
		sb.WriteString(fmt.Sprintf("class %s {\n", typ.Name))

		// Private members by default
		sb.WriteString("private:\n")
		for _, field := range typ.Fields {
			if field.Access != "public" && field.Access != "protected" {
				sb.WriteString(fmt.Sprintf("    %s %s;\n", field.Type, field.Name))
			}
		}

		// Protected members
		hasProtected := false
		for _, field := range typ.Fields {
			if field.Access == "protected" {
				if !hasProtected {
					sb.WriteString("\nprotected:\n")
					hasProtected = true
				}
				sb.WriteString(fmt.Sprintf("    %s %s;\n", field.Type, field.Name))
			}
		}

		// Public members
		sb.WriteString("\npublic:\n")
		// Constructor
		sb.WriteString(fmt.Sprintf("    %s() = default;\n", typ.Name))

		// Public fields
		for _, field := range typ.Fields {
			if field.Access == "public" {
				sb.WriteString(fmt.Sprintf("    %s %s;\n", field.Type, field.Name))
			}
		}

		// Method declarations
		for _, method := range typ.Methods {
			params := make([]string, len(method.Parameters))
			for i, param := range method.Parameters {
				params[i] = fmt.Sprintf("%s %s", param.Type, param.Name)
			}
			returnType := method.ReturnType
			if returnType == "" {
				returnType = "void"
			}
			sb.WriteString(fmt.Sprintf("    %s %s(%s);\n",
				returnType,
				method.Name,
				strings.Join(params, ", ")))
		}

		sb.WriteString("};\n\n")
	}

	sb.WriteString(fmt.Sprintf("\n#endif // %s\n", guardName))
	return sb.String()
}

func (g *CPPGenerator) generateSource(file types.FileConfig) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("#include \"%s.hpp\"\n\n", file.Name))

	// Generate method implementations for each class
	for _, typ := range g.config.Types {
		for _, method := range typ.Methods {
			params := make([]string, len(method.Parameters))
			for i, param := range method.Parameters {
				params[i] = fmt.Sprintf("%s %s", param.Type, param.Name)
			}
			returnType := method.ReturnType
			if returnType == "" {
				returnType = "void"
			}

			sb.WriteString(fmt.Sprintf("%s %s::%s(%s) {\n",
				returnType,
				typ.Name,
				method.Name,
				strings.Join(params, ", ")))

			// Add default return statement
			if returnType != "void" {
				switch returnType {
				case "int":
					sb.WriteString("    return 0;\n")
				case "double":
					sb.WriteString("    return 0.0;\n")
				case "string":
					sb.WriteString("    return \"\";\n")
				case "bool":
					sb.WriteString("    return false;\n")
				default:
					sb.WriteString("    return {};\n")
				}
			}

			sb.WriteString("}\n\n")
		}
	}

	// Generate standalone function implementations
	for _, fn := range file.Functions {
		params := make([]string, len(fn.Parameters))
		for i, param := range fn.Parameters {
			params[i] = fmt.Sprintf("%s %s", param.Type, param.Name)
		}
		returnType := fn.ReturnType
		if returnType == "" {
			returnType = "void"
		}

		sb.WriteString(fmt.Sprintf("%s %s(%s) {\n",
			returnType,
			fn.Name,
			strings.Join(params, ", ")))

		// Add default return statement
		if returnType != "void" {
			switch returnType {
			case "int":
				sb.WriteString("    return 0;\n")
			case "double":
				sb.WriteString("    return 0.0;\n")
			case "string":
				sb.WriteString("    return \"\";\n")
			case "bool":
				sb.WriteString("    return false;\n")
			default:
				sb.WriteString("    return {};\n")
			}
		}

		sb.WriteString("}\n\n")
	}

	return sb.String()
}
