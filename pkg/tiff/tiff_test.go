// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.
package tiff

import (
	"os"
	"testing"
)

// Tests cases for invalid or corrupted files
// We expect to receive an error, not to panic
var tcTiffFileError = []struct {
	filepath string
}{
	{"../../test/data/empty.txt"},
	{"../../test/data/dummy.txt"},
	{"../../test/data/wrong_version.tif"},
	{"../../test/data/wrong_offset0.tif"},
	{"../../test/data/no_ifd0.tif"},
	{"../../test/data/recursive_ifd0.tif"},
	{"../../test/data/wrong_ifd1.tif"},
}

func TestTiffFileError(t *testing.T) {
	for _, tc := range tcTiffFileError {
		f, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening os file failed with err=%s", tc.filepath, err)
		}
		defer f.Close()

		tf, err := Read(f, nil)
		if err == nil || tf != nil {
			t.Errorf("%s: opening tiff file should have failed and returned nil, err=%s, f=%v", tc.filepath, err, tf)
		}
	}
}

// Tests cases for valid TIFF files
var tcTiffFile = []struct {
	filepath string
}{
	{"../../test/data/minimal.tif"},
	{"../../test/data/minimal_with_ifd1.tif"},
}

func TestTiffFile(t *testing.T) {
	for _, tc := range tcTiffFile {
		f, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening os file failed with err=%s", tc.filepath, err)
		}
		defer f.Close()

		tf, err := Read(f, nil)
		if err != nil || tf == nil {
			t.Errorf("%s: opening fails, error=%s, f=%v", tc.filepath, err, tf)
			continue
		}

	}
}
