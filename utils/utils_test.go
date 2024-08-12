package utils

import (
	"io/fs"
	"testing"

	"github.com/alanblins/monodotenv/models"
)

type SpyOs struct {
}

func (myOs *SpyOs) Stat(path string) (fs.FileInfo, error) {
	return nil, nil
}

func (myOs *SpyOs) IsNotExist(error error) bool {
	return true
}

func TestGetValue(t *testing.T) {
	userfile := make(map[string]string)
	secretsFile := models.SecretsYaml{}
	paths := []string{"mypath"}
	workspaces := map[string]string{"local": "value local"}
	ev := models.EnvironmentVariable{
		Key:         "MYKEY",
		Source:      "value",
		Name:        "Any name",
		Description: "Description",
		Paths:       paths,
		Workspaces:  workspaces,
	}
	value, _ := GetValue(ev, "", "local", userfile, secretsFile)
	if value != "value local" {
		t.Fatal("value didn't match")
	}
}

func TestGetValueNotFound(t *testing.T) {
	userfile := make(map[string]string)
	secretsFile := models.SecretsYaml{}
	paths := []string{"mypath"}
	workspaces := map[string]string{"local": "value local"}
	ev := models.EnvironmentVariable{
		Key:         "MYKEY",
		Source:      "value",
		Name:        "Any name",
		Description: "Description",
		Paths:       paths,
		Workspaces:  workspaces,
	}
	_, err := GetValue(ev, "", "stage", userfile, secretsFile)
	if err.Error() != "not found value for the key: MYKEY and workspace: stage" {
		t.Fatal("didn't throw error")
	}
}

func TestGetValueExtend(t *testing.T) {
	userfile := make(map[string]string)
	secretsFile := models.SecretsYaml{}
	paths := []string{"mypath"}
	workspaces := map[string]string{"local": "value local"}
	ev := models.EnvironmentVariable{
		Key:         "MYKEY",
		Source:      "value",
		Name:        "Any name",
		Description: "Description",
		Paths:       paths,
		Workspaces:  workspaces,
	}
	result, _ := GetValue(ev, "local", "stage", userfile, secretsFile)
	if result != "value local" {
		t.Fatal("value didn't match")
	}
}

func TestGetValueUserFile(t *testing.T) {
	userfile := map[string]string{"USERLOCAL": "my user value"}
	secretsFile := models.SecretsYaml{}
	paths := []string{"mypath"}
	workspaces := map[string]string{"local": "USERLOCAL"}
	ev := models.EnvironmentVariable{
		Key:         "MYKEY",
		Source:      "user",
		Name:        "Any name",
		Description: "Description",
		Paths:       paths,
		Workspaces:  workspaces,
	}
	value, _ := GetValue(ev, "", "local", userfile, secretsFile)

	if value != "my user value" {
		t.Fatal("value didn't match")
	}
}

func TestGetValueAESGCM(t *testing.T) {
	passes := []string{"AES256Key-32Characters1234567890", "c835baf3e8b83e5a75a11430"}
	secretsEnvs := map[string][]string{"MYKEY": passes}
	secretsWorkspace := map[string]map[string][]string{"local": secretsEnvs}
	userfile := map[string]string{}
	secretsFile := models.SecretsYaml{
		Secrets: secretsWorkspace,
	}
	paths := []string{"mypath"}
	workspaces := map[string]string{"local": "7e57cc75b26e0e308c67eb81680572639868c51a3cbb"}
	ev := models.EnvironmentVariable{
		Key:         "MYKEY",
		Source:      "aes-gcm",
		Name:        "Any name",
		Description: "Description",
		Paths:       paths,
		Workspaces:  workspaces,
	}
	value, _ := GetValue(ev, "", "local", userfile, secretsFile)

	if value != "mytext" {
		t.Fatal("value didn't match")
	}
}
