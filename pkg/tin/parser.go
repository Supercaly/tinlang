package tin

import (
	"fmt"
)

func parseProgramFromTokens(tokens []token) (program Program) {
	for len(tokens) > 0 {
		switch tokens[0].kind {
		case tokenKindIntLit:
			program = append(program, Instruction{
				Kind:     InstKindPushInt,
				ValueInt: tokens[0].asIntLit,
				token:    tokens[0],
			})
			tokens = tokens[1:]
		case tokenKindKeyword:
			keyword, exist := keywordMap[tokens[0].asKeyword]
			if !exist {
				panic(fmt.Sprintf("%s: unknown keyword '%s'", tokens[0].location, tokens[0].asKeyword))
			}
			program = append(program, Instruction{
				Kind:         IntrinsicKeyword,
				ValueKeyword: keyword,
				token:        tokens[0],
			})
			tokens = tokens[:1]
		case tokenKindWord:
			intrinsic, exist := intrinsicMap[tokens[0].asWord]
			if !exist {
				panic(fmt.Sprintf("%s: unknown intrinsic '%s'", tokens[0].location, tokens[0].asWord))
			}
			program = append(program, Instruction{
				Kind:           InstKindIntrinsic,
				ValueIntrinsic: intrinsic,
				token:          tokens[0],
			})
			tokens = tokens[1:]
		default:
			panic("parseProgramFromTokens: unreachable")
		}
	}
	return program
}
