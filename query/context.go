package query

type SQLContext struct {
	TagMap    map[string]int
	ReqSchema bool
	// For insertion
	WriteStatus uint8

	index int
}

func NewSQLContext() *SQLContext {
	return &SQLContext{
		TagMap: map[string]int{},
	}
}

func (ctx *SQLContext) NextIndex() int {
	ctx.index += 1
	return ctx.index
}

func (ctx *SQLContext) GetTagIndex(tag string) int {
	if i, ok := ctx.TagMap[tag]; ok {
		return i
	} else {
		i := ctx.NextIndex()
		ctx.TagMap[tag] = i
		return i
	}
}

const (
	Regular    uint8 = iota
	ValueOnly  uint8 = iota
	ColumnOnly uint8 = iota
)
