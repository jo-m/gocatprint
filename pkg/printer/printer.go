package printer

import (
	"context"
	"errors"
	"image"
	"strings"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux/gatt"
)

const (
	printerServiceScanUUID = 0xAF30
	printerServiceUUID     = 0xAE30

	printerCharacteristicUUID = 0xAE01
)

// FindOptions represents options for device discovery.
type FindOptions struct {
	// bluetooth device name, ignored if empty
	Name string
	// HCI address, ignored if empty
	Address string
}

// DefaultFindOptions represents the default discovery options.
var DefaultFindOptions = FindOptions{}

type Printer struct {
	client  ble.Client
	profile *ble.Profile

	printerSvc  *ble.Service
	printerChar *ble.Characteristic
}

// Find finds a BLE printer and connects to it.
// Use ctx for timeout.
// Close() after usage.
func Find(ctx context.Context, opts FindOptions) (*Printer, error) {
	filter := func(adv ble.Advertisement) bool {
		if !adv.Connectable() {
			return false
		}

		if opts.Name != "" {
			if !strings.EqualFold(adv.LocalName(), opts.Name) {
				return false
			}
		}

		if opts.Address != "" {
			if !strings.EqualFold(adv.Addr().String(), opts.Address) {
				return false
			}
		}

		for _, s := range adv.Services() {

			if s.Equal(ble.UUID16(printerServiceScanUUID)) {
				return true
			}
		}

		return false
	}

	client, err := ble.Connect(ctx, filter)
	if err != nil {
		return nil, err
	}

	profile, err := client.DiscoverProfile(true)
	if err != nil {
		client.CancelConnection()
		return nil, err
	}

	printerSvc := profile.FindService(ble.NewService(ble.UUID16(printerServiceUUID)))
	if printerSvc == nil {
		client.CancelConnection()
		return nil, errors.New("printer service not found")
	}

	printerChar := profile.FindCharacteristic(ble.NewCharacteristic(ble.UUID16(printerCharacteristicUUID)))
	if printerChar == nil {
		client.CancelConnection()
		return nil, errors.New("printer characteristic not found")
	}

	return &Printer{
		client:  client,
		profile: profile,

		printerSvc:  printerSvc,
		printerChar: printerChar,
	}, nil
}

// Close closes the connection.
func (p *Printer) Close() {
	p.client.ClearSubscriptions()
	p.client.CancelConnection()

	<-p.client.Disconnected()
}

// https://github.com/golang/go/wiki/SliceTricks#batching-with-minimal-allocation
func chunkifyBytes(b []byte, sz int) [][]byte {
	chunks := make([][]byte, 0, (len(b)+sz-1)/sz)

	for sz < len(b) {
		b, chunks = b[sz:], append(chunks, b[0:sz:sz])
	}
	chunks = append(chunks, b)

	return chunks
}

// Print prints an image.
// You may pass it through PrepareImage() beforehand.
func (p *Printer) Print(ctx context.Context, img image.Image) error {
	cmds, err := cmdsPrint(img)
	if err != nil {
		return err
	}

	chunkSz := p.client.Conn().RxMTU() - 3
	chunks := chunkifyBytes(cmds, chunkSz)

	gattC, err := gatt.NewClient(p.client.Conn())
	if err != nil {
		return err
	}
	defer gattC.CancelConnection()
	gattC.Conn().SetContext(ctx)

	for _, b := range chunks {
		err = gattC.WriteCharacteristic(p.printerChar, b, true)
		if err != nil {
			return err
		}
	}

	// TODO hack: try to sleep for as long as printing goes on
	time.Sleep(time.Millisecond*17*time.Duration(img.Bounds().Dy()) + time.Second)

	return nil
}
