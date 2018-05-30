package null

import (
	"github.com/lib/pq"
	"database/sql/driver"
	"encoding/json"
)

/////////
// 
type Int64Array struct {
	pq.Int64Array
	Valid bool
}

func (p *Int64Array) Scan(value interface{}) error {
	if value == nil {
		p.Int64Array, p.Valid = []int64{}, false
		return nil
	}
	p.Valid = true
	return p.Int64Array.Scan(value)
}

// Value implements the driver Valuer interface.
func (p Int64Array) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil // a hack to make sure no array is null in db
	}
	if p.Int64Array == nil {
		return "{}", nil
	}
	return p.Int64Array.Value()
}

func (v Int64Array) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64Array)
	}
	return json.Marshal(nil)
}

func (v *Int64Array) UnmarshalJSON(data []byte) error {
	var val *[]int64
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	if val != nil {
		v.Valid = true
		v.Int64Array = *val
	} else {
		v.Valid = false
	}
	return nil
}


/////////
// 
type Float64Array struct {
	pq.Float64Array
	Valid bool
}

func (p *Float64Array) Scan(value interface{}) error {
	if value == nil {
		p.Float64Array, p.Valid = []float64{}, false
		return nil
	}
	p.Valid = true
	return p.Float64Array.Scan(value)
}

// Value implements the driver Valuer interface.
func (p Float64Array) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil // a hack to make sure no array is null in db
	}
	if p.Float64Array == nil {
		return "{}", nil
	}
	return p.Float64Array.Value()
}

func (v Float64Array) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Float64Array)
	}
	return json.Marshal(nil)
}

func (v *Float64Array) UnmarshalJSON(data []byte) error {
	var val *[]float64
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	if val != nil {
		v.Valid = true
		v.Float64Array = *val
	} else {
		v.Valid = false
	}
	return nil
}


/////////
// 
type BoolArray struct {
	pq.BoolArray
	Valid bool
}

func (p *BoolArray) Scan(value interface{}) error {
	if value == nil {
		p.BoolArray, p.Valid = []bool{}, false
		return nil
	}
	p.Valid = true
	return p.BoolArray.Scan(value)
}

// Value implements the driver Valuer interface.
func (p BoolArray) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil // a hack to make sure no array is null in db
	}
	if p.BoolArray == nil {
		return "{}", nil
	}
	return p.BoolArray.Value()
}

func (v BoolArray) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.BoolArray)
	}
	return json.Marshal(nil)
}

func (v *BoolArray) UnmarshalJSON(data []byte) error {
	var val *[]bool
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	if val != nil {
		v.Valid = true
		v.BoolArray = *val
	} else {
		v.Valid = false
	}
	return nil
}


/////////
// 
type StringArray struct {
	pq.StringArray
	Valid bool
}

func (p *StringArray) Scan(value interface{}) error {
	if value == nil {
		p.StringArray, p.Valid = []string{}, false
		return nil
	}
	p.Valid = true
	return p.StringArray.Scan(value)
}

// Value implements the driver Valuer interface.
func (p StringArray) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil // a hack to make sure no array is null in db
	}
	if p.StringArray == nil {
		return "{}", nil
	}
	return p.StringArray.Value()
}

func (v StringArray) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.StringArray)
	}
	return json.Marshal(nil)
}

func (v *StringArray) UnmarshalJSON(data []byte) error {
	var val *[]string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	if val != nil {
		v.Valid = true
		v.StringArray = *val
	} else {
		v.Valid = false
	}
	return nil
}