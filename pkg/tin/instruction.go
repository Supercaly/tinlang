package tin

type InstKind int

const (
	InstKindPushInt InstKind = iota
	InstKindIntrinsic
)

type Instruction struct {
	Kind           InstKind
	token          token
	ValueInt       int
	ValueIntrinsic Intrinsic
}

type Intrinsic int

const (
	IntrinsicPlus Intrinsic = iota
	IntrinsicMinus
	IntrinsicTimes
	IntrinsicDivMod
)

type Program []Instruction
