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
	keyPtr := (*C.uchar)(unsafe.Pointer(&key[0]))
	pathPtr := (*C.uchar)(unsafe.Pointer(&[]byte(path)[0]))
	pathLen := (C.size_t)(len(path))
	derivationPathPtr := (*C.uchar)(unsafe.Pointer(&[]byte(derivationPath)[0]))
	derivationPathLen := (C.size_t)(len(derivationPath))
	status := C.read_pubkey_from_ledger(
		pathPtr,
		pathLen,
		derivationPathPtr,
		derivationPathLen,
		C.bool(confirmKey),
		keyPtr,
	)
	if status != 0 {
		return nil, fmt.Errorf("error reading pubkey from ledger: %d", uint32(status))
	}

	return
}
