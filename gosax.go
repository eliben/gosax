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
extern void startElementCgo(void*, const xmlChar*, const xmlChar**);
extern void endElementCgo(void*, const xmlChar*);

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
type StartElementFunc func(name string, attrs []string)
type EndElementFunc func(name string)

type SaxCallbacks struct {
	StartDocument StartDocumentFunc
	StartElement  StartElementFunc
	EndElement    EndElementFunc
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

func ParseFile(filename string, cb SaxCallbacks) error {
	var cfilename *C.char = C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	chandler := C.newHandlerStruct()

	if cb.StartDocument != nil {
		chandler.startDocument = C.startDocumentSAXFunc(C.startDocumentCgo)
	} else {
		chandler.startDocument = nil
	}

	if cb.StartElement != nil {
		chandler.startElement = C.startElementSAXFunc(C.startElementCgo)
	} else {
		chandler.startElement = nil
	}

	if cb.EndElement != nil {
		chandler.endElement = C.endElementSAXFunc(C.endElementCgo)
	} else {
		chandler.endElement = nil
	}

	chandler.internalSubset = nil
	chandler.isStandalone = nil
	chandler.hasInternalSubset = nil
	chandler.hasExternalSubset = nil
	chandler.resolveEntity = nil
	chandler.getEntity = nil
	chandler.entityDecl = nil
	chandler.notationDecl = nil
	chandler.attributeDecl = nil
	chandler.elementDecl = nil
	chandler.unparsedEntityDecl = nil
	chandler.setDocumentLocator = nil
	chandler.endDocument = nil
	chandler.reference = nil
	chandler.characters = nil
	chandler.ignorableWhitespace = nil
	chandler.processingInstruction = nil
	chandler.comment = nil
	chandler.warning = nil
	chandler.error = nil
	chandler.fatalError = nil

	user_data := pointer.Save(&cb)
	defer pointer.Unref(user_data)

	rc := C.xmlSAXUserParseFile(&chandler, user_data, cfilename)
	if rc != 0 {
		fmt.Println("xmlSAXUserParseFile returned", rc)
	}

	return nil
}
