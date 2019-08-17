// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.
package tiff

import (
	"os"
	"testing"
)

func TestTiffFileMustReturnError(t *testing.T) {
	tests := []struct {
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

	for _, tc := range tests {
		osf, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening file failed with err=%s", tc.filepath, err)
		}
		defer osf.Close()

		f, err := Read(osf, nil)
		if err == nil || f != nil {
			t.Errorf("%s: reading file should have failed and returned nil, err=%s, f=%v", tc.filepath, err, f)
		}
	}
}

func TestTiffFile(t *testing.T) {
	tests := []struct {
		filepath string
	}{
		{"../../test/data/minimal.tif"},
		{"../../test/data/minimal_with_ifd1.tif"},
	}

	for _, tc := range tests {
		osf, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening file failed with err=%s", tc.filepath, err)
		}
		defer osf.Close()

		f, err := Read(osf, nil)
		if err != nil || f == nil {
			t.Errorf("%s: reading fails, err=%s, f=%v", tc.filepath, err, f)
			continue
		}

	}
}
