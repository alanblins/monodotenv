package models

type EnvironmentVariable struct {
	Key         string            `yaml:"key"`
	Source      string            `yaml:"source"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Workspaces  map[string]string `yaml:"workspaces"`
	Paths       []string          `yaml:"paths"`
}

type Models struct {
	Extends              map[string]string     `yaml:"extends"`
	ProjectFolder        string                `yaml:"projectFolder"`
	EnvironmentVariables []EnvironmentVariable `yaml:"environment_variables"`
}
