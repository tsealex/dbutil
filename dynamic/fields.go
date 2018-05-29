package dynamic

import "reflect"

type abstractField interface {
	Name() string
	Nullable() bool
	Type() reflect.Type
}

// BaseField is immutable.
type BaseField struct {
	name      string
	nullable  bool
	fieldType reflect.Type
}

func (f *BaseField) Name() string {
	return f.name
}

func (f *BaseField) Nullable() bool {
	return f.nullable
}

func (f *BaseField) Type() reflect.Type {
	return f.fieldType
}

// Integer field.
type IntField struct {
	BaseField
	bits uint8
}

func NewIntField(name string, nullable bool, bits int) *IntField {
	if bits < 0 {
		panic("bits must be greater than or equal to 0")
	} else if bits > 64 {
		panic("bits must be less than or equal to 64")
	} else if len(name) == 0 {
		panic("name must not be empty")
	}
	// TODO: Use regex to check if name begins with a lowercase character or
	// TODO: contains any illegal characters.
	f := IntField{}
	f.name = name
	f.nullable = nullable
	if f.nullable {
		// TODO: set fieldType to be 64-bit nullable integer.
		f.bits = 64
	} else {
		// Set fieldType to be the corresponding primitive int type.
		if bits <= 16 {
			f.bits = 16
			f.fieldType = reflect.TypeOf(int16(0))
		} else if bits <= 32 {
			f.bits = 32
			f.fieldType = reflect.TypeOf(int32(0))
		} else {
			f.bits = 64
			f.fieldType = reflect.TypeOf(int64(0))
		}
	}
	return &f
}
