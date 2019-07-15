package gosax

import (
	"fmt"
	"testing"
)

func TestInit(*testing.T) {
	scb := SaxCallbacks{
		StartDocument: func() {
			fmt.Println("Got doc start")
		},
		StartElement: func(name string, attrs []string) {
			fmt.Printf("start elem: %v, attrs %v\n", name, attrs)
		},
	}
	// Just testing that nothing crashed during init
	err := ParseFile("testfiles/fruit.xml", scb)
	fmt.Println(err)
}
