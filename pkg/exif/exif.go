// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package exif

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/vinymeuh/nifuda/pkg/jpeg"
	"github.com/vinymeuh/nifuda/pkg/tiff"
)

/*
Features of the Exif image file specification include the following:
  * The file-recording format is based on existing formats.
    * Compressed files are recorded as JPEG (ISO/IEC 10918-1) with application marker segments (APP1 and APP2) inserted.
	* Uncompressed files are recorded in TIFF rev 6.0 format.
  * Related attribute information for both compressed and uncompressed files is stored in the tag information format defined in TIFF Rev. 6.0.
  * Information specific to the camera system and not defined in TIFF is stored in private tags registered for Exif.
  * Compressed files can record extended data exceeding 64 KBytes by dividing it into multiple APP2 segments.
*/

type File struct {
	jpeg     *jpeg.File
	tiff     *tiff.File // "real" TIFF file or the one embedded in JPEG APP1 segment
	format   FileFormat
	exifTags []tiff.Tag
	gpsTags  []tiff.Tag
}

type FileFormat int

const (
	JPEG FileFormat = iota
	TIFF
	UNKNOWN
)

func (ff FileFormat) String() string {
	switch ff {
	case JPEG:
		return "jpeg"
	case TIFF:
		return "tiff"
	default:
		return "unknown"
	}
}

func Read(rs io.ReadSeeker) (*File, error) {
	f := &File{}

	format := identifyFileFormat(rs)
	switch format {
	case JPEG:
		jpegFile, err := jpeg.Read(rs)
		if err != nil {
			return nil, err
		}
		f.format = JPEG
		f.jpeg = jpegFile
		if f.jpeg.ExifTIFF == nil {
			return nil, errors.New("no Exif APP1 marker found")
		}
		f.tiff, err = tiff.Read(bytes.NewReader(f.jpeg.ExifTIFF.Data[6:]), ExifDictionary)
		if err != nil {
			return nil, err
		}
	case TIFF:
		tiffFile, err := tiff.Read(rs, ExifDictionary)
		if err != nil {
			return nil, err
		}
		f.format = TIFF
		f.tiff = tiffFile
	default:
		return nil, errors.New("not an exif file")
	}

	err := f.parseExif()
	return f, err
}

// a utility function to identify the type of an image file.
func identifyFileFormat(rs io.ReadSeeker) FileFormat {
	b := make([]byte, 2)
	rs.Read(b)
	rs.Seek(0, io.SeekStart)

	switch string(b[0:2]) {
	case "\xff\xd8": // SOI
		return JPEG
	case "II", "MM":
		return TIFF
	default:
		return UNKNOWN
	}
}

func (f *File) parseExif() error {
	for _, tag := range f.tiff.Tags[0] {
		switch tag.TagID {
		case 34665: // Exif IFD
			offset, _ := tag.UInt32(0)
			ifd, err := f.tiff.ReadIFD(offset, ExifDictionary)
			if err != nil {
				return fmt.Errorf("failed to read Exif IFD: %w", err)
			}
			f.exifTags = ifd.Tags
		case 34853: // GPS IFD
			offset, _ := tag.UInt32(0)
			ifd, err := f.tiff.ReadIFD(offset, GPSDictionary)
			if err != nil {
				return fmt.Errorf("failed to read GPS IFD: %w", err)
			}
			f.gpsTags = ifd.Tags
		}

	}
	return nil
}

func (f *File) Tags() map[string][]tiff.Tag {
	t := make(map[string][]tiff.Tag)
	t["ifd0"] = f.tiff.Tags[0]
	t["exif"] = f.exifTags
	t["gps"] = f.gpsTags
	return t
}
