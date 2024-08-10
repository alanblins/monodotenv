/*
Copyright © 2024 Alan Lins <alanblins@gmail.com>
*/
package cmd

import (
	"io/fs"
	"log"
	"os"

	"github.com/alanblins/monodotenv/models"
	"github.com/alanblins/monodotenv/utils"
	"github.com/spf13/cobra"
)

var ForceFlag bool

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var DefaultConfigFile = "monodotenv.yaml"
var DefaultUserFile = ".monodotenv.user.yaml"

type RealMyOs struct {
}

func (myOs *RealMyOs) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}

func (myOs *RealMyOs) IsNotExist(error error) bool {
	return os.IsNotExist(error)
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
		var configYaml models.ConfigYaml
		var userFile map[string]string

		utils.ReadYaml(DefaultConfigFile, &configYaml)
		utils.ReadYaml(DefaultUserFile, &userFile)

		workspace := args[0]
		outputEnvMap := make(map[string]string)

		nonExistingPaths := map[string]bool{}
		envsExisting := map[string]bool{}
		for _, element := range configYaml.EnvironmentVariables {
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
				utils.WriteContent(element, &configYaml, workspace, userFile, outputEnvMap, path)
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
			log.Fatalln("Delete them or you for force overwrite with option -f: mde use <workspace> -f")
		}

		if ForceFlag || len(envsExisting) == 0 {
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
