package tin

type InstKind int

const (
	InstKindPushInt InstKind = iota
	IntrinsicKeyword
	InstKindIntrinsic
)

type Instruction struct {
	Kind           InstKind
	token          token
	ValueInt       int
	ValueKeyword   Keyword
	ValueIntrinsic Intrinsic
}

type Keyword int

const (
	KeywordIf Keyword = iota
	KeywordElse
	KeywordEnd
)

var keywordMap = map[string]Keyword{
	"if":   KeywordIf,
	"else": KeywordElse,
	"end":  KeywordEnd,
}

type Intrinsic int

const (
	IntrinsicPlus Intrinsic = iota
	IntrinsicMinus
	IntrinsicTimes
	IntrinsicDivMod

	IntrinsicGreather
	IntrinsicLess
)

var intrinsicMap = map[string]Intrinsic{
	"+":      IntrinsicPlus,
	"-":      IntrinsicMinus,
	"*":      IntrinsicTimes,
	"divmod": IntrinsicDivMod,
	">":      IntrinsicGreather,
	"<":      IntrinsicLess,
}

type Program []Instruction
