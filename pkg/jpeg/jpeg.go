// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.
package jpeg

import (
	"errors"
	"fmt"
	"io"
)

// A File consists of a sequence of Segments
type File struct {
	rs       io.ReadSeeker
	Segments []*Segment // list of all Segments
}

func Read(rs io.ReadSeeker) (*File, error) {
	f := &File{rs: rs}
	err := f.readSegments()
	if err != nil && len(f.Segments) == 0 {
		return nil, err
	}
	return f, err
}

func (f *File) readSegments() error {
	f.Segments = make([]*Segment, 0)

	// ensure we have a SOI
	s0, err := nextSegment(f.rs)
	if err != nil {
		return fmt.Errorf("unable to read SOI: %w", err)
	}
	if s0.Marker != SOI {
		return errors.New("first segment must be SOI")
	}
	f.Segments = append(f.Segments, s0)

	// next segments
	for {
		s, err := nextSegment(f.rs)
		if err != nil {
			return err
		}
		f.Segments = append(f.Segments, s)

		if s.Marker == EOI || s.Marker == SOS { // FIXME: don't know how to process after SOS marker
			break
		}
	}
	return nil
}
