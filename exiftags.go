// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

// Package exif implements parsing of EXIF tags as defined in EXIF 2.31 specification.
package nifuda

// type ExifTags struct {
// 	ImageWidth  int
// 	ImageLenght int
// }

type ExifTags map[string]Tag
