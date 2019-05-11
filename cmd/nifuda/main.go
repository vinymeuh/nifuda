// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package main

import (
	"os"

	"nifuda/app"
)

// variables defined at build time
var (
	version string
	build   string
)

func main() {
	app := app.New(version, build)
	if app.Run(os.Args) != nil {
		os.Exit(1)
	}
}
