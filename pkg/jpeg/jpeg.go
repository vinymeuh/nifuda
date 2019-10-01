// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package jpeg implements JPEG decoding as defined in JPEG File Interchange Format version 1.02.
package jpeg

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/vinymeuh/nifuda/pkg/tiff"
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
