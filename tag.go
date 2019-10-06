// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"bytes"
	"encoding/binary"
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
	id    uint16 // Tag identifier
	name  string // Tag name as decoded using a TagDictionary
	value TagValue
}

func (t Tag) ID() uint16 {
	return t.id
}

func (t Tag) Type() string {
	return tiffTypes[t.value.tiffType].name
}

func (t Tag) Name() string {
	return t.name
}

func (t Tag) Value() TagValue {
	return t.value
}

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

	switch t.value.tiffType {
	case BYTE:
		var v uint8
		t.value.intValues = make([]int64, t.value.count)
		raw := bytes.NewReader(t.value.raw)
		for i := range t.value.intValues {
			binary.Read(raw, bo, &v)
			t.value.intValues[i] = int64(v)
		}
	case ASCII:
		t.value.strValue = string(t.value.raw[0 : t.value.count-1]) // -1 to remove character '\0'
	case SHORT:
		var v uint16
		t.value.intValues = make([]int64, t.value.count)
		raw := bytes.NewReader(t.value.raw)
		for i := range t.value.intValues {
			binary.Read(raw, bo, &v)
			t.value.intValues[i] = int64(v)
		}
	case LONG:
		var v uint32
		t.value.intValues = make([]int64, t.value.count)
		raw := bytes.NewReader(t.value.raw)
		for i := range t.value.intValues {
			binary.Read(raw, bo, &v)
			t.value.intValues[i] = int64(v)
		}
	case RATIONAL:
		var n, d uint32
		t.value.ratValues = make([][]int64, t.value.count)
		raw := bytes.NewReader(t.value.raw)
		for i := range t.value.ratValues {
			binary.Read(raw, bo, &n)
			binary.Read(raw, bo, &d)
			t.value.ratValues[i] = []int64{int64(n), int64(d)}
		}
	case SBYTE:
		var v int8
		t.value.intValues = make([]int64, t.value.count)
		raw := bytes.NewReader(t.value.raw)
		for i := range t.value.intValues {
			binary.Read(raw, bo, &v)
			t.value.intValues[i] = int64(v)
		}
	case SSHORT:
		var v int16
		t.value.intValues = make([]int64, t.value.count)
		raw := bytes.NewReader(t.value.raw)
		for i := range t.value.intValues {
			binary.Read(raw, bo, &v)
			t.value.intValues[i] = int64(v)
		}
	case SLONG:
		var v int32
		t.value.intValues = make([]int64, t.value.count)
		raw := bytes.NewReader(t.value.raw)
		for i := range t.value.intValues {
			binary.Read(raw, bo, &v)
			t.value.intValues[i] = int64(v)
		}
	case SRATIONAL:
		var n, d int32
		t.value.ratValues = make([][]int64, t.value.count)
		raw := bytes.NewReader(t.value.raw)
		for i := range t.value.ratValues {
			binary.Read(raw, bo, &n)
			binary.Read(raw, bo, &d)
			t.value.ratValues[i] = []int64{int64(n), int64(d)}
		}
	}
}
