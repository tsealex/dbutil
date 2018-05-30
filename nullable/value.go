package nullable

import (
	"encoding/json"
	"database/sql"
	"github.com/lib/pq"
	"time"
)

/////////
//
type String struct {
	sql.NullString
}

func (v String) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	}
	return json.Marshal(nil)
}

func (v *String) UnmarshalJSON(data []byte) error {
	var ptr *string
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}
	if ptr != nil {
		v.Valid = true
		v.String = *ptr
	} else {
		v.Valid = false
	}
	return nil
}

/////////
//
type Int64 struct {
	sql.NullInt64
}

func (v Int64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	}
	return json.Marshal(nil)
}

func (v *Int64) UnmarshalJSON(data []byte) error {
	var ptr *int64
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}
	if ptr != nil {
		v.Valid = true
		v.Int64 = *ptr
	} else {
		v.Valid = false
	}
	return nil
}

/////////
//
type Bool struct {
	sql.NullBool
}

func (v Bool) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Bool)
	}
	return json.Marshal(nil)
}

func (v *Bool) UnmarshalJSON(data []byte) error {
	var ptr *bool
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}
	if ptr != nil {
		v.Valid = true
		v.Bool = *ptr
	} else {
		v.Valid = false
	}
	return nil
}

/////////
//
type Float64 struct {
	sql.NullFloat64
}

func (v Float64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Float64)
	}
	return json.Marshal(nil)
}

func (v *Float64) UnmarshalJSON(data []byte) error {
	var Float64 *float64
	if err := json.Unmarshal(data, &Float64); err != nil {
		return err
	}
	if Float64 != nil {
		v.Valid = true
		v.Float64 = *Float64
	} else {
		v.Valid = false
	}
	return nil
}

///////// TODO: Geography
//
//type Int64 struct {
//	sql.NullInt64
//}
//
//func (v Int64) MarshalJSON() ([]byte, error) {
//	if v.Valid {
//		return json.Marshal(v.Int64)
//	}
//	return json.Marshal(nil)
//}
//
//func (v *Int64) UnmarshalJSON(data []byte) error {
//	var ptr *int64
//	if err := json.Unmarshal(data, &ptr); err != nil {
//		return err
//	}
//	if ptr != nil {
//		v.Valid = true
//		v.Int64 = *ptr
//	} else {
//		v.Valid = false
//	}
//	return nil
//}

///////// TODO: JSONB
//
//type String struct {
//	sql.NullString
//}
//
//func (v String) MarshalJSON() ([]byte, error) {
//	if v.Valid {
//		return json.Marshal(v.String)
//	}
//	return json.Marshal(nil)
//}
//
//func (v *String) UnmarshalJSON(data []byte) error {
//	var str *string
//	if err := json.Unmarshal(data, &str); err != nil {
//		return err
//	}
//	if str != nil {
//		v.Valid = true
//		v.String = *str
//	} else {
//		v.Valid = false
//	}
//	return nil
//}

/////////
//
type Time struct {
	pq.NullTime
}

func (v Time) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	}
	return json.Marshal(nil)
}

func (v *Time) UnmarshalJSON(data []byte) error {
	var ptr *time.Time
	if err := json.Unmarshal(data, &ptr); err != nil {
		return err
	}
	if ptr != nil {
		v.Valid = true
		v.Time = *ptr
	} else {
		v.Valid = false
	}
	return nil
}
