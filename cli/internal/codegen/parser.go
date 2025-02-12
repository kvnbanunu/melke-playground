package codegen

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/kvnbanunu/melke-playground/cli/internal/codegen/types"
)

func ParseConfig(filename string) (*types.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &types.Config{
		Language: "C", // Default language
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
