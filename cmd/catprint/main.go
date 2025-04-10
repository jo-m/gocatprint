/*
Package main (catprint) is a CLI utility to print an image file,
showcasing the catprint advanced API.
It supports dithering and thresholding.

Usage:

	sudo ./catprint ~/image.png
	./catprint --help
*/
package main

import (
	"context"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/rs/zerolog/log"
	"jo-m.ch/go/gocatprint/internal/pkg/logging"
	"jo-m.ch/go/gocatprint/pkg/printer"
)

type flags struct {
	logging.LogConfig

	HCIDevice      int           `arg:"--hci-device" default:"-1" help:"HCI device to use, auto if negative" placeholder:"N"`
	Timeout        time.Duration `arg:"--timeout" default:"10s" help:"how long to allow for discovery and printing" placeholder:"DUR"`
	PrinterName    string        `arg:"--printer-name" default:"" help:"device name to connect to, ignored if empty" placeholder:"NAME"`
	PrinterAddress string        `arg:"--printer-address" default:"" help:"device address to connect to, ignored if empty" placeholder:"ADDR"`

	Threshold bool   `arg:"--threshold" default:"false" help:"use simple thresholding instead of dithering"`
	Preview   string `arg:"--preview" default:"" help:"do not print, just write the (processed) image to the given file" placeholder:"OUT-FILE"`
	Image     string `arg:"positional,required"  help:"image to print, PNG or JPEG, must be 384px wide (unless --scale is passed)" placeholder:"IN-FILE"`
}

func mustSetDefaultDevice(f flags) {
	if f.HCIDevice >= 0 {
		dev, err := linux.NewDevice(ble.OptDeviceID(f.HCIDevice))
		if err != nil {
			log.Panic().Err(err).Msg("cannot instantiate new device")
		}

		log.Debug().Msg("setting default ble device")
		ble.SetDefaultDevice(dev)
		return
	}

	for i := 0; i < 10; i++ {
		dev, err := linux.NewDevice(ble.OptDeviceID(i))
		if err == nil {
			ble.SetDefaultDevice(dev)
			return
		}
	}

	log.Panic().Msg("could not find a HCI device")
}

func mustFindPrinter(f flags) *printer.Printer {
	mustSetDefaultDevice(f)

	ctx, cancel := context.WithTimeout(context.Background(), f.Timeout)
	defer cancel()

	opts := printer.FindOptions{
		Name:    f.PrinterName,
		Address: f.PrinterAddress,
	}

	printer, err := printer.Find(ctx, opts)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	return printer
}

func mustLoadImage(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		log.Panic().Err(err).Msg("could not open image")
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Panic().Err(err).Msg("could not parse image")
	}

	return img
}

func main() {
	f := flags{}
	arg.MustParse(&f)
	logging.MustInit(f.LogConfig)
	log.Debug().Interface("flags", f).Msg("flags")

	log.Info().Msg("loading image..")
	img := mustLoadImage(f.Image)
	img = printer.PrepareImage(img, f.Threshold)

	if f.Preview != "" {
		out, err := os.Create(f.Preview)
		if err != nil {
			log.Panic().Err(err).Send()
		}

		log.Info().Msg("writing image..")
		err = png.Encode(out, img)
		if err != nil {
			log.Panic().Err(err).Send()
		}
		log.Info().Msg("done.")
		return
	}

	log.Info().Msg("connecting to printer..")
	printer := mustFindPrinter(f)
	defer printer.Close()

	log.Info().Msg("printing..")
	err := printer.Print(context.Background(), img)
	if err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Msg("done.")
}
