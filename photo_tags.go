// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import "encoding/binary"

// PhotoTags contains tags from Exif SubIFD.
// Fields are defined in order they appeared in chapter 4.6.5 of Exif 2.31
type PhotoTags struct {
	// A. Tags Relating to Version
	ExifVersion     string
	FlashpixVersion string
	// B. Tag Relating to Image Data Characteristics
	// C. Tags Relating to Image Configuration
	// D. Tags Relating to User Information
	// E. Tag Relating to Related File Information
	// F. Tags Relating to Date and Time
	DateTimeOriginal    string
	DateTimeDigitized   string
	OffsetTime          string
	OffsetTimeOriginal  string
	OffsetTimeDigitized string
	SubSecTime          string
	SubSecTimeOriginal  string
	SubSecTimeDigitized string
	// G. Tags Relating to Picture-Taking Conditions
	ExposureProgram     string
	SpectralSensitivity string
	MeteringMode        string
}

func parseIFDTagsAsPhotoTags(ifd *ifd, bo binary.ByteOrder) PhotoTags {
	var t PhotoTags

	for _, ifdtag := range ifd.tags {
		switch ifdtag.id {
		case 36864: // ExifVersion
			t.ExifVersion = ifdtag.undefinedToString()
		case 40960: // FlashpixVersion
			t.FlashpixVersion = ifdtag.undefinedToString()
		case 36867: // DateTimeOriginal
			t.DateTimeOriginal = ifdtag.asciiToString()
		case 36868: // DateTimeDigitized
			t.DateTimeDigitized = ifdtag.asciiToString()
		case 36880: // OffsetTime
			t.OffsetTime = ifdtag.asciiToString()
		case 36881: // OffsetTimeOriginal
			t.OffsetTimeOriginal = ifdtag.asciiToString()
		case 36882: // OffsetTimeDigitized
			t.OffsetTimeDigitized = ifdtag.asciiToString()
		case 37520: // SubsecTime
			t.SubSecTime = ifdtag.asciiToString()
		case 37521: // SubsecTimeOriginal
			t.SubSecTimeOriginal = ifdtag.asciiToString()
		case 37522: // SubSecTimeDigitized
			t.SubSecTimeDigitized = ifdtag.asciiToString()
		case 34850: // ExposureProgram
			switch ifdtag.shortToUint16(bo)[0] {
			case 0:
				t.ExposureProgram = "not defined"
			case 1:
				t.ExposureProgram = "manual"
			case 2:
				t.ExposureProgram = "normal program"
			case 3:
				t.ExposureProgram = "aperture priority"
			case 4:
				t.ExposureProgram = "shutter priority"
			case 5:
				t.ExposureProgram = "creative program"
			case 6:
				t.ExposureProgram = "action program"
			case 7:
				t.ExposureProgram = "portrait mode"
			case 8:
				t.ExposureProgram = "landscape mode"
			}
		case 34852: // SpectralSensitivity
			t.SpectralSensitivity = ifdtag.asciiToString()
		case 34855: // PhotographicSensitivity (ISO 12232)
		case 34856: // OECF (ISO 14524)
		case 34864: // SensitivityType (ISO 12232)
		case 34865: // StandardOutputSensitivity (ISO 12232)
		case 37383: // MeteringMode
			switch ifdtag.shortToUint16(bo)[0] {
			case 0:
				t.MeteringMode = "unknown"
			case 1:
				t.MeteringMode = "average"
			case 2:
				t.MeteringMode = "center-weighted average"
			case 3:
				t.MeteringMode = "spot"
			case 4:
				t.MeteringMode = "multispot"
			case 5:
				t.MeteringMode = "pattern"
			case 6:
				t.MeteringMode = "partial"
			case 255:
				t.MeteringMode = "other"
			}
		}
	}

	return t
}
