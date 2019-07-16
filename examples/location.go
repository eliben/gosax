package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"unsafe"

	"github.com/eliben/gosax"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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

	err := gosax.ParseFile(flag.Args()[0], scb)
	if err != nil {
		panic(err)
	}

	fmt.Println("counter =", counter)
}
