package internal

import "sync/atomic"

// A AtomicBool is an atomic boolean value.
// The zero value is false.
// Copied from the go 1.19 source code.
type AtomicBool struct {
	v uint32
}

// Load atomically loads and returns the value stored in x.
func (x *AtomicBool) Load() bool { return atomic.LoadUint32(&x.v) != 0 }

// Store atomically stores val into x.
func (x *AtomicBool) Store(val bool) { atomic.StoreUint32(&x.v, b32(val)) }

// Swap atomically stores new into x and returns the previous value.
func (x *AtomicBool) Swap(new bool) (old bool) { return atomic.SwapUint32(&x.v, b32(new)) != 0 }

// CompareAndSwap executes the compare-and-swap operation for the boolean value x.
func (x *AtomicBool) CompareAndSwap(old, new bool) (swapped bool) {
	return atomic.CompareAndSwapUint32(&x.v, b32(old), b32(new))
}

// b32 returns a uint32 0 or 1 representing b.
func b32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}
