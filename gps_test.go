// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"os"
	"reflect"
	"testing"
)

func TestGpsTags(t *testing.T) {
	tests := []struct {
		filepath string
		gps      GpsTags
	}{
		{
			filepath: "./2018-05-14_095545.jpg",
			gps: GpsTags{
				GPSVersionID:    "2.2.0.0",
				GPSLatitudeRef:  "N",
				GPSLongitudeRef: "E",
			},
		},
	}

	for _, tc := range tests {
		// open file
		f, err := os.Open(tc.filepath)
		if err != nil {
			t.Fatalf("%s: opening file failed, error=%s", tc.filepath, err)
		}
		defer f.Close()

		// read exifs
		x, err := Read(f)
		if err != nil || x == nil {
			t.Errorf("%s: reading exifs failed, error=%s, x=%v", tc.filepath, err, x)
			continue
		}

		// test equality on each fields of GpsTags, one by one to have a meaningful message in case of error
		typ := reflect.TypeOf(x.Gps)
		got := reflect.ValueOf(x.Gps)
		want := reflect.ValueOf(tc.gps)

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			gotV := got.FieldByName(field.Name).String()
			wantV := want.FieldByName(field.Name).String()
			if gotV != wantV {
				t.Errorf("%s, Gps.%s: got=%s, want=%s", tc.filepath, field.Name, gotV, wantV)
			}
		}

	}
}
