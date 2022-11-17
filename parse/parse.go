package parse

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FileData struct {
	FilePath   string
	Categories []Category
}

type Category struct {
	// name of the category to be displayed
	Name string

	// string to look for at the beginning of a comment
	// e.g. for ToDos it should be "TODO"
	ParseTarget string

	// slice that stores the comment data for this category
	Comments []Comment
}

type Comment struct {
	Title    string
	Content  string
	Position int
}

func newInlineParser(lang language, categories []Category) *regexp.Regexp {
	targets := []string{}
	for _, category := range categories {
		targets = append(targets, category.ParseTarget)
	}
	targetStr := strings.Join(targets, "|")
	return regexp.MustCompile(
		lang.inline + fmt.Sprintf(" *(?P<title>%s) *(?P<content>.*)", targetStr),
	)
}

func GetExtension(path string) (string, error) {
	ext := filepath.Ext(path)
	if ext == "" {
		return "", fmt.Errorf("file \"%s\" does not have an extension and cannot be parsed", path)
	}
	return ext[1:], nil
}

// Check if a file is valid, supported, and exists
func CheckValid(path string) error {
	// check if the filepath is valid
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	// check if the file extension is valid
	ext, err := GetExtension(path)
	if err != nil {
		return err
	}

	// check if the file extension is supported
	if _, err := GetLanguage(ext); err != nil {
		return err
	}

	return nil
}

// Parse a file line by line for its comments
//
// Returns a FileData containing the parsed file's data
func LineByLine(path string, categories []Category) (FileData, error) {
	filetype, err := GetExtension(path)
	if err != nil {
		return FileData{}, err
	}

	lang, err := GetLanguage(filetype)
	if err != nil {
		return FileData{}, err
	}

	// open file & close it on function end
	file, err := os.Open(path)
	if err != nil {
		return FileData{}, err
	}
	defer file.Close()

	// create a new scanner and comment parser
	scanner := bufio.NewScanner(file)
	inlineParser := newInlineParser(lang, categories)

	// copy the category structs passed in
	categoriesCopy := make([]Category, len(categories))
	copy(categoriesCopy, categories)

	data := FileData{
		FilePath:   path,
		Categories: categoriesCopy,
	}

	catIndexMap := map[string]int{}
	for i, cat := range data.Categories {
		catIndexMap[cat.ParseTarget] = i
	}

	// parse file line by line, adding any TODO or FIXME comments
	pos := 1
	for scanner.Scan() {
		match := inlineParser.FindStringSubmatch(scanner.Text())

		if match != nil {
			comment := Comment{
				Title:    match[1],
				Content:  match[2],
				Position: pos,
			}

			// group comment into its respective category
			catIndex := catIndexMap[comment.Title]
			cat := data.Categories[catIndex]
			data.Categories[catIndex].Comments = append(cat.Comments, comment)
		}
		pos++
	}

	return data, nil
}
