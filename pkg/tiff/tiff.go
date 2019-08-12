// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package tiff implements TIFF decoding as defined in TIFF revision 6.0 specification
package tiff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

/*
TIFF is an image file format built on three kind of structure:
  * a unique Image File Header (IFH)
  * Image File Directories (IFD), each containing information about the image as well as pointers its bitmap data
  * Bitmap Data

Each IFD and its associated bitmap are sometimes called a TIFF subfile.
There is no limit to the number of subfiles a TIFF image file may contain.

IFH contains pointer to the first IFD (IFD0).

A valid TIFF file only require the IFH and IFD0.
*/

// File
type File struct {
	rs io.ReadSeeker
	// Image File Header
	bo      binary.ByteOrder // byte order used within the file
	version uint16           // always "42"
	offset0 uint32           // offset in bytes for IFD0, from the start of the file
	Tags    [][]Tag
}

func Read(rs io.ReadSeeker) (*File, error) {
	f := &File{rs: rs}
	if err := f.readIFH(); err != nil {
		return nil, err
	}
	err := f.readIFDs()
	if err != nil && len(f.Tags) == 0 { // failed to read ifd0
		return nil, err
	}
	return f, err
}

// ReadIFH read the TIFF Header
func (f *File) readIFH() error {
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

// ReadIFDs read all IFDs starting from IFD0 at f.offset0
func (f *File) readIFDs() error {

	ifds := make([]*ifd, 0)

	previous := uint32(0)
	next := f.offset0
	i := 0
	for {
		if next == 0 { // break the infinite loop
			break
		}
		switch next {
		case previous:
			return errors.New("recursive ifd")
		default:
			ifd, err := f.ReadIFD(next)
			if err != nil {
				return fmt.Errorf("failed to read ifd%d: %w", i, err)
			}
			ifds = append(ifds, ifd)
			previous = next
			next = ifd.next
			i++
		}
	}

	// Copy tags to File structure
	f.Tags = make([][]Tag, len(ifds))
	for i, ifd := range ifds {
		f.Tags[i] = ifd.Tags
	}
	return nil
}

// ReadIFD read the IFD that start at offset
func (f *File) ReadIFD(offset uint32) (*ifd, error) {
	f.rs.Seek(int64(offset), io.SeekStart)
	ifd := &ifd{}

	// read the number of entries
	entries := make([]byte, 2)
	if _, err := f.rs.Read(entries); err != nil {
		return ifd, fmt.Errorf("failed to read 2 bytes: %w", err)
	}
	binary.Read(bytes.NewReader(entries), f.bo, &ifd.entries)

	// read the array of Tags
	data := make([]byte, 12*ifd.entries)
	if _, err := f.rs.Read(data); err != nil {
		return ifd, fmt.Errorf("failed to read %d bytes: %w", 12*ifd.entries, err)
	}
	for i := 0; i < int(ifd.entries); i++ {
		tag := Tag{}

		binary.Read(bytes.NewReader(data[12*i:12*i+2]), f.bo, &tag.TagID)
		binary.Read(bytes.NewReader(data[12*i+2:12*i+4]), f.bo, &tag.DataType)
		binary.Read(bytes.NewReader(data[12*i+4:12*i+8]), f.bo, &tag.DataCount)

		length := dataTypes[tag.DataType].size * tag.DataCount
		if length <= 4 {
			tag.DataValue = data[12*i+8 : 12*i+8+int(length)]
		} else {
			binary.Read(bytes.NewReader(data[12*i+8:12*i+12]), f.bo, &tag.DataOffset)
		}
		ifd.Tags = append(ifd.Tags, tag)
	}

	// read offset for next IFD
	next := make([]byte, 4)
	if _, err := f.rs.Read(next); err != nil {
		return ifd, fmt.Errorf("failed to read 4 bytes: %w", err)
	}
	binary.Read(bytes.NewReader(next), f.bo, &ifd.next)

	return ifd, nil
}
