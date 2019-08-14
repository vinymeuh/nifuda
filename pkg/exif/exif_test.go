// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package exif

import (
	"os"
	"testing"
)

var tcExifFileError = []struct {
	filepath string
}{
	{"../../test/data/empty.txt"},
	{"../../test/data/dummy.txt"},
}

func TestExifFileError(t *testing.T) {
	for _, tc := range tcExifFileError {
		f, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening os file failed with err=%s", tc.filepath, err)
		}
		defer f.Close()

		ef, err := Read(f)
		if err == nil || ef != nil {
			t.Errorf("%s: opening should have failed and returned nil, err=%s, f=%v", tc.filepath, err, ef)
		}
	}
}

var tcExifFile = []struct {
	filepath string
	format   FileFormat
	tagIDs   []uint16
}{
	{"../../test/data/minimal.tif", TIFF, []uint16{257}},
	{"../../test/data/TEST_2019-07-21_132615_DSC_0361_DxO_PL2.jpg", JPEG, []uint16{}},
}

func TestExifFile(t *testing.T) {
	for _, tc := range tcExifFile {
		f, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening os file failed with err=%s", tc.filepath, err)
		}
		defer f.Close()

		ef, err := Read(f)
		if err != nil || ef == nil {
			t.Errorf("%s: opening fails, error=%s, f=%v", tc.filepath, err, ef)
			continue
		}

		if ef.format != tc.format {
			t.Errorf("%s: incorrect file format, expected=%s, actual=%s", tc.filepath, tc.format, ef.format)
		}
		// for i, expected := range tc.tagIDs {
		// 	actual := f.tiff.Tags[i].TagID
		// 	if actual != expected {
		// 		t.Errorf("%s: tag ID n°%d is incorrect, expected=%d, actual=%d", tc.filepath, i, expected, actual)
		// 	}
		// }
	}
}
