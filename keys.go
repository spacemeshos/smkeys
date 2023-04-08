package main

/*
#cgo CFLAGS: -I./ed25519_bip32_wasm/
#cgo LDFLAGS: -L./ed25519_bip32_wasm/ -led25519_bip32
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

func main() {
	myString := "m/44'/0'/1'/0'/1'"
	strLen := len(myString)

	// Convert Go string to C-compatible byte array
	strBytes := []byte(myString)

	// Pass the string to Rust
	arrayPtr := C.derive_key_c(
		(*C.uchar)(unsafe.Pointer(&strBytes[0])),
		C.size_t(strLen),
		(*C.uchar)(unsafe.Pointer(&strBytes[0])),
		C.size_t(strLen),
	)
	defer C.free(unsafe.Pointer(arrayPtr))

	const arrayLen = 64

	// Convert the *mut u8 pointer to a Go byte slice
	byteSlice := (*[1 << 30]byte)(unsafe.Pointer(arrayPtr))[:arrayLen:arrayLen]
	bytes := (*[arrayLen]uint8)(unsafe.Pointer(arrayPtr))[:]

	// Access the byte array in Go
	fmt.Println("Received byte array:", byteSlice)
	fmt.Println("Received byte array:", bytes)
}
