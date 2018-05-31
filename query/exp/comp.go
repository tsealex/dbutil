package exp

import (
	"github.com/tsealex/dbutil/query"
	"bytes"
)

type GroupExp struct {
	SubExp Exp
}

func Group(exp Exp) *GroupExp {
	return &GroupExp{SubExp: exp}
}

func (g GroupExp) toSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteByte('(')
	if err = g.SubExp.toSQL(ctx, buf); err != nil {
		return
	}
	buf.WriteByte(')')
	return
}

type BinaryExp struct {
	LeftExp  Exp
	RightExp Exp
	Op       string
}

func Binary(left Exp, op string, right Exp) *BinaryExp {
	return &BinaryExp{LeftExp: left, RightExp: right, Op: op}
}

func (b BinaryExp) toSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	if err = b.LeftExp.toSQL(ctx, buf); err != nil {
		return
	}
	buf.WriteString(b.Op)
	if err = b.RightExp.toSQL(ctx, buf); err != nil {
		return
	}
	return
}

type UnaryExp struct {
	SubExp Exp
	Op     string
	pre    bool
}

func LeftUnary(op string, exp Exp) *UnaryExp {
	return &UnaryExp{SubExp: exp, Op: op, pre: true}
}

func RightUnary(exp Exp, op string) *UnaryExp {
	return &UnaryExp{SubExp: exp, Op: op, pre: false}
}

func (u UnaryExp) toSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	if u.pre {
		buf.WriteString(u.Op)
	}
	if err = u.SubExp.toSQL(ctx, buf); err != nil {
		return
	}
	if !u.pre {
		buf.WriteString(u.Op)
	}
	return
}
