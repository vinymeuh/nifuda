// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.
package tiff

import (
	"encoding/binary"
	"strconv"
	"strings"
	"testing"
)

func TestTagDecode(t *testing.T) {
	tests := []struct {
		id      uint16
		dtType  DataType
		dtCount uint32
		dtRaw   []byte
		// expected values
		name  string
		value string
	}{
		{0, BYTE, 1, []byte{0x7b}, "undefined", "123"},
		{0, ASCII, 6, []byte{'h', 'e', 'l', 'l', 'o', '0'}, "undefined", "hello"},
		{257, SHORT, 2, []byte{0x4, 0x0, 0x0, 0xa}, "ImageLength", "1024 10"},
		{256, LONG, 1, []byte{0x4, 0x0, 0x0, 0x0}, "ImageWidth", "67108864"},
		{0, RATIONAL, 1, []byte{0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x4}, "undefined", "1/4"},
		{0, SBYTE, 2, []byte{0xf8, 0x8}, "undefined", "-8 8"},
		{0, SSHORT, 2, []byte{0xfc, 0x0, 0x4, 0x0}, "undefined", "-1024 1024"},
		{0, SLONG, 2, []byte{0xfc, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0}, "undefined", "-67108864 67108864"},
	}

	for _, tc := range tests {
		tag := Tag{
			id:    tc.id,
			value: TagValue{dataType: tc.dtType, count: tc.dtCount, raw: tc.dtRaw},
		}
		tag.decode(dictExif, binary.BigEndian)

		if tag.ID() != tc.id {
			t.Errorf("[%s]: error decoding tag id, expected=%d, got=%d", tag, tc.id, tag.ID())
		}
		if tag.Type() != dataTypes[tc.dtType].name {
			t.Errorf("[%s]: error decoding tag type, expected=%s, got=%s", tag, dataTypes[tc.dtType].name, tag.Type())
		}
		if tag.Name() != tc.name {
			t.Errorf("[%s]: error decoding tag name, expected=%s, got=%s", tag, tc.name, tag.Name())
		}
		if tag.Value().String() != tc.value {
			t.Errorf("[%s]: error decoding tag value, expected=%s, got=%s", tag, tc.value, tag.Value().String())
		}

		switch tag.value.dataType {
		case BYTE:
			for i, value := range strings.Split(tc.value, " ") {
				expected, _ := strconv.ParseUint(value, 0, 8)
				if tag.Value().UInt8(i) != uint8(expected) {
					t.Errorf("[%s]: error decoding value as uint8, expected=%d, got=%d", tag, tag.Value().UInt8(i), uint8(expected))
				}
			}
		case SHORT:
			for i, value := range strings.Split(tc.value, " ") {
				expected, _ := strconv.ParseUint(value, 0, 16)
				if tag.Value().UInt16(i) != uint16(expected) {
					t.Errorf("[%s]: error decoding value as uint16, expected=%d, got=%d", tag, tag.Value().UInt16(i), uint16(expected))
				}
			}
		case LONG:
			for i, value := range strings.Split(tc.value, " ") {
				expected, _ := strconv.ParseUint(value, 0, 32)
				if tag.Value().UInt32(i) != uint32(expected) {
					t.Errorf("[%s]: error decoding value as uint32, expected=%d, got=%d", tag, tag.Value().UInt32(i), uint32(expected))
				}
			}
		case SBYTE:
			for i, value := range strings.Split(tc.value, " ") {
				expected, _ := strconv.ParseInt(value, 0, 8)
				if tag.Value().Int8(i) != int8(expected) {
					t.Errorf("[%s]: error decoding value as int8, expected=%d, got=%d", tag, tag.Value().Int8(i), int8(expected))
				}
			}
		case SSHORT:
			for i, value := range strings.Split(tc.value, " ") {
				expected, _ := strconv.ParseInt(value, 0, 16)
				if tag.Value().Int16(i) != int16(expected) {
					t.Errorf("[%s]: error decoding value as int16, expected=%d, got=%d", tag, tag.Value().Int16(i), int16(expected))
				}
			}
		case SLONG:
			for i, value := range strings.Split(tc.value, " ") {
				expected, _ := strconv.ParseInt(value, 0, 32)
				if tag.Value().Int32(i) != int32(expected) {
					t.Errorf("[%s]: error decoding value as int32, expected=%d, got=%d", tag, tag.Value().Int32(i), int32(expected))
				}
			}
		}
	}
}
