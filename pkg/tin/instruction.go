package tin

import "fmt"

type InstKind int

const (
	InstKindPushInt InstKind = iota
	InstKindPushString
	InstKindIntrinsic

	InstKindTestCondition
	InstKindElse
	InstKindWhile
	InstKindEnd

	InstKindFunSkip
	InstKindFunDef
	InstKindFunRet
	InstKindFunCall

	InstKindMemPush
)

type Instruction struct {
	Kind           InstKind
	token          token
	ValueInt       int
	ValueString    string
	ValueIntrinsic Intrinsic
	ValueMemory    int
	JmpAddress     int
}

type Intrinsic int

const (
	IntrinsicPlus Intrinsic = iota
	IntrinsicMinus
	IntrinsicTimes
	IntrinsicDivMod

	IntrinsicGreather
	IntrinsicLess
	IntrinsicNotEqual

	IntrinsicDup

	IntrinsicPrint

	IntrinsicSyscall0
	IntrinsicSyscall1
	IntrinsicSyscall2
	IntrinsicSyscall3
	IntrinsicSyscall4
	IntrinsicSyscall5
	IntrinsicSyscall6

	IntrinsicLoad8
	IntrinsicStore8
	IntrinsicLoad32
	IntrinsicStore32
	IntrinsicLoad64
	IntrinsicStore64
)

var intrinsicMap = map[string]Intrinsic{
	"+":        IntrinsicPlus,
	"-":        IntrinsicMinus,
	"*":        IntrinsicTimes,
	"divmod":   IntrinsicDivMod,
	">":        IntrinsicGreather,
	"<":        IntrinsicLess,
	"!=":       IntrinsicNotEqual,
	"dup":      IntrinsicDup,
	"print":    IntrinsicPrint,
	"syscall0": IntrinsicSyscall0,
	"syscall1": IntrinsicSyscall1,
	"syscall2": IntrinsicSyscall2,
	"syscall3": IntrinsicSyscall3,
	"syscall4": IntrinsicSyscall4,
	"syscall5": IntrinsicSyscall5,
	"syscall6": IntrinsicSyscall6,
	"@8":       IntrinsicLoad8,
	"!8":       IntrinsicStore8,
	"@32":      IntrinsicLoad32,
	"!32":      IntrinsicStore32,
	"@64":      IntrinsicLoad64,
	"!64":      IntrinsicStore64,
}

type Program []Instruction

func (i Instruction) String() (out string) {
	// TODO: Better print the instruction
	switch i.Kind {
	case InstKindPushInt:
		out += fmt.Sprint(i.ValueInt)
	case InstKindPushString:
		out += i.ValueString
	case InstKindIntrinsic:
		out += i.ValueIntrinsic.String()
	case InstKindTestCondition:
		out += fmt.Sprintf("(test %d)", i.JmpAddress)
	case InstKindElse:
		out += fmt.Sprintf("(else %d)", i.JmpAddress)
	case InstKindEnd:
		out += fmt.Sprintf("(end %d)", i.JmpAddress)
	case InstKindWhile:
		out += "while"
	case InstKindFunSkip:
		out += fmt.Sprintf("(fskip %d)", i.JmpAddress)
	case InstKindFunDef:
		out += "fdef"
	case InstKindFunRet:
		out += "fret"
	case InstKindFunCall:
		out += fmt.Sprintf("(fcall %d)", i.JmpAddress)
	}
	return out
}

func (ik InstKind) String() string {
	return [...]string{
		"InstKindPushInt",
		"InstKindPushString",
		"InstKindIntrinsic",
		"InstKindTestCondition",
		"InstKindElse",
		"InstKindWhile",
		"InstKindEnd",
		"InstKindFunSkip",
		"InstKindFunDef",
		"InstKindFunRet",
		"InstKindFunCall",
		"InstKindMemPush",
	}[ik]
}

func (i Intrinsic) String() string {
	return [...]string{
		"+",
		"-",
		"*",
		"divmod",
		">",
		"<",
		"!=",
		"dup",
		"print",
		"syscall0",
		"syscall1",
		"syscall2",
		"syscall3",
		"syscall4",
		"syscall5",
		"syscall6",
		"IntrinsicLoad8",
		"IntrinsicStore8",
		"IntrinsicLoad32",
		"IntrinsicStore32",
		"IntrinsicLoad64",
		"IntrinsicStore64",
	}[i]
}
