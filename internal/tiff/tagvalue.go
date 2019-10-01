// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package tiff

import (
	"fmt"
	"math"
	"strings"
)

type TagValue struct {
	dataType  DataType // The scalar type of the data items
	count     uint32   // The number of values in the raw data
	raw       []byte   // The raw value of data
	intValues []int64
	ratValues [][]int64
	strValue  string
}

func (tv TagValue) String() string {
	switch tv.dataType {
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
	switch tv.dataType {
	case BYTE, SBYTE:
		return int8(tv.intValues[i])
	default:
		return 0
	}
}

func (tv TagValue) Int16(i int) int16 {
	switch tv.dataType {
	case SHORT, SSHORT:
		return int16(tv.intValues[i])
	default:
		return 0
	}
}

func (tv TagValue) Int32(i int) int32 {
	switch tv.dataType {
	case LONG, SLONG:
		return int32(tv.intValues[i])
	default:
		return 0
	}
}

func (tv TagValue) UInt8(i int) uint8 {
	switch tv.dataType {
	case BYTE:
		return uint8(tv.intValues[i])
	default:
		return 0
	}
}

func (tv TagValue) UInt16(i int) uint16 {
	switch tv.dataType {
	case SHORT:
		return uint16(tv.intValues[i])
	default:
		return 0
	}
}

func (tv TagValue) UInt32(i int) uint32 {
	switch tv.dataType {
	case LONG:
		return uint32(tv.intValues[i])
	default:
		return 0
	}
}
