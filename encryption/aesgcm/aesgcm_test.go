package aesgcm_test

import (
	"testing"

	"github.com/deweppro/go-sdk/encryption/aesgcm"
	"github.com/deweppro/go-sdk/random"
	"github.com/stretchr/testify/require"
)

func TestUnit_Codec(t *testing.T) {
	rndKey := random.String(32)
	message := []byte("Hello World!")

	c, err := aesgcm.New(rndKey)
	require.NoError(t, err)

	enc1, err := c.Encrypt(message)
	require.NoError(t, err)

	dec1, err := c.Decrypt(enc1)
	require.NoError(t, err)

	require.Equal(t, message, dec1)

	c, err = aesgcm.New(rndKey)
	require.NoError(t, err)

	enc2, err := c.Encrypt(message)
	require.NoError(t, err)

	require.NotEqual(t, enc1, enc2)

	dec2, err := c.Decrypt(enc1)
	require.NoError(t, err)

	require.Equal(t, message, dec2)

}
