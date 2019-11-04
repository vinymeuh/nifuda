// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"encoding/binary"
	"fmt"
)

// GpsTags contains tags from GPS SubIFD.
// Fields are defined in order they appeared in chapter 4.6.6 of Exif 2.31
type GpsTags struct {
	GPSVersionID         string
	GPSLatitudeRef       string
	GPSLatitude          string
	GPSLongitudeRef      string
	GPSLongitude         string
	GPSAltitudeRef       string
	GPSAltitude          float32
	GPSTimeStamp         string
	GPSSatellites        string
	GPSStatus            string
	GPSDOP               float32
	GPSMeasureMode       string
	GPSSpeedRef          string
	GPSSpeed             float32
	GPSTrackRef          string
	GPSTrack             float32
	GPSImgDirectionRef   string
	GPSImgDirection      float32
	GPSMapDatum          string
	GPSDestLatitudeRef   string
	GPSDestLatitude      string
	GPSDestLongitudeRef  string
	GPSDestLongitude     string
	GPSDestBearingRef    string
	GPSDestBearing       float32
	GPSDestDistanceRef   string
	GPSDestDistance      float32
	GPSDateStamp         string
	GPSDifferential      uint16
	GPSHPositioningError float32
}

func parseIFDTagsAsGpsTags(ifd *ifd, bo binary.ByteOrder) GpsTags {
	var t GpsTags

	for _, ifdtag := range ifd.tags {
		switch ifdtag.id {
		case 0: // GPSVersionID
			t.GPSVersionID = intArrayToString(ifdtag.byteToInt(bo), ".")
		case 1: // GPSLatitudeRef
			switch ifdtag.asciiToString() {
			case "N":
				t.GPSLatitudeRef = "North"
			case "S":
				t.GPSLatitudeRef = "South"
			}
		// case 2: // GPSLatitude
		// 	r := ifdtag.rationalToFloat32(bo)
		// 	gps.GPSLatitude = fmt.Sprintf("%2.0f %f' %f\"", r[0], r[1], r[2])
		case 3: // GPSLongitudeRef
			switch ifdtag.asciiToString() {
			case "E":
				t.GPSLongitudeRef = "East"
			case "W":
				t.GPSLongitudeRef = "West"
			}
		// case 4: // GPSLongitude
		// 	r := ifdtag.rationalToFloat32(bo)
		// 	gps.GPSLongitude = fmt.Sprintf("%2.0f %f' %f\"", r[0], r[1], r[2])
		case 5: // GPSAltitudeRef
			switch ifdtag.byteToInt(bo)[0] {
			case 0:
				t.GPSAltitudeRef = "Sea level"
			case 1:
				t.GPSAltitudeRef = "Sea level reference (negative value)"
			}
		case 6: // GPSAltitude
			t.GPSAltitude = ifdtag.rationalToFloat32(bo)[0]
		case 7: // GPSTimeStamp
			r := ifdtag.rationalToFloat32(bo)
			t.GPSTimeStamp = fmt.Sprintf("%02.0f:%02.0f:%02.0fZ", r[0], r[1], r[2])
		case 8: // GPSSatellites
			t.GPSSatellites = ifdtag.asciiToString()
		case 9: // GPSStatus
			switch ifdtag.asciiToString() {
			case "A":
				t.GPSStatus = "Measurement in progress"
			case "V":
				t.GPSStatus = "Measurement interrupted"
			}
		case 10: // GPSMeasureMode
			switch ifdtag.asciiToString() {
			case "2":
				t.GPSMeasureMode = "2-dimensional measurement"
			case "3":
				t.GPSMeasureMode = "3-dimensional measurement"
			}
		case 11: // GPSDOP
			t.GPSDOP = ifdtag.rationalToFloat32(bo)[0]
		case 12: // GPSSpeedRef
			t.GPSSpeedRef = ifdtag.asciiToString()
		case 13: // GPSSpeed
			t.GPSSpeed = ifdtag.rationalToFloat32(bo)[0]
		case 14: // GPSTrackRef
			switch ifdtag.asciiToString() {
			case "M":
				t.GPSTrackRef = "Magnetic direction"
			case "T":
				t.GPSTrackRef = "True direction"
			}
		case 15: // GPSTrack
			t.GPSTrack = ifdtag.rationalToFloat32(bo)[0]
		case 16: // GPSImgDirectionRef
			switch ifdtag.asciiToString() {
			case "M":
				t.GPSImgDirectionRef = "Magnetic direction"
			case "T":
				t.GPSImgDirectionRef = "True direction"
			}
		case 17: // GPSImgDirection
			t.GPSImgDirection = ifdtag.rationalToFloat32(bo)[0]
		case 18: // GPSMapDatum
			t.GPSMapDatum = ifdtag.asciiToString()
		case 19: // GPSDestLatitudeRef
			switch ifdtag.asciiToString() {
			case "N":
				t.GPSDestLatitudeRef = "North"
			case "S":
				t.GPSDestLatitudeRef = "South"
			}
		case 20: // GPSDestLatitude
			r := ifdtag.rationalToFloat32(bo)
			t.GPSDestLatitude = fmt.Sprintf("%2.0f %f' %f\"", r[0], r[1], r[2])
		case 21: // GPSDestLongitudeRef
			switch ifdtag.asciiToString() {
			case "E":
				t.GPSDestLongitudeRef = "East"
			case "W":
				t.GPSDestLongitudeRef = "West"
			}
		case 22: // GPSDestLongitude
			r := ifdtag.rationalToFloat32(bo)
			t.GPSDestLongitude = fmt.Sprintf("%2.0f %f' %f\"", r[0], r[1], r[2])
		case 23: // GPSDestBearingRef
			switch ifdtag.asciiToString() {
			case "M":
				t.GPSDestBearingRef = "Magnetic direction"
			case "T":
				t.GPSDestBearingRef = "True direction"
			}
		case 24: // GPSDestBearing
			t.GPSDestBearing = ifdtag.rationalToFloat32(bo)[0]
		case 25: // GPSDestDistanceRef
			switch ifdtag.asciiToString() {
			case "K":
				t.GPSDestDistanceRef = "Kilometers"
			case "M":
				t.GPSDestDistanceRef = "Miles"
			case "N":
				t.GPSDestDistanceRef = "Nautical miles"
			}
		case 26: // GPSDestDistance
			t.GPSDestDistance = ifdtag.rationalToFloat32(bo)[0]
		case 27: // GPSProcessingMethod
		case 28: // GPSAreaInformation
		case 29: // GPSDateStamp
			t.GPSDateStamp = ifdtag.asciiToString()
		case 30: // GPSDifferential
			t.GPSDifferential = ifdtag.shortToUint16(bo)[0]
		case 31: // GPSHPositioningError
			t.GPSHPositioningError = ifdtag.rationalToFloat32(bo)[0]
		}
	}
	return t
}
