// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package jpeg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// segment is the data structure implementation for JPEG Segment
type segment struct {
	// marker begins with a 0xff byte followed by a byte indicating what kind of marker it is
	marker marker
	// length of the data, includes the two bytes for the length but not the two bytes for the marker (big endian)
	length uint16
	// payload data
	data []byte
}

// marker is the type for JPEG segment markers
type marker byte

// List of JPEG markers (from https://www.disktuna.com/list-of-jpeg-markers/)
const (
	mSOF0  marker = 0xc0
	mSOF1  marker = 0xc1
	mSOF2  marker = 0xc2
	mSOF3  marker = 0xc3
	mDHT   marker = 0xc4
	mSOF5  marker = 0xc5
	mSOF6  marker = 0xc6
	mSOF7  marker = 0xc7
	mJPG   marker = 0xc8
	mSOF9  marker = 0xc9
	mSOF10 marker = 0xca
	mSOF11 marker = 0xcb
	mDAC   marker = 0xcc
	mSOF13 marker = 0xcd
	mSOF14 marker = 0xce
	mSOF15 marker = 0xcf
	mRST0  marker = 0xd0
	mRST1  marker = 0xd1
	mRST2  marker = 0xd2
	mRST3  marker = 0xd3
	mRST4  marker = 0xd4
	mRST5  marker = 0xd5
	mRST6  marker = 0xd6
	mRST7  marker = 0xd7
	mSOI   marker = 0xd8
	mEOI   marker = 0xd9
	mSOS   marker = 0xda
	mDQT   marker = 0xdb
	mDNL   marker = 0xdc
	mDRI   marker = 0xdd
	mDHP   marker = 0xde
	mEXP   marker = 0xdf
	mAPP0  marker = 0xe0
	mAPP1  marker = 0xe1
	mAPP2  marker = 0xe2
	mAPP3  marker = 0xe3
	mAPP4  marker = 0xe4
	mAPP5  marker = 0xe5
	mAPP6  marker = 0xe6
	mAPP7  marker = 0xe7
	mAPP8  marker = 0xe8
	mAPP9  marker = 0xe9
	mAPP10 marker = 0xea
	mAPP11 marker = 0xeb
	mAPP12 marker = 0xec
	mAPP13 marker = 0xed
	mAPP14 marker = 0xee
	mAPP15 marker = 0xef
	mJPG0  marker = 0xf0
	mJPG1  marker = 0xf1
	mJPG2  marker = 0xf2
	mJPG3  marker = 0xf3
	mJPG4  marker = 0xf4
	mJPG5  marker = 0xf5
	mJPG6  marker = 0xf6
	mJPG7  marker = 0xf7
	mJPG8  marker = 0xf8
	mJPG9  marker = 0xf9
	mJPG10 marker = 0xfa
	mJPG11 marker = 0xfb
	mJPG12 marker = 0xfc
	mJPG13 marker = 0xfd
	mCOM   marker = 0xfe
)

func nextSegment(r io.Reader) (*segment, error) {
	buf := make([]byte, 2) // used to read marker & length

	// marker
	if _, err := r.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to read segment marker: %w", err)
	}
	if buf[0] != 0xff {
		return nil, errors.New("invalid segment marker")
	}
	s := &segment{marker: marker(buf[1])}

	if s.marker == mSOI || s.marker == mEOI {
		return s, nil
	}

	// length
	if _, err := r.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to read segment length: %w", err)
	}
	binary.Read(bytes.NewReader(buf), binary.BigEndian, &s.length)

	// data
	s.data = make([]byte, s.length-2)
	read := 0
	for read < int(s.length)-2 {
		n, err := r.Read(s.data[read:])
		if err != nil {
			return nil, fmt.Errorf("failed to read segment data: %w", err)
		}
		read += n
	}

	return s, nil
}
