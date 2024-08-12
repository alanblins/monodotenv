package models

import "io/fs"

type EnvironmentVariable struct {
	Key         string            `yaml:"key"`
	Source      string            `yaml:"source"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Workspaces  map[string]string `yaml:"workspaces"`
	Paths       []string          `yaml:"paths"`
}

type ConfigYaml struct {
	Extends              map[string]string     `yaml:"extends"`
	ProjectFolder        string                `yaml:"projectFolder"`
	EnvironmentVariables []EnvironmentVariable `yaml:"environment_variables"`
}

type EnvironmentPass struct {
}

type SecretsYaml struct {
	Secrets map[string]map[string][]string `yaml:"secrets"`
}

type OsI interface {
	Stat(path string) (fs.FileInfo, error)
	IsNotExist(error error) bool
}
