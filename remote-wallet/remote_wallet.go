package remote_wallet

/*
	#include "remote-wallet.h"
*/
import "C"
import (
	"crypto/ed25519"
	"fmt"
	"unsafe"
)

func ReadPubkeyFromLedger(path, derivationPath string, confirmKey bool) (key *[ed25519.PublicKeySize]byte, err error) {
	key = new([ed25519.PublicKeySize]byte)
	pathStr := C.CString(path)
	defer C.free(unsafe.Pointer(pathStr))
	derivationPathStr := C.CString(derivationPath)
	defer C.free(unsafe.Pointer(derivationPathStr))
	status := C.read_pubkey_from_ledger(
		pathStr,
		derivationPathStr,
		C.bool(confirmKey),
		(*C.uchar)(unsafe.Pointer(&key[0])),
	)
	if status != 0 {
		return nil, fmt.Errorf("error reading pubkey from ledger: %d", uint32(status))
	}

	return
}
