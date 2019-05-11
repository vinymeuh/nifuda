// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package app

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
)

func cmdDevEmptyTIFF(cmd *cli.Cmd) {
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

		var order binary.ByteOrder = binary.LittleEndian
		var buffer bytes.Buffer

		// Tiff Header
		switch order {
		case binary.LittleEndian:
			out.Write([]byte("II"))
		case binary.BigEndian:
			out.Write([]byte("MM"))
		}
		binary.Write(&buffer, order, uint16(42)) // tiff version
		binary.Write(&buffer, order, uint32(8))  // offset ifd0

		// ifd0
		binary.Write(&buffer, order, uint16(1)) // 1 ifd entrie

		binary.Write(&buffer, order, uint16(256)) // id - ImageWidth
		binary.Write(&buffer, order, uint16(4))   // type -  Long
		binary.Write(&buffer, order, uint32(1))   // count
		binary.Write(&buffer, order, uint32(0))   // value

		binary.Write(&buffer, order, uint32(0)) // no next offset
		out.Write(buffer.Bytes())

		out.Close()
	}
}
