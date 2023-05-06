package remote_wallet

/*
	#include <stdarg.h>
	#include <stddef.h>
	#include <stdint.h>
	#include <stdlib.h>
	//#include <ostream.h>
	//#include <new>


	extern uint8_t *read_pubkey_from_ledger(const uint8_t *path,
									 size_t pathlen,
									 const uint8_t *derivation_path,
									 size_t derivation_pathlen);

	/// free the memory allocated and returned by the derive functions by transferring ownership back to
	/// Rust. must be called on each pointer returned by the functions precisely once to ensure safety.
	void free_c(uint8_t *ptr);
*/
import "C"
import (
	"fmt"
	"github.com/spacemeshos/smkeys/bip32"
	"unsafe"
)

func ReadPubkeyFromLedger(path, derivation_path string) (key []byte, err error) {
	pathPtr := (*C.uchar)(unsafe.Pointer(&[]byte(path)[0]))
	pathLen := (C.size_t)(len(path))
	derivation_pathPtr := (*C.uchar)(unsafe.Pointer(&[]byte(derivation_path)[0]))
	derivation_pathLen := (C.size_t)(len(derivation_path))
	arrayPtr := C.read_pubkey_from_ledger(pathPtr, pathLen, derivation_pathPtr, derivation_pathLen)
	if arrayPtr == nil {
		return nil, bip32.PointerErr
	}
	defer C.free_c(arrayPtr)

	// Convert the C pointer to a Go byte slice
	bytes := (*[bip32.ArrayLen]byte)(unsafe.Pointer(arrayPtr))[:]
	key = make([]byte, bip32.ArrayLen)
	bytesCopied := copy(key[:], bytes)
	if bytesCopied != bip32.ArrayLen {
		return nil, fmt.Errorf("error in key length")
	}
	//key = (*C.uint8_t)(unsafe.Pointer(arrayPtr))
	return
}
