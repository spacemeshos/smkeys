package remote_wallet

import "testing"

func TestReadPubkeyFromLedger(t *testing.T) {
	derivationPath := "m/44'/540'/0'/0'/0'"
	path := "usb://ledger"
	key, err := ReadPubkeyFromLedger(path, derivationPath)
	if err != nil {
		t.Errorf("Error in ReadPubkeyFromLedger")
	}
	if key == nil {
		t.Errorf("Error in ReadPubkeyFromLedger")
	}
	t.Log("Got key: ", key)
}
