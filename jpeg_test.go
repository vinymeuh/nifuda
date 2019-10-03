// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"os"
	"testing"
)

func TestJpegFileMustReturnError(t *testing.T) {
	tests := []struct {
		filepath string
	}{
		{"./test/data/empty.txt"},
		{"./test/data/dummy.txt"},
		{"./test/data/nosoi.jpg"},
		{"./test/data/minimal.jpg"},
	}

	for _, tc := range tests {
		osf, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening file failed with err=%s", tc.filepath, err)
		}
		defer osf.Close()

		f, err := jpegRead(osf)
		if err == nil || f != nil {
			t.Errorf("%s: reading file should have failed and returned nil, err=%s, f=%v", tc.filepath, err, f)
		}
	}
}
