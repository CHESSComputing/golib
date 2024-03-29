package lexicon

import (
	"regexp"
	"testing"
)

func initLexicon() {
	LexiconPatterns = make(map[string]LexiconPattern)
	lex := Lexicon{
		Name: "did", Patterns: []string{".*unittest.*"}, Length: 100,
	}
	var patterns []*regexp.Regexp
	for _, pat := range lex.Patterns {
		patterns = append(patterns, regexp.MustCompile(pat))
	}
	lexPattern := LexiconPattern{
		Lexicon:  lex,
		Patterns: patterns,
	}
	LexiconPatterns["did"] = lexPattern
}

// TestPatterns
func TestPatterns(t *testing.T) {
	if LexiconPatterns == nil {
		initLexicon()
	}
	err := CheckPattern("did", "xyz")
	if err == nil {
		t.Error("unable to correct match did pattern")
	}
	err = CheckPattern("did", "unittest")
	if err != nil {
		t.Error(err)
	}
}
