// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"os"
	"testing"
)

func BenchmarkReadExifJpeg(b *testing.B) {
	var (
		filepath = "./2018-05-14_095545.jpg"
		x        *Exif
	)
	for n := 0; n < b.N; n++ {
		f, _ := os.Open(filepath)
		x, _ = Read(f)
	}
	_ = x
}
