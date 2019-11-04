// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
)

// An Image File Directory (IFD) consists of a 2-byte count of the number of directory entries, followed by a
// sequence of 12-byte field entries, followed by a 4-byte offset of the next IFD (or 0 if none).
//
// There must be at least 1 IFD in a TIFF file and each IFD must have at least one entry.
//
// Each TIFF field has an associated Count.
// This means that all fields are actually one-dimensional arrays, even though most fields contain only a single value.
type ifd struct {
	next    uint32   // offset in bytes to the next IFD, from the start of the file. 0 if none
	entries uint16   // number of directory entries
	tags    []ifdTag // list of undecoded tags
}

type ifdTag struct {
	id       uint16 // tag identifier
	tiffType uint16 // tiff type idendifier
	count    uint32 // the number of values in data
	data     []byte // undecoded payload for tag
}

func (it ifdTag) byteToInt(bo binary.ByteOrder) []int {
	b := make([]int, it.count)
	raw := bytes.NewReader(it.data)
	var v uint8
	for i := range b {
		binary.Read(raw, bo, &v)
		b[i] = int(v)
	}
	return b
}

func (it ifdTag) asciiToString() string {
	return string(it.data[0 : it.count-1]) // -1 to remove character '\0'
}

func (it ifdTag) shortToUint16(bo binary.ByteOrder) []uint16 {
	var s uint16
	S := make([]uint16, it.count)
	raw := bytes.NewReader(it.data)
	for i := range S {
		binary.Read(raw, bo, &s)
		S[i] = s
	}
	return S
}

func (it ifdTag) longToUint32(bo binary.ByteOrder) []uint32 {
	var l uint32
	L := make([]uint32, it.count)
	raw := bytes.NewReader(it.data)
	for i := range L {
		binary.Read(raw, bo, &l)
		L[i] = l
	}
	return L
}

func (it ifdTag) rationalToFloat32(bo binary.ByteOrder) []float32 {
	var n, d uint32
	r := make([]float32, it.count)
	raw := bytes.NewReader(it.data)
	for i := range r {
		binary.Read(raw, bo, &n)
		binary.Read(raw, bo, &d)
		r[i] = float32(n / d)
	}
	return r
}

func (it ifdTag) undefinedToString() string {
	return string(it.data[0:it.count])
}

// Helpers
func intArrayToString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}
