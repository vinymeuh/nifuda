// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package tiff

// See page 15 of TIFF Revision 6.0
var dataTypes = map[uint16]struct {
	name string
	size uint32
}{
	1:  {name: "BYTE", size: 1},
	2:  {name: "ASCII", size: 1},
	3:  {name: "SHORT", size: 2},
	4:  {name: "LONG", size: 4},
	5:  {name: "RATIONAL", size: 8},
	6:  {name: "SBYTE", size: 1},
	7:  {name: "UNDEFINED", size: 1},
	8:  {name: "SSHORT", size: 2},
	9:  {name: "SLONG", size: 4},
	10: {name: "SRATIONAL", size: 8},
	11: {name: "FLOAT", size: 4},
	12: {name: "DOUBLE", size: 8},
}

/*
A Tag is a 12-byte record containing TagID identifying the type of information the tag contains and its value.

Each Tag has an associated Count. This means that all fields are actually one-dimensional arrays,
even though most fields contain only a single value.

To save time and space the Offset contains directly the Value instead of pointing to the Value if and only if the value fits into 4 bytes.
*/

type Tag struct {
	TagID      uint16 // Tag identifier
	DataType   uint16 // The scalar type of the data items
	DataCount  uint32 // The number of items in the tag data
	DataOffset uint32 // The byte offset to the data items (4-byte)
	DataValue  []byte // The value data
}

type TagDictionary map[uint16]struct {
	Name string
}

func (t *Tag) Name(d TagDictionary) string {
	return d[t.TagID].Name
}
