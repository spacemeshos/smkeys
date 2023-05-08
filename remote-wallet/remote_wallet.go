package remote_wallet

// TODO(mafa): consider using a header instead

/*
	#include <stdbool.h>
	#include <stdint.h>
	#include <stdlib.h>

	/// read_pubkey_from_ledger reads a pubkey from the ledger device specified by path and
	/// derivation_path. If path is empty, the first ledger device found will be used. If confirm_key
	/// is true, it will prompt the user to confirm the key on the device. It returns
	/// a pointer to the pubkey bytes on success, or null on failure. Note that the caller must free
	/// the returned pointer by passing it back to Rust using sdkutils.free().
	extern uint8_t read_pubkey_from_ledger(const char *path, const char *derivation_path, bool confirm_key, const uint8_t *out);

	// TODO(mafa): having derive_c return an error code (or 0 if OK) and passing in the output buffer (we know the size on both sides)
	// makes handling the FFI on both sides easier
*/
import "C"

import (
	"crypto/ed25519"
	"fmt"
	"unsafe"
)

var (
	ErrUnknown = fmt.Errorf("unknown error")
)

func ReadPubkeyFromLedger(path, derivationPath string, confirmKey bool) ([]byte, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	cDerivationPath := C.CString(derivationPath)
	defer C.free(unsafe.Pointer(cDerivationPath))

	cOut := (C.calloc(ed25519.PublicKeySize, 1))
	defer C.free(cOut)

	retVal := C.read_pubkey_from_ledger(cPath, cDerivationPath, C.bool(confirmKey), (*C.uchar)(cOut))
	switch retVal { // TODO(mafa): error code definitions should be in the header file
	case 0:
		// success
	default:
		return nil, ErrUnknown
	}

	// Convert the C pointer to a Go byte slice
	output := C.GoBytes(cOut, C.int(ed25519.PublicKeySize))
	return output, nil
}
