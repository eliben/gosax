// gosax: Go wrapper for libxml SAX.
//
// C helpers for internal use.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package gosax

/*
#cgo pkg-config: libxml-2.0

#include <libxml/tree.h>
#include <libxml/parser.h>

extern void goStartDocument(void*);
extern void goEndDocument(void*);
extern void goStartElement(void*, const xmlChar*, const xmlChar**, int);
extern void goStartElementNoAttr(void*, const xmlChar*);
extern void goEndElement(void*, const xmlChar*);
extern void goCharacters(void*, const xmlChar*, int);
extern void goCharactersRaw(void*, const xmlChar*, int);

void startDocumentCgo(void* user_data) {
  goStartDocument(user_data);
}

void endDocumentCgo(void* user_data) {
  goEndDocument(user_data);
}

void startElementCgo(void* user_data,
                     const xmlChar* name,
                     const xmlChar** attrs) {
  // The attrs array is terminated with a NULL pointer. To make it usable in
  // Go, we find the length and pass it explicitly to the Go callback.
  int i = 0;
  if (attrs != NULL) {
    while (attrs[i] != NULL) {
      i++;
    }
  }
  goStartElement(user_data, name, attrs, i);
}

void startElementNoAttrCgo(void* user_data,
                           const xmlChar* name,
                           const xmlChar** attrs) {
  goStartElementNoAttr(user_data, name);
}

void endElementCgo(void* user_data, const xmlChar* name) {
  goEndElement(user_data, name);
}

void charactersCgo(void* user_data, const xmlChar* ch, int len) {
  goCharacters(user_data, ch, len);
}

void charactersRawCgo(void* user_data, const xmlChar* ch, int len) {
  goCharactersRaw(user_data, ch, len);
}
*/
import "C"
