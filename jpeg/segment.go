// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.
package jpeg

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

type Segment struct {
	offset uint32
	marker SegmentMarker
	length uint16
	data   []byte
}

type SegmentMarker string

var jifSegmentMarkerName = map[SegmentMarker]string{
	SOI:  "SOI",
	APP0: "APP0",
	APP1: "APP1",
	APP2: "APP2",
	DQT:  "DQT",
	SOF0: "SOF0",
	DHT:  "DHT",
	SOS:  "SOS",
	EOI:  "EOI",
}

const (
	SOI  SegmentMarker = "\xff\xd8"
	APP0 SegmentMarker = "\xff\xe0"
	APP1 SegmentMarker = "\xff\xe1"
	APP2 SegmentMarker = "\xff\xe2"
	DQT  SegmentMarker = "\xff\xdb"
	SOF0 SegmentMarker = "\xff\xc0"
	DHT  SegmentMarker = "\xff\xc4"
	SOS  SegmentMarker = "\xff\xda"
	EOI  SegmentMarker = "\xff\xd9"
)

func parseSegment(rs io.ReadSeeker, offset uint32) (*Segment, error) {
	if _, err := rs.Seek(int64(offset), io.SeekStart); err != nil {
		return nil, errors.Wrapf(err, "failed to seek of %d bytes", offset)
	}

	segment := new(Segment)
	segment.offset = offset

	buffer := make([]byte, 2)

	// read segment
	if _, err := rs.Read(buffer); err != nil {
		return nil, errors.Wrap(err, "failed to read segment marker")
	}
	segment.marker = SegmentMarker(buffer)

	switch segment.marker {
	case SOI, EOI: // return from here for standalone markers
		return segment, nil
	}

	if _, err := rs.Read(buffer); err != nil {
		return segment, errors.Wrap(err, "failed to read 2 bytes")
	}
	err := binary.Read(bytes.NewReader(buffer), binary.BigEndian, &segment.length)
	if err != nil {
		return segment, errors.Wrap(err, "failed to read length of segment's data")
	}

	segment.data = make([]byte, segment.length-2)
	read := 0
	for read < int(segment.length)-2 {
		n, err := rs.Read(segment.data[read:])
		if err != nil {
			return segment, errors.Wrap(err, "failed to read segment's data")
		}
		read += n
	}

	return segment, nil
}
