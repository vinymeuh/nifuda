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
		{"./testdata/errors/empty.txt"},
		{"./testdata/errors/dummy.txt"},
		{"./testdata/errors/nosoi.jpg"},
		{"./testdata/errors/minimal.jpg"},
		{"./testdata/errors/wrong_version.tif"},
		{"./testdata/errors/wrong_offset0.tif"},
		{"./testdata/errors/no_ifd0.tif"},
		//{"./testdata/errors/recursive_ifd0.tif"},  FIXME
		//{"./testdata/errors/wrong_ifd1.tif"}, FIXME
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
