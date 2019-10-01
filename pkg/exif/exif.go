// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package exif implements parsing of EXIF tags as defined in EXIF 2.31 specification.
package exif

import (
	"errors"
	"io"

	"github.com/vinymeuh/nifuda/pkg/jpeg"
	"github.com/vinymeuh/nifuda/pkg/tiff"
)

// Features of the Exif image file specification include the following:
//  * The file-recording format is based on existing formats:
//    * Compressed files are recorded as JPEG (ISO/IEC 10918-1) with application marker segments (APP1 and APP2) inserted.
//	  * Uncompressed files are recorded in TIFF rev 6.0 format.
//  * Related attribute information for both compressed and uncompressed files is stored in the tag information format defined in TIFF Rev. 6.0.
//  * Information specific to the camera system and not defined in TIFF is stored in private tags registered for Exif.
//  * Compressed files can record extended data exceeding 64 KBytes by dividing it into multiple APP2 segments.

// File represents a parsed EXIF file.
type File struct {
	jpeg   *jpeg.File
	tiff   *tiff.File // "real" TIFF file or the one embedded in JPEG APP1 segment
	format fileFormat
}

type fileFormat int

const (
	ffJPEG fileFormat = iota
	ffTIFF
	ffUNKNOWN
)

// Read parses EXIF data from an io.ReadSeeker.
// If successful, the returned file can be used to access EXIF tags.
func Read(rs io.ReadSeeker) (*File, error) {
	f := &File{}
	var err error

	format := identifyFileFormat(rs)
	switch format {
	case ffJPEG:
		f.format = ffJPEG
		f.jpeg, err = jpeg.Read(rs)
		return f, err
	case ffTIFF:
		f.format = ffTIFF
		f.tiff, err = tiff.Read(rs)
		return f, err
	default:
		return nil, errors.New("not an exif file")
	}
}

// Tags returns all tags indexed by tag namespace an tag name.
func (f *File) Tags() map[string]map[string]tiff.Tag {
	tags := make(map[string]map[string]tiff.Tag)
	tags["ifd0"] = make(map[string]tiff.Tag)
	tags["exif"] = make(map[string]tiff.Tag)
	tags["gps"] = make(map[string]tiff.Tag)

	switch f.format {
	case ffJPEG:
		for _, tag := range f.jpeg.Ifd0() {
			tags["ifd0"][tag.Name()] = tag
		}
		for _, tag := range f.jpeg.Exif() {
			tags["exif"][tag.Name()] = tag
		}
		for _, tag := range f.jpeg.Gps() {
			tags["gps"][tag.Name()] = tag
		}
	case ffTIFF:
		for _, tag := range f.tiff.Ifd0() {
			tags["ifd0"][tag.Name()] = tag
		}
		for _, tag := range f.tiff.Exif() {
			tags["exif"][tag.Name()] = tag
		}
		for _, tag := range f.tiff.Gps() {
			tags["gps"][tag.Name()] = tag
		}
	}

	return tags
}

// a utility function to identify the type of an image file.
func identifyFileFormat(rs io.ReadSeeker) fileFormat {
	var a [2]byte
	b := a[:]
	rs.Read(b)
	rs.Seek(0, io.SeekStart)

	switch string(b) {
	case "\xff\xd8": // SOI
		return ffJPEG
	case "II", "MM":
		return ffTIFF
	default:
		return ffUNKNOWN
	}
}
