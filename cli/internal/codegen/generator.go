package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kvnbanunu/melke-playground/cli/internal/codegen/languages"
	"github.com/kvnbanunu/melke-playground/cli/internal/codegen/types"
)

type Generator struct {
	config *types.Config
}

func NewGenerator(config *types.Config) *Generator {
	return &Generator{config: config}
}

func (g *Generator) Generate() error {
	if err := g.createDirectories(); err != nil {
		return err
	}

	var generator interface{ Generate() error }

	switch strings.ToLower(g.config.Language) {
	case "c":
		generator = languages.NewCGenerator(g.config)
	case "c++", "cpp":
		generator = languages.NewCPPGenerator(g.config)
	case "python":
		generator = languages.NewPythonGenerator(g.config)
	case "go":
		generator = languages.NewGoGenerator(g.config)
	case "javascript", "js":
		generator = languages.NewJavaScriptGenerator(g.config)
	case "java":
		generator = languages.NewJavaGenerator(g.config)
	default:
		return fmt.Errorf("unsupported language: %s", g.config.Language)
	}

	return generator.Generate()
}

func (g *Generator) createDirectories() error {
	dirs := []string{
		filepath.Join(g.config.ProjectName, "source", "src"),
		filepath.Join(g.config.ProjectName, "source", "include"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
