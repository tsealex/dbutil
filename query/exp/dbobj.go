package exp

import (
	"github.com/tsealex/dbutil/query"
	"bytes"
)

type ColumnExp struct {
	BaseExp
	Name     string
	Relation *RelationExp
	Quoted   bool
}

func Column(name string) *ColumnExp {
	res := &ColumnExp{Name: name}
	res.Exp = res
	return res
}

func (c *ColumnExp) Assign(exp Exp) *AssignExp {
	res := &AssignExp{Col: c, RightExp: exp}
	res.Exp = res
	return res
}

func (c *ColumnExp) SetRelation(r *RelationExp) *ColumnExp {
	c.Relation = r
	return c
}

func (c *ColumnExp) Quote() *ColumnExp {
	c.Quoted = true
	return c
}

func (c *ColumnExp) Unquote() *ColumnExp {
	c.Quoted = false
	return c
}

func (c ColumnExp) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	name := c.Name
	// TODO: Check name conflict (i.e. if two tables have the same column, make
	// TODO: sure the Table field is specified, else return an error.
	if c.Relation != nil {
		if err = c.Relation.ToSQL(ctx, buf); err != nil {
			return
		}
		buf.WriteByte('.')
	}
	if c.Quoted {
		name = `"` + name + `""`
	}
	buf.WriteString(name)
	return
}

type RelationExp struct {
	BaseExp
	Name   string
	Schema *SchemaExp
	Quoted bool
}

func Relation(name string) *RelationExp {
	res := &RelationExp{Name: name}
	res.Exp = res
	return res
}

func (r *RelationExp) SetSchema(s *SchemaExp) *RelationExp {
	r.Schema = s
	return r
}

func (r *RelationExp) Quote() *RelationExp {
	r.Quoted = true
	return r
}

func (r *RelationExp) Unquote() *RelationExp {
	r.Quoted = false
	return r
}

func (r RelationExp) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	name := r.Name
	// Consult the context to see whether schema name is required here. Note
	// that, in some portion of a query, schema names are not required.
	if ctx.ReqSchema && r.Schema != nil {
		if err = r.Schema.ToSQL(ctx, buf); err != nil {
			return
		}
		buf.WriteByte('.')
	}
	if r.Quoted {
		name = `"` + name + `""`
	}
	buf.WriteString(name)
	return
}

type SchemaExp struct {
	BaseExp
	Name   string
	Quoted bool
}

func Schema(name string) *SchemaExp {
	res := &SchemaExp{Name: name}
	res.Exp = res
	return res
}

func (s *SchemaExp) Quote() *SchemaExp {
	s.Quoted = true
	return s
}

func (s *SchemaExp) Unquote() *SchemaExp {
	s.Quoted = false
	return s
}

func (s SchemaExp) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	name := s.Name
	if s.Quoted {
		name = `"` + name + `""`
	}
	buf.WriteString(name)
	return
}

type AssignExp struct {
	BaseExp
	Col      Exp
	RightExp Exp
}

func Assign(c Exp, e Exp) *AssignExp {
	res := &AssignExp{Col: c, RightExp: e}
	res.Exp = res
	return res
}

func (a AssignExp) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	if ctx.WriteStatus != query.ValueOnly {
		if err = a.Col.ToSQL(ctx, buf); err != nil {
			return
		}
	}
	if ctx.WriteStatus != query.Regular {
		buf.WriteString("=")
	}
	if ctx.WriteStatus != query.ColumnOnly {
		if err = a.RightExp.ToSQL(ctx, buf); err != nil {
			return
		}
	}
	return
}
