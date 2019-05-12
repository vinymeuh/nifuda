// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"

	"github.com/vinymeuh/nifuda"
)

func cmdShow(cmd *cli.Cmd) {
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
	
		for _, namespace := range img.TagsNamespaces() {
			for _, tag := range img.GetTagsFromNamespace(namespace) {
				fmt.Fprintf(os.Stdout, "%-10s %-30s %-10v %s\n", namespace, tag.Name(), tag.Type(), tag.StringValue())
			}
		}
	}
}
