package tin

import "fmt"

var intrinsics = map[string]Intrinsic{
	"+": IntrinsicPlus,
}

func parseProgramFromTokens(tokens []token) (program Program) {
	for len(tokens) > 0 {
		switch tokens[0].kind {
		case tokenKindIntLit:
			program = append(program, Op{
				Type:     OpTypePushInt,
				ValueInt: tokens[0].asIntLit,
				Token:    tokens[0],
			})
			tokens = tokens[1:]
		case tokenKindKeyword:
			panic("parse keyword not implemented")
		case tokenKindWord:
			intrinsic := intrinsics[tokens[0].asWord]
			if intrinsic != -1 {
				program = append(program, Op{
					Type:           OpTypeIntrinsic,
					ValueIntrinsic: intrinsic,
					Token:          tokens[0],
				})
				tokens = tokens[1:]
			} else {
				panic(fmt.Sprintf("unknown intrinsic %s", tokens[0].asWord))
			}
		}
	}
	return program
}
