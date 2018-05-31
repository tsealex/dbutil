package exp

import (
	"github.com/tsealex/dbutil/query"
	"bytes"
)

type FuncExp struct {
	BaseExp
	Name string
	Args []Exp
}

func Func(name string, exps ... Exp) *FuncExp {
	// TODO: type check
	res := &FuncExp{Name: name, Args: exps}
	res.Exp = res
	return res
}

func (f FuncExp) toSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteString(f.Name)
	buf.WriteByte('(')
	for i, arg := range f.Args {
		if i > 0 {
			buf.WriteByte(',')
		}
		if err = arg.toSQL(ctx, buf); err != nil {
			return
		}
	}
	buf.WriteByte(')')
	return
}

type CastExp struct {
	BaseExp
	Type   string
	SubExp Exp
}

func Cast(typeName string, exp Exp) *CastExp {
	// TODO: type check
	res := &CastExp{Type: typeName, SubExp: exp}
	res.Exp = res
	return res
}

func (c CastExp) toSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteByte('(')
	if err = c.SubExp.toSQL(ctx, buf); err != nil {
		return
	}
	buf.WriteString(")::")
	buf.WriteString(c.Type)
	return
}

type AliasExp struct {
	BaseExp
	Name   string
	SubExp Exp
}

func As(name string, exp Exp) *AliasExp {
	// TODO: type check
	res := &AliasExp{Name: name, SubExp: exp}
	res.Exp = res
	return res
}

func (a AliasExp) toSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteByte('(')
	if err = a.SubExp.toSQL(ctx, buf); err != nil {
		return
	}
	buf.WriteString(") AS ")
	buf.WriteString(a.Name)
	return
}
