# nifuda

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/vinymeuh/nifuda?status.svg)](https://godoc.org/github.com/vinymeuh/nifuda)
[![Go Report Card](https://goreportcard.com/badge/github.com/vinymeuh/nifuda)](https://goreportcard.com/report/github.com/vinymeuh/nifuda)
[![Build Status](https://travis-ci.org/vinymeuh/nifuda.svg?branch=master)](https://travis-ci.org/vinymeuh/nifuda)
[![codecov](https://codecov.io/gh/vinymeuh/nifuda/branch/master/graph/badge.svg)](https://codecov.io/gh/vinymeuh/nifuda)

`nifuda` provides a native Go library to read tags from EXIF image files.

## Getting Started

As a example, a very simplistic EXIF reader:

```golang
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/vinymeuh/nifuda"
)

func main() {
    if len(os.Args) != 1 {
        log.Fatalf("Usage: %s EXIF_FILE", os.Args[0])
    }

    f, err := os.Open(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    x, err := nifuda.Read(f)
    if err != nil {
        log.Fatal(err)
    }

    for _, tag := range x.Ifd0 {
        fmt.Printf("ifd0   %-30s   %s\n", tag.Name(), tag.String())
    }
    for _, tag := range x.Exif {
        fmt.Printf("exif   %-30s   %s\n", tag.Name(), tag.String())
    }  
}
```

More examples are availables in [nf-tools repository](https://github.com/vinymeuh/nf-tools).

## References

### Exif

* [CIPA Standards (Exif & DCF)](http://www.cipa.jp/std/std-sec_e.html)
* [Description of Exif file format](http://gvsoft.no-ip.org/exif/exif-explanation.html)

### TIFF

* [TIIF, Revision 6.0 (Library of Congress)](https://www.loc.gov/preservation/digital/formats/fdd/fdd000022.shtml)
* [TIFF at FileFormat.Info](http://www.fileformat.info/format/tiff/index.dir)
* [Adobe TIFF specification](https://www.adobe.io/open/standards/TIFF.html)
* [TIFF File Format FAQ (AWare Systems)](https://www.awaresystems.be/imaging/tiff/faq.html)
* [TIFF Tag Reference](https://www.awaresystems.be/imaging/tiff/tifftags.html)
* [Tags for TIFF, DNG, and Related Specifications](https://www.loc.gov/preservation/digital/formats/content/tiff_tags.shtml)
