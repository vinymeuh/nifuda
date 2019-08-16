// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package tiff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
)

/*
A Tag is a 12-byte record containing TagID identifying the type of information the tag contains and its value.

Each Tag has an associated Count. This means that all fields are actually one-dimensional arrays,
even though most fields contain only a single value.

To save time and space the Offset contains directly the Value instead of pointing to the Value if and only if the value fits into 4 bytes.
*/

type Tag struct {
	TagID     uint16   // Tag identifier
	dataType  dataType // The scalar type of the data items
	dataCount uint32   // The number of items in the tag data
	dataRaw   []byte   // The raw value of data

	Name      string
	intValues []int64
	ratValues [][]int64
	strValue  string
}

// dataType is the TIFF data type as defined in page 15 of TIFF Revision 6.0
type dataType uint16

const (
	BYTE      dataType = 1
	ASCII              = 2
	SHORT              = 3
	LONG               = 4
	RATIONAL           = 5
	SBYTE              = 6
	UNDEFINED          = 7
	SSHORT             = 8
	SLONG              = 9
	SRATIONAL          = 10
	FLOAT              = 11
	DOUBLE             = 12
)

var dataTypes = map[dataType]struct {
	name string
	size uint32
}{
	BYTE:      {name: "BYTE", size: 1},
	ASCII:     {name: "ASCII", size: 1},
	SHORT:     {name: "SHORT", size: 2},
	LONG:      {name: "LONG", size: 4},
	RATIONAL:  {name: "RATIONAL", size: 8},
	SBYTE:     {name: "SBYTE", size: 1},
	UNDEFINED: {name: "UNDEFINED", size: 1},
	SSHORT:    {name: "SSHORT", size: 2},
	SLONG:     {name: "SLONG", size: 4},
	SRATIONAL: {name: "SRATIONAL", size: 8},
	FLOAT:     {name: "FLOAT", size: 4},
	DOUBLE:    {name: "DOUBLE", size: 8},
}

type TagDictionary map[uint16]struct {
	Name string
}

func (t Tag) String() string {
	return fmt.Sprintf("id=%d, type=%s, count=%d, value=%s",
		t.TagID, dataTypes[t.dataType].name, t.dataCount, t.StringValue())
}

func (t *Tag) decode(dict TagDictionary, bo binary.ByteOrder) {
	if dict != nil {
		if val, ok := dict[t.TagID]; ok {
			t.Name = val.Name
		} else {
			t.Name = "unknown"
		}
	} else {
		t.Name = "unknown"
	}

	switch t.dataType {
	case BYTE:
		var v uint8
		t.intValues = make([]int64, t.dataCount)
		raw := bytes.NewReader(t.dataRaw)
		for i := range t.intValues {
			binary.Read(raw, bo, &v)
			t.intValues[i] = int64(v)
		}
	case ASCII:
		t.strValue = string(t.dataRaw[0 : t.dataCount-1]) // -1 to remove character '\0'
	case SHORT:
		var v uint16
		t.intValues = make([]int64, t.dataCount)
		raw := bytes.NewReader(t.dataRaw)
		for i := range t.intValues {
			binary.Read(raw, bo, &v)
			t.intValues[i] = int64(v)
		}
	case LONG:
		var v uint32
		t.intValues = make([]int64, t.dataCount)
		raw := bytes.NewReader(t.dataRaw)
		for i := range t.intValues {
			binary.Read(raw, bo, &v)
			t.intValues[i] = int64(v)
		}
	case RATIONAL:
		var n, d uint32
		t.ratValues = make([][]int64, t.dataCount)
		raw := bytes.NewReader(t.dataRaw)
		for i := range t.ratValues {
			binary.Read(raw, bo, &n)
			binary.Read(raw, bo, &d)
			t.ratValues[i] = []int64{int64(n), int64(d)}
		}
	case SBYTE:
		var v int8
		t.intValues = make([]int64, t.dataCount)
		raw := bytes.NewReader(t.dataRaw)
		for i := range t.intValues {
			binary.Read(raw, bo, &v)
			t.intValues[i] = int64(v)
		}
	case SSHORT:
		var v int16
		t.intValues = make([]int64, t.dataCount)
		raw := bytes.NewReader(t.dataRaw)
		for i := range t.intValues {
			binary.Read(raw, bo, &v)
			t.intValues[i] = int64(v)
		}
	case SLONG:
		var v int32
		t.intValues = make([]int64, t.dataCount)
		raw := bytes.NewReader(t.dataRaw)
		for i := range t.intValues {
			binary.Read(raw, bo, &v)
			t.intValues[i] = int64(v)
		}
	case SRATIONAL:
		var n, d int32
		t.ratValues = make([][]int64, t.dataCount)
		raw := bytes.NewReader(t.dataRaw)
		for i := range t.ratValues {
			binary.Read(raw, bo, &n)
			binary.Read(raw, bo, &d)
			t.ratValues[i] = []int64{int64(n), int64(d)}
		}
	}
}

func (t *Tag) StringValue() string {
	switch t.dataType {
	case ASCII:
		return t.strValue
	case BYTE, SHORT, LONG, SBYTE, SSHORT, SLONG:
		return strings.Trim(strings.Replace(fmt.Sprint(t.intValues), " ", " ", -1), "[]")
	case RATIONAL, SRATIONAL:
		return strings.Trim(strings.Replace(strings.Replace(fmt.Sprint(t.ratValues), " ", "/", -1), "]/[", " ", -1), "[]")
	default:
		// Only 32 firsts characters postfixed with "..."
		return fmt.Sprintf("%s...", string(t.dataRaw[0:int(math.Min(float64(t.dataCount), 32))]))
	}
}

func (t *Tag) UInt8(i int) (uint8, error) {
	switch t.dataType {
	case BYTE:
		return uint8(t.intValues[i]), nil
	default:
		return 0, errors.New("not an uint8")
	}
}

func (t *Tag) UInt16(i int) (uint16, error) {
	switch t.dataType {
	case SHORT:
		return uint16(t.intValues[i]), nil
	default:
		return 0, errors.New("not an uint16")
	}
}

func (t *Tag) UInt32(i int) (uint32, error) {
	switch t.dataType {
	case LONG:
		return uint32(t.intValues[i]), nil
	default:
		return 0, errors.New("not an uint32")
	}
}
