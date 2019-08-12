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

// A Segment
type Segment struct {
	// marker begins with a 0xff byte followed by a byte indicating what kind of marker it is
	Marker Marker
	// length of the data, includes the two bytes for the length but not the two bytes for the marker (big endian)
	length uint16
	// payload data
	Data []byte
}

// Marker is the type for JPEG segment markers
type Marker byte

// List of JPEG markers (from https://www.disktuna.com/list-of-jpeg-markers/)
const (
	SOF0  Marker = 0xc0
	SOF1  Marker = 0xc1
	SOF2  Marker = 0xc2
	SOF3  Marker = 0xc3
	DHT   Marker = 0xc4
	SOF5  Marker = 0xc5
	SOF6  Marker = 0xc6
	SOF7  Marker = 0xc7
	JPG   Marker = 0xc8
	SOF9  Marker = 0xc9
	SOF10 Marker = 0xca
	SOF11 Marker = 0xcb
	DAC   Marker = 0xcc
	SOF13 Marker = 0xcd
	SOF14 Marker = 0xce
	SOF15 Marker = 0xcf
	RST0  Marker = 0xd0
	RST1  Marker = 0xd1
	RST2  Marker = 0xd2
	RST3  Marker = 0xd3
	RST4  Marker = 0xd4
	RST5  Marker = 0xd5
	RST6  Marker = 0xd6
	RST7  Marker = 0xd7
	SOI   Marker = 0xd8
	EOI   Marker = 0xd9
	SOS   Marker = 0xda
	DQT   Marker = 0xdb
	DNL   Marker = 0xdc
	DRI   Marker = 0xdd
	DHP   Marker = 0xde
	EXP   Marker = 0xdf
	APP0  Marker = 0xe0
	APP1  Marker = 0xe1
	APP2  Marker = 0xe2
	APP3  Marker = 0xe3
	APP4  Marker = 0xe4
	APP5  Marker = 0xe5
	APP6  Marker = 0xe6
	APP7  Marker = 0xe7
	APP8  Marker = 0xe8
	APP9  Marker = 0xe9
	APP10 Marker = 0xea
	APP11 Marker = 0xeb
	APP12 Marker = 0xec
	APP13 Marker = 0xed
	APP14 Marker = 0xee
	APP15 Marker = 0xef
	JPG0  Marker = 0xf0
	JPG1  Marker = 0xf1
	JPG2  Marker = 0xf2
	JPG3  Marker = 0xf3
	JPG4  Marker = 0xf4
	JPG5  Marker = 0xf5
	JPG6  Marker = 0xf6
	JPG7  Marker = 0xf7
	JPG8  Marker = 0xf8
	JPG9  Marker = 0xf9
	JPG10 Marker = 0xfa
	JPG11 Marker = 0xfb
	JPG12 Marker = 0xfc
	JPG13 Marker = 0xfd
	COM   Marker = 0xfe
)

func (m Marker) String() string {
	switch m {
	case SOF0:
		return fmt.Sprintf("SOF0 (0xff%x)", uint16(m))
		//FIXME: to be continued, so boring :(
	}
	return fmt.Sprintf("UNKNOWN (0xff%x)", uint16(m))
}

func nextSegment(r io.Reader) (*Segment, error) {
	buf := make([]byte, 2) // used to read marker & length

	// marker
	if _, err := r.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to read segment marker: %w", err)
	}
	if buf[0] != 0xff {
		return nil, errors.New("invalid segment marker")
	}
	s := &Segment{Marker: Marker(buf[1])}

	if s.Marker == SOI || s.Marker == EOI {
		return s, nil
	}

	// length
	if _, err := r.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to read segment length: %w", err)
	}
	binary.Read(bytes.NewReader(buf), binary.BigEndian, &s.length)

	// data
	s.Data = make([]byte, s.length-2)
	read := 0
	for read < int(s.length)-2 {
		n, err := r.Read(s.Data[read:])
		if err != nil {
			return nil, fmt.Errorf("failed to read segment data: %w", err)
		}
		read += n
	}

	return s, nil
}
