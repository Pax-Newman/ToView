package main

import (
	"regexp"
)

type lang struct {
	inline     string
	blockStart string
	blockEnd   string
}

// load and return text data from a given file
func loadFileData(path string) string {
	return ""
}

// find comments in filedata and return a list of them
func findComments(data string, fileType lang) []string {
	comments := []string{}

	if fileType.inline != "" {
		inlinePattern := regexp.MustCompile(fileType.inline + " *(?P<title>TODO|FIXME) *(?P<content>.*)")
		// retreive comment data
		inlineMatches := inlinePattern.FindAllStringSubmatch(data, -1)

		for i := 0; i < len(inlineMatches); i++ {
			match = inlineMatches[i]
		}
	}

	return comments
}

func ParseLine() {

}

func LineByLine() {

}
