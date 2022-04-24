# Gocatprint

Go library to print images on Aliexpress Cat thermo paper printers.
A test binary is included as well.

```
$ make build
$ ./catprint --help

Usage: catprint [--log-pretty] [--log-level LEVEL] [--hci-device N] [--timeout DUR] [--printer-name NAME] [--printer-address ADDR] [--dark-mode] IMG

Positional arguments:
  IMG                    image to print, PNG or JPEG, must be 384 px wide

Options:
  --log-pretty           log pretty [default: true]
  --log-level LEVEL      log level [default: info]
  --hci-device N         HCI device to use [default: 0]
  --timeout DUR          how long to allow for discovery and printing [default: 10s]
  --printer-name NAME    device name to connect to, ignored if empty
  --printer-address ADDR
                         device address to connect to, ignored if empty
  --dark-mode            more contrast, slower speed [default: true]
  --help, -h             display this help and exit
```

Ported to Go from <https://github.com/rbaron/catprinter>.
