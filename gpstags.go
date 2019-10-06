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
	0: {Name: "GPSVersionID"},
	1: {Name: "GPSLatitudeRef"},
	2: {Name: "GPSLatitude"},
	// To be completed ...
	31: {Name: "GPSHPositioningError"},
}
