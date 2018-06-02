package null

import (
	"testing"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestBoolArray_Scan(t *testing.T) {
	instance, err := sql.Open("postgres",
		"dbname=postgres host=localhost port=6000 sslmode=disable")
	if err != nil {
		panic(err)
	}
	Instance := sqlx.NewDb(instance, "postgres")

	tmp := BoolArray{}
	assert.NoError(t, Instance.QueryRow("SELECT ARRAY[false, true]").Scan(&tmp))

	assert.Equal(t, true, tmp.Valid)
	assert.Equal(t, []bool{false, true}, []bool(tmp.BoolArray))

	s := struct{
		B *BoolArray
	}{}
	assert.NoError(t, Instance.Get(&s, "SELECT ARRAY[false, true] AS B"))
	assert.Equal(t, true, s.B.Valid)
	assert.Equal(t, []bool{false, true}, []bool(s.B.BoolArray))

	k := struct{
		B BoolArray
	}{}
	assert.NoError(t, Instance.Get(&k, "SELECT ARRAY[false, true] AS B"))
	assert.Equal(t, true, k.B.Valid)
	assert.Equal(t, []bool{false, true}, []bool(k.B.BoolArray))

}