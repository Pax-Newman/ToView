/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// add config file data here
const defaultConfig = ``
const defaultLanguages = `# Language parsing definitions
[languages]

# To add a new language you can follow this template:
# [languages.$EXT]  --- Replace $EXT with whatever the language's file extension is 
# name = ""         --- The full name of the language
# inline = ""       --- whatever denotes the start of an inline comment
# blockstart = ""   --- whatever denotes the start of a block comment
# blockend = ""     --- whatever denotes the end of a block comment

# If you add support for a language, please submit a pull request
# to https://github.com/Pax-Newman/ToView with your changes

[languages.py]
name = "Python"
inline = "#"
blockstart = ""
blockend = ""

[languages.go]
name = "Go"
inline = "//"
blockstart = ""
blockend = ""`

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("config called")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
