package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"time"

	"github.com/aamcrae/webcam"
)

// CamConfig is the v4l camera config. It contains struct tags compatible with github.com/alexflint/go-arg.
type CamConfig struct {
	PicDev string `arg:"env:PIC_DEV,--pic-dev" default:"/dev/video2" help:"camera video device file path" placeholder:"DEV"`
	// v4l2-ctl --list-formats-ext --device /dev/video2
	PicFormat     string        `arg:"env:PIC_FORMAT,--pic-format" default:"Motion-JPEG" help:"camera preferred image format" placeholder:"STR"`
	PicTimeout    time.Duration `arg:"env:PIC_TIMEOUT,--pic-timeout" default:"2s" help:"how long to give camera time to start" placeholder:"DUR"`
	PicSkipFrames int           `arg:"env:PIC_SKIP_FRAMES,--pic-skip-frames" default:"15" help:"how many frames to skip until picture snap" placeholder:"N"`
}

func chooseFormat(cam *webcam.Webcam, preferred string) (webcam.PixelFormat, string, *webcam.FrameSize, error) {
	fmap := cam.GetSupportedFormats()
	var format webcam.PixelFormat
	var formatStr string
	for f, s := range fmap {
		format = f
		formatStr = s
		if s == preferred {
			break
		}
	}

	if format == 0 {
		return 0, "", nil, errors.New("no format found")
	}

	frameSizes := cam.GetSupportedFrameSizes(format)
	if len(frameSizes) == 0 {
		return 0, "", nil, errors.New("no frame size found")
	}

	return format, formatStr, &frameSizes[0], nil
}

func convertFrame(frame []byte, formatStr string, w, h int) (image.Image, error) {
	switch formatStr {
	case "Motion-JPEG":
		b := bytes.NewBuffer(frame)
		return jpeg.Decode(b)
	default:
		return nil, fmt.Errorf(`unhandled format string "%s"`, formatStr)
	}
}

func Snap(config CamConfig) (image.Image, error) {
	cam, err := webcam.Open(config.PicDev)
	if err != nil {
		return nil, err
	}
	defer cam.Close()

	format, formatStr, size, err := chooseFormat(cam, config.PicFormat)
	if err != nil {
		return nil, err
	}

	f, w, h, _, _, err := cam.SetImageFormat(format, uint32(size.MaxWidth), uint32(size.MaxHeight))
	if err != nil {
		return nil, err
	}
	if f != format {
		return nil, errors.New("unable to choose format")
	}

	err = cam.StartStreaming()
	if err != nil {
		return nil, err
	}

	timeout := uint32(config.PicTimeout.Seconds())
	var frame []byte
	for i := 0; i < config.PicSkipFrames; i++ {
		err = cam.WaitForFrame(timeout)

		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			return nil, errors.New("camera timed out")
		default:
			return nil, err
		}

		frame, err = cam.ReadFrame()
		if err != nil {
			return nil, err
		}
		if len(frame) == 0 {
			return nil, errors.New("received empty frame")
		}
	}

	return convertFrame(frame, formatStr, int(w), int(h))
}
