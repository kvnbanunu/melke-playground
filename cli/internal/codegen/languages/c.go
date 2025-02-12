package languages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kvnbanunu/melke-playground/cli/internal/codegen/types"
)

type CGenerator struct {
	config *types.Config
}

func NewCGenerator(config *types.Config) *CGenerator {
	return &CGenerator{config: config}
}

func (g *CGenerator) Generate() error {
	for _, file := range g.config.Files {
		// Generate header file
		headerPath := filepath.Join(g.config.ProjectName, "source", "include", file.Name+".h")
		headerContent := g.generateHeader(file)
		if err := os.WriteFile(headerPath, []byte(headerContent), 0644); err != nil {
			return err
		}

		// Generate source file
		sourcePath := filepath.Join(g.config.ProjectName, "source", "src", file.Name+".c")
		sourceContent := g.generateSource(file)
		if err := os.WriteFile(sourcePath, []byte(sourceContent), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (g *CGenerator) generateHeader(file types.FileConfig) string {
	var sb strings.Builder

	guardName := strings.ToUpper(file.Name) + "_H"
	sb.WriteString(fmt.Sprintf("#ifndef %s\n", guardName))
	sb.WriteString(fmt.Sprintf("#define %s\n\n", guardName))

	// Generate struct definitions
	for _, typ := range g.config.Types {
		sb.WriteString(fmt.Sprintf("typedef struct %s {\n", typ.Name))
		for _, field := range typ.Fields {
			sb.WriteString(fmt.Sprintf("    %s %s;\n", field.Type, field.Name))
		}
		sb.WriteString(fmt.Sprintf("} %s;\n\n", typ.Name))
	}

	// Generate function declarations
	for _, fn := range file.Functions {
		params := make([]string, len(fn.Parameters))
		for i, param := range fn.Parameters {
			params[i] = fmt.Sprintf("%s %s", param.Type, param.Name)
		}
		returnType := fn.ReturnType
		if returnType == "" {
			returnType = "void"
		}
		sb.WriteString(fmt.Sprintf("%s %s(%s);\n",
			returnType,
			fn.Name,
			strings.Join(params, ", ")))
	}

	sb.WriteString(fmt.Sprintf("\n#endif // %s\n", guardName))
	return sb.String()
}

func (g *CGenerator) generateSource(file types.FileConfig) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("#include \"%s.h\"\n\n", file.Name))

	// Generate function implementations
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

		// Add placeholder return statement if needed
		if returnType != "void" {
			switch {
			case strings.Contains(returnType, "int"):
				sb.WriteString("    return 0;\n")
			case strings.Contains(returnType, "char"):
				sb.WriteString("    return '\\0';\n")
			case strings.Contains(returnType, "*"):
				sb.WriteString("    return NULL;\n")
			default:
				sb.WriteString("    return 0;\n")
			}
		}

		sb.WriteString("}\n\n")
	}

	return sb.String()
}
