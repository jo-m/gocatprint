package printer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBs(t *testing.T) {
	signed := []int8{81, 120, -93, 0, 1, 0, 0, 0, -1}
	unsigned := bs(signed)

	assert.Equal(t, []byte{'Q', 'x', 0xa3, 0x00, 0x01, 0x00, 0x00, 0x00, 0xff}, unsigned)
}

func TestCheckSum(t *testing.T) {
	assert.Equal(t, byte(161), checkSum(cmdLatticeStart, 5, 12))
}

func TestUnquote(t *testing.T) {
	assert.Equal(t, []byte{'Q', 'x', 0xbd, 0xff}, unquote(`Qx\xbd\xff`))
}

func TestCmdFeedPaper(t *testing.T) {
	assert.Equal(t, []byte{'Q', 'x', 0xbd, 0x00, 0x01, 0x00, 0x14, 'l', 0xff}, cmdFeedPaper(20))

	assert.Equal(t, unquote(`Qx\xbd\x00\x01\x00\x14l\xff`), cmdFeedPaper(20))
	assert.Equal(t, unquote(`Qx\xbd\x00\x01\x00xo\xff`), cmdFeedPaper(120))
	assert.Equal(t, unquote(`Qx\xbd\x00\x01\x00\xff\xf3\xff`), cmdFeedPaper(-1))
}

func TestCmdSetEnergy(t *testing.T) {
	assert.Equal(t, unquote(`Qx\xaf\x00\x02\x00\x00l\x00\xff`), cmdSetEnergy(20))
	assert.Equal(t, unquote(`Qx\xaf\x00\x02\x00\x00o\x00\xff`), cmdSetEnergy(120))
	assert.Equal(t, unquote(`Qx\xaf\x00\x02\x00\x00v\x00\xff`), cmdSetEnergy(200))
	assert.Equal(t, unquote(`Qx\xaf\x00\x02\x00\xff$\x00\xff`), cmdSetEnergy(-1))
	assert.Equal(t, unquote(`Qx\xaf\x00\x02\x00\xff\n\x00\xff`), cmdSetEnergy(-100))
}
