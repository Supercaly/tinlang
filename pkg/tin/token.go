package tin

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type tokenKind int

const (
	tokenKindWord tokenKind = iota
	tokenKindKeyword
	tokenKindIntLit
	tokenKindStringLit
)

type token struct {
	kind     tokenKind
	value    string
	location fileLocation
}

const (
	spaceRegexStr     string = `^\s`
	commentRegexStr   string = `^#.*`
	intLitRegexStr    string = `^\d+`
	stringLitRegexStr string = `^"([^"\\]|\\.)*"`
	keywordRegexStr   string = `^(if|else|end|while|do|def|include|memory)`
)

func tokenizeSource(source string, fileName string) (out []token) {
	location := fileLocation{fileName: fileName}

	spaceRegex, err := regexp.Compile(spaceRegexStr)
	if err != nil {
		panic(err)
	}
	commentRegex, err := regexp.Compile(commentRegexStr)
	if err != nil {
		panic(err)
	}
	intLitRegex, err := regexp.Compile(intLitRegexStr)
	if err != nil {
		panic(err)
	}
	stringLitRegex, err := regexp.Compile(stringLitRegexStr)
	if err != nil {
		panic(err)
	}
	keywordRegex, err := regexp.Compile(keywordRegexStr)
	if err != nil {
		panic(err)
	}

	for len(source) > 0 {
		if spaceRegex.MatchString(source) {
			switch source[0] {
			case ' ':
				location.col++
			case '\n':
				location.col = 0
				location.row++
			case '\r':
				location.col = 0
			default:
				// TODO: manage all whitespace characters
				panic(fmt.Sprintf("%s: unsupported whitespace character %v", location, source[0]))
			}
			source = source[1:]
		} else if commentRegex.MatchString(source) {
			idxs := commentRegex.FindIndex([]byte(source))
			if idxs == nil {
				panic(fmt.Sprintf("%s: cannot find the end of a comment", location))
			}
			source = source[idxs[1]:]
			// TODO: Comments don't increment the location
			// This is not a big deal since at the end of a comment ther's always a new line
		} else if intLitRegex.MatchString(source) {
			idxs := intLitRegex.FindIndex([]byte(source))
			if idxs == nil {
				panic(fmt.Sprintf("%s: cannot find the end of an integer literal", location))
			}
			intStr := source[:idxs[1]]
			source = source[idxs[1]:]
			out = append(out, token{
				kind:     tokenKindIntLit,
				value:    intStr,
				location: location,
			})
			location.col += idxs[1]
		} else if stringLitRegex.MatchString(source) {
			idxs := stringLitRegex.FindIndex([]byte(source))
			if idxs == nil {
				panic(fmt.Sprintf("%s: cannot find the end of an string literal", location))
			}
			// TODO: Unsupported multi-line strings
			str := source[:idxs[1]]
			source = source[idxs[1]:]
			out = append(out, token{
				kind:     tokenKindStringLit,
				value:    str[1 : len(str)-1],
				location: location,
			})
			location.col += idxs[1]
		} else if keywordRegex.MatchString(source) {
			idxs := keywordRegex.FindIndex([]byte(source))
			if idxs == nil {
				panic(fmt.Sprintf("%s: cannot find the end of a keyword", location))
			}
			out = append(out, token{
				kind:     tokenKindKeyword,
				value:    source[:idxs[1]],
				location: location,
			})
			source = source[idxs[1]:]
			location.col += idxs[1]
		} else {
			idx := strings.IndexFunc(source, unicode.IsSpace)
			var word string
			if idx == -1 {
				word = source
				source = ""
			} else {
				word = source[:idx]
				source = source[idx:]
			}
			out = append(out, token{
				kind:     tokenKindWord,
				value:    word,
				location: location,
			})
			location.col += len(word)
		}
	}
	return out
}

func (t tokenKind) String() string {
	return [...]string{
		"tokenKindWord",
		"tokenKindKeyword",
		"tokenKindIntLit",
		"tokenKindStringLit",
	}[t]
}

func (t token) String() string {
	return fmt.Sprintf("(%s, %s, %s)", t.location, t.kind, t.value)
}
