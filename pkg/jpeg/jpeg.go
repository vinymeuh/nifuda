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
	ExifTIFF *Segment // embedded TIFF file for Exif tags
}

func Read(rs io.ReadSeeker) (*File, error) {
	f := &File{rs: rs}
	err := f.readSegments()
	if err != nil {
		return nil, err
	}
	return f, err
}

func (f *File) readSegments() error {
	// ensure we have a SOI
	s0, err := nextSegment(f.rs)
	if err != nil {
		return fmt.Errorf("unable to read SOI: %w", err)
	}
	if s0.Marker != SOI {
		return errors.New("first segment must be SOI")
	}

	// next segments until we have found APP1 Exif
	for {
		s, err := nextSegment(f.rs)
		if err != nil {
			return err
		}

		if s.Marker == APP1 && string(s.Data[0:6]) == ExifIdentifier {
			f.ExifTIFF = s
			break
		}

		if s.Marker == EOI || s.Marker == SOS { // don't know how to process after SOS marker
			break
		}
	}
	return nil
}
