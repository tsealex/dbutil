package dbtype

import (
	"encoding/json"
	"errors"
	"database/sql/driver"
)

type Jsonb map[string]interface{}

// Reference:
// http://coussej.github.io/2016/02/16/Handling-JSONB-in-Go-Structs/

func (j Jsonb) Value() (driver.Value, error) {
	res, err := json.Marshal(j)
	return res, err
}

func (j *Jsonb) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*j, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}

// Return defaultValue if key does not exist in the Jsonb object.the
func (j *Jsonb) GetOrDefault(key string, defaultValue interface{}) interface{} {
	if value, in := (*j)[key]; in {
		return value
	}
	return defaultValue
}
