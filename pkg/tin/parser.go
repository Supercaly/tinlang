package tin

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

const (
	maxIncludeLevel int = 100
)

type parser struct {
	ip           int
	ipStack      []int
	funStack     map[string]int
	includeLevel int
}

func (p *parser) parseProgramFromTokens(tokens []token) (program Program) {
	if p.funStack == nil {
		p.funStack = make(map[string]int)
	}

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
			p.ip++
		case tokenKindStringLit:
			program = append(program, Instruction{
				Kind:        InstKindPushString,
				ValueString: tokens[0].value,
				token:       tokens[0],
			})
			tokens = tokens[1:]
			p.ip++
		case tokenKindKeyword:
			keyword := tokens[0].value
			switch keyword {
			case "if":
				program = append(program, Instruction{
					Kind:  InstKindTestCondition,
					token: tokens[0],
				})
				tokens = tokens[1:]
				p.ipStack = append(p.ipStack, p.ip)
				p.ip++
			case "else":
				if len(p.ipStack) == 0 {
					panic("cannot parse the else of a non existing if")
				}
				if_addr := p.ipStack[len(p.ipStack)-1]
				p.ipStack = p.ipStack[:len(p.ipStack)-1]
				p.ipStack = append(p.ipStack, p.ip)
				program = append(program, Instruction{
					Kind:  InstKindElse,
					token: tokens[0],
				})
				tokens = tokens[1:]
				program[if_addr].JmpAddress = p.ip + 1
				p.ip++
			case "while":
				program = append(program, Instruction{
					Kind:  InstKindWhile,
					token: tokens[0],
				})
				tokens = tokens[1:]
				p.ipStack = append(p.ipStack, p.ip)
				p.ip++
			case "do":
				if len(p.ipStack) == 0 {
					panic("'do' used without a preceding 'while'")
				}
				program = append(program, Instruction{
					Kind:  InstKindTestCondition,
					token: tokens[0],
				})
				tokens = tokens[1:]
				p.ipStack = append(p.ipStack, p.ip)
				p.ip++
			case "end":
				if len(p.ipStack) == 0 {
					panic("'end' used without a preceding 'if', 'while', 'def'")
				}
				prec_addr := p.ipStack[len(p.ipStack)-1]
				p.ipStack = p.ipStack[:len(p.ipStack)-1]

				inst := program[prec_addr]
				switch inst.Kind {
				case InstKindTestCondition:
					endJmpAddr := p.ip + 1
					if len(p.ipStack) > 0 && program[p.ipStack[len(p.ipStack)-1]].Kind == InstKindWhile {
						endJmpAddr = p.ipStack[len(p.ipStack)-1]
						p.ipStack = p.ipStack[1:]
					} else {
						program[prec_addr].JmpAddress = p.ip + 1
					}
					program = append(program, Instruction{
						Kind:       InstKindEnd,
						token:      tokens[0],
						JmpAddress: endJmpAddr,
					})
					tokens = tokens[1:]
				case InstKindElse:
					program[prec_addr].JmpAddress = p.ip + 1
					program = append(program, Instruction{
						Kind:       InstKindEnd,
						token:      tokens[0],
						JmpAddress: p.ip + 1,
					})
					tokens = tokens[1:]
				case InstKindFunSkip:
					program[prec_addr].JmpAddress = p.ip + 1
					program = append(program, Instruction{
						Kind:  InstKindFunRet,
						token: tokens[0],
					})
					tokens = tokens[1:]
				default:
					panic(fmt.Sprintf("unexpected keyword '%s' as the preceder of 'end'", inst.Kind))
				}
				p.ip++
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
				p.ipStack = append(p.ipStack, p.ip)
				p.ip++

				if len(tokens) == 0 {
					panic("'def' used without a name")
				}
				funName := tokens[0].value
				tokens = tokens[1:]
				p.funStack[funName] = p.ip
				p.ip++
			case "include":
				tokens = tokens[1:]
				if len(tokens) == 0 || tokens[0].kind != tokenKindStringLit {
					panic("expected location after include")
				}
				includePath := tokens[0].value
				tokens = tokens[1:]

				if p.includeLevel+1 > maxIncludeLevel {
					panic("max include lever reached")
				}
				p.includeLevel++
				source, err := ioutil.ReadFile(includePath)
				if err != nil {
					panic(err)
				}
				tokens := tokenizeSource(string(source), includePath)
				program = append(program, p.parseProgramFromTokens(tokens)...)
				p.includeLevel--
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
				p.ip++
			} else {
				// a function call
				if fun_addr, ok := p.funStack[tokens[0].value]; ok {
					program = append(program, Instruction{
						Kind:       InstKindFunCall,
						JmpAddress: fun_addr,
						token:      tokens[0],
					})
					tokens = tokens[1:]
					p.ip++
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
