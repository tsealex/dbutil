package exp

import (
	"github.com/tsealex/dbutil/query"
	"bytes"
	"fmt"
	"strconv"
)

type Exp interface {
	ToSQL(*query.SQLContext, *bytes.Buffer) error
}

type BaseExp struct {
	Exp
}

func (b *BaseExp) Add(exp Exp) *BinaryExp {
	return Binary(b, "+", exp)
}

func (b *BaseExp) Sub(exp Exp) *BinaryExp {
	return Binary(b, "-", exp)
}

func (b *BaseExp) Mul(exp Exp) *BinaryExp {
	return Binary(b, "*", exp)
}

func (b *BaseExp) Div(exp Exp) *BinaryExp {
	return Binary(b, "/", exp)
}

func (b *BaseExp) Mod(exp Exp) *BinaryExp {
	return Binary(b, "%", exp)
}

func (b *BaseExp) Expo(exp Exp) *BinaryExp {
	return Binary(b, "%", exp)
}

func (b *BaseExp) Is(exp *LiteralExp) *BinaryExp {
	return Binary(b, " IS ", exp)
}

func (b *BaseExp) IsNot(exp *LiteralExp) *BinaryExp {
	return Binary(b, " IS NOT ", exp)
}

func (b *BaseExp) Gt(exp Exp) *BinaryExp {
	return Binary(b, ">", exp)
}

func (b *BaseExp) Lt(exp Exp) *BinaryExp {
	return Binary(b, "<", exp)
}

func (b *BaseExp) Gte(exp Exp) *BinaryExp {
	return Binary(b, ">=", exp)
}

func (b *BaseExp) Lte(exp Exp) *BinaryExp {
	return Binary(b, "<=", exp)
}

func (b *BaseExp) RightShift(exp Exp) *BinaryExp {
	return Binary(b, ">>", exp)
}

func (b *BaseExp) LeftShift(exp Exp) *BinaryExp {
	return Binary(b, "<<", exp)
}

func (b *BaseExp) Concat(exp Exp) *BinaryExp {
	return Binary(b, "||", exp)
}

func (b *BaseExp) Union(exp Exp) *BinaryExp {
	return Binary(b, "|", exp)
}

func (b *BaseExp) Overlap(exp Exp) *BinaryExp {
	return Binary(b, "&&", exp)
}

func (b *BaseExp) Intersect(exp Exp) *BinaryExp {
	return Binary(b, "&", exp)
}

func (b *BaseExp) Contain(exp Exp) *BinaryExp {
	return Binary(b, "@>", exp)
}

func (b *BaseExp) ContainedBy(exp Exp) *BinaryExp {
	return Binary(b, "<@", exp)
}

func (b *BaseExp) Eq(exp Exp) *BinaryExp {
	return Binary(b, "=", exp)
}

func (b *BaseExp) NotEq(exp Exp) *BinaryExp {
	return Binary(b, "<>", exp)
}

func (b *BaseExp) Match(exp Exp, caseSens bool) *BinaryExp {
	if caseSens {
		return Binary(b, "~", exp)
	} else {
		return Binary(b, "~*", exp)
	}
}

func (b *BaseExp) NotMatch(exp Exp, caseSens bool) *BinaryExp {
	if caseSens {
		return Binary(b, "!~", exp)
	} else {
		return Binary(b, "!~*", exp)
	}
}

func (b *BaseExp) Like(pattern *LiteralExp) *BinaryExp {
	return Binary(b, " ~~ ", pattern)
}

func (b *BaseExp) ILike(pattern *LiteralExp) *BinaryExp {
	return Binary(b, " ~~* ", pattern)
}

func (b *BaseExp) NotLike(pattern *LiteralExp) *BinaryExp {
	return Binary(b, " !~~ ", pattern)
}

func (b *BaseExp) NotILike(pattern *LiteralExp) *BinaryExp {
	return Binary(b, " !~~* ", pattern)
}

func (b *BaseExp) SimilarTo(pattern *LiteralExp) *BinaryExp {
	return Binary(b, " SIMILAR TO ", pattern)
}

func (b BaseExp) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	return fmt.Errorf("unimplemented error")
}

type LiteralExp struct {
	BaseExp
	Value interface{}
	isExp bool
}

func Literal(value interface{}) (res *LiteralExp) {
	// TODO: type check, return nil if not a primitive type
	res = &LiteralExp{Value: value}
	res.isExp = false
	res.Exp = res
	return res
}

func Expression(value string) (res *LiteralExp) {
	// TODO: type check, return nil if not a primitive type
	res = &LiteralExp{Value: value}
	res.isExp = true
	res.Exp = res
	return res
}

func (l LiteralExp) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	if str, ok := l.Value.(string); !l.isExp && ok {
		_, err = buf.WriteString(fmt.Sprintf(`'%s'`, str))
	} else {
		_, err = buf.WriteString(fmt.Sprint(l.Value))
	}
	return
}

var All = Expression("*")

type ArrayExp struct {
	LiteralExp
	Values []interface{}
}

func Array(values ... interface{}) *ArrayExp {
	// TODO: type check, return nil if no a primitive type nor Exp instance.
	res := &ArrayExp{Values: values}
	res.Exp = res
	return res
}

func (a ArrayExp) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteString("ARRAY")
	buf.WriteByte('[')
	for i, val := range a.Values {
		if i > 0 {
			buf.WriteByte(',')
		}
		if exp, ok := val.(Exp); ok {
			if err = exp.ToSQL(ctx, buf); err != nil {
				return
			}
		} else {
			if str, ok := val.(string); ok {
				_, err = buf.WriteString(fmt.Sprintf(`'%s'`, str))
			} else {
				_, err = buf.WriteString(fmt.Sprint(val))
			}
		}
	}
	buf.WriteByte(']')
	return
}

type UnbindExp struct {
	BaseExp
	Tag string // Optional
}

func Unbind() *UnbindExp {
	return &UnbindExp{}
}

func TaggedUnbind(tag string) (res *UnbindExp) {
	res = &UnbindExp{Tag: tag}
	res.Exp = res
	return res
}

func (u UnbindExp) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
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

type GroupExp struct {
	BaseExp
	SubExp Exp
}

func Group(exp Exp) *GroupExp {
	res := &GroupExp{SubExp: exp}
	res.Exp = res
	return res
}

func (g GroupExp) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteByte('(')
	if err = g.SubExp.ToSQL(ctx, buf); err != nil {
		return
	}
	buf.WriteByte(')')
	return
}

type BinaryExp struct {
	BaseExp
	LeftExp  Exp
	RightExp Exp
	Op       string
}

func Binary(left Exp, op string, right Exp) *BinaryExp {
	res := &BinaryExp{LeftExp: left, RightExp: right, Op: op}
	res.Exp = res
	return res
}

func (b BinaryExp) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteByte('(')
	if err = b.LeftExp.ToSQL(ctx, buf); err != nil {
		return
	}
	buf.WriteString(b.Op)
	if err = b.RightExp.ToSQL(ctx, buf); err != nil {
		return
	}
	buf.WriteByte(')')
	return
}

type UnaryExp struct {
	BaseExp
	SubExp Exp
	Op     string
	pre    bool
}

func LeftUnary(op string, exp Exp) *UnaryExp {
	res := &UnaryExp{SubExp: exp, Op: op, pre: true}
	res.Exp = res
	return res
}

func RightUnary(exp Exp, op string) *UnaryExp {
	res := &UnaryExp{SubExp: exp, Op: op, pre: false}
	res.Exp = res
	return res
}

func Minus(exp Exp) *UnaryExp {
	return LeftUnary("-", exp)
}

func Not(exp Exp) *UnaryExp {
	return LeftUnary(" NOT ", exp)
}

func (u UnaryExp) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	if u.pre {
		buf.WriteString(u.Op)
	}
	if err = u.SubExp.ToSQL(ctx, buf); err != nil {
		return
	}
	if !u.pre {
		buf.WriteString(u.Op)
	}
	return
}

type CondExp struct {
	BaseExp
	Exps []Exp
	Op   string
}

func Cond(op string, exps ... Exp) *CondExp {
	res := &CondExp{Exps: exps, Op: op}
	res.Exp = res
	return res
}

func (c *CondExp) And(exps ... Exp) *CondExp {
	var tmp = make([]Exp, len(exps) + 1)
	tmp = append(tmp, c)
	tmp = append(tmp, exps...)
	return And(tmp...)
}

func And(exps ... Exp) *CondExp {
	return Cond(" AND ", exps...)
}

func Or(exps ... Exp) *CondExp {
	return Cond(" OR ", exps...)
}

func (c CondExp) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteByte('(')
	for i, exp := range c.Exps {
		if i > 0 {
			buf.WriteString(c.Op)
		}
		if err = exp.ToSQL(ctx, buf); err != nil {
			return
		}
	}
	buf.WriteByte(')')
	return
}

// TODO: Add Subquery Expression
