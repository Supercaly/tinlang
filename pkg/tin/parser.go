package tin

import "fmt"

var intrinsics = map[string]Intrinsic{
	"+": IntrinsicPlus,
}

func parseProgramFromTokens(tokens []Token) (program Program) {
	for len(tokens) > 0 {
		switch tokens[0].Type {
		case TokenTypeIntLit:
			program = append(program, Op{
				Type:     OpTypePushInt,
				ValueInt: tokens[0].AsIntLit,
				Token:    tokens[0],
			})
			tokens = tokens[1:]
		case TokenTypeKeyword:
			panic("parse keyword not implemented")
		case TokenTypeWord:
			intrinsic := intrinsics[tokens[0].AsWord]
			if intrinsic != -1 {
				program = append(program, Op{
					Type:           OpTypeIntrinsic,
					ValueIntrinsic: intrinsic,
					Token:          tokens[0],
				})
				tokens = tokens[1:]
			} else {
				panic(fmt.Sprintf("unknown intrinsic %s", tokens[0].AsWord))
			}
		}
	}
	return program
}
