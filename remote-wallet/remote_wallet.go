package remote_wallet

/*
	#include "remote_wallet.h"
*/
import "C"
import (
	"crypto/ed25519"
	"fmt"
	"unsafe"
)

func ReadPubkeyFromLedger(path, derivationPath string, confirmKey bool) (key []byte, err error) {
	// Allocate a buffer for the output
	cKey := C.calloc(C.size_t(ed25519.PublicKeySize), 1)
	defer C.free(cKey)
	pathStr := C.CString(path)
	defer C.free(unsafe.Pointer(pathStr))
	derivationPathStr := C.CString(derivationPath)
	defer C.free(unsafe.Pointer(derivationPathStr))
	status := C.read_pubkey_from_ledger(
		pathStr,
		derivationPathStr,
		C.bool(confirmKey),
		(*C.uchar)(cKey),
	)
	if status != 0 {
		return nil, fmt.Errorf("error reading pubkey from ledger: %d", uint32(status))
	}
	key = C.GoBytes(cKey, C.int(ed25519.PublicKeySize))

	return
}
