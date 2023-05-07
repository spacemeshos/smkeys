package common

/*
	#include <stdint.h>

	/// free the memory allocated and returned by extern functions by transferring ownership back to
	/// Rust. must be called on each pointer returned by the functions precisely once to ensure safety.
	extern void freeptr(uint8_t *ptr);
*/
import "C"
import (
	"fmt"
)

// Common errors used in various packages
var (
	KeyLenErr  = fmt.Errorf("invalid key length")
	PathErr    = fmt.Errorf("empty path supplied")
	PointerErr = fmt.Errorf("error in wrapped function, got nil pointer")
)

type CUChar = *C.uchar

func FreeCPointer(ptr CUChar) {
	if ptr != nil {
		C.freeptr(ptr)
	}
}
