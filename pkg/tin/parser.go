package tin

import (
	"fmt"
	"strconv"
)

func parseProgramFromTokens(tokens []token) (program Program) {
	var ipStack []int
	var ip int = 0
	var funStack map[string]int = make(map[string]int)

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
			keyword := tokens[0].value
			switch keyword {
			case "if":
				program = append(program, Instruction{
					Kind:  InstKindTestCondition,
					token: tokens[0],
				})
				tokens = tokens[1:]
				ipStack = append(ipStack, ip)
				ip++
			case "else":
				if len(ipStack) == 0 {
					panic("cannot parse the else of a non existing if")
				}
				if_addr := ipStack[len(ipStack)-1]
				ipStack = ipStack[:len(ipStack)-1]
				ipStack = append(ipStack, ip)
				program = append(program, Instruction{
					Kind:  InstKindElse,
					token: tokens[0],
				})
				tokens = tokens[1:]
				program[if_addr].JmpAddress = ip + 1
				ip++
			case "while":
				program = append(program, Instruction{
					Kind:  InstKindWhile,
					token: tokens[0],
				})
				tokens = tokens[1:]
				ipStack = append(ipStack, ip)
				ip++
			case "do":
				if len(ipStack) == 0 {
					panic("'do' used without a preceding 'while'")
				}
				program = append(program, Instruction{
					Kind:  InstKindTestCondition,
					token: tokens[0],
				})
				tokens = tokens[1:]
				ipStack = append(ipStack, ip)
				ip++
			case "end":
				if len(ipStack) == 0 {
					panic("'end' used without a preceding 'if', 'while', 'def'")
				}
				prec_addr := ipStack[len(ipStack)-1]
				ipStack = ipStack[:len(ipStack)-1]

				inst := program[prec_addr]
				switch inst.Kind {
				case InstKindTestCondition:
					endJmpAddr := ip + 1
					if len(ipStack) > 0 && program[ipStack[len(ipStack)-1]].Kind == InstKindWhile {
						endJmpAddr = ipStack[len(ipStack)-1]
						ipStack = ipStack[1:]
					} else {
						program[prec_addr].JmpAddress = ip + 1
					}
					program = append(program, Instruction{
						Kind:       InstKindEnd,
						token:      tokens[0],
						JmpAddress: endJmpAddr,
					})
					tokens = tokens[1:]
				case InstKindElse:
					program[prec_addr].JmpAddress = ip + 1
					program = append(program, Instruction{
						Kind:       InstKindEnd,
						token:      tokens[0],
						JmpAddress: ip + 1,
					})
					tokens = tokens[1:]
				case InstKindFunSkip:
					program[prec_addr].JmpAddress = ip + 1
					program = append(program, Instruction{
						Kind:  InstKindFunRet,
						token: tokens[0],
					})
					tokens = tokens[1:]
				default:
					panic(fmt.Sprintf("unexpected keyword '%s' as the preceder of 'end'", inst.Kind))
				}
				ip++
			case "def":
				program = append(program, Instruction{
					Kind:  InstKindFunSkip,
					token: tokens[0],
				})
				program = append(program, Instruction{
					Kind:  InstKindFunDef,
					token: tokens[0],
				})
				tokens = tokens[1:]
				ipStack = append(ipStack, ip)
				ip++

				if len(tokens) == 0 {
					panic("'def' used without a name")
				}
				funName := tokens[0].value
				tokens = tokens[1:]
				funStack[funName] = ip
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
				// a function call
				if fun_addr, ok := funStack[tokens[0].value]; ok {
					program = append(program, Instruction{
						Kind:       InstKindFunCall,
						JmpAddress: fun_addr,
						token:      tokens[0],
					})
					tokens = tokens[1:]
					ip++
				} else {
					panic(fmt.Sprintf("%s: unknown word '%s'", tokens[0].location, tokens[0].value))
				}

			}
		default:
			panic("there is a problem with 'parseProgramFromTokens' because this should be unreachable")
		}
	}
	return program
}
