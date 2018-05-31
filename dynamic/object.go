// It is important to create structs dynamically at runtime, especially because
// SQL queries may be generated dynamically. This is what this package aims to
// achieve, to allow developers to create a struct dynamically, thus to be able
// to use instances of that struct to collect query results.
package dynamic

import (
	"reflect"
	"fmt"
)

type abstractObject interface {
	Type() reflect.Type
	CreateInstance() ObjectPointer
}

type Object struct {
	structRepr      reflect.Type
	structSliceRepr reflect.Type
}

func NewObject(fields ... abstractField) *Object {
	obj := Object{}
	structFields := make([]reflect.StructField, len(fields))
	for i, field := range fields {
		name := field.Name()
		structFields[i] = reflect.StructField{
			Name:      name,
			Type:      field.Type(),
			Anonymous: false,
		}
	}
	obj.structRepr = reflect.StructOf(structFields)
	obj.structSliceRepr = reflect.SliceOf(obj.structRepr)
	return &obj
}

func (o *Object) Type() reflect.Type {
	return o.structRepr
}

// Returns a pointer to an instance of this Object.
func (o *Object) CreateInstance() ObjectPointer {
	return reflect.New(o.structRepr).Interface()
}

// Returns a pointer to slice of instances of this Object.
func (o *Object) CreateSlice() SlicePointer {
	return reflect.New(o.structSliceRepr).Interface()
}

func GetField(ptr ObjectPointer, name string) (interface{}, error) {
	if val := reflect.ValueOf(ptr); val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("ptr is not a pointer")
	} else if val = val.Elem(); val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ptr is not a pointer to a struct instance")
	} else if val = val.FieldByName(name); !val.IsValid() {
		return nil, fmt.Errorf("%s is not a field of the struct", name)
	} else {
		return val.Interface(), nil
	}
}


func GetElem(ptr SlicePointer, i int) (ObjectPointer, error) {
	if val := reflect.ValueOf(ptr); val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("ptr is not a pointer")
	} else if val = val.Elem(); val.Kind() != reflect.Slice {
		return nil, fmt.Errorf("ptr is not a slice")
	} else if i >= val.Len() {
		return nil, fmt.Errorf("i (%d) is out of range", i)
	} else if val = val.Index(i); val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ptr is not a pointer to a struct slice")
	} else if !val.CanAddr() {
		return nil, fmt.Errorf("pointed struct instance is unaddressable")
	} else {
		return val.Addr().Interface(), nil
	}
}


func GetLen(ptr SlicePointer) (int, error) {
	if val := reflect.ValueOf(ptr); val.Kind() != reflect.Ptr {
		return 0, fmt.Errorf("ptr is not a pointer")
	} else if val = val.Elem(); val.Kind() != reflect.Slice {
		return 0, fmt.Errorf("ptr is not a slice")
	} else {
		return val.Len(), nil
	}
}

type SlicePointer interface{}

type ObjectPointer interface{}
