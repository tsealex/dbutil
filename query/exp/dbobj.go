package exp

import (
	"github.com/tsealex/dbutil/query"
	"bytes"
)

type ColumnExp struct {
	Name     string
	Relation *RelationExp
	Quoted   bool
}

func Column(name string) *ColumnExp {
	return &ColumnExp{Name: name}
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

func (c ColumnExp) toSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	name := c.Name
	// TODO: Check name conflict (i.e. if two tables have the same column, make
	// TODO: sure the Table field is specified, else return an error.
	if c.Relation != nil {
		if err = c.Relation.toSQL(ctx, buf); err != nil {
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
	Name   string
	Schema *SchemaExp
	Quoted bool
}

func Relation(name string) *RelationExp {
	return &RelationExp{Name: name}
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

func (r RelationExp) toSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	name := r.Name
	// Consult the context to see whether schema name is required here. Note
	// that, in some portion of a query, schema names are not required.
	if ctx.ReqSchema && r.Schema != nil {
		if err = r.Schema.toSQL(ctx, buf); err != nil {
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
	Name   string
	Quoted bool
}

func Schema(name string) *SchemaExp {
	return &SchemaExp{Name: name}
}

func (s *SchemaExp) Quote() *SchemaExp {
	s.Quoted = true
	return s
}

func (s *SchemaExp) Unquote() *SchemaExp {
	s.Quoted = false
	return s
}

func (s SchemaExp) toSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	name := s.Name
	if s.Quoted {
		name = `"` + name + `""`
	}
	buf.WriteString(name)
	return
}
