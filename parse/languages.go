package parse

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
