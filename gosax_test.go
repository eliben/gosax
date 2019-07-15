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
		EndElement: func(name string) {
			fmt.Printf("end elem: %v\n", name)
		},
	}
	// Just testing that nothing crashed
	err := ParseFile("testfiles/fruit.xml", scb)
	fmt.Println("err", err)
}
