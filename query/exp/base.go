package exp

import (
	"github.com/tsealex/dbutil/query"
	"bytes"
	"fmt"
	"strconv"
)

type Exp interface {
	toSQL(*query.SQLContext, *bytes.Buffer) error
}

type LiteralExp struct {
	Value interface{}
}

func Literal(value interface{}) *LiteralExp {
	// TODO: type check, return nil if no a primitive type
	return &LiteralExp{Value: value}
}

func (l LiteralExp) toSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	if str, ok := l.Value.(string); ok {
		_, err = buf.WriteString(fmt.Sprintf(`"%s"`, str))
	} else {
		_, err = buf.WriteString(fmt.Sprint(l.Value))
	}
	return
}

type UnbindExp struct {
	Tag string // Optional
}

func Unbind() *UnbindExp {
	return &UnbindExp{}
}

func TaggedUnbind(tag string) *UnbindExp {
	return &UnbindExp{Tag: tag}
}

func (u UnbindExp) toSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	// TODO: Handle other DBMS cases. This only supports Postgres right now.
	var i int
	if u.Tag != "" {
		i = ctx.GetTagIndex(u.Tag)
	} else {
		i = ctx.NextIndex()
	}
	buf.WriteByte('$')
	buf.WriteString(strconv.FormatInt(int64(i), 10))
	return
}


