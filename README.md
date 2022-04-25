# Gocatprint

![Demo](demo.gif)

Go library to print images on cheap thermo paper printers from Aliexpress.
A CLI tool `catprint` is available as well.

Install it via `go install github.com/jo-m/gocatprint/cmd/catprint@latest`.

There is a simple API in `pkg/simple` and a more advanced one in `pkg/printer`.

Ported to Go from <https://github.com/rbaron/catprinter>.

```
$ make build
$ ./catprint --help
Usage: catprint [--log-pretty] [--log-level LEVEL] [--hci-device N] [--timeout DUR] [--printer-name NAME] [--printer-address ADDR] [--fast-mode] [--threshold] [--preview OUT-FILE] IN-FILE

Positional arguments:
  IN-FILE                image to print, PNG or JPEG, must be 384px wide (unless --scale is passed)

Options:
  --log-pretty           log pretty [default: true]
  --log-level LEVEL      log level [default: info]
  --hci-device N         HCI device to use, auto if negative [default: -1]
  --timeout DUR          how long to allow for discovery and printing [default: 10s]
  --printer-name NAME    device name to connect to, ignored if empty
  --printer-address ADDR
                         device address to connect to, ignored if empty
  --fast-mode            less contrast, higher printer speed [default: false]
  --threshold            use simple thresholding instead of dithering [default: false]
  --preview OUT-FILE     do not print, just write the (processed) image to the given file
  --help, -h             display this help and exit
```
