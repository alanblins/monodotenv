/*
Copyright © 2024 Alan Lins <alanblins@gmail.com>
*/
package utils

import (
	"errors"
	"log"
	"os"

	"github.com/alanblins/monodotenv/models"
	"gopkg.in/yaml.v3"
)

func IsFileExist(path string, myOs models.OsI) (bool, error) {
	_, error := myOs.Stat(path)
	if error == nil {
		return true, nil
	}
	return false, error

}

func ValidateEnvironmanetVariable(environmentVariable models.EnvironmentVariable, workspace string, extendWorkspace string, userFile map[string]string) (bool, error) {

	return true, nil
}

func GetValue(environmentVariable models.EnvironmentVariable, userfile map[string]string, extendWorkspace string, workspace string) (string, error) {
	value, workspaceExist := environmentVariable.Workspaces[workspace]
	if !workspaceExist {
		if extendWorkspace != "" {
			value, workspaceExist = environmentVariable.Workspaces[extendWorkspace]
			if !workspaceExist {
				return "", errors.New("not found value for the key: " + environmentVariable.Key + " and workspace: " + workspace)
			}
		} else {
			return "", errors.New("not found value for the key: " + environmentVariable.Key + " and workspace: " + workspace)
		}
	}

	if environmentVariable.Source == "value" {
		return value, nil
	}

	if environmentVariable.Source == "user" {
		if userfile == nil {
			return "", errors.New("no .monodotenv.user.yaml file found")
		}
		val, ok := userfile[value]
		if !ok {
			return "", errors.New(value + "not found in .monodotenv.user.yaml file.")
		}
		return val, nil
	}
	return "", errors.New("source unknown")
}

func ReadYaml[T models.ConfigYaml | map[string]string](filepath string, yamlDecoded *T) {
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &yamlDecoded)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func WriteContent(environmentVariable models.EnvironmentVariable, configYaml *models.ConfigYaml, workspace string, userFileYaml map[string]string, outputEnvMap map[string]string, path string) {
	extendWorkspace := configYaml.Extends[workspace]
	value, errorReadValue := GetValue(environmentVariable, userFileYaml, extendWorkspace, workspace)
	if errorReadValue != nil {
		log.Fatalln(errorReadValue)
	}
	content := outputEnvMap[path] + environmentVariable.Key + "=" + value + "\n"
	outputEnvMap[path] = content
}
