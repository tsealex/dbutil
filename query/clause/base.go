package clause

import (
	"github.com/tsealex/dbutil/query"
	"bytes"
	"github.com/tsealex/dbutil/query/exp"
)

type Clause interface {
	ToSQL(*query.SQLContext, *bytes.Buffer) error
}

type HavingClause struct {
	cond *exp.CondExp
}

func Having(cond ... interface{}) *HavingClause {
	h := HavingClause{}
	tmp := make([]exp.Exp, len(cond))
	for i, e := range cond {
		if t := getExp(e); t != nil {
			tmp[i] = t
		}
	}
	h.cond = exp.And(tmp...)
	return &h
}

func (hc *HavingClause) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteString("HAVING ")
	if err = hc.cond.ToSQL(ctx, buf); err != nil {
		return
	}
	return
}

type GroupByClause struct {
	cols []exp.Exp
}

func GroupBy(cols ... interface{}) *GroupByClause {
	gb := GroupByClause{}
	for _, col := range cols {
		e := getExp(col)
		gb.cols = append(gb.cols, e)
	}
	return &gb
}

func (oc *GroupByClause) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteString("GROUP BY ")
	if err = concatExps(",", ctx, buf, oc.cols); err != nil {
		return
	}
	return
}

type OrderByClause struct {
	cols []exp.Exp
}

func Order() *OrderByClause {
	return &OrderByClause{}
}

func (oc *OrderByClause) By(order string, cols ... interface{}) *OrderByClause {
	for _, col := range cols {
		e := getExp(col)
		oc.cols = append(oc.cols, exp.RightUnary(e, order))
	}
	return oc
}

const (
	// Orders.
	ASC  = "ASC"
	DESC = "DESC"
)

func (oc *OrderByClause) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteString("ORDER BY ")
	if err = concatExps(",", ctx, buf, oc.cols); err != nil {
		return
	}
	return
}

type OnConflictClause struct {
	cols    []exp.Exp
	write   []exp.Exp
	nothing bool
}

func (oc *OnConflictClause) Write(col exp.Exp, val interface{}) *OnConflictClause {
	colExp := getExp(col)
	rightVal := getExp(val)
	if colExp == nil || rightVal == nil {
		return oc
	}
	oc.write = append(oc.write, exp.Assign(col, rightVal))
	return oc
}

func (oc *OnConflictClause) DoNothing() *OnConflictClause {
	oc.nothing = true
	return oc
}

func (oc *OnConflictClause) ToSQL(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	buf.WriteString("ON CONFLICT (")
	if err = concatExps(",", ctx, buf, oc.cols); err != nil {
		return
	}
	buf.WriteByte(')')
	if oc.nothing || len(oc.write) == 0 {
		buf.WriteString(" DO NOTHING")
	} else {
		buf.WriteString(" DO UPDATE SET ")
		if err = concatExps(",", ctx, buf, oc.write); err != nil {
			return
		}
	}
	return
}

func OnConflict(cols ... interface{}) *OnConflictClause {
	res := OnConflictClause{}
	for _, col := range cols {
		if e := getExp(col); e != nil {
			res.cols = append(res.cols, e)
		}
	}
	return &res
}

type SQLRecipe struct {
	read  []exp.Exp    // Any but AssignExp, CondExp, RelationExp and SchemaExp
	write []exp.Exp    // AssignExp
	cond  *exp.CondExp // CondExp
	addlClauses []Clause
}

func SQL() *SQLRecipe {
	return &SQLRecipe{}
}

func (r *SQLRecipe) AddClause(clauses ... Clause) *SQLRecipe {
	r.addlClauses = append(r.addlClauses, clauses...)
	return r
}

func (r *SQLRecipe) Read(exps ... interface{}) *SQLRecipe {
	tmp := make([]exp.Exp, len(exps))
	for i, e := range exps {
		if t := getExp(e); t != nil {
			tmp[i] = t
		}
	}
	r.read = append(r.read, tmp...)
	return r
}

func (r *SQLRecipe) Write(col exp.Exp, val interface{}) *SQLRecipe {
	colExp := getExp(col)
	rightVal := getExp(val)
	if colExp == nil || rightVal == nil {
		return r
	}
	r.write = append(r.write, exp.Assign(col, rightVal))
	return r
}

func (r *SQLRecipe) Where(cond ... interface{}) *SQLRecipe {
	tmp := make([]exp.Exp, len(cond))
	for i, e := range cond {
		if t := getExp(e); t != nil {
			tmp[i] = t
		}
	}
	if r.cond != nil {
		r.cond = r.cond.And(tmp...)
	} else {
		r.cond = exp.And(tmp...)
	}
	return r
}

func (r *SQLRecipe) Select(tables ... interface{}) (q string, err error) {
	var tableExps = make([]exp.Exp, len(tables))
	for i, tb := range tables {
		tableExps[i] = getExp(tb)
	}
	buf := &bytes.Buffer{}
	ctx := query.NewSQLContext()
	buf.WriteString("SELECT ")
	if err = concatExps(",", ctx, buf, r.read); err != nil {
		return
	}
	if len(tableExps) > 0 {
		buf.WriteString(" FROM ")
		if err = concatExps(",", ctx, buf, tableExps); err != nil {
			return
		}
	}
	if r.cond != nil {
		buf.WriteString(" WHERE ")
		if err = r.cond.ToSQL(ctx, buf); err != nil {
			return
		}
	}
	// Parses additional clauses.
	r.parseClauses(ctx, buf)
	q = buf.String()
	return
}

func (r *SQLRecipe) Update(table exp.Exp) (q string, err error) {
	buf := &bytes.Buffer{}
	ctx := query.NewSQLContext()
	buf.WriteString("UPDATE ")
	if err = table.ToSQL(ctx, buf); err != nil {
		return
	}
	if len(r.write) > 0 {
		buf.WriteString(" SET ")
		if err = concatExps(",", ctx, buf, r.write); err != nil {
			return
		}
	}
	if r.cond != nil {
		buf.WriteString(" WHERE ")
		if err = r.cond.ToSQL(ctx, buf); err != nil {
			return
		}
	}
	if len(r.read) > 0 {
		buf.WriteString(" RETURNING ")
		if err = concatExps(",", ctx, buf, r.read); err != nil {
			return
		}
	}
	// Parses additional clauses.
	r.parseClauses(ctx, buf)
	q = buf.String()
	return
}

func (r *SQLRecipe) Delete(table exp.Exp) (q string, err error) {
	buf := &bytes.Buffer{}
	ctx := query.NewSQLContext()
	buf.WriteString("DELETE FROM ")
	if err = table.ToSQL(ctx, buf); err != nil {
		return
	}
	if r.cond != nil {
		buf.WriteString(" WHERE ")
		if err = r.cond.ToSQL(ctx, buf); err != nil {
			return
		}
	}
	if len(r.read) > 0 {
		buf.WriteString(" RETURNING ")
		if err = concatExps(",", ctx, buf, r.read); err != nil {
			return
		}
	}
	// Parses additional clauses.
	r.parseClauses(ctx, buf)
	q = buf.String()
	return
}

func (r *SQLRecipe) Insert(table exp.Exp) (q string, err error) {
	buf := &bytes.Buffer{}
	ctx := query.NewSQLContext()
	buf.WriteString("INSERT INTO ")
	if err = table.ToSQL(ctx, buf); err != nil {
		return
	}
	buf.WriteByte('\n')
	if len(r.write) > 0 {
		ctx.WriteStatus = query.ColumnOnly
		buf.WriteByte('(')
		if err = concatExps(",", ctx, buf, r.write); err != nil {
			return
		}
		ctx.WriteStatus = query.ValueOnly
		buf.WriteString(") VALUES (")
		if err = concatExps(",", ctx, buf, r.write); err != nil {
			return
		}
		buf.WriteByte(')')
		ctx.WriteStatus = query.Regular
	}
	if len(r.read) > 0 {
		buf.WriteString(" RETURNING ")
		if err = concatExps(",", ctx, buf, r.read); err != nil {
			return
		}
	}
	// Parses additional clauses.
	r.parseClauses(ctx, buf)
	q = buf.String()
	return
}

func (r *SQLRecipe) parseClauses(ctx *query.SQLContext, buf *bytes.Buffer) (err error) {
	for _, clause := range r.addlClauses {
		buf.WriteByte(' ')
		if err = clause.ToSQL(ctx, buf); err != nil {
			return
		}
	}
	return
}

func concatExps(str string, ctx *query.SQLContext, buf *bytes.Buffer, exps []exp.Exp) (err error) {
	for i, e := range exps {
		if i > 0 {
			buf.WriteString(str)
		}
		if err = e.ToSQL(ctx, buf); err != nil {
			return
		}
	}
	return
}

func getExp(e interface{}) exp.Exp {
	switch e.(type) {
	case exp.Exp:
		return e.(exp.Exp)
	case string:
		return exp.Expression(e.(string))
	}
	// TODO: Handle other invalid type.
	return nil
}
