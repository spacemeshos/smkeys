package remote_wallet

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadPubkeyFromLedger(t *testing.T) {
	derivationPath := "m/44'/540'/0'/0'/0'"
	path := "usb://ledger"
	key, err := ReadPubkeyFromLedger(path, derivationPath)
	require.NoError(t, err)
	require.NotNil(t, key)
	t.Log("Got key: ", key)
}
