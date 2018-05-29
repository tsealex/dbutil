// It is important to create structs dynamically at runtime, especially because
// SQL queries may be generated dynamically. This is what this package aims to
// achieve, to allow developers to create a struct dynamically, thus to be able
// to use instances of that struct to collect query results.
package dynamic

import (
	"reflect"
	"strings"
	"fmt"
	"bytes"
)

type abstractObject interface {
	GetField(v reflect.Value, field string) interface{}
	Type() reflect.Type
	CreateInstance() interface{}
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
		tag := fmt.Sprintf(`json:"%s" db:"%s"`,
			toSnakeCase(name), strings.ToLower(name))
		structFields[i] = reflect.StructField{
			Name:      name,
			Type:      field.Type(),
			Tag:       reflect.StructTag(tag),
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
func (o *Object) CreateInstance() interface{} {
	return reflect.New(o.structRepr).Interface()
}

// Returns a pointer to slice of instances of this Object.
func (o *Object) CreateSlice() interface{} {
	return reflect.New(o.structSliceRepr).Interface()
}

// Helpers:
// Reference: https://gist.github.com/elwinar/14e1e897fdbe4d3432e1
func toSnakeCase(in string) string {
	length := len(in)

	out := bytes.Buffer{}
	for i := 0; i < length; i++ {
		if i > 0 && isUpper(in[i]) && ((i+1 < length &&
			isLower(in[i+1])) || isLower(in[i-1])) {
			out.WriteByte('_')
		}
		out.WriteByte(toLower(in[i]))
	}
	return out.String()
}

func isLower(c uint8) bool {
	return c >= 'a' && c <= 'z'
}

func isUpper(c uint8) bool {
	return c >= 'A' && c <= 'Z'
}

func toLower(c uint8) uint8 {
	if !isLower(c) {
		return c + 32
	}
	return c
}
