// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"bytes"
	"encoding/binary"
)

// An Image File Directory (IFD) consists of a 2-byte count of the number of directory entries, followed by a
// sequence of 12-byte field entries, followed by a 4-byte offset of the next IFD (or 0 if none).
//
// There must be at least 1 IFD in a TIFF file and each IFD must have at least one entry.
//
// Each TIFF field has an associated Count.
// This means that all fields are actually one-dimensional arrays, even though most fields contain only a single value.
type tiffIFD struct {
	next    uint32   // offset in bytes to the next IFD, from the start of the file. 0 if none
	entries uint16   // number of directory entries
	tags    []ifdTag // list of undecoded tags
}

type ifdTag struct {
	id       uint16 // tag identifier
	tiffType uint16 // tiff type idendifier
	count    uint32 // the number of values in data
	data     []byte // undecoded payload for tag
}

// TIFF types as defined in page 15 of TIFF Revision 6.0
const (
	ttBYTE      uint16 = 1
	ttASCII            = 2
	ttSHORT            = 3
	ttLONG             = 4
	ttRATIONAL         = 5
	ttSBYTE            = 6
	ttUNDEFINED        = 7
	ttSSHORT           = 8
	ttSLONG            = 9
	ttSRATIONAL        = 10
	ttFLOAT            = 11
	ttDOUBLE           = 12
)

var tiffTypes = map[uint16]struct {
	name string
	size uint32
}{
	ttBYTE:      {name: "BYTE", size: 1},
	ttASCII:     {name: "ASCII", size: 1},
	ttSHORT:     {name: "SHORT", size: 2},
	ttLONG:      {name: "LONG", size: 4},
	ttRATIONAL:  {name: "RATIONAL", size: 8},
	ttSBYTE:     {name: "SBYTE", size: 1},
	ttUNDEFINED: {name: "UNDEFINED", size: 1},
	ttSSHORT:    {name: "SSHORT", size: 2},
	ttSLONG:     {name: "SLONG", size: 4},
	ttSRATIONAL: {name: "SRATIONAL", size: 8},
	ttFLOAT:     {name: "FLOAT", size: 4},
	ttDOUBLE:    {name: "DOUBLE", size: 8},
}

func (it ifdTag) byteToInt(bo binary.ByteOrder) []int {
	b := make([]int, it.count)
	raw := bytes.NewReader(it.data)
	var v uint8
	for i := range b {
		binary.Read(raw, bo, &v)
		b[i] = int(v)
	}
	return b
}

func (it ifdTag) asciiToString() string {
	return string(it.data[0 : it.count-1]) // -1 to remove character '\0'
}

func (it ifdTag) shortToUint16(bo binary.ByteOrder) []uint16 {
	var s uint16
	S := make([]uint16, it.count)
	raw := bytes.NewReader(it.data)
	for i := range S {
		binary.Read(raw, bo, &s)
		S[i] = s
	}
	return S
}

func (it ifdTag) longToUint32(bo binary.ByteOrder) []uint32 {
	var l uint32
	L := make([]uint32, it.count)
	raw := bytes.NewReader(it.data)
	for i := range L {
		binary.Read(raw, bo, &l)
		L[i] = l
	}
	return L
}

func (it ifdTag) rationalToFloat32(bo binary.ByteOrder) []float32 {
	var n, d uint32
	r := make([]float32, it.count)
	raw := bytes.NewReader(it.data)
	for i := range r {
		binary.Read(raw, bo, &n)
		binary.Read(raw, bo, &d)
		r[i] = float32(n / d)
	}
	return r
}

func (it ifdTag) undefinedToString() string {
	return string(it.data[0:it.count])
}

/*** ***/
// func (raw ifdTag) decode(dict tagDictionary, bo binary.ByteOrder) Tag {
// 	tag := Tag{raw: raw}
// 	if dict != nil {
// 		if val, ok := dict[raw.id]; ok {
// 			tag.name = val.Name
// 		} else {
// 			tag.name = "undefined"
// 		}
// 	} else {
// 		tag.name = "undefined"
// 	}

// 	switch raw.tiffType {
// 	case ttBYTE:
// 		var v uint8
// 		tag.intValues = make([]int64, raw.count)
// 		raw := bytes.NewReader(raw.data)
// 		for i := range tag.intValues {
// 			binary.Read(raw, bo, &v)
// 			tag.intValues[i] = int64(v)
// 		}
// 	case ttASCII:
// 		tag.strValue = string(raw.data[0 : raw.count-1]) // -1 to remove character '\0'
// 	case ttSHORT:
// 		var v uint16
// 		tag.intValues = make([]int64, raw.count)
// 		raw := bytes.NewReader(raw.data)
// 		for i := range tag.intValues {
// 			binary.Read(raw, bo, &v)
// 			tag.intValues[i] = int64(v)
// 		}
// 	case ttLONG:
// 		var v uint32
// 		tag.intValues = make([]int64, raw.count)
// 		raw := bytes.NewReader(raw.data)
// 		for i := range tag.intValues {
// 			binary.Read(raw, bo, &v)
// 			tag.intValues[i] = int64(v)
// 		}
// 	case ttRATIONAL:
// 		var n, d uint32
// 		tag.ratValues = make([][]int64, raw.count)
// 		raw := bytes.NewReader(raw.data)
// 		for i := range tag.ratValues {
// 			binary.Read(raw, bo, &n)
// 			binary.Read(raw, bo, &d)
// 			tag.ratValues[i] = []int64{int64(n), int64(d)}
// 		}
// 	case ttSBYTE:
// 		var v int8
// 		tag.intValues = make([]int64, raw.count)
// 		raw := bytes.NewReader(raw.data)
// 		for i := range tag.intValues {
// 			binary.Read(raw, bo, &v)
// 			tag.intValues[i] = int64(v)
// 		}
// 	case ttSSHORT:
// 		var v int16
// 		tag.intValues = make([]int64, raw.count)
// 		raw := bytes.NewReader(raw.data)
// 		for i := range tag.intValues {
// 			binary.Read(raw, bo, &v)
// 			tag.intValues[i] = int64(v)
// 		}
// 	case ttSLONG:
// 		var v int32
// 		tag.intValues = make([]int64, raw.count)
// 		raw := bytes.NewReader(raw.data)
// 		for i := range tag.intValues {
// 			binary.Read(raw, bo, &v)
// 			tag.intValues[i] = int64(v)
// 		}
// 	case ttSRATIONAL:
// 		var n, d int32
// 		tag.ratValues = make([][]int64, raw.count)
// 		raw := bytes.NewReader(raw.data)
// 		for i := range tag.ratValues {
// 			binary.Read(raw, bo, &n)
// 			binary.Read(raw, bo, &d)
// 			tag.ratValues[i] = []int64{int64(n), int64(d)}
// 		}
// 	}

// 	return tag
// }
