// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package nifuda implements parsing of EXIF tags as defined in EXIF 2.31 specification.
package nifuda

import (
	"errors"
	"io"
)

type Exif struct {
	Ifd0 ExifTags
	Exif ExifTags
	//	Gps  GpsTags
}

// Read parses EXIF data from an io.ReadSeeker.
// If successful, the returned file can be used to access EXIF tags.
func Read(rs io.ReadSeeker) (*Exif, error) {

	var a [2]byte
	b := a[:]
	rs.Read(b)
	rs.Seek(0, io.SeekStart)

	switch string(b) {
	case "\xff\xd8": // SOI
		x, err := jpegRead(rs)
		return x, err
	case "II", "MM":
		x, err := tiffRead(rs)
		return x, err
	default:
		return nil, errors.New("not an exif file")
	}
}
