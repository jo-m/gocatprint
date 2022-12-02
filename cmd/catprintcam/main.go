/*
Package main (catprintcamp) is a CLI utility to grab a picture from a webcam and print it,
showcasing the simple API.
Pictures are automatically dithered.

Usage:

	sudo ./catprintcam --pic-dev /dev/video0
	./catprintcam --help
*/
package main

import (
	_ "image/gif"
	_ "image/jpeg"

	"github.com/alexflint/go-arg"
	"github.com/jo-m/gocatprint/internal/pkg/logging"
	"github.com/jo-m/gocatprint/pkg/simple"
	"github.com/rs/zerolog/log"
)

type flags struct {
	logging.LogConfig
	CamConfig
}

func main() {
	f := flags{}
	arg.MustParse(&f)
	logging.MustInit(f.LogConfig)
	log.Debug().Interface("flags", f).Msg("flags")

	log.Info().Msg("snapping picture..")
	im, err := Snap(f.CamConfig)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	log.Info().Msg("printing..")
	err = simple.Print(im, false)
	if err != nil {
		log.Panic().Err(err).Send()
	}
	log.Info().Msg("done.")
}
