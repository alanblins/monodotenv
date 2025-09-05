/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io/fs"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "monodotenv",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var DefaultConfigFile = "monodotenv.yaml"
var DefaultUserFile = ".monodotenv.user.yaml"
var DefaultSecretsFile = ".monodotenv.secrets.yaml"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type RealMyOs struct {
}

func (myOs *RealMyOs) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(path)
}

func (myOs *RealMyOs) IsNotExist(error error) bool {
	return os.IsNotExist(error)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.multenv.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
