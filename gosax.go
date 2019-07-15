package gosax

import (
	"sync"
	"unsafe"
)

/*
#cgo pkg-config: libxml-2.0

#include <libxml/tree.h>
#include <libxml/parser.h>
#include <libxml/parserInternals.h>
*/
import "C"
import "github.com/eliben/gosax/pointer"

// Used to ensure that xmlInitParser is only called once.
var initOnce sync.Once

func init() {
	initOnce.Do(func() {
		C.xmlInitParser()
	})
}

type StartDocumentFunc func()
type StartElementFunc func(name string, attrs []string)

type SaxCallbacks struct {
	StartDocument StartDocumentFunc
	StartElement  StartElementFunc
}

//export goStartDocument
func goStartDocument(user_data unsafe.Pointer) {
	gcb := pointer.Restore(user_data).(*SaxCallbacks)
	gcb.StartDocument()
}

//export goStartElement
func goStartElement(user_data unsafe.Pointer, name *C.char, attrs **C.char) {

}

func ParseFile(filename string, cb SaxCallbacks) error {

	return nil
}
