package simple

import (
	"context"
	"errors"
	"image"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"jo-m.ch/go/gocatprint/pkg/printer"
)

/*
Package simple provides a simple wrapper around the printer API.
*/

func setBleDefaultDevice() error {
	for i := 0; i < 10; i++ {
		dev, err := linux.NewDevice(ble.OptDeviceID(i))
		if err == nil {
			ble.SetDefaultDevice(dev)
			return nil
		}
	}

	return errors.New("could not find a HCI device")
}

// Print tries to print the give image via the first Bluetooth device found.
// If threshold is true, simple thresholding is used instead of dithering.
func Print(img image.Image, threshold bool) error {
	img = printer.PrepareImage(img, threshold)

	err := setBleDefaultDevice()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	printer, err := printer.Find(ctx, printer.DefaultFindOptions)
	if err != nil {
		return err
	}
	defer printer.Close()

	return printer.Print(context.Background(), img)
}
