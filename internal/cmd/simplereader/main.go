package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/vinymeuh/nifuda"
)

func main() {
	if len(os.Args) != 2 {
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

	v := reflect.ValueOf(x.Image)
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		fmt.Printf("Image.%-30s   %v\n", vt.Field(i).Name, v.Field(i).Interface())
	}

	v = reflect.ValueOf(x.Photo)
	vt = v.Type()
	for i := 0; i < vt.NumField(); i++ {
		fmt.Printf("Photo.%-30s   %v\n", vt.Field(i).Name, v.Field(i).Interface())
	}
}
