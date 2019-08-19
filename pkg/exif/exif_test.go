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
	type testTag struct {
		namespace string
		name      string
		value     string
	}
	tests := []struct {
		filepath string
		format   FileFormat
		tags     []testTag
	}{
		{"../../test/data/minimal.tif", TIFF, []testTag{
			{"ifd0", "ImageLength", "1"},
		}},
		{"../../test/data/TEST_2019-07-21_132615_DSC_0361_DxO_PL2.jpg", JPEG, []testTag{
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
