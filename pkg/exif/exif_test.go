// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package exif

import (
	"os"
	"testing"
)

func TestExifFileError(t *testing.T) {
	tests := []struct {
		filepath string
	}{
		{"../../test/data/empty.txt"},
		{"../../test/data/dummy.txt"},
	}

	for _, tc := range tests {
		osf, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening file failed with err=%s", tc.filepath, err)
		}
		defer osf.Close()

		f, err := Read(osf)
		if err == nil || f != nil {
			t.Errorf("%s: reading file should have failed and returned nil, err=%s, f=%v", tc.filepath, err, f)
		}
	}
}

func TestExifFile(t *testing.T) {
	tests := []struct {
		filepath string
		format   FileFormat
		tagIDs   []uint16
	}{
		{"../../test/data/minimal.tif", TIFF, []uint16{257}},
		{"../../test/data/TEST_2019-07-21_132615_DSC_0361_DxO_PL2.jpg", JPEG, []uint16{}},
	}

	for _, tc := range tests {
		f, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening file failed with err=%s", tc.filepath, err)
		}
		defer f.Close()

		ef, err := Read(f)
		if err != nil || ef == nil {
			t.Errorf("%s: reading fails, error=%s, f=%v", tc.filepath, err, ef)
			continue
		}

		if ef.format != tc.format {
			t.Errorf("%s: incorrect file format, expected=%v, got=%v", tc.filepath, tc.format, ef.format)
		}

		// for i, expected := range tc.tagIDs {
		// 	actual := f.tiff.Tags[i].TagID
		// 	if actual != expected {
		// 		t.Errorf("%s: tag ID n°%d is incorrect, expected=%d, actual=%d", tc.filepath, i, expected, actual)
		// 	}
		// }
	}
}
