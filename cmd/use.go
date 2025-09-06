/*
Copyright © 2024 Alan Lins <alanblins@gmail.com>
*/
package cmd

import (
	"log"

	"github.com/alanblins/monodotenv/models"
	"github.com/alanblins/monodotenv/utils"
	"github.com/spf13/cobra"
)

var ForceFlag bool
var SuffixFlag string

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [environment]",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Short: "Create .env files to multiple folders for an specific environment",
	Long: `Add the environment and environment variables into monodotenv.yaml. Ex:

		environment_variables:
		- name: Server URL
		description: Url of the web server
		key: SERVER_URL
		source: value
		environments:
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
		var configYaml models.ConfigYaml
		var userFile map[string]string
		var secretsFile models.SecretsYaml

		err := utils.ReadYaml(DefaultConfigFile, &configYaml)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		utils.ReadYaml(DefaultUserFile, &userFile)
		utils.ReadYaml(DefaultSecretsFile, &secretsFile)

		environment := args[0]
		outputEnvMap := make(map[string]string)

		nonExistingPaths := map[string]bool{}
		envsExisting := map[string]bool{}
		for _, element := range configYaml.EnvironmentVariables {

			if element.Paths == nil {
				var path = "./"
				envPath := "./.env"

				exist, err := utils.IsFileExist(path, &RealMyOs{})
				if !exist || err != nil {
					nonExistingPaths[path] = true
				}
				exist, err = utils.IsFileExist(envPath, &RealMyOs{})
				if exist && err == nil {
					envsExisting[envPath] = true
				}
				utils.WriteContent(element, &configYaml, environment, outputEnvMap, path, userFile, secretsFile)
			} else {
				for _, path := range element.Paths {
					envPath := path + "/.env"

					exist, err := utils.IsFileExist(path, &RealMyOs{})
					if !exist || err != nil {
						nonExistingPaths[path] = true
					}
					exist, err = utils.IsFileExist(envPath, &RealMyOs{})
					if exist && err == nil {
						envsExisting[envPath] = true
					}
					utils.WriteContent(element, &configYaml, environment, outputEnvMap, path, userFile, secretsFile)
				}
			}
		}

		if len(nonExistingPaths) > 0 {
			log.Println("The following paths don't exist: ")
			for path := range nonExistingPaths {
				log.Println(path)
			}
			log.Fatalln("Please create these folders.")
		}
		if !ForceFlag && len(envsExisting) > 0 {
			log.Println("There are existing .env files on the following directories:")
			for path := range envsExisting {
				log.Println(path)
			}
			log.Fatalln("Delete them or you for force overwrite with option -f: mde use <environment> -f")
		}

		if ForceFlag || len(envsExisting) == 0 {
			for path, output := range outputEnvMap {
				outputBytes := []byte(output)
				err := utils.WriteFile(path, outputBytes, SuffixFlag)
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
	useCmd.PersistentFlags().StringVarP(&SuffixFlag, "suffix", "s", "", "suffix for .env files, ex: -s production will create .env.production")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
