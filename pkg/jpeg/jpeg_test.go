// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.
package jpeg

import (
	"os"
	"testing"
)

// Tests cases for invalid or corrupted files
// We expect to receive an error, not to panic
var tcJpegFileError = []struct {
	filepath string
}{
	{"../../test/data/empty.txt"},
	{"../../test/data/dummy.txt"},
	{"../../test/data/nosoi.jpg"},
}

func TestJpegFileError(t *testing.T) {
	for _, tc := range tcJpegFileError {
		f, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening os file failed with err=%s", tc.filepath, err)
		}
		defer f.Close()

		jf, err := Read(f)
		if err == nil || jf != nil {
			t.Errorf("%s: opening should have failed and returned nil, err=%s, f=%v", tc.filepath, err, f)
		}
	}
}

var tcJpegFile = []struct {
	filepath string
	hasExif  bool
}{
	{"../../test/data/minimal.jpg", false},
	{"../../test/data/TEST_2019-07-21_132615_DSC_0361_DxO_PL2.jpg", true},
}

func TestJpegFile(t *testing.T) {
	for _, tc := range tcJpegFile {
		f, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening os file failed with err=%s", tc.filepath, err)
		}
		defer f.Close()

		jf, err := Read(f)
		if err != nil {
			t.Errorf("%s: opening fails, error=%s", tc.filepath, err)
			continue
		}

		if tc.hasExif == false && jf.ExifSubTIFF() != nil {
			t.Errorf("%s: invalid ExifSubTIFF detected", tc.filepath)
		}
		if tc.hasExif == true && jf.ExifSubTIFF() == nil {
			t.Errorf("%s: should have detected ExifSubTIFF", tc.filepath)
		}

	}
}
