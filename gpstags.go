// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package exif implements parsing of EXIF tags as defined in EXIF 2.31 specification.
package nifuda

type GpsTags struct {
}

// gpsTags contains GPS tags definitions
var dictGps = TagDictionary{
	/*****************************/
	/* GPS Attribute Information */
	/*****************************/
	// A. Tags Relating to GPS
	0: {Name: "GPSVersionID"},
	1: {Name: "GPSLatitudeRef"},
	2: {Name: "GPSLatitude"},
	// To be completed ...
	31: {Name: "GPSHPositioningError"},
}
