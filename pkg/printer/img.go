package printer

import (
	"errors"
	"image"
	"image/draw"
)

const (
	PrintWidth = 384

	feedPaper = 15
)

func encodeRunLenRepetition(n uint, val bool) []int8 {
	v := int8(0)
	if val {
		v = 1
	}

	ret := []int8{}
	for n > 0x7f {
		ret = append(ret, 0x7f|(v<<7))
		n -= 0x7f
	}

	if n > 0 {
		ret = append(ret, (v<<7)|int8(n))
	}

	return ret
}

func runLengthEncode(row []bool) []int8 {
	ret := []int8{}
	count := uint(0)
	lastVal := false

	for i, val := range row {
		if i != 0 && val == lastVal {
			count++
		} else {
			ret = append(ret, encodeRunLenRepetition(count, lastVal)...)
			count = 1
		}
		lastVal = val
	}

	if count > 0 {
		ret = append(ret, encodeRunLenRepetition(count, lastVal)...)
	}
	return ret
}

func byteEncode(row []bool) []int8 {
	bitEncode := func(chunkStart, bitIx int) int8 {
		if row[chunkStart+bitIx] {
			return 1 << bitIx
		} else {
			return 0
		}
	}

	ret := []int8{}
	for chunkStart := 0; chunkStart < len(row); chunkStart += 8 {
		var b int8 = 0
		for bitIx := 0; bitIx < 8; bitIx++ {
			b |= bitEncode(chunkStart, bitIx)
		}
		ret = append(ret, b)
	}

	return ret
}

func cmdPrintRow(row []bool) []byte {
	// try to use run-length compression on the image data.
	encoded := runLengthEncode(row)

	// if the resulting compression takes more than PRINT_WIDTH // 8, it means
	// it's not worth it. So we fallback to a simpler, fixed-length encoding.
	if len(encoded) > PrintWidth/8 {
		encoded = byteEncode(row)

		ret := []int8{
			81,
			120,
			-94,
			0,
			int8(len(encoded)),
			0,
		}
		ret = append(ret, encoded...)
		ret = append(ret, 0)
		ret = append(ret, -1)

		retConv := bs(ret)
		retConv[len(retConv)-2] = checkSum(retConv, 6, len(encoded))
		return retConv
	}

	ret := []int8{
		81,
		120,
		-65,
		0,
		int8(len(encoded)),
		0,
	}
	ret = append(ret, encoded...)
	ret = append(ret, 0)
	ret = append(ret, -1)
	retConv := bs(ret)
	retConv[len(retConv)-2] = checkSum(retConv, 6, len(encoded))
	return retConv
}

func imgToBitmap(img image.Image) ([][]bool, error) {
	if img.Bounds().Dx() != PrintWidth {
		return nil, errors.New("invalid size")
	}

	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, gray.Bounds(), img, img.Bounds().Min, draw.Src)

	ret := make([][]bool, gray.Bounds().Dy())
	for y := 0; y < gray.Bounds().Dy(); y++ {
		row := make([]bool, gray.Bounds().Dx())
		for x := 0; x < gray.Bounds().Dx(); x++ {
			row[x] = gray.GrayAt(x, y).Y <= 127 // invert
		}

		ret[y] = row
	}

	return ret, nil
}

func cmdsImg(img image.Image) ([]byte, error) {
	bmp, err := imgToBitmap(img)
	if err != nil {
		return nil, err
	}

	ret := []byte{}
	for _, row := range bmp {
		ret = append(ret, cmdPrintRow(row)...)
	}

	return ret, nil
}

func cmdsPrint(img image.Image, textMode bool) ([]byte, error) {
	imgCmds, err := cmdsImg(img)
	if err != nil {
		return nil, err
	}

	ret := []byte{}

	ret = append(ret, cmdGetDevState...)
	ret = append(ret, cmdSetQuality200DPI...)
	if textMode {
		ret = append(ret, cmdPrintText...)
	} else {
		ret = append(ret, cmdPrintImg...)
	}
	ret = append(ret, cmdLatticeStart...)

	ret = append(ret, imgCmds...)

	ret = append(ret, cmdFeedPaper(feedPaper)...)
	ret = append(ret, cmdSetPaper...)
	ret = append(ret, cmdSetPaper...)
	ret = append(ret, cmdSetPaper...)
	ret = append(ret, cmdLatticeEnd...)
	ret = append(ret, cmdGetDevState...)

	return ret, nil
}
