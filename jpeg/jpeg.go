// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.
package jpeg

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math"

	"github.com/pkg/errors"

	"github.com/vinymeuh/nifuda/tag"

	"github.com/vinymeuh/nifuda/exif"
	"github.com/vinymeuh/nifuda/tiff"
)

// JPEG Interchange Format segment marker
// See https://en.wikipedia.org/wiki/JPEG#JPEG_files

type Jpeg struct {
	segments []*Segment
	tiffTags map[string]map[string]tag.Tag
}

func (j *Jpeg) TagsNamespaces() []string {
	keys := make([]string, 0)
	for key := range j.tiffTags {
		keys = append(keys, key)
	}
	return keys
}

func (j *Jpeg) GetTagsFromNamespace(namespace string) []tag.Tag {
	tags := make([]tag.Tag, 0)
	for _, tag := range j.tiffTags[namespace] {
		tags = append(tags, tag)
	}
	return tags
}

func (j *Jpeg) GetTag(namespace string, name string) tag.Tag {
	return j.tiffTags[namespace][name]
}

// Parse is the jpeg's factory.
func Parse(rs io.ReadSeeker) (*Jpeg, error) {

	// checks that stream starts with a SOI segment
	segment, err := parseSegment(rs, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read SOI segment")
	}
	if segment.marker != SOI {
		return nil, errors.New("SOI marker not found")
	}

	// initialize jpeg structure with first segment
	j := &Jpeg{
		segments: make([]*Segment, 0, 16),
		tiffTags: map[string]map[string]tag.Tag{},
	}
	j.segments = append(j.segments, segment)

	// reads each successive segments
	previous := 0
	for {
		segment, err := parseSegment(rs, j.segments[previous].offset+uint32(j.segments[previous].length)+2)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read segment %d", previous+1)
		}
		j.segments = append(j.segments, segment)

		// does this segment has interesting tags ?
		switch segment.marker {
		case APP1:
			if string(segment.data[0:6]) == exif.IdentifierCode {
				subTiff, err := tiff.Parse(bytes.NewReader(segment.data[6:]))
				if err != nil {
					return j, errors.Wrap(err, "failed to read embedded Exif segment")
				}
				for _, namespace := range subTiff.TagsNamespaces() {
					j.tiffTags[namespace] = map[string]tag.Tag{}
					for _, tag := range subTiff.GetTagsFromNamespace(namespace) {
						j.tiffTags[namespace][tag.Name()] = tag
					}
				}

			}
		case SOS: // don't know how to process after SOS marker
			return j, nil
		case EOI: // End Of Image, the "normal" exit
			return j, nil
		}
		previous++ // next segment
	}
}

// PrintStructure is a function used for test & debugging purposes
// Should not be used by applications.
func (j *Jpeg) PrintStructure(wout io.Writer) {
	var markerHex, markerName string
	var data []byte

	fmt.Fprintf(wout, " %7s | %-11s | %6s | %s\n", "address", "marker", "length", "data")
	for _, segment := range j.segments {
		markerHex = hex.EncodeToString([]byte(string(segment.marker)))
		markerName = jifSegmentMarkerName[SegmentMarker(segment.marker)]
		switch segment.marker {
		case SOI, SOS:
			fmt.Fprintf(wout, " %7d | 0x%4s %-4s\n", segment.offset, markerHex, markerName)
		case APP0, APP1:
			// Print only 32 firsts characters
			// Control characters are replaced by dots
			data = segment.data[0:int(math.Min(float64(segment.length-2), 32))]
			dataDumped := make([]byte, len(data))
			for _, b := range data {
				if b < 32 || b > 126 {
					dataDumped = append(dataDumped, '.')
				} else {
					dataDumped = append(dataDumped, b)
				}
			}
			fmt.Fprintf(wout, " %7d | 0x%4s %-4s | %6d | %s\n", segment.offset, markerHex, markerName, segment.length, dataDumped)
		default:
			fmt.Fprintf(wout, " %7d | 0x%4s %-4s | %6d |\n", segment.offset, markerHex, markerName, segment.length)
		}
	}
}
