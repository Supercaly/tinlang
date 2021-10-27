package tin

type OpType int

const (
	OpTypePushInt OpType = iota
	OpTypeIntrinsic
)

type Op struct {
	Type           OpType
	Token          token
	ValueInt       int
	ValueIntrinsic Intrinsic
}

type Intrinsic int

const (
	IntrinsicPlus Intrinsic = iota
)

type Program []Op
