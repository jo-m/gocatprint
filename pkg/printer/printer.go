package printer

import (
	"context"
	"errors"
	"image"
	"strings"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux/gatt"
	"github.com/rs/zerolog/log"
)

const (
	printerServiceScanUUID = 0xAF30
	printerServiceUUID     = 0xAE30

	printerCharacteristicUUID = 0xAE01
)

type FindOptions struct {
	// bluetooth device name, ignored if empty
	Name string
	// HCI address, ignored if empty
	Address string
}

var DefaultFindOptions = FindOptions{}

type Printer struct {
	client ble.Client
	char   *ble.Characteristic
}

func findCharacteristic(client ble.Client) (*ble.Characteristic, error) {
	log.Debug().Msg("discovering services")

	svcs, err := client.DiscoverServices(nil)
	if err != nil {
		return nil, err
	}
	for _, s := range svcs {
		logger := log.With().Str("svcUuid", s.UUID.String()).Logger()
		logger.Debug().Msg("found service")

		if !s.UUID.Equal(ble.UUID16(printerServiceUUID)) {
			logger.Trace().Msg("no match")
			continue
		}

		chars, err := client.DiscoverCharacteristics(nil, s)
		if err != nil {
			return nil, err
		}
		for _, c := range chars {
			logger.Trace().Str("charUuid", c.UUID.String()).Msg("found characteristic")

			if !c.UUID.Equal(ble.UUID16(printerCharacteristicUUID)) {
				logger.Trace().Msg("no match")
				continue
			}

			logger.Debug().Msg("found")
			return c, nil
		}
	}

	return nil, errors.New("characteristic not found")
}

func Find(ctx context.Context, opts FindOptions) (*Printer, error) {
	filter := func(adv ble.Advertisement) bool {
		logger := log.With().
			Str("addr", adv.Addr().String()).
			Int("rssi", adv.RSSI()).
			Str("name", adv.LocalName()).
			Logger()

		logger.Debug().Msg("received advertisement")

		if !adv.Connectable() {
			logger.Trace().Msg("not connectable")
			return false
		}

		if opts.Name != "" {
			if strings.ToUpper(adv.LocalName()) != strings.ToUpper(opts.Name) {
				logger.Trace().Msg("name mismatch")
				return false
			}
		}

		if opts.Address != "" {
			if strings.ToUpper(adv.Addr().String()) != strings.ToUpper(opts.Address) {
				logger.Trace().Msg("address mismatch")
				return false
			}
		}

		for _, s := range adv.Services() {
			logger.Trace().Str("service", s.String()).Msg("found service")

			if s.Equal(ble.UUID16(printerServiceScanUUID)) {
				logger.Debug().Msg("found printer advertisement")

				return true
			}
		}

		return false
	}

	client, err := ble.Connect(ctx, filter)
	if err != nil {
		return nil, err
	}

	log.Info().Str("addr", client.Addr().String()).Msg("connected")

	char, err := findCharacteristic(client)
	if err != nil {
		client.CancelConnection()
		return nil, err
	}

	return &Printer{
		client: client,
		char:   char,
	}, nil
}

func (p *Printer) Close() {
	log.Debug().Msg("Close()")

	p.client.ClearSubscriptions()
	p.client.CancelConnection()

	<-p.client.Disconnected()
	log.Debug().Msg("connection closed")
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

func (p *Printer) Print(ctx context.Context, img image.Image, darkMode bool) error {
	cmds, err := cmdsPrint(img, darkMode)
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

	// fmt.Println(gattC.Addr())
	// fmt.Println(gattC.DiscoverCharacteristics())
	// fmt.Println(gattC.DiscoverProfile(true))

	// return nil

	// TODO subscribe
	// err = gattC.Subscribe(nil, false, func(req []byte) {
	// 	log.Info().Msg("received sub")
	// })
	// if err != nil {
	// 	return err
	// }

	// return nil

	for _, b := range chunks {
		log.Trace().Int("len", len(b)).Msg("writing chunk")
		err = gattC.WriteCharacteristic(p.char, b, true)
		if err != nil {
			return err
		}
	}

	// TODO fix hack
	time.Sleep(time.Second * 5)

	return nil
}
