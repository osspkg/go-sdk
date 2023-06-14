package jwt_test

import (
	"testing"
	"time"

	"github.com/osspkg/go-sdk/auth/jwt"
	"github.com/stretchr/testify/require"
)

type demoJwtPayload struct {
	ID int `json:"id"`
}

func TestUnit_NewJWT(t *testing.T) {
	conf := make([]jwt.Config, 0)
	conf = append(conf, jwt.Config{ID: "789456", Key: "123456789123456789123456789123456789", Algorithm: jwt.AlgHS256})
	j, err := jwt.New(conf)
	require.NoError(t, err)

	payload1 := demoJwtPayload{ID: 159}
	token, err := j.Sign(&payload1, time.Hour)
	require.NoError(t, err)

	payload2 := demoJwtPayload{}
	head1, err := j.Verify(token, &payload2)
	require.NoError(t, err)

	require.Equal(t, payload1, payload2)

	head2, err := j.Verify(token, &payload2)
	require.NoError(t, err)
	require.Equal(t, head1, head2)
}
