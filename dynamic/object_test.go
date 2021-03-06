package dynamic

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/tsealex/dbutil/null"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

func TestNewObject(t *testing.T) {
	instance, err := sql.Open("postgres",
		"dbname=postgres host=localhost port=6000 sslmode=disable")
	if err != nil {
		panic(err)
	}
	Instance := sqlx.NewDb(instance, "postgres")

	intField1 := NewIntField("IntOne", false, true, 64)
	intField2 := NewIntField("IntTwo", true, true, 32)
	obj := NewObject(intField1, intField2)
	assert.NotNil(t, obj.CreateEditable())
	assert.NotNil(t, obj.CreateEditableSlice())

	ptr := obj.CreateInstance()
	assert.NotNil(t, ptr)

	assert.NoError(t, Instance.Get(ptr,
		"SELECT 1 AS IntOne, 2 AS IntTwo"))
	assert.NotNil(t, ptr)

	b, err := json.Marshal(ptr)
	assert.NoError(t, err)
	assert.Equal(t, `{"IntOne":1,"IntTwo":2}`, string(b))

	v, err := GetField(ptr, "IntOne")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), v)

	_, err = GetField("hello", "")
	assert.Error(t, err)
	_, err = GetField(obj.CreateSlice(), "")
	assert.Error(t, err)
	_, err = GetField(ptr, "FieldNotExist")
	assert.Error(t, err)

	s, ok := ptr.(*struct{
		IntOne int64
		IntTwo null.Int64
	})
	assert.True(t, ok)
	assert.Equal(t, int64(2), s.IntTwo.Int64)

	ptr = obj.CreateSlice()
	assert.NotNil(t, ptr)
	assert.NoError(t, Instance.Select(ptr,
		"SELECT 1 AS IntOne, 2 AS IntTwo"))

	l, err := GetLen(ptr)
	assert.NoError(t, err)
	assert.Equal(t, l, 1)
	_, err = GetElem(ptr, 1)
	assert.Error(t, err)
	_, err = GetElem("hello", 1)
	assert.Error(t, err)
	_, err = GetElem(obj.CreateInstance(), 1)
	assert.Error(t, err)
	elem, err := GetElem(ptr, 0)
	assert.NoError(t, err)

	v, err = GetField(elem, "IntOne")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), v)

	b, err = json.Marshal(ptr)
	assert.NoError(t, err)
	assert.Equal(t, `[{"IntOne":1,"IntTwo":2}]`, string(b))

	e := obj.CreateEditable()
	assert.NoError(t, json.Unmarshal([]byte(`{"IntTwo":2}`), e))
	f, err := GetField(e, "IntTwo")
	assert.NoError(t, err)
	i := f.(*null.Int64)
	assert.NotNil(t, i)
	assert.True(t, i.Valid)
	assert.Equal(t, int64(2), i.Int64)

	f, err = GetField(e, "IntOne")
	assert.NoError(t, err)
	assert.Nil(t, f)

}