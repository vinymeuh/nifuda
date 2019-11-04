// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import "encoding/binary"

// ImageTags contains tags from first IFD (IFD0).
// Fields are defined in order they appeared in chapter 4.6.4 of Exif 2.31
type ImageTags struct {
	// A. Tags relating to image data structure
	ExifIFD             uint32
	GpsIFD              uint32
	InteroperabilityIFD uint32
	//ImageWidth SHORT or LONG :(
	//ImageLength SHORT or LONG :(
	BitsPerSample             uint16
	Compression               string
	PhotometricInterpretation string
	SamplesPerPixel           uint16
	PlanarConfiguration       string
	YCbCrPositioning          string
	ResolutionUnit            string
	// B. Tags relating to recording offset
	// C. Tags relating to image data characteristics
	// D. Other tags
	DateTime         string
	ImageDescription string
	Make             string
	Model            string
	Software         string
	Artist           string
	Copyright        string
}

func parseIFDTagsAsImageTags(ifd *ifd, bo binary.ByteOrder) ImageTags {
	var t ImageTags

	for _, ifdtag := range ifd.tags {
		switch ifdtag.id {
		case 34665: // Exif IFD
			t.ExifIFD = ifdtag.longToUint32(bo)[0]
		case 34853: // GPS IFD
			t.GpsIFD = ifdtag.longToUint32(bo)[0]
		case 258: // BitsPerSample
			t.BitsPerSample = ifdtag.shortToUint16(bo)[0]
		case 259: // Compression
			switch ifdtag.shortToUint16(bo)[0] {
			case 1:
				t.Compression = "uncompressed"
			case 6:
				t.Compression = "JPEG compression"
			}
		case 262: // PhotometricInterpretation
			switch ifdtag.shortToUint16(bo)[0] {
			case 2:
				t.PhotometricInterpretation = "RGB"
			case 6:
				t.PhotometricInterpretation = "YCbCr"
			}
		case 274: // Orientation
		case 277: // SamplesPerPixel
			t.SamplesPerPixel = ifdtag.shortToUint16(bo)[0]
		case 284: // PlanarConfiguration
			switch ifdtag.shortToUint16(bo)[0] {
			case 1:
				t.PlanarConfiguration = "chunky format"
			case 2:
				t.PlanarConfiguration = "planar format"
			}
		case 531: // YCbCrPositioning
			switch ifdtag.shortToUint16(bo)[0] {
			case 1:
				t.YCbCrPositioning = "centered"
			case 2:
				t.YCbCrPositioning = "co-sited"
			}
		case 296: // ResolutionUnit
			switch ifdtag.shortToUint16(bo)[0] {
			case 2:
				t.ResolutionUnit = "inches"
			case 3:
				t.ResolutionUnit = "centimeters"
			}
		case 306: // DateTime
			t.DateTime = ifdtag.asciiToString()
		case 270: // ImageDescription
			t.ImageDescription = ifdtag.asciiToString()
		case 271: // Make
			t.Make = ifdtag.asciiToString()
		case 272: // Model
			t.Model = ifdtag.asciiToString()
		case 305: // Software
			t.Software = ifdtag.asciiToString()
		case 315: // Artist
			t.Artist = ifdtag.asciiToString()
		case 33432: // Copyright
			t.Copyright = ifdtag.asciiToString()
		}
	}

	return t
}
