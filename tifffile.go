// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// TIFF is an image file format built on three kind of structure:
//   * a unique Image File Header (IFH)
//   * Image File Directories (IFD), each containing information about the image as well as pointers its bitmap data
//   * Bitmap Data
//
// Each IFD and its associated bitmap are sometimes called a TIFF subfile.
// There is no limit to the number of subfiles a TIFF image file may contain.
//
// IFH contains pointer to the first IFD (IFD0).
// A valid TIFF file only require the IFH and IFD0.

// File represents a parsed TIFF file.
type tiffFile struct {
	rs      io.ReadSeeker
	bo      binary.ByteOrder // byte order used within the file
	version uint16           // always "42"
	offset0 uint32           // offset in bytes for IFD0, from the start of the file
}

// Parses TIFF data from an io.ReadSeeker.
func tiffRead(rs io.ReadSeeker) (*Exif, error) {
	x := &Exif{}
	f := &tiffFile{rs: rs}
	if err := f.readIFH(); err != nil {
		return nil, err
	}

	err := f.readIFD0(x)
	if err != nil { // && len(f.ifd0) == 0 { // failed to read ifd0
		return nil, err
	}

	return x, err
}

// readIFH reads the TIFF Header
func (f *tiffFile) readIFH() error {
	header := make([]byte, 8)
	if _, err := f.rs.Read(header); err != nil {
		return fmt.Errorf("failed to read 8 bytes: %w", err)
	}

	// retrieve byte order indication
	// Legal values are "II" (little-endian) and "MM" (big-endian)
	switch string(header[0:2]) {
	case "II": // Intel little-endian (0x4949)
		f.bo = binary.LittleEndian
	case "MM": // Motorola big-endian (0x4D4D)
		f.bo = binary.BigEndian
	default:
		return errors.New("invalid tiff byte order indication")
	}

	// validate tiff version
	binary.Read(bytes.NewReader(header[2:4]), f.bo, &f.version)
	if f.version != 42 {
		return errors.New("invalid tiff version")
	}

	// read offset for IFD0
	binary.Read(bytes.NewReader(header[4:8]), f.bo, &f.offset0)
	if f.offset0 < 8 { // ifd0 can not be located in the first 8 bytes used by IFH
		return errors.New("invalid offset for ifd0")
	}

	return nil
}

// readID0 reads the first IFD and decode it as Exif data
func (f *tiffFile) readIFD0(x *Exif) error {
	// IFD0
	ifd0, err := f.readIFD(f.offset0)
	if err != nil {
		return err
	}
	x.Ifd0 = parseIFDTagsAsExifTag(ifd0, f.bo)

	// Exif IFD
	if x.Ifd0.ExifIFD > 0 {
		exifIFD, err := f.readIFD(x.Ifd0.ExifIFD)
		if err != nil {
			return err
		}
		x.Exif = parseIFDTagsAsExifTag(exifIFD, f.bo)
	}

	// GPS IFD
	if x.Ifd0.GpsIFD > 0 {
		gpsIFD, err := f.readIFD(x.Ifd0.GpsIFD)
		if err != nil {
			return err
		}
		x.Gps = parseIFDTagsAsGpsTag(gpsIFD, f.bo)
	}

	return nil
}

// readIFD read the IFD starting at offset
func (f *tiffFile) readIFD(offset uint32) (*tiffIFD, error) {
	f.rs.Seek(int64(offset), io.SeekStart)
	ifd := tiffIFD{}

	// read the number of entries
	entries := make([]byte, 2)
	if _, err := f.rs.Read(entries); err != nil {
		return &ifd, fmt.Errorf("failed to read 2 bytes: %w", err)
	}
	binary.Read(bytes.NewReader(entries), f.bo, &ifd.entries)

	// read the data
	data := make([]byte, 12*ifd.entries)
	if _, err := f.rs.Read(data); err != nil {
		return &ifd, fmt.Errorf("failed to read %d bytes: %w", 12*ifd.entries, err)
	}

	// read offset for next IFD
	next := make([]byte, 4)
	if _, err := f.rs.Read(next); err != nil {
		return &ifd, fmt.Errorf("failed to read 4 bytes: %w", err)
	}
	binary.Read(bytes.NewReader(next), f.bo, &ifd.next)

	// parse raw tags (after offset because of possible nested Seek)
	ifd.tags = make([]ifdTag, ifd.entries)
	tag := ifdTag{}
	for i := 0; i < int(ifd.entries); i++ {

		binary.Read(bytes.NewReader(data[12*i:12*i+2]), f.bo, &tag.id)
		binary.Read(bytes.NewReader(data[12*i+2:12*i+4]), f.bo, &tag.tiffType)
		binary.Read(bytes.NewReader(data[12*i+4:12*i+8]), f.bo, &tag.count)

		length := tiffTypes[tag.tiffType].size * tag.count
		if length <= 4 {
			tag.data = data[12*i+8 : 12*i+8+int(length)]
		} else {
			var offset uint32
			binary.Read(bytes.NewReader(data[12*i+8:12*i+12]), f.bo, &offset)
			if _, err := f.rs.Seek(int64(offset), io.SeekStart); err != nil {
				//return ifd, fmt.Errorf("failed to seek of %d bytes: %w", offset, err)
				return &ifd, err
			}
			tag.data = make([]byte, length)
			if _, err := f.rs.Read(tag.data); err != nil {
				//return ifd, fmt.Errorf("failed to read field value: %w", err)
				return &ifd, err
			}
		}
		ifd.tags[i] = tag
	}

	return &ifd, nil
}
