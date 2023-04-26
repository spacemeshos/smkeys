package bip32

// TODO: is it possible to read the C header info from the header file rather than hardcoding it here?

/*
	#cgo CFLAGS: -I../deps
	#cgo LDFLAGS: -L../deps -led25519_bip32
	#include <stdint.h>

	/// derive_c does the same thing as the above function, but is intended for use over the CFFI.
	/// it adds error handling in order to be friendlier to the FFI caller: in case of an error, it
	/// prints the error and returns a null pointer.
	/// note that the caller must free() the returned memory as it's not managed/freed here.
	extern uint8_t *derive_c(const uint8_t *seed, size_t seedlen, const uint8_t *path, size_t pathlen);

	/// free the memory allocated and returned by the derive functions by transferring ownership back to
	/// Rust. must be called on each pointer returned by the functions precisely once to ensure safety.
	extern void derive_free_c(uint8_t *ptr);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// ArrayLen is the expected length of the return value in bytes: the first 32 bytes are the public key
// and the second 32 bytes are the private key.
// TODO: can we read this from somewhere?
const ArrayLen = 64

// Bip39SeedLen is the expected length of a BIP39-compatible seed. This is specified in the BIP itself.
const Bip39SeedLen = 64

var (
	pathErr    = fmt.Errorf("empty path supplied")
	badSeedLen = fmt.Errorf("invalid seed length")
	pointerErr = fmt.Errorf("error in wrapped function, got nil pointer")
)

// Derive wraps the underlying CFFI function. It derives a new keypair from a path and a seed.
func Derive(path string, seed []byte) (key *[ArrayLen]byte, err error) {
	pathLen := len(path)
	seedLen := len(seed)

	// empty path and empty seed will both cause upstream errors, go ahead and catch them here.
	// note: we don't attempt to actually parse the path here. if it contains less than two elements
	// that will also cause upstream errors.
	if pathLen < 1 {
		return nil, pathErr
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
		return nil, pointerErr
	}
	defer C.derive_free_c(arrayPtr)

	// Convert the *mut u8 pointer to a Go byte slice
	bytes := (*[ArrayLen]byte)(unsafe.Pointer(arrayPtr))[:]
	key = new([ArrayLen]byte)
	bytesCopied := copy(key[:], bytes)
	if bytesCopied != ArrayLen {
		return nil, fmt.Errorf("error in key length")
	}
	return
}
