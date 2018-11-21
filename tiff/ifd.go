// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package tiff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/vinymeuh/nifuda/exif"
)

/*
An Image File Directory (IFD) consists of a 2-byte count of the number of directory entries, followed by a
 sequence of 12-byte field entries, followed by a 4-byte offset of the next IFD (or 0 if none).

There must be at least 1 IFD in a TIFF file and each IFD must have at least one entry.

A TIFF field is a logical entity consisting of TIFF tag and its value. This logical concept is implemented as an IFD Entry,
plus the actual value if it doesn’t fit into the value/offset part, the last 4 bytes of the IFD Entry.
The terms TIFF field and IFD entry are interchangeable in most contexts.

To save time and space the Value Offset contains the Value instead of pointing to the Value if and only if the Value fits into 4 bytes.

Each TIFF field has an associated Count. This means that all fields are actually one-dimensional arrays, even though most fields contain only a single value.
*/

type ifd struct {
	offset  uint32 // offset in bytes of the IFD, from the start of the file
	entries uint16 // number of directory entries
	next    uint32 // offset in bytes to the next IFD, from the start of the file. 0 if none
	data    []byte // 12-byte IFD entry
}

type ifdEntry struct {
	id        uint16
	name      string
	fieldType uint16
	bytes     []byte
}

type ifdEntryTypeDefinition struct {
	name string
	size uint32
}

// See page 15 of TIFF Revision 6.0
var ifdEntryTypes = map[uint16]ifdEntryTypeDefinition{
	1:  {name: "Byte", size: 1},
	2:  {name: "Ascii", size: 1},
	3:  {name: "Short", size: 2},
	4:  {name: "Long", size: 4},
	5:  {name: "Rational", size: 8},
	6:  {name: "SByte", size: 1},
	7:  {name: "Undefined", size: 1},
	8:  {name: "SShort", size: 2},
	9:  {name: "SLong", size: 4},
	10: {name: "SRational", size: 8},
	11: {name: "Float", size: 4},
	12: {name: "Double", size: 8},
}

func parseIFD(namespace string, rs io.ReadSeeker, offset uint32, bom binary.ByteOrder) (*ifd, map[string]map[string]*tiffTag, error) {
	if _, err := rs.Seek(int64(offset), io.SeekStart); err != nil {
		return nil, nil, errors.Wrapf(err, "failed to seek of %d bytes", offset)
	}

	// read the IFD
	ifd := &ifd{offset: offset}

	bEntries := make([]byte, 2)
	if _, err := rs.Read(bEntries); err != nil {
		return ifd, nil, errors.Wrap(err, "failed to read 2 bytes")
	}
	if err := binary.Read(bytes.NewReader(bEntries), bom, &ifd.entries); err != nil {
		return ifd, nil, errors.Wrap(err, "failed to read number of entries")
	}

	ifd.data = make([]byte, 12*ifd.entries)
	if _, err := rs.Read(ifd.data); err != nil { // TODO - ReadAll !?
		return ifd, nil, errors.Wrapf(err, "failed to read %d bytes", 12*ifd.entries)
	}

	bNext := make([]byte, 4)
	if _, err := rs.Read(bNext); err != nil {
		return ifd, nil, errors.Wrap(err, "failed to read 4 bytes")
	}
	if err := binary.Read(bytes.NewReader(bNext), bom, &ifd.next); err != nil {
		return ifd, nil, errors.Wrap(err, "failed to read offset for next ifd")
	}

	// extract tags from the IFD
	T := map[string]map[string]*tiffTag{}
	T[namespace] = map[string]*tiffTag{}

	var fieldID uint16
	var fieldType uint16
	var fieldCount uint32
	var fieldValueOffset uint32

	for i := 0; i < int(ifd.entries); i++ {
		// can we ignore error because read from []byte ?
		binary.Read(bytes.NewReader(ifd.data[12*i:12*i+2]), bom, &fieldID)
		binary.Read(bytes.NewReader(ifd.data[12*i+2:12*i+4]), bom, &fieldType)
		binary.Read(bytes.NewReader(ifd.data[12*i+4:12*i+8]), bom, &fieldCount)

		tag := new(tiffTag)
		tag.name = exif.Dictionary[fieldID].Name
		if tag.name == "" {
			tag.name = fmt.Sprintf("%d", fieldID)
		}

		// Whether the Value fits within 4 bytes is determined by the Type and Count of the field
		fieldLength := ifdEntryTypes[fieldType].size * fieldCount
		if fieldLength <= 4 {
			// If the Value is shorter than 4 bytes, it is left-justified within the 4-byte Value Offset,
			tag.bytes = ifd.data[12*i+8 : 12*i+8+int(fieldLength)]
		} else {
			if err := binary.Read(bytes.NewReader(ifd.data[12*i+8:12*i+12]), bom, &fieldValueOffset); err != nil {
				return ifd, T, errors.Wrap(err, "failed to read offset for field value")
			}
			tag.bytes = make([]byte, fieldLength)
			if _, err := rs.Seek(int64(fieldValueOffset), io.SeekStart); err != nil {
				return ifd, T, errors.Wrapf(err, "failed to seek of %d bytes", int64(fieldValueOffset))
			}
			if _, err := rs.Read(tag.bytes); err != nil {
				return ifd, T, errors.Wrap(err, "failed to read field value")
			}
		}
		tag.convertTiffBytesValue(fieldType, fieldCount, bom) // TODO
		T[namespace][tag.Name()] = tag
	}

	return ifd, T, nil
}
