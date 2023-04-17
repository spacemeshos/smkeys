package bip32

/*
   #cgo CFLAGS: -I..
   #cgo LDFLAGS: -L.. -led25519_bip32
   #include <stdlib.h>
   #include <stdint.h>

   // Declare the Rust function "derive_key_c" here
   extern uint8_t *derive_key_c(const uint8_t *seed,
   					  uintptr_t seedlen,
   					  const uint8_t *path,
   					  uintptr_t pathlen);

   // Declare the C standard library function "free" here
   extern void free(void*);

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

// DeriveKey wraps the underlying CFFI function. It accepts a path and a seed and returns a derived keypair.
func DeriveKey(path string, seed []byte) (key *[ArrayLen]byte, err error) {
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
	arrayPtr := C.derive_key_c(
		(*C.uchar)(unsafe.Pointer(&seed[0])),
		C.size_t(seedLen),
		(*C.uchar)(unsafe.Pointer(&pathBytes[0])),
		C.size_t(pathLen),
	)
	if arrayPtr == nil {
		return nil, pointerErr
	}
	defer C.free(unsafe.Pointer(arrayPtr))

	// Convert the *mut u8 pointer to a Go byte slice
	bytes := (*[ArrayLen]byte)(unsafe.Pointer(arrayPtr))[:]
	key = new([ArrayLen]byte)
	bytesCopied := copy(key[:], bytes)
	if bytesCopied != ArrayLen {
		return nil, fmt.Errorf("error in key length")
	}
	return
}
