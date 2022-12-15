package parse

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Pax-Newman/toview/internal/configuration"
	"github.com/Pax-Newman/toview/internal/filehelpers"
)

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

type Language struct {
	Name       string
	Inline     string
	BlockStart string
	BlockEnd   string
}

type LanguageConfig struct {
	Languages map[string]Language
}

// defines an error when trying to accces a filetype that is not yet supported
type NotSupportedError struct {
	Filetype string
}

func (e NotSupportedError) Error() string {
	return fmt.Sprintf("file extension \"%s\" not currently supported", e.Filetype)
}

type Parser struct {
	languages map[string]Language
	parsers   map[string]*regexp.Regexp
}

func InitParser() (*Parser, error) {
	languages, err := configuration.UnmarshalFromPath[LanguageConfig]("languages.toml")
	if err != nil {
		return nil, err
	}

	parser := Parser{
		languages: languages.Languages,
		parsers:   map[string]*regexp.Regexp{},
	}
	return &parser, nil
}

func (p *Parser) newInlineParser(lang Language, categories []Category) *regexp.Regexp {
	// Create a regex string containing all category target strings e.g. TODO|FIXME
	targets := []string{}
	for _, category := range categories {
		targets = append(targets, category.ParseTarget)
	}
	targetStr := strings.Join(targets, "|")

	// add the new language regex parser to the Parser object's parsers map
	p.parsers[lang.Name] = regexp.MustCompile(
		lang.Inline + fmt.Sprintf(" *(?P<title>%s) *(?P<content>.*)", targetStr),
	)
	return p.parsers[lang.Name]
}

// Returns an error if the filetype is not defined in languages var, nil otherwise
func (p *Parser) GetLanguage(filetype string) (Language, error) {
	lang, exists := p.languages[filetype]

	if !exists {
		return Language{}, NotSupportedError{filetype}
	}
	return lang, nil
}

// Parse a file line by line for its comments
//
// Returns a FileData containing the parsed file's data
func (p *Parser) LineByLine(path string, categories []Category) (FileData, error) {
	filetype, err := filehelpers.GetExtension(path)
	if err != nil {
		return FileData{}, err
	}

	lang, err := p.GetLanguage(filetype)
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
	inlineParser, exists := p.parsers[lang.Name]
	if !exists {
		// create a new parser if we don't already have one for this language
		inlineParser = p.newInlineParser(lang, categories)
	}

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
