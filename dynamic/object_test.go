package dynamic

import (
	"testing"
	"fmt"
	"github.com/tsealex/dbutil"
	"encoding/json"
)

func TestNewObject(t *testing.T) {
	intField1 := NewIntField("IntOne", false, 64)
	intField2 := NewIntField("IntTwo", false, 32)
	obj := NewObject(intField1, intField2)
	ptr := obj.CreateInstance()
	fmt.Println(ptr)

	fmt.Println(dbutil.Instance.Get(ptr,
		"SELECT 1 AS IntOne, 2 AS IntTwo"))
	fmt.Println(ptr)

	b, err := json.Marshal(ptr)
	fmt.Println(err)
	fmt.Println(string(b))

	ptr = obj.CreateSlice()
	fmt.Println(ptr)
	fmt.Println(dbutil.Instance.Select(ptr,
		"SELECT 1 AS IntOne, 2 AS IntTwo"))
	fmt.Println(ptr)

	b, err = json.Marshal(ptr)
	fmt.Println(err)
	fmt.Println(string(b))
}