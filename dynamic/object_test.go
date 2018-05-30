package dynamic

import (
	"testing"
	"fmt"
	"github.com/tsealex/dbutil"
	"encoding/json"
	"reflect"
	"database/sql"
)

func TestNewObject(t *testing.T) {
	intField1 := NewIntField("IntOne", false, 64)
	intField2 := NewIntField("IntTwo", true, 32)
	obj := NewObject(intField1, intField2)
	ptr := obj.CreateInstance()
	fmt.Println(ptr)

	fmt.Println(dbutil.Instance.Get(ptr,
		"SELECT 1 AS IntOne, 2 AS IntTwo"))
	fmt.Println(ptr)

	b, err := json.Marshal(ptr)
	fmt.Println(err)
	fmt.Println(string(b))

	fmt.Println(reflect.ValueOf(ptr).Elem().FieldByName("IntTwo").
		Interface().(sql.NullInt64))

	fmt.Println("------------------------")
	fmt.Println(GetField(ptr, "IntOne"))
	fmt.Println(GetField("hello", ""))
	fmt.Println(GetField(obj.CreateSlice(), ""))
	fmt.Println(GetField(ptr, "FieldNotExist"))



	//fmt.Println(ptr.(*struct{
	//	IntOne int64
	//	IntTwo int32
	//}).IntOne)

	ptr = obj.CreateSlice()
	fmt.Println(ptr)
	fmt.Println(dbutil.Instance.Select(ptr,
		"SELECT 1 AS IntOne, 2 AS IntTwo"))
	fmt.Println(ptr)

	fmt.Println("------------------------")
	fmt.Println(GetLen(ptr))
	fmt.Println(GetElem(ptr, 1))
	fmt.Println(GetElem("hello", 1))
	fmt.Println(GetElem(obj.CreateInstance(), 1))
	elem, err := GetElem(ptr, 0)
	fmt.Println(elem, err)
	fmt.Println(GetField(elem, "IntOne"))

	b, err = json.Marshal(ptr)
	fmt.Println(err)
	fmt.Println(string(b))

	fmt.Println(">>>>>>>>>>>>>>>>>>>>")
	s := reflect.New(obj.structRepr)
	b, err = json.Marshal(s.Interface())
	fmt.Println(err)
	fmt.Println(string(b))



	//var v = struct {
	//	IntOne int64
	//	IntTwo int32
	//}{1, 2}
	//b, err = json.Marshal(v)
	//fmt.Println(err)
	//fmt.Println(string(b))
}