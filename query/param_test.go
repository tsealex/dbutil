package query

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/tsealex/dbutil/dynamic"
	"encoding/json"
	"github.com/tsealex/dbutil/null"
)

func TestPrepareParameters(t *testing.T) {
	names := []string{"Hello", "Two"}
	arg0 := struct {
		Two int
	}{Two: 4}
	arg1 := map[interface{}]string{
		"Hello": "World",
		2: "Okay",
		"OhNo": "Yeah",
	}
	arg2 := map[int]string{1: "No"}
	assert.Equal(t, []interface{}{"World", 4},
		*PrepareParameters(&names, arg0, &arg0, arg2, &arg1))
}

func TestPrepareParameters2(t *testing.T) {
	instance, err := sql.Open("postgres",
		"dbname=postgres host=localhost port=6000 sslmode=disable")
	if err != nil {
		panic(err)
	}
	Instance := sqlx.NewDb(instance, "postgres")
	query := `
	CREATE TABLE testing (
		i INT,
		f FLOAT,
		b BOOLEAN,
		s VARCHAR NOT NULL 
	)`
	_, err = Instance.Exec(query)
	defer Instance.Exec("DROP TABLE testing CASCADE")
	assert.NoError(t, err)
	obj := dynamic.NewObject(
		dynamic.NewIntField("I", true, true, 32),
		dynamic.NewFloatField("F", true, false, 32),
		dynamic.NewBoolField("B", true, true),
		dynamic.NewStringField("S", false, true),
	)
	ptr := obj.CreateInstance()
	assert.NoError(t, json.Unmarshal([]byte(`{"I":32,"B":false,"S":"not null"}`), ptr))

	params := PrepareParameters(&[]string{"I","B","S"}, ptr)
	query = `INSERT INTO testing (i,b,s) VALUES ($1,$2,$3)`
	_, err = Instance.Exec(query, *params...)
	assert.NoError(t, err)

	ptr2 := obj.CreateSlice()
	assert.NoError(t, Instance.Select(ptr2, `SELECT * FROM testing`))

	l, err := dynamic.GetLen(ptr2)
	assert.NoError(t, err)
	assert.Equal(t, 1, l)

	elem, err := dynamic.GetElem(ptr2, 0)
	assert.NoError(t, err)

	i, err := dynamic.GetField(elem, "I")
	assert.NoError(t, err)
	assert.Equal(t, int64(32), i.(null.Int64).Int64)

	i, err = dynamic.GetField(elem, "B")
	assert.NoError(t, err)
	assert.Equal(t, true, i.(null.Bool).Valid)
	assert.Equal(t, false, i.(null.Bool).Bool)

	i, err = dynamic.GetField(elem, "F")
	assert.NoError(t, err)
	assert.Equal(t, false, i.(null.Float64).Valid)

	i, err = dynamic.GetField(elem, "S")
	assert.NoError(t, err)
	assert.Equal(t, "not null", i.(string))

}
