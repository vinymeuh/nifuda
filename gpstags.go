// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"encoding/binary"
)

type GpsTags struct {
	GPSVersionID        string
	GPSLatitudeRef      string
	GPSLongitudeRef     string
	GPSSpeedRef         string
	GPSTrackRef         string
	GPSImgDirectionRef  string
	GPSMapDatum         string
	GPSDestLatitudeRef  string
	GPSDestLongitudeRef string
	GPSDestBearingRef   string
	GPSDestDistanceRef  string
	GPSDateStamp        string
}

func parseIFDTagsAsGpsTag(ifd *tiffIFD, bo binary.ByteOrder) GpsTags {
	var gps GpsTags
	for _, ifdtag := range ifd.tags {
		switch ifdtag.id {
		case 0: // GPSVersionID
			gps.GPSVersionID = intArrayToString(ifdtag.ttByte(bo), ".")
		case 1: // GPSLatitudeRef
			switch ifdtag.ttAscii() {
			case "N":
				gps.GPSLatitudeRef = "North"
			case "S":
				gps.GPSLatitudeRef = "South"
			}
		case 2: // GPSLatitude
		case 3: // GPSLongitudeRef
			switch ifdtag.ttAscii() {
			case "E":
				gps.GPSLongitudeRef = "East"
			case "W":
				gps.GPSLongitudeRef = "West"
			}
		case 4: // GPSLongitude
		case 5: // GPSAltitudeRef
		case 6: // GPSAltitude
		case 7: // GPSTimeStamp
		case 8: // GPSSatellites
		case 9: // GPSStatus
		case 10: // GPSMeasureMode
		case 11: // GPSDOP
		case 12: // GPSSpeedRef
			gps.GPSSpeedRef = ifdtag.ttAscii()
		case 13: // GPSSpeed
		case 14: // GPSTrackRef
			switch ifdtag.ttAscii() {
			case "M":
				gps.GPSTrackRef = "Magnetic direction"
			case "T":
				gps.GPSTrackRef = "True direction"
			}
		case 15: // GPSTrack
		case 16: // GPSImgDirectionRef
			switch ifdtag.ttAscii() {
			case "M":
				gps.GPSImgDirectionRef = "Magnetic direction"
			case "T":
				gps.GPSImgDirectionRef = "True direction"
			}
		case 17: // GPSImgDirection
		case 18: // GPSMapDatum
			gps.GPSMapDatum = ifdtag.ttAscii()
		case 19: // GPSDestLatitudeRef
			switch ifdtag.ttAscii() {
			case "N":
				gps.GPSDestLatitudeRef = "North"
			case "S":
				gps.GPSDestLatitudeRef = "South"
			}
		case 20: // GPSDestLatitude
		case 21: // GPSDestLongitudeRef
			switch ifdtag.ttAscii() {
			case "E":
				gps.GPSDestLongitudeRef = "East"
			case "W":
				gps.GPSDestLongitudeRef = "West"
			}
		case 22: // GPSDestLongitude
		case 23: // GPSDestBearingRef
			switch ifdtag.ttAscii() {
			case "M":
				gps.GPSDestBearingRef = "Magnetic direction"
			case "T":
				gps.GPSDestBearingRef = "True direction"
			}
		case 24: // GPSDestBearing
		case 25: // GPSDestDistanceRef
			switch ifdtag.ttAscii() {
			case "K":
				gps.GPSDestDistanceRef = "Kilometers"
			case "M":
				gps.GPSDestDistanceRef = "Miles"
			case "N":
				gps.GPSDestDistanceRef = "Nautical miles"
			}
		case 26: // GPSDestDistance
		case 27: // GPSProcessingMethod
		case 28: // GPSAreaInformation
		case 29: // GPSDateStamp
			gps.GPSDateStamp = ifdtag.ttAscii()
		case 30: // GPSDifferential
		case 31: // GPSHPositioningError
		}
	}
	return gps
}
