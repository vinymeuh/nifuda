// Copyright 2018 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package tag

// Tag is the interface for all types of image's metadata.
//
// Used as a return type for most of nifuda.Image's methods, so has its own package to avoid import cycle.
type Tag interface {
	Name() string
	StringValue() string
	Type() string
}
