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
			filepath: "./testdata/TEST_2018-05-14_095545.jpg",
			gps: GpsTags{
				GPSVersionID:   "2.2.0.0",
				GPSLatitudeRef: "North",
				//GPSLatitude:        "35° 1' 1.03\"",
				GPSLongitudeRef: "East",
				//GPSLongitude:       "135° 46' 59.69\"",
				GPSAltitudeRef:     "Sea level",
				GPSAltitude:        102,
				GPSTimeStamp:       "00:55:44Z",
				GPSImgDirectionRef: "Magnetic direction",
				GPSImgDirection:    307,
				GPSMapDatum:        "WGS-84",
				GPSDateStamp:       "2018:05:14",
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
		gpsType := reflect.TypeOf(x.Gps)
		got := reflect.ValueOf(x.Gps)
		want := reflect.ValueOf(tc.gps)

		for i := 0; i < gpsType.NumField(); i++ {
			field := gpsType.Field(i)
			gotV := got.FieldByName(field.Name)
			wantV := want.FieldByName(field.Name)
			switch gotV.Kind() {
			case reflect.Float64:
				if gotV.Float() != wantV.Float() {
					t.Errorf("%s, Gps.%s: got=%f, want=%f", tc.filepath, field.Name, gotV.Float(), wantV.Float())
				}
			case reflect.String:
				if gotV.String() != wantV.String() {
					t.Errorf("%s, Gps.%s: got=%s, want=%s", tc.filepath, field.Name, gotV.String(), wantV.String())
				}
			}
		}

	}
}

func BenchmarkReadGpsTags(b *testing.B) {
	var (
		filepath = "./testdata/TEST_2018-05-14_095545.jpg"
		x        *Exif
		g        GpsTags
	)
	for n := 0; n < b.N; n++ {
		f, _ := os.Open(filepath)
		x, _ = Read(f)
		g = x.Gps
	}
	_ = g
}
