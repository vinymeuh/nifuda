// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

//
// Create tiff files mainly for error test cases
//
// Run with: go run testtiff.go
//
package main

import (
	"log"
	"os"
)

func ifh() []byte {
	return []byte{
		0x4d, 0x4d, // MM - big endian
		0x00, 0x2a, // 42
		0x00, 0x00, 0x00, 0x08, // offset0=8
	}
}

func ifd() []byte {
	return []byte{
		// number of entries (2 bits)
		0x00, 0x01, // 1
		// entry (12 bits)
		0x01, 0x01, // Tag ID = 257 ImageLength
		0x00, 0x03, // Data Type 3 = SHORT => Value size = 2 bits
		0x00, 0x00, 0x00, 0x01, // Data Count = 1
		0x00, 0x01, 0x00, 0x00, // Value = 1
		// next offset (4 bits)
		0x00, 0x00, 0x00, 0x00,
	}
}

func createFile(filepath string, data ...[]byte) {
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, d := range data {
		_, err = f.Write(d)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	createFile("data/minimal.tif", ifh(), ifd())

	ifd0 := ifd()
	ifd0[17] = 0x26
	ifd1 := ifd()
	createFile("data/minimal_with_ifd1.tif", ifh(), ifd0, ifd1)

	header := ifh()
	header[3] = 0x00
	createFile("data/wrong_version.tif", header)

	header = ifh()
	header[7] = 0x09
	createFile("data/wrong_offset0.tif", header, ifd())

	header = ifh()
	header[7] = 0x00
	createFile("data/no_ifd0.tif", header)

	ifd0 = ifd()
	ifd0[17] = 0x08
	createFile("data/recursive_ifd0.tif", ifh(), ifd0)

	ifd0 = ifd()
	ifd0[17] = 0x26
	createFile("data/wrong_ifd1.tif", ifh(), ifd0, ifd()[:4])

}
