// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"fmt"
	"math"
	"strings"
)

type TagValue struct {
	tiffType  tiffType // The scalar type of the data items
	count     uint32   // The number of values in the raw data
	raw       []byte   // The raw value of data
	intValues []int64
	ratValues [][]int64
	strValue  string
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

func (tv TagValue) String() string {
	switch tv.tiffType {
	case ASCII:
		return tv.strValue
	case BYTE, SHORT, LONG, SBYTE, SSHORT, SLONG:
		return strings.Trim(strings.Replace(fmt.Sprint(tv.intValues), " ", " ", -1), "[]")
	case RATIONAL, SRATIONAL:
		return strings.Trim(strings.Replace(strings.Replace(fmt.Sprint(tv.ratValues), " ", "/", -1), "]/[", " ", -1), "[]")
	default:
		// Only 32 firsts characters postfixed with "..."
		return fmt.Sprintf("%s...", string(tv.raw[0:int(math.Min(float64(tv.count), 32))]))
	}
}

func (tv TagValue) Int8(i int) int8 {
	switch tv.tiffType {
	case BYTE, SBYTE:
		return int8(tv.intValues[i])
	default:
		return 0
	}
}

func (tv TagValue) Int16(i int) int16 {
	switch tv.tiffType {
	case SHORT, SSHORT:
		return int16(tv.intValues[i])
	default:
		return 0
	}
}

func (tv TagValue) Int32(i int) int32 {
	switch tv.tiffType {
	case LONG, SLONG:
		return int32(tv.intValues[i])
	default:
		return 0
	}
}

func (tv TagValue) UInt8(i int) uint8 {
	switch tv.tiffType {
	case BYTE:
		return uint8(tv.intValues[i])
	default:
		return 0
	}
}

func (tv TagValue) UInt16(i int) uint16 {
	switch tv.tiffType {
	case SHORT:
		return uint16(tv.intValues[i])
	default:
		return 0
	}
}

func (tv TagValue) UInt32(i int) uint32 {
	switch tv.tiffType {
	case LONG:
		return uint32(tv.intValues[i])
	default:
		return 0
	}
}
