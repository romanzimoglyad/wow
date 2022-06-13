package pow

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestCheckValidHash(t *testing.T) {
	assert.False(t, CheckValidHash("asdasd", 2))
	assert.True(t, CheckValidHash("00008877f1f3f71b02c560beb02339d39884281eea156fa01a86d012442c21bf", 4))
	assert.True(t, CheckValidHash("00123", 2))
}

func TestDoWork(t *testing.T) {
	hb := HashBlock{
		Ver:      1,
		Bits:     4,
		Date:     time.Time{}.Unix(),
		Resource: "client1",
		Rand:     "rnd",
		Counter:  0,
	}
	err := hb.DoWork(100)
	require.Error(t, err)
	err = hb.DoWork(1000000000)
	require.NoError(t, err)
	assert.Equal(t, 62461, hb.Counter)
}
