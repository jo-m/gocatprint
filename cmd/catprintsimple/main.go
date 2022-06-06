package main

import (
	"flag"
	"image"
	"os"

	_ "image/gif"
	_ "image/jpeg"

	"github.com/jo-m/gocatprint/pkg/simple"
)

/*
Usage:
	./catprintsimple -img pkg/printer/testdata/swan.jpg
*/

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
