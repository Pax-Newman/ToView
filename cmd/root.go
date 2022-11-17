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

func renderCategory(cmd *cobra.Command, category string, items []parse.Comment) string {
	// get flags
	flagAll, _ := cmd.Flags().GetBool("all")

	// init the string we'll build and return
	renderStr := ""

	// check if category is empty
	if len(items) <= 0 {
		if flagAll {
			// report anyways if --all is set
			renderStr += fmt.Sprintf("## %s\n", category)
			renderStr += "### No comments to report\n"
		}
		return renderStr
	}
	// if there are items then populate our render string with them
	renderStr += fmt.Sprintf("## %s\n", category)
	for _, item := range items {
		renderStr += fmt.Sprintf(" - __%d:__ %s\n", item.Position, item.Content)
	}
	return renderStr
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "toview filepath ...",
	Short: "A utility to find and display ToDo items in your code",
	Long: `toview is a CLI utility to parse files for TODO and FIXME
	comments, rendering them in customizable markdown`,

	Args: cobra.MatchAll(
		// ensure there is at least one arg
		cobra.MinimumNArgs(1),
		// ensure all of the args are valid and supported files
		func(cmd *cobra.Command, args []string) error {
			debug, _ := cmd.Flags().GetBool("debug")
			if debug {
				fmt.Printf("Args: %v\n", args)
				cmd.DebugFlags()
			}
			ignore_unsupported, _ := cmd.Flags().GetBool("ignore-unsupported")

			for _, arg := range args {
				if err := parse.CheckValid(arg); err != nil && !ignore_unsupported {
					return err
				}
			}
			return nil
		},
	),

	Run: func(cmd *cobra.Command, args []string) {
		// check if the debug flag has been set
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			cobra.CompErrorln(err.Error())
		}

		// FIXME implement --flagAll into the rendering
		flagAll, _ := cmd.Flags().GetBool("all")

		datas := []parse.FileData{}
		// parse data for each file in args

		// TODO move this to a config?
		// TODO add a flag to define additional catergories?
		categories := []parse.Category{
			{
				Name:        "To Do",
				ParseTarget: "TODO",
				Comments:    []parse.Comment{},
			},
			{
				Name:        "Fix Me",
				ParseTarget: "FIXME",
				Comments:    []parse.Comment{},
			},
		}

		for _, path := range args {
			data, err := parse.LineByLine(path, categories)
			if err != nil && debug {
				cobra.CompErrorln(err.Error())
			}
			datas = append(datas, data)
		}

		// init the string that we will render and display
		renderStr := ""

		// TODO split rendering into multiple functions, consider breaking into another file
		// prepare data from each file for the render
		for _, data := range datas {
			// skip if the struct is empty
			// this would occur if there was an error while parsing one of the files
			if reflect.ValueOf(data).IsZero() {
				continue
			}

			// TODO consider if there should be a config for reporting the relative path instead
			// report the filename

			hasItems := []parse.Category{}
			for _, category := range data.Categories {
				if len(category.Comments) > 0 {
					hasItems = append(hasItems, category)
				}
			}
			if len(hasItems) <= 0 {
				if flagAll {
					renderStr += fmt.Sprintf("# %s\n", filepath.Base(data.FilePath))
					renderStr += "### No comments to report on yet\n"
				}
				continue
			}
			renderStr += fmt.Sprintf("# %s\n", filepath.Base(data.FilePath))

			for _, category := range data.Categories {
				renderStr += renderCategory(cmd, category.Name, category.Comments)
			}
		}

		// FIXME handle the err from the render method?
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
	rootCmd.Flags().BoolP("debug", "d", false, "Enable debugging output")
	rootCmd.Flags().BoolP("ignore-unsupported", "i", false, "Skips any unsupported files without stopping execution")
	rootCmd.Flags().BoolP("all", "a", false, "Displays all files and categories even if empty")

	// TODO add flag for "quiet output" i.e only show true output, no loading bars
}
