// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"

	"github.com/vinymeuh/nifuda/jpeg"
)

func cmdDevEmptyJPEG(cmd *cli.Cmd) {
	var (
		imgFile = cmd.StringArg("FILE", "", "Name of the file to create")
	)

	cmd.Spec = "FILE"

	cmd.Action = func() {

		// do not override file if already exists
		out, err := os.OpenFile(*imgFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			fmt.Println(err)
			cli.Exit(2)
		}

		out.Write([]byte(jpeg.SOI))
		out.Write([]byte(jpeg.EOI))

		out.Close()
	}
}
