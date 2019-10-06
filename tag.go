// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

type TagDictionary map[uint16]struct {
	Name string
}

/*
A Tag is a 12-byte record containing TagID identifying the type of information the tag contains and its value.

Each Tag has an associated Count. This means that all fields are actually one-dimensional arrays,
even though most fields contain only a single value.

To save time and space the Offset contains directly the Value instead of pointing to the Value if and only if the value fits into 4 bytes.
*/

type Tag struct {
	id   uint16 // Tag identifier
	name string // Tag name as decoded using a TagDictionary
	//value     TagValue
	raw       rawTag
	intValues []int64
	ratValues [][]int64
	strValue  string
}

func (t Tag) ID() uint16 {
	return t.id
}

func (t Tag) Type() string {
	return tiffTypes[t.raw.tiffType].name
}

func (t Tag) Name() string {
	return t.name
}

// func (t Tag) Value() TagValue {
// 	return t.value
// }

func (t *Tag) decode(dict TagDictionary, bo binary.ByteOrder) {
	if dict != nil {
		if val, ok := dict[t.id]; ok {
			t.name = val.Name
		} else {
			t.name = "undefined"
		}
	} else {
		t.name = "undefined"
	}

	switch t.raw.tiffType {
	case BYTE:
		var v uint8
		t.intValues = make([]int64, t.raw.count)
		raw := bytes.NewReader(t.raw.data)
		for i := range t.intValues {
			binary.Read(raw, bo, &v)
			t.intValues[i] = int64(v)
		}
	case ASCII:
		t.strValue = string(t.raw.data[0 : t.raw.count-1]) // -1 to remove character '\0'
	case SHORT:
		var v uint16
		t.intValues = make([]int64, t.raw.count)
		raw := bytes.NewReader(t.raw.data)
		for i := range t.intValues {
			binary.Read(raw, bo, &v)
			t.intValues[i] = int64(v)
		}
	case LONG:
		var v uint32
		t.intValues = make([]int64, t.raw.count)
		raw := bytes.NewReader(t.raw.data)
		for i := range t.intValues {
			binary.Read(raw, bo, &v)
			t.intValues[i] = int64(v)
		}
	case RATIONAL:
		var n, d uint32
		t.ratValues = make([][]int64, t.raw.count)
		raw := bytes.NewReader(t.raw.data)
		for i := range t.ratValues {
			binary.Read(raw, bo, &n)
			binary.Read(raw, bo, &d)
			t.ratValues[i] = []int64{int64(n), int64(d)}
		}
	case SBYTE:
		var v int8
		t.intValues = make([]int64, t.raw.count)
		raw := bytes.NewReader(t.raw.data)
		for i := range t.intValues {
			binary.Read(raw, bo, &v)
			t.intValues[i] = int64(v)
		}
	case SSHORT:
		var v int16
		t.intValues = make([]int64, t.raw.count)
		raw := bytes.NewReader(t.raw.data)
		for i := range t.intValues {
			binary.Read(raw, bo, &v)
			t.intValues[i] = int64(v)
		}
	case SLONG:
		var v int32
		t.intValues = make([]int64, t.raw.count)
		raw := bytes.NewReader(t.raw.data)
		for i := range t.intValues {
			binary.Read(raw, bo, &v)
			t.intValues[i] = int64(v)
		}
	case SRATIONAL:
		var n, d int32
		t.ratValues = make([][]int64, t.raw.count)
		raw := bytes.NewReader(t.raw.data)
		for i := range t.ratValues {
			binary.Read(raw, bo, &n)
			binary.Read(raw, bo, &d)
			t.ratValues[i] = []int64{int64(n), int64(d)}
		}
	}
}

func (t Tag) String() string {
	switch t.raw.tiffType {
	case ASCII:
		return t.strValue
	case BYTE, SHORT, LONG, SBYTE, SSHORT, SLONG:
		return strings.Trim(strings.Replace(fmt.Sprint(t.intValues), " ", " ", -1), "[]")
	case RATIONAL, SRATIONAL:
		return strings.Trim(strings.Replace(strings.Replace(fmt.Sprint(t.ratValues), " ", "/", -1), "]/[", " ", -1), "[]")
	default:
		// Only 32 firsts characters postfixed with "..."
		return fmt.Sprintf("%s...", string(t.raw.data[0:int(math.Min(float64(len(t.raw.data)), 32))]))
	}
}

func (t Tag) Int8(i int) int8 {
	switch t.raw.tiffType {
	case BYTE, SBYTE:
		return int8(t.intValues[i])
	default:
		return 0
	}
}

func (t Tag) Int16(i int) int16 {
	switch t.raw.tiffType {
	case SHORT, SSHORT:
		return int16(t.intValues[i])
	default:
		return 0
	}
}

func (t Tag) Int32(i int) int32 {
	switch t.raw.tiffType {
	case LONG, SLONG:
		return int32(t.intValues[i])
	default:
		return 0
	}
}

func (t Tag) UInt8(i int) uint8 {
	switch t.raw.tiffType {
	case BYTE:
		return uint8(t.intValues[i])
	default:
		return 0
	}
}

func (t Tag) UInt16(i int) uint16 {
	switch t.raw.tiffType {
	case SHORT:
		return uint16(t.intValues[i])
	default:
		return 0
	}
}

func (t Tag) UInt32(i int) uint32 {
	switch t.raw.tiffType {
	case LONG:
		return uint32(t.intValues[i])
	default:
		return 0
	}
}
