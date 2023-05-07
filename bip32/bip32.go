package bip32

// TODO: is it possible to read the C header info from the header file rather than hardcoding it here?

/*
	#include <stdint.h>

	/// derive_c does the same thing as the above function, but is intended for use over the CFFI.
	/// it adds error handling in order to be friendlier to the FFI caller: in case of an error, it
	/// prints the error and returns a null pointer.
	/// note that the caller must free() the returned memory as it's not managed/freed here.
	extern uint8_t *derive_c(const uint8_t *seed, size_t seedlen, const uint8_t *path, size_t pathlen);
*/
import "C"
import (
	"crypto/ed25519"
	"fmt"
	"github.com/spacemeshos/smkeys/common"
	"unsafe"
)

// Bip39SeedLen is the expected length of a BIP39-compatible seed. This is specified in the BIP itself.
const Bip39SeedLen = 64

var (
	badSeedLen = fmt.Errorf("invalid seed length")
)

// Derive wraps the underlying CFFI function. It derives a new keypair from a path and a seed.
func Derive(path string, seed []byte) (key *[ed25519.PrivateKeySize]byte, err error) {
	pathLen := len(path)
	seedLen := len(seed)

	// empty path and empty seed will both cause upstream errors, go ahead and catch them here.
	// note: we don't attempt to actually parse the path here. if it contains less than two elements
	// that will also cause upstream errors.
	if pathLen < 1 {
		return nil, common.PathErr
	}
	if seedLen != Bip39SeedLen {
		return nil, badSeedLen
	}

	// Convert Go string to C-compatible byte array
	pathBytes := []byte(path)

	// Pass the string to Rust
	arrayPtr := C.derive_c(
		(*C.uchar)(unsafe.Pointer(&seed[0])),
		C.size_t(seedLen),
		(*C.uchar)(unsafe.Pointer(&pathBytes[0])),
		C.size_t(pathLen),
	)
	if arrayPtr == nil {
		return nil, common.PointerErr
	}
	defer common.FreeCPointer(common.CUChar(arrayPtr))

	// Convert the *mut u8 pointer to a Go byte slice
	bytes := (*[ed25519.PrivateKeySize]byte)(unsafe.Pointer(arrayPtr))[:]
	key = new([ed25519.PrivateKeySize]byte)
	bytesCopied := copy(key[:], bytes)
	if bytesCopied != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("error in key length")
	}
	return
}
