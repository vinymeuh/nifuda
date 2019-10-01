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
		//{"../../test/data/recursive_ifd0.tif"},
		//{"../../test/data/wrong_ifd1.tif"},
	}

	for _, tc := range tests {
		osf, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening file failed with err=%s", tc.filepath, err)
		}
		defer osf.Close()

		f, err := Read(osf)
		if err == nil {
			t.Errorf("%s: reading file should have failed, err=%s, f=%v", tc.filepath, err, f)
		}
	}
}

func BenchmarkReadTiff(b *testing.B) {
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
