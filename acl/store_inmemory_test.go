package acl_test

import (
	"testing"

	"github.com/deweppro/go-sdk/acl"

	"github.com/stretchr/testify/require"
)

func TestUnit_NewInMemoryStorage(t *testing.T) {
	opt := acl.OptionInMemoryStorageSetupData(map[string]string{
		"u1": "123",
		"u2": "456",
	})
	store := acl.NewInMemoryStorage(opt)
	require.NotNil(t, store)

	val, err := store.FindACL("u1")
	require.NoError(t, err)
	require.Equal(t, "123", val)

	val, err = store.FindACL("u2")
	require.NoError(t, err)
	require.Equal(t, "456", val)

	val, err = store.FindACL("u3")
	require.Error(t, err)
	require.Equal(t, "", val)

	err = store.ChangeACL("u2", "789")
	require.NoError(t, err)

	val, err = store.FindACL("u2")
	require.NoError(t, err)
	require.Equal(t, "789", val)

	val, err = store.FindACL("u5")
	require.Error(t, err)
	require.Equal(t, "", val)

	err = store.ChangeACL("u5", "333")
	require.NoError(t, err)

	val, err = store.FindACL("u5")
	require.NoError(t, err)
	require.Equal(t, "333", val)
}
