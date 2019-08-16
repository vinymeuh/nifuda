// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package tiff

type TagDictionary map[uint16]struct {
	Name string
}

// BaselineTags contains Baseline TIFF tags definitions in order they appeared in TIFF 6.0 Specification
var BaselineTags = TagDictionary{
	262: {Name: "PhotometricInterpretation"},
	259: {Name: "Compression"},
	257: {Name: "ImageLength"},
	256: {Name: "ImageWidth"},
	// to be continued, but not sure it is very usefull ...
}
