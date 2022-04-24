package printer

import (
	"image"
	"image/color"

	"github.com/makeworld-the-better-one/dither/v2"
	"golang.org/x/image/draw"
)

// PrepareImage prepares an image (scaling, thresholding, dithering)
// so that it can be passed to printer.Print().
func PrepareImage(img image.Image, threshold bool) image.Image {
	// scale
	ratio := float64(PrintWidth) / float64(img.Bounds().Max.X)
	scaled := image.NewRGBA(image.Rect(0, 0, PrintWidth, int(ratio*float64(img.Bounds().Max.Y))))
	draw.BiLinear.Scale(scaled, scaled.Rect, img, img.Bounds(), draw.Over, nil)
	img = scaled

	// grayscale
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, gray.Bounds(), img, img.Bounds().Min, draw.Src)
	img = gray

	if threshold {
		// simple thresholding
		access := img.(*image.Gray)
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				if access.GrayAt(x, y).Y > 127 {
					access.Set(x, y, color.White)
				} else {
					access.Set(x, y, color.Black)
				}
			}
		}
	} else {
		// dither
		palette := []color.Color{
			color.Black,
			color.White,
		}

		d := dither.NewDitherer(palette)
		d.Matrix = dither.FloydSteinberg

		img = d.Dither(img)
	}

	return img
}
