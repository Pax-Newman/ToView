package render

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/Pax-Newman/toview/internal/parse"
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

func RenderFile(cmd *cobra.Command, data parse.FileData) string {
	flagAll, _ := cmd.Flags().GetBool("all")
	renderStr := ""
	// skip if the struct is empty
	// this would occur if there was an error while parsing one of the files
	if reflect.ValueOf(data).IsZero() {
		return ""
	}

	// TODO consider if there should be a config for reporting the relative path instead
	// report the filename

	// filter out empty categories
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
		return ""
	}
	renderStr += fmt.Sprintf("# %s\n", filepath.Base(data.FilePath))

	for _, category := range data.Categories {
		renderStr += RenderCategory(cmd, category.Name, category.Comments)
	}
	return renderStr
}
