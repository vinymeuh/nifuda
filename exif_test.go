// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"os"
	"testing"
)

func TestExifFileError(t *testing.T) {
	tests := []struct {
		filepath string
	}{
		{"./test/data/empty.txt"},
		{"./test/data/dummy.txt"},
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

func TestExifRead(t *testing.T) {
	type testTag struct {
		name  string
		value string
	}
	tests := []struct {
		filepath string
		ifd0     []testTag
		exif     []testTag
	}{
		{filepath: "./test/data/minimal.tif",
			ifd0: []testTag{
				{"ImageLength", "1"},
			}},
		{filepath: "./test/data/TEST_2019-07-21_132615_DSC_0361_DxO_PL2.jpg",
			ifd0: []testTag{
				{"Make", "NIKON CORPORATION"},
			},
			exif: []testTag{
				{"LensMake", "NIKON"},
			}},
	}

	for _, tc := range tests {
		f, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening file failed with err=%s", tc.filepath, err)
		}
		defer f.Close()

		x, err := Read(f)
		if err != nil || x == nil {
			t.Errorf("%s: reading fails, error=%s, x=%v", tc.filepath, err, x)
			continue
		}

		// namespace Ifd0
		for _, tctag := range tc.ifd0 {
			if _, ok := x.Ifd0[tctag.name]; ok == false {
				t.Errorf("%s, %v: tag name not found in Ifd0", tc.filepath, tctag)
				continue
			}

			got := x.Ifd0[tctag.name].Value().String()
			if got != tctag.value {
				t.Errorf("%s, Ifd0.%v: error with tag value, expected=%s, got=%s", tc.filepath, tctag, tctag.value, got)
			}
		}

		// namespace Exif
		for _, tctag := range tc.exif {
			if _, ok := x.Exif[tctag.name]; ok == false {
				t.Errorf("%s, %v: tag name not found in Exif", tc.filepath, tctag)
				continue
			}

			got := x.Exif[tctag.name].Value().String()
			if got != tctag.value {
				t.Errorf("%s, Exif.%v: error with tag value, expected=%s, got=%s", tc.filepath, tctag, tctag.value, got)
			}
		}
	}
}

func BenchmarkReadExifJpeg(b *testing.B) {
	var (
		filepath = "./test/data/TEST_2019-07-21_132615_DSC_0361_DxO_PL2.jpg"
		x        *Exif
	)
	for n := 0; n < b.N; n++ {
		f, _ := os.Open(filepath)
		x, _ = Read(f)
	}
	_ = x
}
