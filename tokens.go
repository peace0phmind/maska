package mask

import "regexp"

type Token struct {
	Pattern   *regexp.Regexp
	Multiple  bool
	Optional  bool
	Repeated  bool
	Transform func(char string) string
}

type Tokens map[string]Token

var defaultTokens = Tokens{
	"#": {Pattern: regexp.MustCompile(`[0-9]`)},
	"@": {Pattern: regexp.MustCompile(`[a-zA-Z]`)},
	"*": {Pattern: regexp.MustCompile(`[a-zA-Z0-9]`)},
}
