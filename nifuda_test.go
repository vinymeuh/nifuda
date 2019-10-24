// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

type testData struct {
	Filepath string
	Exif     Exif
}

func TestNifuda(t *testing.T) {
	testcases := []string{
		"./testdata/TEST_2018-05-14_095545.json",
		"./testdata/TEST_2019-07-21_132615_DSC_0361_DxO_PL2.json",
	}

	for _, tc := range testcases {
		// load test data
		data, err := ioutil.ReadFile(tc)
		if err != nil {
			t.Fatalf("%s: opening file failed, error=%s", tc, err)
		}

		td := testData{}
		err = json.Unmarshal([]byte(data), &td)
		if err != nil {
			t.Fatalf("%s: unmarshaling test data failed, error=%s", tc, err)
		}

		// read exif data
		f, err := os.Open(td.Filepath)
		if err != nil {
			t.Fatalf("%s: opening file failed, error=%s", td.Filepath, err)
		}
		defer f.Close()

		x, err := Read(f)
		if err != nil || x == nil {
			t.Errorf("%s: reading exifs failed, error=%s, x=%v", td.Filepath, err, x)
			continue
		}

		// test IFD0 Tags
		testEachFields(t, td.Filepath,
			reflect.TypeOf(x.Ifd0),
			reflect.ValueOf(x.Ifd0),
			reflect.ValueOf(td.Exif.Ifd0),
		)

		// test Exif Tags
		testEachFields(t, td.Filepath,
			reflect.TypeOf(x.Exif),
			reflect.ValueOf(x.Exif),
			reflect.ValueOf(td.Exif.Exif),
		)

		// test GPS Tags
		testEachFields(t, td.Filepath,
			reflect.TypeOf(x.Gps),
			reflect.ValueOf(x.Gps),
			reflect.ValueOf(td.Exif.Gps),
		)
	}
}

func testEachFields(t *testing.T, filepath string, theType reflect.Type, got reflect.Value, want reflect.Value) {
	for i := 0; i < theType.NumField(); i++ {
		field := theType.Field(i)
		gotV := got.FieldByName(field.Name)
		wantV := want.FieldByName(field.Name)

		switch gotV.Kind() {
		case reflect.Float64:
			if gotV.Float() != wantV.Float() {
				t.Errorf("%s, %s.%s: got=%f, want=%f", filepath, theType.Name(), field.Name, gotV.Float(), wantV.Float())
			}
		case reflect.String:
			if gotV.String() != wantV.String() {
				t.Errorf("%s, %s.%s: got=%s, want=%s", filepath, theType.Name(), field.Name, gotV.String(), wantV.String())
			}
		case reflect.Uint16:
			if gotV.Uint() != wantV.Uint() {
				t.Errorf("%s, %s.%s: got=%d, want=%d", filepath, theType.Name(), field.Name, gotV.Uint(), wantV.Uint())
			}
		}
	}
}
