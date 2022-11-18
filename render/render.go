package render

import (
	"fmt"

	"github.com/Pax-Newman/toview/parse"
	"github.com/spf13/cobra"
)

func RenderCategory(cmd *cobra.Command, category string, items []parse.Comment) string {
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
