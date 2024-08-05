/*
Copyright © 2024 Alan Lins <alanblins@gmail.com>
*/
package cmd

import (
	"log"
	"os"

	"github.com/alanblins/monodotenv/models"
	"github.com/alanblins/monodotenv/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ForceFlag bool

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [workspace]",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Create .env files for an specific workspace",
	Long: `Add the workspace into environment variables on multenv.yaml. Ex:

		environment_variables:
		- name: Server URL
		description: Url of the web server
		key: SERVER_URL
		source: value
		workspaces:
			stage: https://stage.myserver.com
			local: http://localhost:1000
		paths:
		- packages/app

	Run the command to populate stage values
	monodotenv use stage

	Run the command to populate local values
	monodotenv use local
	`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckConfigFile()
		workspace := args[0]

		var model models.Models
		var userFile map[string]string

		yamlFile, err := os.ReadFile("monodotenv.yaml")
		yamlUserFile, errUserFile := os.ReadFile(".monodotenv.user.yaml")

		if err != nil {
			log.Fatalf("yamlFile.Get err #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, &model)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}

		if errUserFile == nil {
			errUserFile = yaml.Unmarshal(yamlUserFile, &userFile)
			if errUserFile != nil {
				log.Fatalf("Unmarshal: %v", err)
			}
		}

		currentWorkspace := workspace

		outputEnvMap := make(map[string]string)

		nonExistingPaths := map[string]bool{}
		containsNotExistingPaths := false

		envsExisting := map[string]bool{}
		containsExistingEnvs := false
		for _, element := range model.EnvironmentVariables {
			for _, path := range element.Paths {
				var exist, _ = utils.IsFileExist(path)
				if !exist {
					nonExistingPaths[path] = true
					containsNotExistingPaths = true
				}
				envPath := path + "/.env"
				var envExist, _ = utils.IsFileExist(envPath)
				if envExist {
					envsExisting[envPath] = true
					containsExistingEnvs = true
				}

				extendWorkspace := model.Extends[currentWorkspace]
				extendValue := ""
				if extendWorkspace != "" {
					extendValue = element.Workspaces[extendWorkspace]
				}

				value, errorReadValue := utils.GetValue(element.Workspaces[currentWorkspace], element.Key, element.Source, userFile, extendValue)
				if errorReadValue != nil {
					log.Fatalln(errorReadValue)
				}
				content := outputEnvMap[path] + element.Key + "=" + value + "\n"
				outputEnvMap[path] = content
			}
		}

		if containsNotExistingPaths {
			log.Println("The following paths don't exist: ")
			for path := range nonExistingPaths {
				log.Println(path)
			}
			log.Fatalln("Please create these folders.")
		}
		if !ForceFlag && containsExistingEnvs {
			log.Println("There are existing .env files on the following directories:")
			for path := range envsExisting {
				log.Println(path)
			}
			log.Fatalln("Delete them or you for force overwrite with option -f")
		}

		if ForceFlag || !containsExistingEnvs {
			for path, output := range outputEnvMap {
				outputBytes := []byte(output)
				err := os.WriteFile(path+"/.env", outputBytes, 0644)
				check(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	useCmd.PersistentFlags().BoolVarP(&ForceFlag, "force", "f", false, "overwrites existing .env files")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
