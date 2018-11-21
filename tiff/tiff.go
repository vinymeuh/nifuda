// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package tiff implements TIFF decoding as defined in TIFF revision 6.0 specification
package tiff

/*
TIFF is an image file format.

A TIFF file begins with an 8-byte image file header that points to an image file directory (IFD).
An image file directory contains information about the image, as well as pointers to the actual image data.
*/

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/vinymeuh/nifuda/tag"
)

// Tiff implements nifuda.Image for TIFF format.
type Tiff struct {
	order    binary.ByteOrder               // byte order used within the file. Legal values are "II" (little-endian) and "MM" (big-endian)
	version  uint16                         // always "42"
	offset0  uint32                         // offset in bytes of the first IFD, from the start of the file
	ifds     []*ifd                         // list of all IFDs
	tiffTags map[string]map[string]*tiffTag // all Tags from file, indexed by their namespace then their name
}

func (t *Tiff) TagsNamespaces() []string {
	keys := make([]string, 0)
	for key := range t.tiffTags {
		keys = append(keys, key)
	}
	return keys
}

func (t *Tiff) GetTagsFromNamespace(namespace string) []tag.Tag {
	tags := make([]tag.Tag, 0)
	for _, tag := range t.tiffTags[namespace] {
		tags = append(tags, tag)
	}
	return tags
}

func (t *Tiff) GetTag(namespace string, name string) tag.Tag {
	if tag, ok := t.tiffTags[namespace][name]; ok {
		return tag
	}
	return nil
}

// Parse is the tiff's factory.
func Parse(rs io.ReadSeeker) (*Tiff, error) {
	t := new(Tiff)
	t.ifds = make([]*ifd, 1, 8)
	t.tiffTags = map[string]map[string]*tiffTag{}

	// read tiff header
	rs.Seek(0, io.SeekStart)
	header := make([]byte, 8)
	if _, err := rs.Read(header); err != nil {
		return nil, errors.Wrap(err, "failed to read 8 bytes")
	}

	// retrieve byte order indication
	switch string(header[0:2]) {
	case "II": // Intel little-endian (0x4949)
		t.order = binary.LittleEndian
	case "MM": // Motorola big-endian (0x4D4D)
		t.order = binary.BigEndian
	default:
		return nil, errors.New("invalid tiff byte order indication")
	}

	// validate tiff version
	if err := binary.Read(bytes.NewReader(header[2:4]), t.order, &t.version); err != nil {
		return t, errors.Wrap(err, "failed to read tiff version")
	}
	if t.version != 42 {
		return t, errors.New("invalid tiff version")
	}

	// read IFD0
	if err := binary.Read(bytes.NewReader(header[4:8]), t.order, &t.offset0); err != nil {
		return t, errors.Wrap(err, "failed to read offset for ifd0")
	}
	ifd, ifdTags, err := parseIFD("ifd0", rs, t.offset0, t.order)
	if err != nil {
		return t, errors.Wrap(err, "failed to read ifd0")
	}
	t.ifds[0] = ifd
	t.tiffTags["ifd0"] = map[string]*tiffTag{}
	t.tiffTags["ifd0"] = ifdTags["ifd0"]

	// Private Sub-IFD
	if tag := t.tiffTags["ifd0"]["Exif IFD"]; tag != nil {
		exifIFD, exifTags, err := parseIFD("exif", rs, uint32(tag.intValues[0]), t.order)
		if err != nil {
			return t, errors.Wrap(err, "failed to read Exif IFD")
		}
		t.ifds = append(t.ifds, exifIFD)
		t.tiffTags["exif"] = map[string]*tiffTag{}
		t.tiffTags["exif"] = exifTags["exif"]
	}

	// next IFDs
	i := 0
	for {
		switch t.ifds[i].next {
		case 0:
			return t, nil
		case t.ifds[i].offset:
			return t, errors.New("recursive ifd")
		default: // read next
			i++
			ifd, ifdTags, err := parseIFD(fmt.Sprintf("ifd%d", i), rs, t.ifds[i-1].next, t.order)
			if err != nil {
				return t, errors.Wrapf(err, "failed to read ifd%d", i)
			}
			t.ifds = append(t.ifds, ifd)
			t.tiffTags[fmt.Sprintf("ifd%d", i)] = map[string]*tiffTag{}
			t.tiffTags[fmt.Sprintf("ifd%d", i)] = ifdTags[fmt.Sprintf("ifd%d", i)]
		}
	}
}

// PrintStructure is a function used for test & debugging purposes
// Should not be used by applications.
func (t *Tiff) PrintStructure(out io.Writer) {
	switch t.order {
	case binary.LittleEndian:
		fmt.Fprintln(out, "BYTE ORDER INDICATION: \"II\" (little-endian)")
	case binary.BigEndian:
		fmt.Fprintln(out, "BYTE ORDER INDICATION: \"MM\" (big-endian)")
	}

	fmt.Fprintf(out, " %7s | %-3s | %-7s | %-8s |\n", "address", "ifd", "entries", "next")
	for n, ifd := range t.ifds {
		fmt.Fprintf(out, " %7d | %3d | %7d | %8d |\n", ifd.offset, n, ifd.entries, ifd.next)
	}
}
