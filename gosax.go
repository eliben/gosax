package gosax

import (
	"fmt"
	"sync"
)

import "github.com/eliben/gosax/pointer"

/*
#cgo pkg-config: libxml-2.0

#include <libxml/tree.h>
#include <libxml/parser.h>
#include <libxml/parserInternals.h>
*/
import "C"

// Used to ensure that xmlInitParser is only called once.
var initOnce sync.Once

func init() {
	initOnce.Do(func() {
		C.xmlInitParser()
	})

	up := pointer.Save(3)
	fmt.Println(up)
}
