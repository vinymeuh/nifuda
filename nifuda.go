// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package nifuda

import (
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/vinymeuh/nifuda/tag"

	"github.com/vinymeuh/nifuda/jpeg"
	"github.com/vinymeuh/nifuda/tiff"
)

// Image is the top-level interface of nifuda module which define the methods to access image's tags.
//
// Image's tags must be identified by their namespace (Exif, XMP, GPS, ...) and their name.
type Image interface {
	// TagsNamespaces returns all namespaces.
	TagsNamespaces() []string
	// GetTagsFromNamespace returns all image's tags for a given namespace.
	// Returns a empty array if namespace not found or contains no tags.
	GetTagsFromNamespace(namespace string) []tag.Tag
	// GetTag returns an image's tag uniquely identified by its namespace and its name.
	// Returns nil if tag not found.
	GetTag(namespace string, name string) tag.Tag
}

// Parse method is "the generic constructor of Image".
//
// Most applications will use it to create a new Image from a supported image format. The alternative
// will be to create directly an Image from a type implementing the Image interface.
//
// Currently, only TIFF and JPEG images are supported.
func Parse(rs io.ReadSeeker) (Image, error) {
	b := make([]byte, 16)
	if _, err := rs.Read(b); err != nil {
		return nil, errors.Wrap(err, "failed to read the 16 first bytes from header")
	}
	switch string(b[0:2]) {
	case string(jpeg.SOI):
		return jpeg.Parse(rs)
	case "II", "MM":
		return tiff.Parse(rs)
	}
	return nil, errors.New("unknown image format")
}

// ParseFromFile is a convenient function to create an Image from a File instead of an ReaderSeeker.
//
// It only manages the open/close of the file before to call Parse().
func ParseFromFile(filePath string) (Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open %q", filePath)
	}
	defer f.Close()
	return Parse(f)
}
