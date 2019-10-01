// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package jpeg implements JPEG decoding as defined in JPEG File Interchange Format version 1.02.
package jpeg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/vinymeuh/nifuda/internal/tiff"
)

// File represents a parsed JPEG file.
type File struct {
	*tiff.File // embedded TIFF file for Exif tags
}

// Read parses JPEG from an io.ReadSeeker to retrieve Exif tags.
// Returns an error if no Exif data found.
func Read(rs io.ReadSeeker) (*File, error) {
	// ensure we have a SOI
	s0, err := nextSegment(rs)
	if err != nil {
		return nil, fmt.Errorf("unable to read SOI: %w", err)
	}
	if s0.marker != mSOI {
		return nil, errors.New("first segment must be SOI")
	}

	// next segments until we have found APP1 Exif
	f := File{}
	for {
		s, err := nextSegment(rs)
		if err != nil {
			return nil, err
		}

		if s.marker == mAPP1 && string(s.data[0:6]) == "Exif\x00\x00" {
			f.File, err = tiff.Read(bytes.NewReader(s.data[6:]))
			return &f, err
		}

		if s.marker == mEOI || s.marker == mSOS { // don't know how to process after SOS marker
			return nil, errors.New("no Exif data found")
		}
	}
}

// segment is the data structure implementation for JPEG Segment
type segment struct {
	marker byte   // marker begins with a 0xff byte followed by a byte indicating what kind of marker it is
	length uint16 // length of the data, includes the two bytes for the length but not the two bytes for the marker (big endian)
	data   []byte // payload data
}

// Restricted list to the used JPEG markers.
// For full list see https://www.disktuna.com/list-of-jpeg-markers/
const (
	mSOI  byte = 0xd8
	mEOI  byte = 0xd9
	mSOS  byte = 0xda
	mAPP1 byte = 0xe1
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
	s := &segment{marker: buf[1]}

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
