package query

import (
	"database/sql"
	"github.com/tsealex/dbutil/dynamic"
)

type Queryable interface {
	QueryOnce(...interface{}) (dynamic.ObjectPointer, error)
	Query(...interface{}) (dynamic.SlicePointer, error)
	Exec(...interface{}) (sql.Result, error)
}

