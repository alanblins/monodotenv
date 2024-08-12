/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/alanblins/monodotenv/utils"
	"github.com/spf13/cobra"
)

var DecryptFlag bool

// encCmd represents the doc command
var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "Encrypt password",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if DecryptFlag {
			cipherText := args[0]
			keyString := args[1]
			nonce := args[2]
			textDecrypted := utils.GCMDecrypter(keyString, cipherText, nonce)
			fmt.Println(textDecrypted)
		} else {
			textString := args[0]
			keyString := args[1]
			nonceHex := ""
			if len(args) > 2 {
				nonceHex = args[2]
			}
			cipherText, nonce := utils.GCMEncrypter(keyString, textString, nonceHex)
			fmt.Printf("Ciphertext: %s\n", cipherText)
			fmt.Printf("key: %s\n", keyString)
			fmt.Printf("nonce: %s\n", nonce)
		}
	},
}

func init() {
	rootCmd.AddCommand(encCmd)
	encCmd.PersistentFlags().BoolVarP(&DecryptFlag, "decrypt", "d", false, "decrypt")
}
