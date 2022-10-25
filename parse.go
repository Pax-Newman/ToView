package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type language struct {
	inline     string
	blockStart string
	blockEnd   string
}

type item struct {
	title    string
	content  string
	position int
}

// load and return text data from a given file
func loadFileData(path string) string {
	return ""
}

// find comments in filedata and return a list of them
func findComments(data string, lang language) []string {
	comments := []string{}

	if lang.inline != "" {
		inlinePattern := regexp.MustCompile(lang.inline + " *(?P<title>TODO|FIXME) *(?P<content>.*)")
		// retreive comment data
		inlineMatches := inlinePattern.FindAllStringSubmatch(data, -1)

		for i := 0; i < len(inlineMatches); i++ {
			match := inlineMatches[i]
			comment := item{
				title:    match[1],
				content:  match[2],
				position: -1,
			}
		}
	}

	return comments
}

func newInlineParser(lang language) *regexp.Regexp {
	return regexp.MustCompile(
		lang.inline + " *(?P<title>TODO|FIXME) *(?P<content>.*)",
	)
}

// Parse a file line by line for its comments
// Returns a slice where index 0 = slice of TODO items, 1 = slice of FIXME items
func LineByLine(path string) [][]item {
	// open file & close it on function end
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	filetype := filepath.Ext(path)

	// FIXME get language from filetype
	lang := language{}

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

	return comments
}
