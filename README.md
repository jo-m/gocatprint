# Gocatprint

![Demo](demo.gif)

Go library to print images on cheap thermo paper printers from Aliexpress.
A CLI tool `catprint` is available as well.

Install it via `go install github.com/jo-m/gocatprint/cmd/catprint@latest`.

There is a simple API in `pkg/simple` and a more advanced one in `pkg/printer`.

See the binaries in `cmd/` for usage.

Ported to Go from <https://github.com/rbaron/catprinter>.

Needs root to run (for Bluetooth), there is probably a way to fix this via Linux system permissions.

```
$ make build
$ ./catprint --help
Usage: catprint [--log-pretty] [--log-level LEVEL] [--hci-device N] [--timeout DUR] [--printer-name NAME] [--printer-address ADDR] [--threshold] [--preview OUT-FILE] IN-FILE

Positional arguments:
  IN-FILE                image to print, PNG or JPEG, must be 384px wide (unless --scale is passed)

Options:
  --log-pretty           log pretty [default: true, env: LOG_PRETTY]
  --log-level LEVEL      log level [default: info, env: LOG_LEVEL]
  --hci-device N         HCI device to use, auto if negative [default: -1]
  --timeout DUR          how long to allow for discovery and printing [default: 10s]
  --printer-name NAME    device name to connect to, ignored if empty
  --printer-address ADDR
                         device address to connect to, ignored if empty
  --threshold            use simple thresholding instead of dithering [default: false]
  --preview OUT-FILE     do not print, just write the (processed) image to the given file
  --help, -h             display this help and exit

$ sudo ./catprint ~/image.png
```

## Webcam example

```
$ make build
$ ./catprintcam --help
Usage: catprintcam [--log-pretty] [--log-level LEVEL] [--pic-dev DEV] [--pic-format STR] [--pic-timeout DUR] [--pic-skip-frames N]

Options:
  --log-pretty           log pretty [default: true, env: LOG_PRETTY]
  --log-level LEVEL      log level [default: info, env: LOG_LEVEL]
  --pic-dev DEV          camera video device file path [default: /dev/video2, env: PIC_DEV]
  --pic-format STR       camera preferred image format [default: Motion-JPEG, env: PIC_FORMAT]
  --pic-timeout DUR      how long to give camera time to start [default: 2s, env: PIC_TIMEOUT]
  --pic-skip-frames N    how many frames to skip until picture snap [default: 15, env: PIC_SKIP_FRAMES]
  --help, -h             display this help and exit

$ sudo ./catprintcam --pic-dev /dev/video0
```
