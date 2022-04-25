package main

import (
	"context"
	"image"
	"image/png"
	"os"
	"time"

	_ "image/gif"
	_ "image/jpeg"

	"github.com/alexflint/go-arg"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/jo-m/gocatprint/pkg/printer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type flags struct {
	LogPretty bool   `arg:"--log-pretty" default:"true" help:"log pretty"`
	LogLevel  string `arg:"--log-level" default:"info" help:"log level" placeholder:"LEVEL"`

	HCIDevice      int           `arg:"--hci-device" default:"-1" help:"HCI device to use, auto if negative" placeholder:"N"`
	Timeout        time.Duration `arg:"--timeout" default:"10s" help:"how long to allow for discovery and printing" placeholder:"DUR"`
	PrinterName    string        `arg:"--printer-name" default:"" help:"device name to connect to, ignored if empty" placeholder:"NAME"`
	PrinterAddress string        `arg:"--printer-address" default:"" help:"device address to connect to, ignored if empty" placeholder:"ADDR"`

	FastMode  bool   `arg:"--fast-mode" default:"false" help:"less contrast, higher printer speed"`
	Threshold bool   `arg:"--threshold" default:"false" help:"use simple thresholding instead of dithering"`
	Preview   string `arg:"--preview" default:"" help:"do not print, just write the (processed) image to the given file" placeholder:"OUT-FILE"`
	Image     string `arg:"positional,required"  help:"image to print, PNG or JPEG, must be 384px wide (unless --scale is passed)" placeholder:"IN-FILE"`
}

func mustSetupLogging(f flags) {
	level, err := zerolog.ParseLevel(f.LogLevel)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	if f.LogPretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Logger = log.Logger.Level(level).With().Timestamp().Caller().Logger()
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
	mustSetupLogging(f)
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
	err := printer.Print(context.Background(), img, !f.FastMode)
	if err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Msg("done.")
}
