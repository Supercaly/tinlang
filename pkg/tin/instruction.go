package tin

import "fmt"

type InstKind int

const (
	InstKindPushInt InstKind = iota
	InstKindPushString
	InstKeyword
	InstKindIntrinsic
)

type Instruction struct {
	Kind           InstKind
	token          token
	ValueInt       int
	ValueString    string
	ValueKeyword   Keyword
	ValueIntrinsic Intrinsic
}

type KeywordKind int

const (
	KeywordKindIf KeywordKind = iota
	KeywordKindElse
	KeywordKindEnd
	KeywordKindWhile
	KeywordKindDo
)

var keywordMap = map[string]KeywordKind{
	"if":    KeywordKindIf,
	"else":  KeywordKindElse,
	"end":   KeywordKindEnd,
	"while": KeywordKindWhile,
	"do":    KeywordKindDo,
}

type Keyword struct {
	Kind       KeywordKind
	JmpAddress int
	HasJmp     bool
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
)

var intrinsicMap = map[string]Intrinsic{
	"+":      IntrinsicPlus,
	"-":      IntrinsicMinus,
	"*":      IntrinsicTimes,
	"divmod": IntrinsicDivMod,
	">":      IntrinsicGreather,
	"<":      IntrinsicLess,
	"!=":     IntrinsicNotEqual,
	"dup":    IntrinsicDup,
	"print":  IntrinsicPrint,
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
	case InstKeyword:
		out += i.ValueKeyword.Kind.String()
	}
	return out
}

func (ik InstKind) String() string {
	return [...]string{
		"InstKindPushInt",
		"InstKindPushString",
		"InstKeyword",
		"InstKindIntrinsic",
	}[ik]
}

func (k KeywordKind) String() string {
	return [...]string{
		"if",
		"else",
		"end",
		"while",
		"do",
	}[k]
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
	}[i]
}
