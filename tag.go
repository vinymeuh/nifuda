// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
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
	name      string // Tag name as decoded using a TagDictionary
	raw       rawTag
	intValues []int64
	ratValues [][]int64
	strValue  string
}

func (t Tag) ID() uint16 {
	return t.raw.id
}

func (t Tag) Type() string {
	return tiffTypes[t.raw.tiffType].name
}

func (t Tag) Name() string {
	return t.name
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
