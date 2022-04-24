package printer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkifyBytes(t *testing.T) {
	assert.Equal(t,
		[][]byte{[]byte{1, 2, 3}, []byte{4, 5, 6}, []byte{7, 8, 9}, []byte{10, 11}},
		chunkifyBytes([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, 3))
}
