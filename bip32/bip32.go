package bip32

/*
	#include "ed25519_bip32.h"
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
	badSeedLen = fmt.Errorf("invalid seed length")
	ffiErr     = fmt.Errorf("error deriving key")
	pathErr    = fmt.Errorf("invalid path")
)

// Derive wraps the underlying CFFI function. It derives a new keypair from a path and a seed.
func Derive(path string, seed []byte) (key []byte, err error) {
	// Allocate a buffer for the output
	cKey := C.calloc(C.size_t(ed25519.PrivateKeySize), 1)
	defer C.free(cKey)
	pathStr := C.CString(path)
	defer C.free(unsafe.Pointer(pathStr))
	seedLen := len(seed)

	// empty path and empty seed will both cause upstream errors, go ahead and catch them here.
	// note: we don't attempt to actually parse the path here. if it contains less than two elements
	// that will also cause upstream errors.
	if len(path) < 1 {
		return nil, pathErr
	}
	if seedLen != Bip39SeedLen {
		return nil, badSeedLen
	}

	// Pass the string to Rust
	status := C.derive_c(
		(*C.uchar)(unsafe.Pointer(&seed[0])),
		C.size_t(seedLen),
		pathStr,
		(*C.uchar)(cKey),
	)
	if status != 0 {
		return nil, ffiErr
	}
	key = C.GoBytes(cKey, C.int(ed25519.PrivateKeySize))
	return
}
