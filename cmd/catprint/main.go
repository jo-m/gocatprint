package main

import (
	"context"
	"image"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/jo-m/gocatprint/pkg/printer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type flags struct {
	LogPretty bool   `arg:"--log-pretty" default:"true" help:"log pretty"`
	LogLevel  string `arg:"--log-level" default:"info" help:"log level" placeholder:"LEVEL"`

	HCIDevice      int           `arg:"--hci-device" default:"0" help:"HCI device to use" placeholder:"N"`
	Timeout        time.Duration `arg:"--timeout" default:"10s" help:"how long to allow for discovery and printing" placeholder:"DUR"`
	PrinterName    string        `arg:"--printer-name" default:"" help:"device name to connect to, ignored if empty" placeholder:"NAME"`
	PrinterAddress string        `arg:"--printer-address" default:"" help:"device address to connect to, ignored if empty" placeholder:"ADDR"`

	DarkMode bool   `arg:"--dark-mode" default:"true" help:"more contrast, slower speed"`
	Image    string `arg:"positional,required"  help:"image to print, PNG or JPEG, must be 384 px wide" placeholder:"IMG"`
}

func setupLogging(f flags) {
	level, err := zerolog.ParseLevel(f.LogLevel)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	if f.LogPretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Logger = log.Logger.Level(level).With().Timestamp().Caller().Logger()
}

func mustReadImage(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		log.Panic().Err(err).Msg("could not open image")
	}
	defer f.Close()

	data, _, err := image.Decode(f)
	if err != nil {
		log.Panic().Err(err).Msg("could not parse image")
	}

	return data
}

func mustSetDefaultDevice(f flags) {
	dev, err := linux.NewDevice(ble.OptDeviceID(f.HCIDevice))
	if err != nil {
		log.Panic().Err(err).Msg("cannot instantiate new device")
	}
	log.Debug().Msg("setting default ble device")
	ble.SetDefaultDevice(dev)
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

func main() {
	f := flags{}
	arg.MustParse(&f)
	setupLogging(f)
	log.Debug().Interface("flags", f).Msg("flags")

	log.Info().Msg("loading image..")
	img := mustReadImage(f.Image)

	log.Info().Msg("connecting to printer..")
	printer := mustFindPrinter(f)
	defer printer.Close()

	log.Info().Msg("printing..")
	err := printer.Print(context.Background(), img, f.DarkMode)
	if err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Msg("done.")
}
