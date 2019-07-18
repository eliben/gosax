package gosax

import (
	"fmt"
	"strings"
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
extern void startElementNoAttrCgo(void*, const xmlChar*, const xmlChar**);
extern void endElementCgo(void*, const xmlChar*);
extern void charactersCgo(void*, const xmlChar*, int);
extern void charactersRawCgo(void*, const xmlChar*, int);

// Since this structure contains pointers, take extra care to zero it out
// before passing it to Go code.
static inline xmlSAXHandler newHandlerStruct() {
	xmlSAXHandler h = {0};
	return h;
}

// Wrap a C macro in a function callable from Go.
static inline xmlError* getLastError() {
	return xmlGetLastError();
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

// SaxCallbacks collects callback functions to invoke on SAX events. Only
// populate callbacks you're interested in - callbacks left as nil will not
// be registered with the C layer and may save processing time.
// Some callbacks override others for optimization purposes - check the comments
// for more information.
type SaxCallbacks struct {
	// StartDocument is invoked on the "start document" event.
	StartDocument StartDocumentFunc

	// EndDocument is invoked on the "end document" event
	EndDocument EndDocumentFunc

	// StartElement is invoked whenever the beginning of a new element is found.
	// name will be the element name, and attrs a slice of attributes where
	// attribute names alternate with values. For example, given the element
	// <elem foo="bar" id="100"> the callback will get name="elem" and
	// attrs=["foo", "bar", "id", "100"].
	StartElement StartElementFunc

	// StartElementNoAttr will override StartElement, if set. When you don't
	// care about the attributes of an element, use this one - it will be faster
	// because it doesn't have to do attribute unpacking, which is expensive.
	StartElementNoAttr StartElementNoAttrFunc

	// EndElement is invoked at the end of parsing an element (after closing tag
	// has been processed), with name being the element name.
	EndElement EndElementFunc

	// Characters is invoked on character data inside elements. contents is the
	// data, as string. Note that this callback may be invoked multiple times
	// within a single tag.
	Characters CharactersFunc

	// CharactersRaw will override Characters, if set. It doesn't translate XML
	// data into a Go string, but leaves it as an opaque pair of (ch, chlen),
	// which you could use UnpackString to convert to a string if needed. This
	// could be a useful optimization if you're only occasionally interested in
	// the contents of character data.
	CharactersRaw CharactersRawFunc
}

type StartDocumentFunc func()
type EndDocumentFunc func()
type StartElementFunc func(name string, attrs []string)
type StartElementNoAttrFunc func(name string)
type EndElementFunc func(name string)
type CharactersFunc func(contents string)
type CharactersRawFunc func(ch unsafe.Pointer, chlen int)

func UnpackString(ch unsafe.Pointer, chlen int) string {
	return C.GoStringN((*C.char)(ch), C.int(chlen))
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
	// StartElementNoAttr overrides StartElement
	if cb.StartElementNoAttr != nil {
		SAXhandler.startElement = C.startElementSAXFunc(C.startElementNoAttrCgo)
	}

	if cb.EndElement != nil {
		SAXhandler.endElement = C.endElementSAXFunc(C.endElementCgo)
	}

	if cb.Characters != nil {
		SAXhandler.characters = C.charactersSAXFunc(C.charactersCgo)
	}
	// CharactersRaw overrides Characters
	if cb.CharactersRaw != nil {
		SAXhandler.characters = C.charactersSAXFunc(C.charactersRawCgo)
	}

	user_data := pointer.Save(&cb)
	defer pointer.Unref(user_data)

	rc := C.xmlSAXUserParseFile(&SAXhandler, user_data, cfilename)
	if rc != 0 {
		xmlErr := C.getLastError()
		msg := strings.TrimSpace(C.GoString(xmlErr.message))
		return fmt.Errorf("line %v: error: %v", xmlErr.line, msg)
	}

	return nil
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

//export goStartElementNoAttr
func goStartElementNoAttr(user_data unsafe.Pointer, name *C.char) {
	gcb := pointer.Restore(user_data).(*SaxCallbacks)
	gcb.StartElementNoAttr(C.GoString(name))
}

//export goEndElement
func goEndElement(user_data unsafe.Pointer, name *C.char) {
	gcb := pointer.Restore(user_data).(*SaxCallbacks)
	gcb.EndElement(C.GoString(name))
}

//export goCharacters
func goCharacters(user_data unsafe.Pointer, ch *C.char, chlen C.int) {
	gcb := pointer.Restore(user_data).(*SaxCallbacks)
	gcb.Characters(C.GoStringN(ch, chlen))
}

//export goCharactersRaw
func goCharactersRaw(user_data unsafe.Pointer, ch *C.char, chlen C.int) {
	gcb := pointer.Restore(user_data).(*SaxCallbacks)
	gcb.CharactersRaw(unsafe.Pointer(ch), int(chlen))
}
