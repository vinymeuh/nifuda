// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

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
