package remote_wallet

/*
	#include <stdint.h>

	extern uint8_t *read_pubkey_from_ledger(const uint8_t *path,
									 size_t pathlen,
									 const uint8_t *derivation_path,
									 size_t derivation_pathlen);
*/
import "C"
import (
	"crypto/ed25519"
	"github.com/spacemeshos/smkeys/common"
	"unsafe"
)

func ReadPubkeyFromLedger(path, derivationPath string) (key []byte, err error) {
	pathPtr := (*C.uchar)(unsafe.Pointer(&[]byte(path)[0]))
	pathLen := (C.size_t)(len(path))
	derivationPathPtr := (*C.uchar)(unsafe.Pointer(&[]byte(derivationPath)[0]))
	derivationPathLen := (C.size_t)(len(derivationPath))
	arrayPtr := C.read_pubkey_from_ledger(pathPtr, pathLen, derivationPathPtr, derivationPathLen)
	if arrayPtr == nil {
		return nil, common.PointerErr
	}
	defer common.FreeCPointer(arrayPtr)

	// Convert the C pointer to a Go byte slice
	bytes := (*[ed25519.PublicKeySize]byte)(unsafe.Pointer(arrayPtr))[:]
	key = make([]byte, ed25519.PublicKeySize)
	bytesCopied := copy(key[:], bytes)
	if bytesCopied != ed25519.PublicKeySize {
		return nil, common.KeyLenErr
	}
	//key = (*C.uint8_t)(unsafe.Pointer(arrayPtr))
	return
}
