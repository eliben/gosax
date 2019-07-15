package gosax

/*
#cgo pkg-config: libxml-2.0

#include <libxml/tree.h>
#include <libxml/parser.h>

extern void goStartDocument(void*);
extern void goStartElement(void*, const xmlChar*, const xmlChar**);

void startDocumentCgo(void* user_data) {
  goStartDocument(user_data);
}

void startElementCgo(void* user_data,
                     const xmlChar* name,
										 const xmlChar** attrs) {
  goStartElement(user_data, name, attrs);
}
*/
import "C"
