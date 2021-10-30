package tin

import (
	"fmt"
)

func parseProgramFromTokens(tokens []token) (program Program) {
	var ipStack []int
	var ip int = 0

	for len(tokens) > 0 {
		switch tokens[0].kind {
		case tokenKindIntLit:
			program = append(program, Instruction{
				Kind:     InstKindPushInt,
				ValueInt: tokens[0].asIntLit,
				token:    tokens[0],
			})
			tokens = tokens[1:]
			ip++
		case tokenKindKeyword:
			keyword, exist := keywordMap[tokens[0].asKeyword]
			if !exist {
				panic(fmt.Sprintf("%s: unknown keyword '%s'", tokens[0].location, tokens[0].asKeyword))
			}
			program = append(program, Instruction{
				Kind:         InstKeyword,
				ValueKeyword: Keyword{Kind: keyword},
				token:        tokens[0],
			})
			tokens = tokens[1:]

			switch keyword {
			case KeywordKindIf:
				ipStack = append(ipStack, ip)
			case KeywordKindElse:
				if len(ipStack) == 0 {
					panic("cannot parse the else of a non existing if")
				}
				if_addr := ipStack[len(ipStack)-1]
				ipStack = ipStack[:len(ipStack)-1]
				ipStack = append(ipStack, ip)
				program[if_addr].ValueKeyword.JmpAddress = ip + 1
			case KeywordKindEnd:
				if len(ipStack) == 0 {
					panic("cannot parse the end of a non existing block")
				}
				inst_addr := ipStack[len(ipStack)-1]
				ipStack = ipStack[:len(ipStack)-1]

				inst := program[inst_addr]
				if inst.Kind != InstKeyword {
					panic("unexpected instruction parsing end block")
				}
				switch program[inst_addr].ValueKeyword.Kind {
				case KeywordKindIf:
					program[inst_addr].ValueKeyword.JmpAddress = ip + 1
				case KeywordKindElse:
					// do nothing
				default:
					panic("unexpected keyword parsing end block")
				}
			}

			ip++
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
			ip++
		default:
			panic("there is a problem with 'parseProgramFromTokens' because this should be unreachable")
		}
	}
	return program
}
