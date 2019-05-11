// Copyright 2019 VinyMeuh. All rights reserved.
// Use of the source code is governed by a MIT-style license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"runtime"

	cli "github.com/jawher/mow.cli"
)

var app *App

type App struct {
	*cli.Cli
}

const (
	appName = "nifuda"
	appDesc = "A tool to manage my photo library"
)

func New(version string, build string) *App {
	if app == nil {
		app = &App{
			Cli: cli.App(appName, appDesc),
		}

		// top-level global options
		app.Version("v version", fmt.Sprintf("%s %s built %s for %s/%s with %s",
			appName, version, build, runtime.GOOS, runtime.GOARCH, runtime.Version()))

		// subcommands
		app.Command("show", "Display image metadata", cmdNotYetImplemented)

		app.Command("dev", "Miscellaneous tools for development and debug", func(cmd *cli.Cmd) {
			cmd.Command("empty-jpeg", "Create an empty JPEG file", cmdDevEmptyJPEG)
			cmd.Command("empty-tiff", "Create an empty TIFF file", cmdDevEmptyTIFF)
			cmd.Command("ps", "Print structure of a file", cmdDevPrintStructure)
		})
	}
	return app
}

func cmdNotYetImplemented(cmd *cli.Cmd) {
	cmd.Action = func() {
		fmt.Println("not yet implemented")
	}
}
