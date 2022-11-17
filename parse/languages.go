package parse

import (
	"fmt"
)

// TODO consider adding colors to each language for rendering
type language struct {
	name       string
	inline     string
	blockStart string
	blockEnd   string
}

var languages = map[string]language{
	"py": {
		name:       "Python",
		inline:     "#",
		blockStart: "",
		blockEnd:   "",
	},
	"go": {
		name:       "Go",
		inline:     "//",
		blockStart: "",
		blockEnd:   "",
	},
}

// defines an error when trying to accces a filetype that is not yet supported
type NotSupportedError struct {
	Filetype string
}

func (e NotSupportedError) Error() string {
	return fmt.Sprintf("file extension \"%s\" not currently supported", e.Filetype)
}

// Returns an error if the filetype is not defined in languages var, nil otherwise
func GetLanguage(filetype string) (language, error) {
	lang := languages[filetype]
	if lang.name == "" {
		return language{}, NotSupportedError{filetype}
	}
	return lang, nil
}
