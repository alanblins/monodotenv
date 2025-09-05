/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/alanblins/monodotenv/models"
	"github.com/alanblins/monodotenv/utils"

	"github.com/spf13/cobra"
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc [workspace]",
	Short: "Create markdown doc for environment variables",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var configYaml models.ConfigYaml
		var userFile map[string]string
		var secretsFile models.SecretsYaml

		header := []string{"Key", "Name", "Description", "Path"}

		err := utils.ReadYaml(DefaultConfigFile, &configYaml)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		utils.ReadYaml(DefaultUserFile, &userFile)
		utils.ReadYaml(DefaultSecretsFile, &secretsFile)

		mapWorskpaces := make(map[string]string)

		workspace := ""
		if len(args) > 0 {
			workspace = args[0]
		}
		if workspace == "" {
			for _, ev := range configYaml.EnvironmentVariables {
				for keyWorkspace := range ev.Workspaces {
					mapWorskpaces[keyWorkspace] = keyWorkspace
				}
			}
		} else {
			mapWorskpaces[workspace] = workspace
		}

		workspaces := []string{}

		for keyWorkspace := range mapWorskpaces {
			header = append(header, keyWorkspace)
			workspaces = append(workspaces, keyWorkspace)
		}

		contents := [][]string{}

		for _, element := range configYaml.EnvironmentVariables {
			paths := []string{"./"}
			if element.Paths != nil {
				paths = element.Paths
			}
			for _, path := range paths {
				contents = utils.WriteContentDocLine(contents, element, &configYaml, workspaces, path, userFile, secretsFile)
			}
		}

		columnWidth := []int{}

		for _, value := range header {
			columnWidth = append(columnWidth, len(value))
		}

		for _, valueArr := range contents {
			for index, value := range valueArr {
				columnWidth[index] = max(len(value), columnWidth[index])
			}
		}

		headerMark := renderTextLine(header, columnWidth)
		headerDashMark := renderDashes(columnWidth)

		fmt.Println(headerMark)
		fmt.Println(headerDashMark)

		for _, valueArr := range contents {
			contentMark := renderTextLine(valueArr, columnWidth)
			fmt.Println(contentMark)
		}

	},
}

func renderTextLine(contents []string, columnWidth []int) string {
	content := "|"
	for index, value := range contents {
		repeat := columnWidth[index] - len(value)
		content = content + " " + value + strings.Repeat(" ", repeat) + " |"
	}
	return content
}

func renderDashes(columnWidth []int) string {
	headerDashMark := "|"
	for _, value := range columnWidth {
		headerDashMark = headerDashMark + strings.Repeat("-", value+1) + " |"
	}
	return headerDashMark
}

func init() {
	rootCmd.AddCommand(docCmd)
	// docCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
