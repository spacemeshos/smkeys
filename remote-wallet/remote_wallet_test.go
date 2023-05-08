package remote_wallet

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadPubkeyFromLedger(t *testing.T) {
	derivationPath := "m/44'/540'/0'/0'/0'"
	path := "usb://ledger"
	key, err := ReadPubkeyFromLedger(path, derivationPath, true)
	require.NoError(t, err, "Got error attempting to read pubkey from ledger, "+
		"please make sure the device is connected and unlocked and the correct app is open")
	require.NotNil(t, key)
	t.Log("Got key: ", key)
}
