// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

type GpsTags struct {
}

// gpsTags contains GPS tags definitions
var dictGps = tagDictionary{
	/*****************************/
	/* GPS Attribute Information */
	/*****************************/
	// A. Tags Relating to GPS
	0:  {Name: "GPSVersionID"},
	1:  {Name: "GPSLatitudeRef"},
	2:  {Name: "GPSLatitude"},
	3:  {Name: "GPSLongitudeRef"},
	4:  {Name: "GPSLongitude"},
	5:  {Name: "GPSAltitudeRef"},
	6:  {Name: "GPSAltitude"},
	7:  {Name: "GPSTimeStamp"},
	8:  {Name: "GPSSatellites"},
	9:  {Name: "GPSStatus"},
	10: {Name: "GPSMeasureMode"},
	11: {Name: "GPSDOP"},
	12: {Name: "GPSSpeedRef"},
	13: {Name: "GPSSpeed"},
	14: {Name: "GPSTrackRef"},
	15: {Name: "GPSTrack"},
	16: {Name: "GPSTrack"},
	17: {Name: "GPSImgDirection"},
	18: {Name: "GPSMapDatum"},
	19: {Name: "GPSDestLatitudeRef"},
	20: {Name: "GPSDestLatitude"},
	21: {Name: "GPSDestLongitudeRef"},
	22: {Name: "GPSDestLongitude"},
	23: {Name: "GPSDestBearingRef"},
	24: {Name: "GPSDestBearing"},
	25: {Name: "GPSDestDistanceRef"},
	26: {Name: "GPSDestDistance"},
	27: {Name: "GPSProcessingMethod"},
	28: {Name: "GPSAreaInformation"},
	29: {Name: "GPSDateStamp"},
	30: {Name: "GPSDifferential"},
	31: {Name: "GPSHPositioningError"},
}
