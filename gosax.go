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

extern void startDocumentCgo(void*);
extern void endDocumentCgo(void*);
extern void startElementCgo(void*, const xmlChar*, const xmlChar**);
extern void endElementCgo(void*, const xmlChar*);
extern void charactersCgo(void*, const xmlChar*, int);

// Since this structure contains pointers, take extra care to zero it out
// before passing it to Go code.
static inline xmlSAXHandler newHandlerStruct() {
	xmlSAXHandler h = {0};
	return h;
}
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
type EndDocumentFunc func()
type StartElementFunc func(name string, attrs []string)
type EndElementFunc func(name string)
type CharactersFunc func(contents string)

type SaxCallbacks struct {
	StartDocument StartDocumentFunc
	EndDocument   EndDocumentFunc
	StartElement  StartElementFunc
	EndElement    EndElementFunc
	Characters    CharactersFunc
}

//export goStartDocument
func goStartDocument(user_data unsafe.Pointer) {
	gcb := pointer.Restore(user_data).(*SaxCallbacks)
	gcb.StartDocument()
}

//export goEndDocument
func goEndDocument(user_data unsafe.Pointer) {
	gcb := pointer.Restore(user_data).(*SaxCallbacks)
	gcb.EndDocument()
}

//export goStartElement
func goStartElement(user_data unsafe.Pointer, name *C.char, attrs **C.char, attrlen C.int) {
	// TODO: optimization opportunity: provide callback for start elements w/o
	// attributes, to save all this work here and in the C wrapper.
	gcb := pointer.Restore(user_data).(*SaxCallbacks)
	length := int(attrlen)
	var goattrs []string
	if length > 0 {
		tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(attrs))[:length:length]
		goattrs = make([]string, length)
		for i, s := range tmpslice {
			goattrs[i] = C.GoString(s)
		}
	}
	gcb.StartElement(C.GoString(name), goattrs)
}

//export goEndElement
func goEndElement(user_data unsafe.Pointer, name *C.char) {
	gcb := pointer.Restore(user_data).(*SaxCallbacks)
	gcb.EndElement(C.GoString(name))
}

//export goCharacters
func goCharacters(user_data unsafe.Pointer, ch *C.char, chlen C.int) {
	// TODO: optimization opportunity: GoStringN copies the data, so we may want
	// the user to get raw *C.char and C.int instead and selectively copy to a Go
	// string when in the right parsing stage.
	gcb := pointer.Restore(user_data).(*SaxCallbacks)
	gcb.Characters(C.GoStringN(ch, chlen))
}

func ParseFile(filename string, cb SaxCallbacks) error {
	var cfilename *C.char = C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	// newHandlerStruct zeroes out all the pointers; we assign only those that
	// are passed as non-nil in SaxCallbacks.
	SAXhandler := C.newHandlerStruct()

	if cb.StartDocument != nil {
		SAXhandler.startDocument = C.startDocumentSAXFunc(C.startDocumentCgo)
	}

	if cb.EndDocument != nil {
		SAXhandler.endDocument = C.endDocumentSAXFunc(C.endDocumentCgo)
	}

	if cb.StartElement != nil {
		SAXhandler.startElement = C.startElementSAXFunc(C.startElementCgo)
	}

	if cb.EndElement != nil {
		SAXhandler.endElement = C.endElementSAXFunc(C.endElementCgo)
	}

	if cb.Characters != nil {
		SAXhandler.characters = C.charactersSAXFunc(C.charactersCgo)
	}

	user_data := pointer.Save(&cb)
	defer pointer.Unref(user_data)

	// TODO: more real error handling -- actually report the parsing error
	rc := C.xmlSAXUserParseFile(&SAXhandler, user_data, cfilename)
	if rc != 0 {
		fmt.Println("xmlSAXUserParseFile returned", rc)
	}

	return nil
}
