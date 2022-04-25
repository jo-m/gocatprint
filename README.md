# Gocatprint

![Demo](demo.gif)

Go library to print images on cheap thermo paper printers from Aliexpress.
A CLI tool `catprint` is available as well.

Install it via `go install github.com/jo-m/gocatprint/cmd/catprint@latest`.

Ported to Go from <https://github.com/rbaron/catprinter>.

```
$ make build
$ ./catprint --help

Usage: catprint [--log-pretty] [--log-level LEVEL] [--hci-device N] [--timeout DUR] [--printer-name NAME] [--printer-address ADDR] [--fast-mode] [--threshold] [--preview IMG-FILE] IMG-FILE

Positional arguments:
  IMG-FILE               image to print, PNG or JPEG, must be 384px wide (unless --scale is passed)

Options:
  --log-pretty           log pretty [default: true]
  --log-level LEVEL      log level [default: info]
  --hci-device N         HCI device to use [default: 0]
  --timeout DUR          how long to allow for discovery and printing [default: 10s]
  --printer-name NAME    device name to connect to, ignored if empty
  --printer-address ADDR
                         device address to connect to, ignored if empty
  --fast-mode            less contrast, higher printer speed [default: false]
  --threshold            use simple thresholding instead of dithering [default: false]
  --preview IMG-FILE     do not print, just write the (processed) image to the given file
  --help, -h             display this help and exit
```
