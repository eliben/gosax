package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unsafe"

	"github.com/eliben/gosax"
)

func main() {
	counter := 0
	inLocation := false

	scb := gosax.SaxCallbacks{
		StartElement: func(name string, attrs []string) {
			if name == "location" {
				inLocation = true
			} else {
				inLocation = false
			}
		},

		// this overrides StartElement
		StartElementNoAttr: func(name string) {
			if name == "location" {
				inLocation = true
			} else {
				inLocation = false
			}
		},

		EndElement: func(name string) {
			inLocation = false
		},

		Characters: func(contents string) {
			if inLocation && strings.Contains(contents, "Africa") {
				counter++
			}
		},

		// this overrides Characters
		CharactersRaw: func(ch unsafe.Pointer, chlen int) {
			if inLocation {
				if strings.Contains(gosax.UnpackString(ch, chlen), "Africa") {
					counter++
				}
			}
		},
	}

	var xctx *gosax.XmlParserCtxt

	reader := bufio.NewReader(os.Stdin)
	for {
		chunk := make([]byte, 4096)
		_, err := reader.Read(chunk)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		if xctx == nil {
			xctx, err = gosax.CreatePushParser(string(chunk), scb)
			if err != nil {
				log.Fatal(err)
			}
			defer xctx.Close()
			continue
		}

		err = xctx.ParseChunk(string(chunk), false)
		if err != nil && err != gosax.ERROR_DOCUMENT_END {
			log.Fatal(err)
		}
	}

	// in case the document is very small parsechunk has to be called once empty
	err := xctx.ParseChunk("", true)
	if err != nil && err != gosax.ERROR_DOCUMENT_END {
		log.Fatal(err)
	}

	fmt.Printf("counter = %d\n", counter)
}
