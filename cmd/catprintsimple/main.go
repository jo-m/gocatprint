/*
Package main (catprint) is a CLI utility to print an image file,
showcasing the catprint simple API.
The image is automatically dithered.

Usage:

	sudo ./catprintsimple -img pkg/printer/testdata/swan.jpg
*/
package main

import (
	"flag"
	"image"
	"os"

	_ "image/gif"
	_ "image/jpeg"

	"github.com/jo-m/gocatprint/pkg/simple"
)

func main() {
	imgPath := flag.String("img", "", "image file to print")
	flag.Parse()

	file, err := os.Open(*imgPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	err = simple.Print(img, false)
	if err != nil {
		panic(err)
	}
}
