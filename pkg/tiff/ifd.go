// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package tiff

/*
An Image File Directory (IFD) consists of a 2-byte count of the number of directory entries, followed by a
 sequence of 12-byte field entries, followed by a 4-byte offset of the next IFD (or 0 if none).

There must be at least 1 IFD in a TIFF file and each IFD must have at least one entry.

Each TIFF field has an associated Count. This means that all fields are actually one-dimensional arrays, even though most fields contain only a single value.
*/

type ifd struct {
	entries uint16 // number of directory entries
	Tags    []Tag  // array of 12-byte Tags entries
	next    uint32 // offset in bytes to the next IFD, from the start of the file. 0 if none
}
