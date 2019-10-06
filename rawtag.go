// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"bytes"
	"encoding/binary"
)

type rawTag struct {
	id       uint16   // tag identifier
	tiffType tiffType // tiff type idendifier
	count    uint32   // the number of values in data
	data     []byte   // undecoded payload for tag
}

// tiffType is the TIFF data type as defined in page 15 of TIFF Revision 6.0
type tiffType uint16

const (
	BYTE      tiffType = 1
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

var tiffTypes = map[tiffType]struct {
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

func (raw rawTag) decode(dict TagDictionary, bo binary.ByteOrder) Tag {
	tag := Tag{raw: raw}
	if dict != nil {
		if val, ok := dict[raw.id]; ok {
			tag.name = val.Name
		} else {
			tag.name = "undefined"
		}
	} else {
		tag.name = "undefined"
	}

	switch raw.tiffType {
	case BYTE:
		var v uint8
		tag.intValues = make([]int64, raw.count)
		raw := bytes.NewReader(raw.data)
		for i := range tag.intValues {
			binary.Read(raw, bo, &v)
			tag.intValues[i] = int64(v)
		}
	case ASCII:
		tag.strValue = string(raw.data[0 : raw.count-1]) // -1 to remove character '\0'
	case SHORT:
		var v uint16
		tag.intValues = make([]int64, raw.count)
		raw := bytes.NewReader(raw.data)
		for i := range tag.intValues {
			binary.Read(raw, bo, &v)
			tag.intValues[i] = int64(v)
		}
	case LONG:
		var v uint32
		tag.intValues = make([]int64, raw.count)
		raw := bytes.NewReader(raw.data)
		for i := range tag.intValues {
			binary.Read(raw, bo, &v)
			tag.intValues[i] = int64(v)
		}
	case RATIONAL:
		var n, d uint32
		tag.ratValues = make([][]int64, raw.count)
		raw := bytes.NewReader(raw.data)
		for i := range tag.ratValues {
			binary.Read(raw, bo, &n)
			binary.Read(raw, bo, &d)
			tag.ratValues[i] = []int64{int64(n), int64(d)}
		}
	case SBYTE:
		var v int8
		tag.intValues = make([]int64, raw.count)
		raw := bytes.NewReader(raw.data)
		for i := range tag.intValues {
			binary.Read(raw, bo, &v)
			tag.intValues[i] = int64(v)
		}
	case SSHORT:
		var v int16
		tag.intValues = make([]int64, raw.count)
		raw := bytes.NewReader(raw.data)
		for i := range tag.intValues {
			binary.Read(raw, bo, &v)
			tag.intValues[i] = int64(v)
		}
	case SLONG:
		var v int32
		tag.intValues = make([]int64, raw.count)
		raw := bytes.NewReader(raw.data)
		for i := range tag.intValues {
			binary.Read(raw, bo, &v)
			tag.intValues[i] = int64(v)
		}
	case SRATIONAL:
		var n, d int32
		tag.ratValues = make([][]int64, raw.count)
		raw := bytes.NewReader(raw.data)
		for i := range tag.ratValues {
			binary.Read(raw, bo, &n)
			binary.Read(raw, bo, &d)
			tag.ratValues[i] = []int64{int64(n), int64(d)}
		}
	}

	return tag
}
