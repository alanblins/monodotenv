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
	value, _ := GetValue(ev, userfile, "", "local")
	if value != "value local" {
		t.Fatal("value didn't match")
	}
}

func TestGetValueNotFound(t *testing.T) {
	userfile := make(map[string]string)
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
	_, err := GetValue(ev, userfile, "", "stage")
	if err.Error() != "not found value for the key: MYKEY and workspace: stage" {
		t.Fatal("didn't throw error")
	}
}

func TestGetValueExtend(t *testing.T) {
	userfile := make(map[string]string)
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
	result, _ := GetValue(ev, userfile, "local", "stage")
	if result != "value local" {
		t.Fatal("value didn't match")
	}
}

func TestGetValueUserFile(t *testing.T) {
	userfile := map[string]string{"USERLOCAL": "my user value"}
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
	value, _ := GetValue(ev, userfile, "", "local")

	if value != "value local" {
		t.Fatal("value didn't match")
	}
}
