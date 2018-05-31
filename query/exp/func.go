package exp

import (
	"github.com/tsealex/dbutil/query"
	"bytes"
)

type FuncExp struct {
	Name string
	Args []Exp
}

func Func(name string, exps ... Exp) *FuncExp {
	// TODO: type check
	return &FuncExp{Name: name, Args: exps}
}

func (f FuncExp) toSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteString(f.Name)
	buf.WriteByte('(')
	for i, arg := range f.Args {
		if err = arg.toSQL(ctx, buf); err != nil {
			return
		}
		if i > 0 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(')')
	return
}

type CastExp struct {
	Type   string
	SubExp Exp
}

func Cast(typeName string, exp Exp) *CastExp {
	// TODO: type check
	return &CastExp{Type: typeName, SubExp: exp}
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
	Name   string
	SubExp Exp
}

func As(name string, exp Exp) *AliasExp {
	// TODO: type check
	return &AliasExp{Name: name, SubExp: exp}
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
