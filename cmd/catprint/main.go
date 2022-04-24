package main

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/alexflint/go-arg"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/jo-m/gocatprint/pkg/printer"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
)

type flags struct {
	LogPretty bool   `arg:"--log-pretty" default:"true" help:"log pretty"`
	LogLevel  string `arg:"--log-level" default:"info" help:"log level" placeholder:"LEVEL"`

	HCIDevice      int           `arg:"--hci-device" default:"0" help:"HCI device to use" placeholder:"N"`
	Timeout        time.Duration `arg:"--timeout" default:"10s" help:"how long to allow for discovery and printing" placeholder:"DUR"`
	PrinterName    string        `arg:"--printer-name" default:"" help:"device name to connect to, ignored if empty" placeholder:"NAME"`
	PrinterAddress string        `arg:"--printer-address" default:"" help:"device address to connect to, ignored if empty" placeholder:"ADDR"`

	DarkMode bool   `arg:"--dark-mode" default:"true" help:"more contrast, slower speed"`
	NoScale  bool   `arg:"--no-scale" default:"false" help:"do not scale input image, must be provided with 384px width"`
	NoDither bool   `arg:"--no-dither" default:"false" help:"do not apply dithering to the image, use simple thresholding instead"`
	Preview  string `arg:"--preview" default:"" help:"do not print, just write the (processed) image to the given file" placeholder:"IMG-FILE"`
	Image    string `arg:"positional,required"  help:"image to print, PNG or JPEG, must be 384px wide (unless --scale is passed)" placeholder:"IMG-FILE"`
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

func mustPrepareImage(f flags) image.Image {
	file, err := os.Open(f.Image)
	if err != nil {
		log.Panic().Err(err).Msg("could not open image")
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Panic().Err(err).Msg("could not parse image")
	}

	if !f.NoScale {
		ratio := float64(printer.PrintWidth) / float64(img.Bounds().Max.X)

		scaled := image.NewRGBA(image.Rect(0, 0, printer.PrintWidth, int(ratio*float64(img.Bounds().Max.Y))))
		draw.BiLinear.Scale(scaled, scaled.Rect, img, img.Bounds(), draw.Over, nil)

		img = scaled
	}

	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, gray.Bounds(), img, img.Bounds().Min, draw.Src)
	img = gray

	if !f.NoDither {
		// dither
		palette := []color.Color{
			color.Black,
			color.White,
		}

		d := dither.NewDitherer(palette)
		d.Matrix = dither.FloydSteinberg

		img = d.Dither(img)
	} else {
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
	}

	return img
}

func main() {
	f := flags{}
	arg.MustParse(&f)
	setupLogging(f)
	log.Debug().Interface("flags", f).Msg("flags")

	log.Info().Msg("loading & preparing image..")
	img := mustPrepareImage(f)

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
	err := printer.Print(context.Background(), img, f.DarkMode)
	if err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Msg("done.")
}
