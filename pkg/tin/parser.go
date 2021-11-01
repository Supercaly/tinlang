package tin

import (
	"fmt"
	"strconv"
)

func parseProgramFromTokens(tokens []token) (program Program) {
	var ipStack []int
	var ip int = 0
	var defStack map[string]int = make(map[string]int)

	for len(tokens) > 0 {
		switch tokens[0].kind {
		case tokenKindIntLit:
			intVal, err := strconv.ParseInt(tokens[0].value, 10, 64)
			if err != nil {
				panic(err)
			}
			program = append(program, Instruction{
				Kind:     InstKindPushInt,
				ValueInt: int(intVal),
				token:    tokens[0],
			})
			tokens = tokens[1:]
			ip++
		case tokenKindStringLit:
			program = append(program, Instruction{
				Kind:        InstKindPushString,
				ValueString: tokens[0].value,
				token:       tokens[0],
			})
			tokens = tokens[1:]
			ip++
		case tokenKindKeyword:
			keyword, exist := keywordMap[tokens[0].value]
			if !exist {
				panic(fmt.Sprintf("%s: unknown keyword '%s'", tokens[0].location, tokens[0].value))
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
				ip++
			case KeywordKindElse:
				if len(ipStack) == 0 {
					panic("cannot parse the else of a non existing if")
				}
				if_addr := ipStack[len(ipStack)-1]
				ipStack = ipStack[:len(ipStack)-1]
				ipStack = append(ipStack, ip)
				program[if_addr].ValueKeyword.HasJmp = true
				program[if_addr].ValueKeyword.JmpAddress = ip + 1
				ip++
			case KeywordKindWhile:
				ipStack = append(ipStack, ip)
				ip++
			case KeywordKindDo:
				if len(ipStack) == 0 {
					panic("'do' used without a preceding 'while'")
				}
				while_addr := ipStack[len(ipStack)-1]
				ipStack = ipStack[:len(ipStack)-1]

				inst := program[while_addr]
				if inst.Kind != InstKeyword || inst.ValueKeyword.Kind != KeywordKindWhile {
					panic("unexpected preceder parsing 'do'")
				}
				program[ip].ValueKeyword.HasJmp = true
				program[ip].ValueKeyword.JmpAddress = while_addr
				ipStack = append(ipStack, ip)
				ip++
			case KeywordKindEnd:
				if len(ipStack) == 0 {
					panic("'end' used without a preceding 'if' or 'while'")
				}
				inst_addr := ipStack[len(ipStack)-1]
				ipStack = ipStack[:len(ipStack)-1]

				inst := program[inst_addr]
				if inst.Kind != InstKeyword {
					panic("unexpected preceder parsing 'end'")
				}
				switch program[inst_addr].ValueKeyword.Kind {
				case KeywordKindIf:
					program[inst_addr].ValueKeyword.HasJmp = true
					program[inst_addr].ValueKeyword.JmpAddress = ip + 1
				case KeywordKindElse:
					// do nothing
				case KeywordKindDo:
					program[ip].ValueKeyword.HasJmp = true
					program[ip].ValueKeyword.JmpAddress = program[inst_addr].ValueKeyword.JmpAddress
					program[inst_addr].ValueKeyword.HasJmp = true
					program[inst_addr].ValueKeyword.JmpAddress = ip + 1
				case KeywordKindDef:
					program[inst_addr].ValueKeyword.HasJmp = true
					program[inst_addr].ValueKeyword.JmpAddress = ip + 1
					program[ip].ValueKeyword.IsRet = true
				default:
					panic(fmt.Sprintf("unexpected keyword '%s' as the preceder of 'end'", inst.ValueKeyword.Kind))
				}
				ip++
			case KeywordKindDef:
				ipStack = append(ipStack, ip)
				ip++
			default:
				panic(fmt.Sprintf("unknown keyword '%s'", keyword))
			}
		case tokenKindWord:
			intrinsic, exist := intrinsicMap[tokens[0].value]
			// word is an intrinsic
			if exist {
				program = append(program, Instruction{
					Kind:           InstKindIntrinsic,
					ValueIntrinsic: intrinsic,
					token:          tokens[0],
				})
				tokens = tokens[1:]
				ip++
			} else {
				if len(ipStack) > 0 {
					// check if word is the name of a function
					def_addr := ipStack[len(ipStack)-1]
					inst := program[def_addr]
					if inst.Kind != InstKeyword || program[ipStack[len(ipStack)-1]].ValueKeyword.Kind != KeywordKindDef {
						panic("boh")
					}

					program = append(program, Instruction{
						Kind:         InstKeyword,
						ValueKeyword: Keyword{Kind: KeywordKindDefName},
						token:        tokens[0],
					})
					defStack[tokens[0].value] = ip
					tokens = tokens[1:]
					ip++
				} else {
					// a function call
					if fun_addr, ok := defStack[tokens[0].value]; ok {
						program = append(program, Instruction{
							Kind:         InstKindFunCall,
							ValueFunCall: fun_addr,
							token:        tokens[0],
						})
						tokens = tokens[1:]
						ip++
					} else {
						panic(fmt.Sprintf("%s: unknown intrinsic '%s'", tokens[0].location, tokens[0].value))
					}
				}
			}
		default:
			panic("there is a problem with 'parseProgramFromTokens' because this should be unreachable")
		}
	}
	return program
}
