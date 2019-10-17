// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"encoding/binary"
	"fmt"
)

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

func parseIFDTagsAsGpsTag(ifd *tiffIFD, bo binary.ByteOrder) GpsTags {
	var gps GpsTags
	for _, ifdtag := range ifd.tags {
		switch ifdtag.id {
		case 0: // GPSVersionID
			gps.GPSVersionID = intArrayToString(ifdtag.byteToInt(bo), ".")
		case 1: // GPSLatitudeRef
			switch ifdtag.asciiToString() {
			case "N":
				gps.GPSLatitudeRef = "North"
			case "S":
				gps.GPSLatitudeRef = "South"
			}
		// case 2: // GPSLatitude
		// 	r := ifdtag.rationalToFloat32(bo)
		// 	gps.GPSLatitude = fmt.Sprintf("%2.0f %f' %f\"", r[0], r[1], r[2])
		case 3: // GPSLongitudeRef
			switch ifdtag.asciiToString() {
			case "E":
				gps.GPSLongitudeRef = "East"
			case "W":
				gps.GPSLongitudeRef = "West"
			}
		// case 4: // GPSLongitude
		// 	r := ifdtag.rationalToFloat32(bo)
		// 	gps.GPSLongitude = fmt.Sprintf("%2.0f %f' %f\"", r[0], r[1], r[2])
		case 5: // GPSAltitudeRef
			switch ifdtag.byteToInt(bo)[0] {
			case 0:
				gps.GPSAltitudeRef = "Sea level"
			case 1:
				gps.GPSAltitudeRef = "Sea level reference (negative value)"
			}
		case 6: // GPSAltitude
			gps.GPSAltitude = ifdtag.rationalToFloat32(bo)[0]
		case 7: // GPSTimeStamp
			r := ifdtag.rationalToFloat32(bo)
			gps.GPSTimeStamp = fmt.Sprintf("%02.0f:%02.0f:%02.0fZ", r[0], r[1], r[2])
		case 8: // GPSSatellites
			gps.GPSSatellites = ifdtag.asciiToString()
		case 9: // GPSStatus
			switch ifdtag.asciiToString() {
			case "A":
				gps.GPSStatus = "Measurement in progress"
			case "V":
				gps.GPSStatus = "Measurement interrupted"
			}
		case 10: // GPSMeasureMode
			switch ifdtag.asciiToString() {
			case "2":
				gps.GPSMeasureMode = "2-dimensional measurement"
			case "3":
				gps.GPSMeasureMode = "3-dimensional measurement"
			}
		case 11: // GPSDOP
			gps.GPSDOP = ifdtag.rationalToFloat32(bo)[0]
		case 12: // GPSSpeedRef
			gps.GPSSpeedRef = ifdtag.asciiToString()
		case 13: // GPSSpeed
			gps.GPSSpeed = ifdtag.rationalToFloat32(bo)[0]
		case 14: // GPSTrackRef
			switch ifdtag.asciiToString() {
			case "M":
				gps.GPSTrackRef = "Magnetic direction"
			case "T":
				gps.GPSTrackRef = "True direction"
			}
		case 15: // GPSTrack
			gps.GPSTrack = ifdtag.rationalToFloat32(bo)[0]
		case 16: // GPSImgDirectionRef
			switch ifdtag.asciiToString() {
			case "M":
				gps.GPSImgDirectionRef = "Magnetic direction"
			case "T":
				gps.GPSImgDirectionRef = "True direction"
			}
		case 17: // GPSImgDirection
			gps.GPSImgDirection = ifdtag.rationalToFloat32(bo)[0]
		case 18: // GPSMapDatum
			gps.GPSMapDatum = ifdtag.asciiToString()
		case 19: // GPSDestLatitudeRef
			switch ifdtag.asciiToString() {
			case "N":
				gps.GPSDestLatitudeRef = "North"
			case "S":
				gps.GPSDestLatitudeRef = "South"
			}
		case 20: // GPSDestLatitude
			r := ifdtag.rationalToFloat32(bo)
			gps.GPSDestLatitude = fmt.Sprintf("%2.0f %f' %f\"", r[0], r[1], r[2])
		case 21: // GPSDestLongitudeRef
			switch ifdtag.asciiToString() {
			case "E":
				gps.GPSDestLongitudeRef = "East"
			case "W":
				gps.GPSDestLongitudeRef = "West"
			}
		case 22: // GPSDestLongitude
			r := ifdtag.rationalToFloat32(bo)
			gps.GPSDestLongitude = fmt.Sprintf("%2.0f %f' %f\"", r[0], r[1], r[2])
		case 23: // GPSDestBearingRef
			switch ifdtag.asciiToString() {
			case "M":
				gps.GPSDestBearingRef = "Magnetic direction"
			case "T":
				gps.GPSDestBearingRef = "True direction"
			}
		case 24: // GPSDestBearing
			gps.GPSDestBearing = ifdtag.rationalToFloat32(bo)[0]
		case 25: // GPSDestDistanceRef
			switch ifdtag.asciiToString() {
			case "K":
				gps.GPSDestDistanceRef = "Kilometers"
			case "M":
				gps.GPSDestDistanceRef = "Miles"
			case "N":
				gps.GPSDestDistanceRef = "Nautical miles"
			}
		case 26: // GPSDestDistance
			gps.GPSDestDistance = ifdtag.rationalToFloat32(bo)[0]
		case 27: // GPSProcessingMethod
		case 28: // GPSAreaInformation
		case 29: // GPSDateStamp
			gps.GPSDateStamp = ifdtag.asciiToString()
		case 30: // GPSDifferential
			gps.GPSDifferential = ifdtag.shortToUint16(bo)[0]
		case 31: // GPSHPositioningError
			gps.GPSHPositioningError = ifdtag.rationalToFloat32(bo)[0]
		}
	}
	return gps
}
