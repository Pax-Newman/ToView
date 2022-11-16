package parse

import (
	"bufio"
	"errors"
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

// Parse a file line by line for its comments
// Returns a slice where index 0 = slice of TODO items, 1 = slice of FIXME items
func LineByLine(path string) (FileData, error) {
	filetype := filepath.Ext(path)
	// check if the file has an extension
	if filetype == "" {
		errMsg := fmt.Sprintf("file \"%s\" does not have an extension and cannot be parsed\n", path)
		return FileData{}, errors.New(errMsg)
	} else {
		// remove the dot from the file extension
		filetype = filetype[1:]
	}

	lang := languages[filetype]
	if lang.name == "" {
		errMsg := fmt.Sprintf("file extension \"%s\" not defined in languages.go\n", filetype)
		return FileData{}, errors.New(errMsg)
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
