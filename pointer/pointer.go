// package pointer helps pass arbitrary Go values to C code and back.
//
// Taken from https://github.com/mattn/go-pointer/ with small modifications.
// mattn's MIT license applies.
package pointer

// #include <stdlib.h>
import "C"
import (
	"sync"
	"unsafe"
)

var (
	mutex sync.RWMutex
	store = map[unsafe.Pointer]interface{}{}
)

func Save(v interface{}) unsafe.Pointer {
	if v == nil {
		return nil
	}

	var ptr unsafe.Pointer = C.malloc(C.size_t(1))
	if ptr == nil {
		panic("C.malloc failed")
	}

	mutex.Lock()
	store[ptr] = v
	mutex.Unlock()

	return ptr
}

func Restore(ptr unsafe.Pointer) (v interface{}) {
	if ptr == nil {
		return nil
	}

	mutex.RLock()
	v = store[ptr]
	mutex.RUnlock()
	return
}

func Unref(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}

	mutex.Lock()
	delete(store, ptr)
	mutex.Unlock()

	C.free(ptr)
}
