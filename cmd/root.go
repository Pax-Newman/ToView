/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/Pax-Newman/toview/parse"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ToView",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	// TODO write custom validator to accept any number of valid file paths
	// TODO add --ignore-unsupported flag to skip unsupported filetypes
	Args: cobra.ArbitraryArgs,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)

		datas := []parse.FileData{}
		// parse data for each file in args
		for _, path := range args {
			data, err := parse.LineByLine(path)
			if err != nil {
				cobra.CompError(err.Error())
			}
			datas = append(datas, data)
		}

		// init the string that we will render and display
		renderStr := ""

		// prepare data from each file for the render
		for _, data := range datas {
			// skip if the struct is empty
			// this would occur if there was an error while parsing one of the files
			if reflect.ValueOf(data).IsZero() {
				continue
			}

			// report the filename
			// TODO consider if there should be a config for reporting the relative path instead
			renderStr += fmt.Sprintf("# %s\n", filepath.Base(data.FilePath))

			// check if there's anything to report in the file
			if len(data.ToDo) <= 0 && len(data.FixMe) <= 0 {
				renderStr += "### No comments available to report\n"
				continue
			}

			// check for and display TODOs in the file
			renderStr += "## To Do\n"
			if len(data.ToDo) > 0 {
				for _, todo := range data.ToDo {
					renderStr += fmt.Sprintf(" - __%d:__ %s\n", todo.Position, todo.Content)
				}
			} else {
				renderStr += "### No ToDos to report\n"
			}

			// check for and display FIXMEs in the file
			renderStr += "## Fix Me\n"
			if len(data.FixMe) > 0 {
				for _, fixme := range data.FixMe {
					renderStr += fmt.Sprintf(" - __%d:__ %s\n", fixme.Position, fixme.Content)
				}
			} else {
				renderStr += "### No FixMes to report\n"
			}
		}

		// FIXME handle the err from the render method
		// TODO allow users to set their own custom style in a config
		out, _ := glamour.Render(renderStr, "dark")
		fmt.Println(out)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ToView.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// TODO add flag for "quiet output" i.e only show true output, no loading bars
}
