package tin

import "fmt"

type InstKind int

const (
	InstKindPushInt InstKind = iota
	InstKeyword
	InstKindIntrinsic
)

type Instruction struct {
	Kind           InstKind
	token          token
	ValueInt       int
	ValueKeyword   Keyword
	ValueIntrinsic Intrinsic
}

type KeywordKind int

const (
	KeywordKindIf KeywordKind = iota
	KeywordKindElse
	KeywordKindEnd
)

var keywordMap = map[string]KeywordKind{
	"if":   KeywordKindIf,
	"else": KeywordKindElse,
	"end":  KeywordKindEnd,
}

type Keyword struct {
	Kind       KeywordKind
	JmpAddress int
}

type Intrinsic int

const (
	IntrinsicPlus Intrinsic = iota
	IntrinsicMinus
	IntrinsicTimes
	IntrinsicDivMod

	IntrinsicGreather
	IntrinsicLess

	IntrinsicPrint
)

var intrinsicMap = map[string]Intrinsic{
	"+":      IntrinsicPlus,
	"-":      IntrinsicMinus,
	"*":      IntrinsicTimes,
	"divmod": IntrinsicDivMod,
	">":      IntrinsicGreather,
	"<":      IntrinsicLess,
	"print":  IntrinsicPrint,
}

type Program []Instruction

func (i Instruction) String() (out string) {
	// TODO: Better print the instruction
	switch i.Kind {
	case InstKindPushInt:
		out += fmt.Sprint(i.ValueInt)
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
		"InstKeyword",
		"InstKindIntrinsic",
	}[ik]
}

func (k KeywordKind) String() string {
	return [...]string{
		"if",
		"else",
		"end",
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
		"print",
	}[i]
}
