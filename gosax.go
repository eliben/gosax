package gosax

import (
	"fmt"
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
func goStartElement(user_data unsafe.Pointer, name *C.char, attrs **C.char, attrlen C.int) {
	gcb := pointer.Restore(user_data).(*SaxCallbacks)
	length := int(attrlen)
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(attrs))[:length:length]
	goattrs := make([]string, length)
	for i, s := range tmpslice {
		goattrs[i] = C.GoString(s)
	}
	gcb.StartElement(C.GoString(name), goattrs)

}

func ParseFile(filename string, cb SaxCallbacks) error {
	var cfilename *C.char = C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	cHandler := C.xmlSAXHandler{}

	if cb.StartDocument != nil {
		cHandler.startDocument = C.startDocumentSAXFunc(startDocumentCgo)
	} else {
		cHandler.startDocument = nil
	}

	if cb.StartElementFunc != nil {
		cHandler.startElemnt = C.startElementSAXFunc(startElementCgo)
	} else {
		cHandler.startElement = nil
	}

	user_data := pointer.Save(cb)
	defer pointer.Unref(user_data)

	rc := c.xmlSAXUserParseFile(&cHandler, user_data, cfilename)
	if rc != 0 {
		fmt.Println("xmlSAXUserParseFile returned", rc)
	}

	return nil
}
