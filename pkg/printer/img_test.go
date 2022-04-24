package printer

import (
	"bytes"
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "image/png"
)

func TestEncodeRunLenRepetition(t *testing.T) {
	assert.Equal(t, unquote(`\x8a`), bs(encodeRunLenRepetition(10, true)))
	assert.Equal(t, unquote(`d`), bs(encodeRunLenRepetition(100, false)))
	assert.Equal(t, unquote(`\xff`), bs(encodeRunLenRepetition(127, true)))
	assert.Equal(t, unquote(`\x7fI`), bs(encodeRunLenRepetition(200, false)))
	assert.Equal(t, unquote(`\xff\xfb`), bs(encodeRunLenRepetition(250, true)))
}

func TestRunLengthEncode(t *testing.T) {
	assert.Equal(t, unquote(`\x04\x89\x01\x81\x02\x82`), bs(runLengthEncode([]bool{false, false, false, false, true, true, true, true, true, true, true, true, true, false, true, false, false, true, true})))
}

func TestByteEncode(t *testing.T) {
	assert.Equal(t, unquote(`\xf0_`), bs(byteEncode([]bool{false, false, false, false, true, true, true, true, true, true, true, true, true, false, true, false})))
}

func TestImgToBitmap(t *testing.T) {
	f, err := os.Open("testdata/test.png")
	assert.NoError(t, err)
	defer f.Close()

	imgData, _, err := image.Decode(f)
	assert.NoError(t, err)

	bmp, err := imgToBitmap(imgData)
	assert.NoError(t, err)
	mat := bytes.Buffer{}
	for _, row := range bmp {
		for _, v := range row {
			if v {
				mat.WriteString("1 ")
			} else {
				mat.WriteString("0 ")
			}
		}
		mat.WriteString("\n")
	}

	// dumped from Python version
	truth, err := os.ReadFile("testdata/test.mat")
	assert.NoError(t, err)

	assert.Equal(t, truth, mat.Bytes())
}

func TestCmdsImg(t *testing.T) {
	f, err := os.Open("testdata/test.png")
	assert.NoError(t, err)
	defer f.Close()

	imgData, _, err := image.Decode(f)
	assert.NoError(t, err)

	cmds, err := cmdsImg(imgData)
	assert.NoError(t, err)

	// dumped from Python version
	truth, err := os.ReadFile("testdata/test.cmds.bin")
	assert.NoError(t, err)
	assert.Equal(t, truth, cmds)
}
