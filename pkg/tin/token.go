package tin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type TokenType int

const (
	TokenTypeWord TokenType = iota
	TokenTypeKeyword
	TokenTypeIntLit
)

type Token struct {
	Type      TokenType
	AsWord    string
	AsKeyword string
	AsIntLit  int
}

func tokenizeSource(source string) (out []Token) {
	spaceRegex, err := regexp.Compile(`^\s`)
	if err != nil {
		panic(err)
	}
	commentRegex, err := regexp.Compile(`^#.*`)
	if err != nil {
		panic(err)
	}
	intLitRegex, err := regexp.Compile(`^\d+`)
	if err != nil {
		panic(err)
	}
	keywordRegex, err := regexp.Compile(`^(if|else)`)
	if err != nil {
		panic(err)
	}

	for len(source) > 0 {
		if spaceRegex.MatchString(source) {
			source = source[1:]
		} else if commentRegex.MatchString(source) {
			idxs := commentRegex.FindIndex([]byte(source))
			if idxs == nil {
				panic("cannot find the end of a comment")
			}
			source = source[idxs[1]:]
		} else if intLitRegex.MatchString(source) {
			idxs := intLitRegex.FindIndex([]byte(source))
			if idxs == nil {
				panic("cannot find the end of an integer literal")
			}
			intLit, err := strconv.ParseInt(source[:idxs[1]], 10, 64)
			if err != nil {
				panic(err)
			}
			source = source[idxs[1]:]
			out = append(out, Token{
				Type:     TokenTypeIntLit,
				AsIntLit: int(intLit),
			})
		} else if keywordRegex.MatchString(source) {
			idxs := keywordRegex.FindIndex([]byte(source))
			if idxs == nil {
				panic("cannot find the end of a keyword")
			}
			out = append(out, Token{
				Type:      TokenTypeKeyword,
				AsKeyword: source[:idxs[1]],
			})
			source = source[idxs[1]:]
		} else {
			idx := strings.IndexRune(source, ' ')
			var word string
			if idx == -1 {
				word = source
			} else {
				word = source[:idx]
				source = source[idx:]
			}
			out = append(out, Token{
				Type:   TokenTypeWord,
				AsWord: word,
			})
		}
	}
	return out
}

func (t TokenType) String() string {
	return [...]string{
		"TokenTypeWord",
		"TokenTypeKeyword",
		"TokenTypeIntLit",
	}[t]
}

func (t Token) String() (out string) {
	out += "("
	out += fmt.Sprintf("%s, ", t.Type)
	switch t.Type {
	case TokenTypeWord:
		out += t.AsWord
	case TokenTypeKeyword:
		out += t.AsKeyword
	case TokenTypeIntLit:
		out += fmt.Sprint(t.AsIntLit)
	}
	out += ")"
	return out
}
