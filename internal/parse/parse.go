package parse

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Pax-Newman/toview/internal/configuration"
	"github.com/Pax-Newman/toview/internal/filehelpers"
	"github.com/spf13/viper"
)

type Language struct {
	Name       string
	Inline     string
	BlockStart string
	BlockEnd   string
}

// defines an error when trying to accces a filetype that is not yet supported
type NotSupportedError struct {
	Filetype string
}

func (e NotSupportedError) Error() string {
	return fmt.Sprintf("file extension \"%s\" not currently supported", e.Filetype)
}

// Holds data relating to the data parsed from a file
type FileData struct {
	FilePath   string
	Categories []Category
}

// Holds data relating to group of comments of a specific type
type Category struct {
	// name of the category to be displayed
	Name string

	// string to look for at the beginning of a comment
	// e.g. for ToDos it should be "TODO"
	ParseTarget string

	// slice that stores the comment data for this category
	Comments []Comment
}

// Holds data relating to a comment
type Comment struct {
	Title    string
	Content  string
	Position int
}

type LanguageConfig struct {
	Languages map[string]Language
}

var languages LanguageConfig

// Initialize the language definitions from the config file
func Init() error {
	// config, err := configuration.LoadConfig("languages.toml")
	// if err != nil {
	// 	return err
	// }
	// config.Unmarshal(&languages)
	// return nil
	var err error
	languages, err = configuration.UnmarshalFromConfig[LanguageConfig]("languages.toml")

	return err
}

func newInlineParser(lang Language, categories []Category) *regexp.Regexp {
	targets := []string{}
	for _, category := range categories {
		targets = append(targets, category.ParseTarget)
	}
	targetStr := strings.Join(targets, "|")
	return regexp.MustCompile(
		lang.Inline + fmt.Sprintf(" *(?P<title>%s) *(?P<content>.*)", targetStr),
	)
}

// Check if a file is valid, supported, and exists
func CheckValid(path string) error {
	// check if the filepath is valid
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	// check if the file extension is valid
	ext, err := filehelpers.GetExtension(path)
	if err != nil {
		return err
	}

	// check if the file extension is supported
	if _, err := GetLanguage(ext); err != nil {
		return err
	}

	return nil
}

// Returns an error if the filetype is not defined in languages var, nil otherwise
func GetLanguage(filetype string) (Language, error) {
	lang := languages.Languages[filetype]

	if lang.Name == "" {
		return Language{}, NotSupportedError{filetype}
	}
	return lang, nil
}

// Parse a file line by line for its comments
//
// Returns a FileData containing the parsed file's data
func LineByLine(path string, categories []Category, languages *viper.Viper) (FileData, error) {
	filetype, err := filehelpers.GetExtension(path)
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
		// TODO handle languages that have start/end comment markers
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
