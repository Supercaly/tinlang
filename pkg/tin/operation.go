package tin

type OpType int

const (
	OpTypePushInt OpType = iota
	OpTypeIntrinsic
	OpTypeIf
)

type Op struct {
	Type           OpType
	Token          Token
	ValueInt       int
	ValueIntrinsic Intrinsic
}

type Intrinsic int

const (
	IntrinsicPlus Intrinsic = iota
)

type Program []Op
