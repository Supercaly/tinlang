package tin

import (
	"fmt"
)

var intrinsicMap = map[string]Intrinsic{
	"+":      IntrinsicPlus,
	"-":      IntrinsicMinus,
	"*":      IntrinsicTimes,
	"divmod": IntrinsicDivMod,
}

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
			panic(fmt.Sprintf("%s: parse keyword not implemented", tokens[0].location))
		case tokenKindWord:
			intrinsic, exist := intrinsicMap[tokens[0].asWord]
			if exist {
				program = append(program, Instruction{
					Kind:           InstKindIntrinsic,
					ValueIntrinsic: intrinsic,
					token:          tokens[0],
				})
				tokens = tokens[1:]
			} else {
				panic(fmt.Sprintf("%s: unknown intrinsic '%s'", tokens[0].location, tokens[0].asWord))
			}
		default:
			panic("parseProgramFromTokens: unreachable")
		}
	}
	return program
}
