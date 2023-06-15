package models

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	textPrefix = "SRID=4326;POINT("
)

type Point4326 struct {
	Lat float64 `json:"lat"` // vi tuyen (pgis: y)
	Lng float64 `json:"lng"` // kinh tuyen:(pgis: x)
} // @name models.Point4326

func (p *Point4326) String() string {
	return fmt.Sprintf("SRID=4326;POINT(%v %v)", p.Lng, p.Lat)
}

// Scan :
func (p *Point4326) Scan(val interface{}) error {
	var s = ""
	switch t := val.(type) {
	case string:
		s = val.(string)
	case []byte:
		s = string(val.([]byte))
	default:
		fmt.Println("Positon4326 scan input type invalid ", t)
		return errors.New("type is invalid")
	}
	idx := strings.Index(s, textPrefix)
	if idx >= 0 {
		return p.fromEWKT(s)
	}
	return p.fromEWKB(s)
}

func (p Point4326) Value() (driver.Value, error) {
	return p.String(), nil
}

func (p *Point4326) MakePoint() string {
	return fmt.Sprintf("ST_SetSRID(ST_MakePoint(%f, %f), 4326)", p.Lng, p.Lat)
}

func (p *Point4326) fromEWKT(val string) error {
	s := val[len(textPrefix):]
	s = s[:len(s)-1]
	// fmt.Println("s : ", s)
	ss := strings.Split(strings.TrimSpace(s), " ")
	// fmt.Println("SS: ", ss)
	if len(ss) != 2 {
		return fmt.Errorf("format of %s is invalid", val)
	}
	var err error
	p.Lng, err = strconv.ParseFloat(ss[0], 64)
	if nil != err {
		return err
	}
	p.Lat, err = strconv.ParseFloat(ss[1], 64)
	if nil != err {
		return err
	}

	return nil
}

func (p *Point4326) fromEWKB(val string) error {
	b, err := hex.DecodeString(val)
	if err != nil {
		return err
	}

	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("invalid byte order %d", wkbByteOrder)
	}

	var wkbGeometryType uint64
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

// select *,
// ST_DistanceSphere(states."location",
// ST_SetSRID(ST_MakePoint(105.834160, 21.037763), 4326)) as distance
// from states
// where ST_DWithin(states."location", ST_SetSRID(ST_MakePoint(105.834160, 21.027763), 4326), 1)
// order by ST_Distance(states."location", ST_SetSRID(ST_MakePoint(105.834160, 21.027763), 4326));
