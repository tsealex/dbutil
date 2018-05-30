package dbtype

import (
	"fmt"
	"database/sql/driver"
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

// Reference:
// https://github.com/nferruzzi/gormGIS
type Point struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func (p *Point) String() string {
	return fmt.Sprintf("(Lat: %v, Long: %v)", p.Lat, p.Long)
}

func (p *Point) Scan(v interface{}) error {
	b, err := hex.DecodeString(string(v.([]uint8)))
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	var wkbOrder uint8
	var byteOrder binary.ByteOrder
	if err := binary.Read(r, binary.LittleEndian, &wkbOrder); err != nil {
		return err
	} else if wkbOrder == 0 {
		byteOrder = binary.BigEndian
	} else if wkbOrder == 1 {
		byteOrder = binary.LittleEndian
	} else {
		return fmt.Errorf("invalid WKB representation")
	}

	var wkbType uint64
	if err := binary.Read(r, byteOrder, &wkbType); err != nil {
		return err
	}
	// TODO: check if wkbType is a Point
	// Reference:
	// https://en.wikipedia.org/wiki/Well-known_text#Well-known_binary

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}
	return nil
}

func (p Point) Value() (driver.Value, error) {
	// PostGis stores geography points in long and then lat format
	// https://postgis.net/2013/08/18/tip_lon_lat/NullFloat64Array
	return fmt.Sprintf("SRID=4326;POINT(%v %v)", p.Long, p.Lat), nil
}
