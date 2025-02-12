package types

// Config represents the structure of the YAML configuration file
type Config struct {
	Language    string       `yaml:"language"`
	ProjectName string       `yaml:"projectName"`
	Types       []TypeConfig `yaml:"types"`
	Files       []FileConfig `yaml:"files"`
}

type TypeConfig struct {
	Name    string           `yaml:"name"`
	Fields  []FieldConfig    `yaml:"fields"`
	Methods []FunctionConfig `yaml:"methods"`
}

type FieldConfig struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Access string `yaml:"access"`
}

type FileConfig struct {
	Name      string           `yaml:"name"`
	Functions []FunctionConfig `yaml:"functions"`
}

type FunctionConfig struct {
	Name       string            `yaml:"name"`
	Parameters []ParameterConfig `yaml:"parameters"`
	ReturnType string            `yaml:"returnType"`
	Access     string            `yaml:"access"`
}

type ParameterConfig struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}
