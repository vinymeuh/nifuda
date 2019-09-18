// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package exif implements parsing of EXIF tags as defined in EXIF 2.31 specification.
package exif

import (
	"bytes"
	"errors"
	"fmt"
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
	jpeg     *jpeg.File
	tiff     *tiff.File // "real" TIFF file or the one embedded in JPEG APP1 segment
	format   fileFormat
	exifTags []tiff.Tag
	gpsTags  []tiff.Tag
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

	format := identifyFileFormat(rs)
	switch format {
	case ffJPEG:
		jpegFile, err := jpeg.Read(rs)
		if err != nil {
			return nil, err
		}
		f.format = ffJPEG
		f.jpeg = jpegFile
		subTIFF := f.jpeg.ExifSubTIFF()
		if subTIFF == nil {
			return nil, errors.New("no Exif APP1 marker found")
		}
		f.tiff, err = tiff.Read(bytes.NewReader(subTIFF), exifTags)
		if err != nil {
			return nil, err
		}
	case ffTIFF:
		tiffFile, err := tiff.Read(rs, exifTags)
		if err != nil {
			return nil, err
		}
		f.format = ffTIFF
		f.tiff = tiffFile
	default:
		return nil, errors.New("not an exif file")
	}

	err := f.parseExif()
	return f, err
}

// Tags returns all tags indexed by tag namespace an tag name.
func (f *File) Tags() map[string]map[string]tiff.Tag {
	tags := make(map[string]map[string]tiff.Tag)

	tags["ifd0"] = make(map[string]tiff.Tag)
	for _, tag := range f.tiff.Tags[0] {
		tags["ifd0"][tag.Name()] = tag
	}

	tags["exif"] = make(map[string]tiff.Tag)
	for _, tag := range f.exifTags {
		tags["exif"][tag.Name()] = tag
	}

	tags["gps"] = make(map[string]tiff.Tag)
	for _, tag := range f.gpsTags {
		tags["gps"][tag.Name()] = tag
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

// a utility function to read EXIF and GPS IFD after hava read IFD0
func (f *File) parseExif() error {
	for _, tag := range f.tiff.Tags[0] {
		switch tag.ID() {
		case tagExifIfd:
			ifd, err := f.tiff.ReadIFD(tag.Value().UInt32(0), exifTags)
			if err != nil {
				return fmt.Errorf("failed to read Exif IFD: %w", err)
			}
			f.exifTags = ifd.Tags
		case tagGpsIfd:
			ifd, err := f.tiff.ReadIFD(tag.Value().UInt32(0), gpsTags)
			if err != nil {
				return fmt.Errorf("failed to read GPS IFD: %w", err)
			}
			f.gpsTags = ifd.Tags
		}
	}
	return nil
}
