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

// listCmd represents the use command
var listCmd = &cobra.Command{
	Use:   "list [environment]",
	Short: "List environments available on monodotenv.yaml",
	Long: `List environments available on monodotenv.yaml.
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

		mapEnvironments := make(map[string]string)
		if len(args) == 0 {
			for _, ev := range configYaml.EnvironmentVariables {
				for keyEnvironment := range ev.Environments {
					mapEnvironments[keyEnvironment] = keyEnvironment
				}
			}

			for keyEnvironment := range mapEnvironments {
				println(keyEnvironment)
			}
			return
		}
		environment := args[0]
		outputEnvMap := make(map[string]string)

		nonExistingPaths := map[string]bool{}
		for _, element := range configYaml.EnvironmentVariables {
			if element.Paths == nil {
				var path = "./"
				utils.WriteContent(element, &configYaml, environment, outputEnvMap, path, userFile, secretsFile, true)
			} else {
				for _, path := range element.Paths {
					utils.WriteContent(element, &configYaml, environment, outputEnvMap, path, userFile, secretsFile, true)
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

		for path, output := range outputEnvMap {
			outputBytes := []byte(output)
			utils.DryWriteFile(path, outputBytes, SuffixFlag)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	listCmd.PersistentFlags().StringVarP(&SuffixFlag, "suffix", "s", "", "suffix for .env files, ex: -s production will create .env.production")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
