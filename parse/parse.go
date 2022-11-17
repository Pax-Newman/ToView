package parse

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type FileData struct {
	FilePath string
	ToDo     []Comment
	FixMe    []Comment
}

type Comment struct {
	Title    string
	Content  string
	Position int
}

func newInlineParser(lang language) *regexp.Regexp {
	return regexp.MustCompile(
		lang.inline + " *(?P<title>TODO|FIXME) *(?P<content>.*)",
	)
}

func GetExtension(path string) (string, error) {
	ext := filepath.Ext(path)
	if ext == "" {
		return "", fmt.Errorf("file \"%s\" does not have an extension and cannot be parsed", path)
	}
	return ext[1:], nil
}

// Parse a file line by line for its comments
//
// Returns a FileData containing the parsed file's data
func LineByLine(path string) (FileData, error) {
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
	inlineParser := newInlineParser(lang)

	// parse file line by line, adding any TODO or FIXME comments
	TODOs := []Comment{}
	FIXMEs := []Comment{}
	pos := 1
	for scanner.Scan() {
		match := inlineParser.FindStringSubmatch(scanner.Text())

		if match != nil {
			comment := Comment{
				Title:    match[1],
				Content:  match[2],
				Position: pos,
			}
			switch comment.Title {
			case "TODO":
				TODOs = append(TODOs, comment)
			case "FIXME":
				FIXMEs = append(FIXMEs, comment)
			}
		}
		pos++
	}

	data := FileData{
		FilePath: path,
		ToDo:     TODOs,
		FixMe:    FIXMEs,
	}

	return data, nil
}
