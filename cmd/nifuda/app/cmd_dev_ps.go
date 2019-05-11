// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"

	"github.com/vinymeuh/nifuda"
	"github.com/vinymeuh/nifuda/jpeg"
	"github.com/vinymeuh/nifuda/tiff"
)

func cmdDevPrintStructure(cmd *cli.Cmd) {
	var (
		imgFile = cmd.StringArg("FILE", "", "Image file")
	)

	cmd.Spec = "FILE"

	cmd.Action = func() {
		img, err := nifuda.ParseFromFile(*imgFile)
		if err != nil {
			fmt.Println(err)
			cli.Exit(2)
		}

		switch img.(type) {
		case *jpeg.Jpeg:
			img := img.(*jpeg.Jpeg)
			img.PrintStructure(os.Stdout)
		case *tiff.Tiff:
			img := img.(*tiff.Tiff)
			img.PrintStructure(os.Stdout)
		}
	}
}
