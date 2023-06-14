package acl_test

import (
	"testing"

	acl2 "github.com/osspkg/go-sdk/acl"
	"github.com/stretchr/testify/require"
)

func TestUnit_NewACL(t *testing.T) {
	store := acl2.NewInMemoryStorage()
	acl := acl2.NewACL(store, 3)

	email := "demo@example.com"

	t.Log("user not exist")

	levels, err := acl.GetAll(email)
	require.Error(t, err)
	require.Nil(t, levels)

	require.Error(t, acl.Set(email, 10, 1))

	t.Log("user exist")

	require.NoError(t, store.ChangeACL(email, ""))

	require.Error(t, acl.Set(email, 10, 1))

	levels, err = acl.GetAll(email)
	require.NoError(t, err)
	require.Equal(t, []uint8{0, 0, 0}, levels)

	require.NoError(t, acl.Set(email, 2, 10))

	levels, err = acl.GetAll(email)
	require.NoError(t, err)
	require.Equal(t, []uint8{0, 0, 9}, levels)
}
