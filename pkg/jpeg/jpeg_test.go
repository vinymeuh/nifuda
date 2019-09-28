// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.
package jpeg

import (
	"os"
	"testing"
)

func TestJpegFileMustReturnError(t *testing.T) {
	tests := []struct {
		filepath string
	}{
		{"../../test/data/empty.txt"},
		{"../../test/data/dummy.txt"},
		{"../../test/data/nosoi.jpg"},
		{"../../test/data/minimal.jpg"},
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

func TestJpegFile(t *testing.T) {
	tests := []struct {
		filepath string
		hasExif  bool
	}{
		{"../../test/data/TEST_2019-07-21_132615_DSC_0361_DxO_PL2.jpg", true},
	}

	for _, tc := range tests {
		osf, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening file failed with err=%s", tc.filepath, err)
		}
		defer osf.Close()

		f, err := Read(osf)
		if err != nil || f == nil {
			t.Errorf("%s: reading fails, err=%s, f=%v", tc.filepath, err, f)
			continue
		}

		if tc.hasExif == false && f.ExifSubTIFF() != nil {
			t.Errorf("%s: invalid ExifSubTIFF detected", tc.filepath)
		}
		if tc.hasExif == true && f.ExifSubTIFF() == nil {
			t.Errorf("%s: should have detected ExifSubTIFF", tc.filepath)
		}
	}
}

func BenchmarkReadJpeg(b *testing.B) {
	var (
		filepath = "../../test/data/TEST_2019-07-21_132615_DSC_0361_DxO_PL2.jpg"
		f        *File
	)
	for n := 0; n < b.N; n++ {
		osf, _ := os.Open(filepath)
		f, _ = Read(osf)
	}
	_ = f
}
