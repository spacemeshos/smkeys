package common

/*
	#include <stdint.h>

	/// free the memory allocated and returned by the derive functions by transferring ownership back to
	/// Rust. must be called on each pointer returned by the functions precisely once to ensure safety.
	extern void derive_free_c(uint8_t *ptr);
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

func FreeCPointer(ptr *C.uint8_t) {
	if ptr != nil {
		C.derive_free_c(ptr)
	}
}
