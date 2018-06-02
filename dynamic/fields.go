package dynamic

import (
	"reflect"
	"github.com/lib/pq"
	"time"
	"github.com/tsealex/dbutil/null"
	"github.com/tsealex/dbutil/dbtype"
)

type abstractField interface {
	Name() string
	Nullable() bool
	Editable() bool
	Type() reflect.Type
}

// BaseField is immutable.
type BaseField struct {
	name      string
	nullable  bool
	editable  bool
	fieldType reflect.Type
}

func (f *BaseField) Name() string {
	return f.name
}

func (f *BaseField) Nullable() bool {
	return f.nullable
}

func (f *BaseField) Editable() bool {
	return f.editable
}

func (f *BaseField) Type() reflect.Type {
	return f.fieldType
}

// Integer field.
type IntField struct {
	BaseField
	bits uint8
}

func NewIntField(name string, nullable bool, editable bool, bits int) *IntField {
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
	f.editable = editable
	if f.nullable {
		// Set fieldType to be 64-bit null integer.
		f.bits = 64
		f.fieldType = reflect.TypeOf(null.Int64{})
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

type FloatField struct {
	BaseField
	bits uint8
}

func NewFloatField(name string, nullable bool, editable bool, bits int) *FloatField {
	if bits < 0 {
		panic("bits must be greater than or equal to 0")
	} else if bits > 64 {
		panic("bits must be less than or equal to 64")
	} else if len(name) == 0 {
		panic("name must not be empty")
	}
	// TODO: Use regex to check if name begins with a lowercase character or
	// TODO: contains any illegal characters.
	f := FloatField{}
	f.name = name
	f.nullable = nullable
	f.editable = editable
	if f.nullable {
		// TODO: set fieldType to be 64-bit null float.
		f.bits = 64
		f.fieldType = reflect.TypeOf(null.Float64{})
	} else {
		// Set fieldType to be the corresponding primitive int type.
		if bits <= 32 {
			f.bits = 32
			f.fieldType = reflect.TypeOf(float32(0))
		} else {
			f.bits = 64
			f.fieldType = reflect.TypeOf(float64(0))
		}
	}
	return &f
}

type StringField struct {
	BaseField
}

func NewStringField(name string, nullable bool, editable bool) *StringField {
	if len(name) == 0 {
		panic("name must not be empty")
	}
	// TODO: Use regex to check if name begins with a lowercase character or
	// TODO: contains any illegal characters.
	f := StringField{}
	f.name = name
	f.nullable = nullable
	f.editable = editable
	if f.nullable {
		// TODO: set fieldType to be null string.
		f.fieldType = reflect.TypeOf(null.Float64{})
	} else {
		// Set fieldType to be the corresponding primitive int type.
		f.fieldType = reflect.TypeOf(string(""))
	}
	return &f
}

type BoolField struct {
	BaseField
}

func NewBoolField(name string, nullable bool, editable bool) *BoolField {
	if len(name) == 0 {
		panic("name must not be empty")
	}
	// TODO: Use regex to check if name begins with a lowercase character or
	// TODO: contains any illegal characters.
	f := BoolField{}
	f.name = name
	f.nullable = nullable
	f.editable = editable
	if f.nullable {
		// TODO: set fieldType to be null bool.
		f.fieldType = reflect.TypeOf(null.Bool{})
	} else {
		// Set fieldType to be the corresponding primitive int type.
		f.fieldType = reflect.TypeOf(false)
	}
	return &f
}

type TimeField struct {
	BaseField
}

func NewTimeField(name string, nullable bool, editable bool) *TimeField {
	if len(name) == 0 {
		panic("name must not be empty")
	}
	// TODO: Use regex to check if name begins with a lowercase character or
	// TODO: contains any illegal characters.
	f := TimeField{}
	f.name = name
	f.nullable = nullable
	f.editable = editable
	if f.nullable {
		// TODO: set fieldType to be null bool.
		f.fieldType = reflect.TypeOf(pq.NullTime{})
	} else {
		// Set fieldType to be the corresponding primitive int type.
		f.fieldType = reflect.TypeOf(time.Time{})
	}
	return &f
}

type PointField struct {
	BaseField
}

func NewPointField(name string, nullable bool, editable bool) *PointField {
	if len(name) == 0 {
		panic("name must not be empty")
	}
	// TODO: Use regex to check if name begins with a lowercase character or
	// TODO: contains any illegal characters.
	f := PointField{}
	f.name = name
	f.nullable = nullable
	f.editable = editable
	if f.nullable {
		// TODO: set fieldType to be null bool.
		f.fieldType = reflect.TypeOf(null.Point{})
	} else {
		// Set fieldType to be the corresponding primitive int type.
		f.fieldType = reflect.TypeOf(dbtype.Point{})
	}
	return &f
}

type JsonbField struct {
	BaseField
}

func NewJsonbField(name string, nullable bool, editable bool) *JsonbField {
	if len(name) == 0 {
		panic("name must not be empty")
	}
	// TODO: Use regex to check if name begins with a lowercase character or
	// TODO: contains any illegal characters.
	f := JsonbField{}
	f.name = name
	f.nullable = nullable
	f.editable = editable
	if f.nullable {
		// TODO: set fieldType to be null bool.
		f.fieldType = reflect.TypeOf(null.Jsonb{})
	} else {
		// Set fieldType to be the corresponding primitive int type.
		f.fieldType = reflect.TypeOf(dbtype.Jsonb{})
	}
	return &f
}

// Array types.
type StringArrayField struct {
	BaseField
}

func NewStringArrayField(name string, editable bool) *StringArrayField {
	if len(name) == 0 {
		panic("name must not be empty")
	}
	f := StringArrayField{}
	f.name = name
	f.nullable = true
	f.editable = editable
	f.fieldType = reflect.TypeOf(null.StringArray{})
	return &f
}

type IntArrayField struct {
	BaseField
}

func NewIntArrayField(name string, editable bool) *IntArrayField {
	if len(name) == 0 {
		panic("name must not be empty")
	}
	f := IntArrayField{}
	f.name = name
	f.nullable = true
	f.editable = editable
	f.fieldType = reflect.TypeOf(null.Int64Array{})
	return &f
}

type FloatArrayField struct {
	BaseField
}

func NewFloatArrayField(name string, editable bool) *FloatArrayField {
	if len(name) == 0 {
		panic("name must not be empty")
	}
	f := FloatArrayField{}
	f.name = name
	f.nullable = true
	f.editable = editable
	f.fieldType = reflect.TypeOf(null.Float64Array{})
	return &f
}

type BoolArrayField struct {
	BaseField
}

func NewBoolArrayField(name string, editable bool) *BoolArrayField {
	if len(name) == 0 {
		panic("name must not be empty")
	}
	f := BoolArrayField{}
	f.name = name
	f.nullable = true
	f.editable = editable
	f.fieldType = reflect.TypeOf(null.BoolArray{})
	return &f
}
