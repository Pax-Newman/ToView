package parse

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type item struct {
	title    string
	content  string
	position int
}

func newInlineParser(lang language) *regexp.Regexp {
	return regexp.MustCompile(
		lang.inline + " *(?P<title>TODO|FIXME) *(?P<content>.*)",
	)
}

// Parse a file line by line for its comments
// Returns a slice where index 0 = slice of TODO items, 1 = slice of FIXME items
func LineByLine(path string) ([][]item, error) {
	// open file & close it on function end
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	filetype := filepath.Ext(path)
	lang := languages[filetype]
	if lang.name == "" {
		errMsg := fmt.Sprintf("file extension \".%s\" not defined in languages.go", filetype)
		return nil, errors.New(errMsg)
	}

	// create a new scanner and comment parser
	scanner := bufio.NewScanner(file)
	inlineParser := newInlineParser(lang)

	// parse file line by line, adding any TODO or FIXME comments
	TODOs := []item{}
	FIXMEs := []item{}
	pos := 1
	for scanner.Scan() {
		match := inlineParser.FindStringSubmatch(scanner.Text())

		comment := item{
			title:    match[1],
			content:  match[2],
			position: pos,
		}
		switch comment.title {
		case "TODO":
			TODOs = append(TODOs, comment)
		case "FIXME":
			FIXMEs = append(FIXMEs, comment)
		}
		pos++
	}

	comments := [][]item{TODOs, FIXMEs}

	return comments, nil
}
