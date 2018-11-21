// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package tiff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

// tiffTag implements interface tag.Tag
type tiffTag struct {
	id       uint16
	name     string
	tiffType uint16
	bytes    []byte
	// ?
	intValues      []int64
	rationalValues [][]int64
	stringValue    string
}

func (tag tiffTag) Name() string { return tag.name }

func (tag tiffTag) StringValue() string {
	switch tag.tiffType {
	case 1, 3, 4, 6, 8, 9:
		return strings.Trim(strings.Replace(fmt.Sprint(tag.intValues), " ", " ", -1), "[]")
	case 5, 10:
		return strings.Trim(strings.Replace(strings.Replace(fmt.Sprint(tag.rationalValues), " ", "/", -1), "]/[", " ", -1), "[]")
	case 2:
		return tag.stringValue
	default:
		// Only 32 firsts characters postfixed with "..."
		return fmt.Sprintf("%s...", string(tag.bytes[0:int(math.Min(float64(len(tag.bytes)), 32))]))
	}
}

func (tag tiffTag) Type() string {
	return ifdEntryTypes[tag.tiffType].name
}

func (tag *tiffTag) convertTiffBytesValue(tiffType uint16, count uint32, order binary.ByteOrder) {
	switch tiffType {
	case 1:
		var v uint8
		tag.tiffType = tiffType
		tag.intValues = make([]int64, int(count))
		raw := bytes.NewReader(tag.bytes)
		for i := 0; i < int(count); i++ {
			binary.Read(raw, order, &v) // error ignored
			tag.intValues[i] = int64(v)
		}
	case 2:
		tag.tiffType = tiffType
		tag.stringValue = string(tag.bytes[0 : len(tag.bytes)-1]) // -1 to remove character '\0'
	case 3:
		var v uint16
		tag.tiffType = tiffType
		tag.intValues = make([]int64, int(count))
		raw := bytes.NewReader(tag.bytes)
		for i := 0; i < int(count); i++ {
			binary.Read(raw, order, &v) // error ignored
			tag.intValues[i] = int64(v)
		}
	case 4:
		var v uint32
		tag.tiffType = tiffType
		tag.intValues = make([]int64, int(count))
		raw := bytes.NewReader(tag.bytes)
		for i := 0; i < int(count); i++ {
			binary.Read(raw, order, &v) // error ignored
			tag.intValues[i] = int64(v)
		}
	case 5:
		var n, d uint32
		tag.tiffType = tiffType
		tag.rationalValues = make([][]int64, int(count))
		raw := bytes.NewReader(tag.bytes)
		for i := 0; i < int(count); i++ {
			binary.Read(raw, order, &n) // error ignored
			binary.Read(raw, order, &d) // error ignored
			tag.rationalValues[i] = []int64{int64(n), int64(d)}
		}
	case 6:
		var v int8
		tag.tiffType = tiffType
		tag.intValues = make([]int64, int(count))
		raw := bytes.NewReader(tag.bytes)
		for i := 0; i < int(count); i++ {
			binary.Read(raw, order, &v) // error ignored
			tag.intValues[i] = int64(v)
		}
	case 7, 11, 12:
		tag.tiffType = tiffType
	case 8:
		var v int16
		tag.tiffType = tiffType
		tag.intValues = make([]int64, int(count))
		raw := bytes.NewReader(tag.bytes)
		for i := 0; i < int(count); i++ {
			binary.Read(raw, order, &v) // error ignored
			tag.intValues[i] = int64(v)
		}
	case 9:
		var v int32
		tag.tiffType = tiffType
		tag.intValues = make([]int64, int(count))
		raw := bytes.NewReader(tag.bytes)
		for i := 0; i < int(count); i++ {
			binary.Read(raw, order, &v) // error ignored
			tag.intValues[i] = int64(v)
		}
	case 10:
		var n, d int32
		tag.tiffType = tiffType
		tag.rationalValues = make([][]int64, int(count))
		raw := bytes.NewReader(tag.bytes)
		for i := 0; i < int(count); i++ {
			binary.Read(raw, order, &n) // error ignored
			binary.Read(raw, order, &d) // error ignored
			tag.rationalValues[i] = []int64{int64(n), int64(d)}
		}
	}

}
