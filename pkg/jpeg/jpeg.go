// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package jpeg implements JPEG decoding as defined in JPEG File Interchange Format version 1.02.
package jpeg

import (
	"errors"
	"fmt"
	"io"
)

// File represents a parsed JPEG file.
type File struct {
	rs          io.ReadSeeker // the JPEG data stream
	exifSubTIFF []byte        // embedded TIFF file for Exif tags
}

// Read parses JPEG data from an io.ReadSeeker.
// Note that parsing is restricted to APP1 Exif segment, all others segments are currently ignored.
func Read(rs io.ReadSeeker) (*File, error) {
	f := &File{rs: rs}
	err := f.readSegments()
	if err != nil {
		return nil, err
	}
	return f, err
}

// ExifSubTIFF returns the embedded Exif TIFF file.
func (f *File) ExifSubTIFF() []byte {
	return f.exifSubTIFF
}

func (f *File) readSegments() error {
	// ensure we have a SOI
	s0, err := nextSegment(f.rs)
	if err != nil {
		return fmt.Errorf("unable to read SOI: %w", err)
	}
	if s0.marker != mSOI {
		return errors.New("first segment must be SOI")
	}

	// next segments until we have found APP1 Exif
	for {
		s, err := nextSegment(f.rs)
		if err != nil {
			return err
		}

		if s.marker == mAPP1 && string(s.data[0:6]) == "Exif\x00\x00" {
			f.exifSubTIFF = s.data[6:]
			break
		}

		if s.marker == mEOI || s.marker == mSOS { // don't know how to process after SOS marker
			break
		}
	}
	return nil
}
