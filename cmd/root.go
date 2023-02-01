/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Pax-Newman/toview/internal/configuration"
	"github.com/Pax-Newman/toview/internal/filehelpers"
	"github.com/Pax-Newman/toview/internal/parse"
	"github.com/Pax-Newman/toview/internal/render"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

type config struct {
	Categories map[string]parse.Category
	Styles     map[string]render.Style
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
			// check if the debug flag is on
			debug, _ := cmd.Flags().GetBool("debug")
			if debug {
				fmt.Printf("Args: %v\n", args)
				cmd.DebugFlags()
			}
			ignore_unsupported, _ := cmd.Flags().GetBool("ignore-unsupported")

			// init a parser so we can check what languages are suppported
			parser, err := parse.InitParser()
			if err != nil {
				return err
			}

			for _, arg := range args {
				// check if the file exists and has an extension
				if err := filehelpers.CheckValid(arg); err != nil {
					return err
				}
				// check if the file is supported
				ext, _ := filehelpers.GetExtension(arg)
				if _, err := parser.GetLanguage(ext); err != nil && !ignore_unsupported {
					return err
				}
			}
			return nil
		},
	),

	RunE: func(cmd *cobra.Command, args []string) error {

		// check if the debug flag has been set
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			cobra.CompErrorln(err.Error())
		}

		// Load our config file
		// TODO create a config file if one doesn't exist
		config, err := configuration.UnmarshalFromPath[config]("config.toml")
		if err != nil {
			cobra.CompErrorln(err.Error())
		}

		// Load our categories from config
		categories := []parse.Category{}
		for _, cat := range config.Categories {
			categories = append(categories, cat)
		}
		filedatas := []parse.FileData{}

		parser, err := parse.InitParser()
		if err != nil {
			return err
		}

		// parse each file given to us
		for _, path := range args {
			data, err := parser.LineByLine(path, categories)
			// print parsing errors if our debug flag is set
			if err != nil && debug {
				cobra.CompErrorln(err.Error())
			}
			filedatas = append(filedatas, data)
		}

		// init the string that we will render and display
		renderStr := ""

		// prepare data from each file for the render
		for _, data := range filedatas {
			renderStr += render.RenderFile(cmd, data)
		}

		// FIXME handle the err from the render method?
		// TODO allow users to set their own custom style in a config
		// TODO consider replacing glamour with lipgloss for custom styling

		raw, _ := cmd.Flags().GetBool("raw")
		if raw {
			fmt.Println(renderStr)
		} else {
			out, _ := glamour.Render(renderStr, "dark")
			fmt.Println(out)
		}

		return nil
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
	rootCmd.Flags().BoolP("debug", "d", false, "Enable debugging output")
	rootCmd.Flags().BoolP("ignore-unsupported", "i", false, "Skips any unsupported files without stopping execution")
	rootCmd.Flags().BoolP("all", "a", false, "Displays all files and categories even if empty")

	// TODO add flag for "quiet output" i.e only show true output, no loading bars

	// TODO add a flag for raw output without any styling
	rootCmd.Flags().BoolP("raw", "r", false, "Output raw text without styling")

	// TODO add a command to output to a file
	// TODO add a flag to render each file's data to a seperate md file
	// i.e. toview render -s/--seperate main.go commands.go ---> main.md commands.md
	// TODO add --style flag for specifying which style to use
}
