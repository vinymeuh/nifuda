// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

//
// Create jpeg files mainly for error test cases
//
// Run with: go run testjpeg.go
//
package main

import (
	"log"
	"os"
)

func soi() []byte {
	return []byte{
		0xff, 0xd8,
	}
}

func eoi() []byte {
	return []byte{
		0xff, 0xd9,
	}
}

func app0() []byte {
	return []byte{
		0xff, 0xe0,
		0x00, 0x10, // segment length
		0x4a, 0x46, 0x49, 0x46, 0x00, // JFIF
		0x01, 0x02, // version 1.2
		0x01,       // unit
		0x00, 0x60, // x-density
		0x00, 0x60, // y-density
		0x00, 0x00, // thumbnail width x height
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
	createFile("data/nosoi.jpg", eoi())

	createFile("data/minimal.jpg", soi(), app0(), eoi())
}
