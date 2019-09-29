// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package exif

import (
	"bytes"
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
	type testTag struct {
		namespace string
		name      string
		value     string
	}
	tests := []struct {
		filepath string
		format   fileFormat
		tags     []testTag
	}{
		{"../../test/data/minimal.tif", ffTIFF, []testTag{
			{"ifd0", "ImageLength", "1"},
		}},
		{"../../test/data/TEST_2019-07-21_132615_DSC_0361_DxO_PL2.jpg", ffJPEG, []testTag{
			{"ifd0", "Make", "NIKON CORPORATION"},
			{"exif", "LensMake", "NIKON"},
		}},
	}

	for _, tc := range tests {
		osf, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening file failed with err=%s", tc.filepath, err)
		}
		defer osf.Close()

		f, err := Read(osf)
		if err != nil || f == nil {
			t.Errorf("%s: reading fails, error=%s, f=%v", tc.filepath, err, f)
			continue
		}

		if f.format != tc.format {
			t.Errorf("%s: incorrect file format, expected=%v, got=%v", tc.filepath, tc.format, f.format)
		}

		ftags := f.Tags()
		for _, tctag := range tc.tags {
			if _, ok := ftags[tctag.namespace]; ok == false {
				t.Errorf("%s, %v: tag namespace not found", tc.filepath, tctag)
				continue
			}

			if _, ok := ftags[tctag.namespace][tctag.name]; ok == false {
				t.Errorf("%s, %v: tag name not found in namespace", tc.filepath, tctag)
				continue
			}

			got := ftags[tctag.namespace][tctag.name].Value().String()
			if got != tctag.value {
				t.Errorf("%s, %v: error with tag value, expected=%s, got=%s", tc.filepath, tctag, tctag.value, got)
			}
		}
	}
}

var format fileFormat // avoid compiler optimisation
func BenchmarkIdentifyFileFormat(b *testing.B) {
	var (
		benchData = []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00}
		ff        fileFormat
	)
	rs := bytes.NewReader(benchData)
	for n := 0; n < b.N; n++ {
		ff = identifyFileFormat(rs)
	}
	format = ff
}

func BenchmarkReadExifJpeg(b *testing.B) {
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

func BenchmarkReadExifTiff(b *testing.B) {
	var (
		filepath = "../../test/data/TEST_2019-07-21_132615_DSC_0361.NEF"
		f        *File
	)
	for n := 0; n < b.N; n++ {
		osf, _ := os.Open(filepath)
		f, _ = Read(osf)
	}
	_ = f
}
