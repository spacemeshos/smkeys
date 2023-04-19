package bip32

// TODO: is it possible to read the C header info from the header file rather than hardcoding it here?

/*
	#cgo CFLAGS: -I..
	#cgo LDFLAGS: -L.. -led25519_bip32
	#include <stdint.h>

	/// derive_c does the same thing as the above function, but is intended for use over the CFFI.
	/// it adds error handling in order to be friendlier to the FFI caller: in case of an error, it
	/// prints the error and returns a null pointer.
	/// note that the caller must free() the returned memory as it's not managed/freed here.
	extern uint8_t *derive_c(const uint8_t *seed, size_t seedlen, const uint8_t *path, size_t pathlen);

	/// derive_child_c derives a new child key from a seed and a single hardened path element.
	/// the childidx always refers to a hardened path element, as we do not support non-hardened paths.
	/// note that the caller must free() the returned memory as it's not managed/freed here.
	extern uint8_t *derive_child_c(const uint8_t *seed, size_t seedlen, uint32_t childidx);

	/// free the memory allocated and returned by the derive functions by transferring ownership back to
	/// Rust. must be called on each pointer returned by the functions precisely once to ensure safety.
	extern void derive_free_c(uint8_t *ptr);

	/// from_seed_c derives a new extended secret key from a seed
	extern uint8_t *from_seed_c(const uint8_t *seed, size_t seedlen);
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

var (
	pathErr    = fmt.Errorf("empty path supplied")
	seedErr    = fmt.Errorf("empty seed supplied")
	pointerErr = fmt.Errorf("error in wrapped function, got nil pointer")
)

// FromSeed wraps the underlying CFFI function. It derives a new keypair from a seed.
func FromSeed(seed []byte) (key *[ArrayLen]byte, err error) {
	seedLen := len(seed)

	// empty seed will both cause downstream errors, go ahead and catch it here.
	// TODO: do we want to allow empty seed?
	if seedLen < 1 {
		return nil, pathErr
	}

	// Pass the string to Rust
	arrayPtr := C.from_seed_c(
		(*C.uchar)(unsafe.Pointer(&seed[0])),
		C.size_t(seedLen),
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

// Derive wraps the underlying CFFI function. It derives a new keypair from a path and a seed.
func Derive(path string, seed []byte) (key *[ArrayLen]byte, err error) {
	pathLen := len(path)
	seedLen := len(seed)

	// empty path and empty seed will both cause downstream errors, go ahead and catch them here.
	// TODO: do we want to allow empty seed?
	if pathLen < 1 {
		return nil, pathErr
	}
	if seedLen < 1 {
		return nil, seedErr
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

// DeriveChild wraps the underlying CFFI function. It derives a new keypair from a seed and a single path component.
func DeriveChild(seed []byte, childIdx uint32) (key *[ArrayLen]byte, err error) {
	seedLen := len(seed)

	// empty seed will both cause downstream errors, go ahead and catch it here.
	// TODO: do we want to allow empty seed?
	if seedLen < 1 {
		return nil, pathErr
	}

	// Pass the string to Rust
	arrayPtr := C.derive_child_c(
		(*C.uchar)(unsafe.Pointer(&seed[0])),
		C.size_t(seedLen),
		C.uint(childIdx),
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
