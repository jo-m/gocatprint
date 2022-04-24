package main

import (
	"context"
	"image"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/jo-m/goprint/pkg/printer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "image/png"
)

type flags struct {
	LogPretty bool   `arg:"--log-pretty" default:"true" help:"log pretty"`
	LogLevel  string `arg:"--log-level" default:"info" help:"log level" placeholder:"LEVEL"`

	Timeout        time.Duration `arg:"-t,--timeout" default:"10s" help:"how long to allow for discovery and printing" placeholder:"DUR"`
	PrinterName    string        `arg:"--printer-name" default:"" help:"device name to connect to, ignored if empty" placeholder:"NAME"`
	PrinterAddress string        `arg:"--printer-address" default:"" help:"device address to connect to, ignored if empty" placeholder:"ADDR"`

	Image    string `arg:"positional,required"  help:"image to print, PNG or JPEG, must be 384 px wide" placeholder:"IMG"`
	DarkMode bool   `arg:"--dark-mode" default:"true" help:"more contrast, slower speed"`
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

func mustSetDefaultDevice() {
	dev, err := linux.NewDevice()
	if err != nil {
		log.Panic().Err(err).Msg("cannot instantiate new device")
	}
	log.Info().Msg("setting default ble device")
	ble.SetDefaultDevice(dev)
}

func mustReadImage(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		log.Panic().Err(err).Msg("could not open image")
	}
	defer f.Close()

	data, _, err := image.Decode(f)

	return data
}

func main() {
	f := flags{}
	err := arg.Parse(&f)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	setupLogging(f)

	log.Debug().Interface("flags", f).Msg("flags")

	img := mustReadImage("pkg/printer/testdata/test.png")

	// TODO
	mustSetDefaultDevice()

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
	defer printer.Close()

	printer.Print(ctx, img, true)
	if err != nil {
		log.Panic().Err(err).Send()
	}
}
