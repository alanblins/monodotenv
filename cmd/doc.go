/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Create markdown doc for environment variables",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("doc called")
		header := []string{"Key", "Name", "Description"}
		header = append(header, "Stage")
		header = append(header, "Local")

		content := [][]string{}
		content1 := []string{"someke", "some key", "key sample bla bal", "stagevalue", "localvalue"}

		content = append(content, content1)

		columnWidth := []int{}

		for _, valueArr := range content {
			for index, value := range valueArr {
				columnWidth = append(columnWidth, 0)
				if columnWidth[index] == 0 || len(value) > columnWidth[index] {
					columnWidth[index] = len(value)
				}
			}
		}

		for index, value := range header {
			if columnWidth[index] == 0 || len(value) > columnWidth[index] {
				columnWidth[index] = len(value)
			}
		}

		headerMark := "|"
		for index, value := range header {
			repeat := columnWidth[index] - len(value)
			headerMark = headerMark + " " + value + strings.Repeat(" ", repeat) + " |"
		}
		headerDashMark := "|"
		for _, value := range columnWidth {
			headerDashMark = headerDashMark + strings.Repeat("-", value+1) + " |"
		}

		fmt.Println(headerMark)
		fmt.Println(headerDashMark)

		for _, valueArr := range content {
			contentMark := "|"
			for index, value := range valueArr {
				repeat := columnWidth[index] - len(value)
				contentMark = contentMark + " " + value + strings.Repeat(" ", repeat) + " |"
			}
			fmt.Println(contentMark)
		}

	},
}

func init() {
	rootCmd.AddCommand(docCmd)
	// docCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
