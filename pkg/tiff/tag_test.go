// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.
package tiff

import (
	"encoding/binary"
	"testing"
)

func TestTagDecode(t *testing.T) {
	tests := []struct {
		dtType   dataType
		dtCount  uint32
		dtRaw    []byte
		expected string
	}{
		{BYTE, 1, []byte{0x7b}, "123"},
		{ASCII, 6, []byte{'h', 'e', 'l', 'l', 'o', '0'}, "hello"},
		{SHORT, 2, []byte{0x4, 0x0, 0x0, 0xa}, "1024 10"},
		{LONG, 1, []byte{0x4, 0x0, 0x0, 0x0}, "67108864"},
		{RATIONAL, 1, []byte{0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x4}, "1/4"},
	}

	for _, tc := range tests {
		tag := Tag{dataType: tc.dtType, dataCount: tc.dtCount, dataRaw: tc.dtRaw}
		tag.decode(nil, binary.BigEndian)
		if tag.StringValue() != tc.expected {
			t.Errorf("error decoding tag [%s], expected value=%s", tag, tc.expected)
		}
	}
}
