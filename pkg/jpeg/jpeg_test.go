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
	markers  []Marker
}{
	{"../../test/data/minimal.jpg", []Marker{SOI, APP0, EOI}},
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

		for i, expected := range tc.markers {
			if jf.Segments[i].Marker != expected {
				t.Errorf("%s: wrong segment marker, expected=%s, actual=%s", tc.filepath, expected, jf.Segments[i].Marker) //FIXME: formatage segment marker
			}
		}
	}
}
