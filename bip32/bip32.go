package bip32

// TODO: is it possible to read the C header info from the header file rather than hardcoding it here?
// TODO (mafa): yes, we do that for libgpu-post and for libpostrs

/*
	#include <stdint.h>
	#include <stdlib.h>

	/// derive_c does the same thing as the above function, but is intended for use over the CFFI. - TODO(mafa): needs better description
	/// it adds error handling in order to be friendlier to the FFI caller: in case of an error, it
	/// prints the error and returns a null pointer. - TODO(mafa): imo it's easier to handle with a return error code
	extern uint8_t derive_c(const uint8_t *seed, const char *path, uint8_t *out);

	// TODO(mafa): having derive_c return an error code (or 0 if OK) and passing in the output buffer (we know the size on both sides)
	// makes handling the FFI on both sides easier
*/
import "C"

import (
	"crypto/ed25519"
	"fmt"
	"unsafe"
)

// Bip39SeedLen is the expected length of a BIP39-compatible seed. This is specified in the BIP itself.
const Bip39SeedLen = 64

var (
	ErrInvalidPath     = fmt.Errorf("invalid path")
	ErrBadSeed         = fmt.Errorf("invalid seed length")
	ErrNonHardenedPath = fmt.Errorf("non-hardened path not supported")
	ErrUnknown         = fmt.Errorf("unknown error")
)

// Derive wraps the underlying CFFI function. It derives a new keypair from a path and a seed.
func Derive(path string, seed []byte) ([]byte, error) {
	// convert Go string to C-compatible byte array
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	cSeed := C.CBytes(seed)
	defer C.free(cSeed)

	// Allocate a buffer for the output
	cOut := (C.calloc(ed25519.PrivateKeySize, 1))
	defer C.free(cOut)

	// Pass the string to Rust
	retVal := C.derive_c((*C.uchar)(cSeed), cPath, (*C.uchar)(cOut))
	switch retVal { // TODO(mafa): error code definitions should be in the header file
	case 0:
		// all good
	case 1:
		return nil, ErrInvalidPath
	case 2:
		return nil, ErrBadSeed
	case 3:
		return nil, ErrNonHardenedPath
	default:
		return nil, ErrUnknown
	}

	// Convert the *mut u8 pointer to a Go byte slice
	output := C.GoBytes(cOut, C.int(ed25519.PrivateKeySize))
	return output, nil
}
