package gosax

import (
	"testing"
)

func TestInit(*testing.T) {
	// Test that nothing crashed in init()
}

func TestBasic(t *testing.T) {
	var plantId string
	var startDoc bool
	var endDoc bool

	scb := SaxCallbacks{
		StartDocument: func() {
			startDoc = true
		},
		EndDocument: func() {
			endDoc = true
		},
		StartElement: func(name string, attrs []string) {
			if name == "plant" {
				if len(attrs) < 2 {
					t.Errorf("want len(attrs) at least 2, got %v", len(attrs))
				}
				if attrs[0] != "id" {
					t.Errorf("want 'id' attr, got %v", attrs[0])
				}
				plantId = attrs[1]
			}
		},
		EndElement: func(name string) {
		},
	}

	err := ParseFile("testfiles/fruit.xml", scb)
	if err != nil {
		panic(err)
	}

	if plantId != "27" {
		t.Errorf("want plant id %v, got %v", 27, plantId)
	}

	if !startDoc {
		t.Errorf("want doc start, found none")
	}

	if !endDoc {
		t.Errorf("want doc end, found none")
	}
}

func TestCharacters(t *testing.T) {
	m := make(map[string]bool)
	scb := SaxCallbacks{
		Characters: func(contents string) {
			m[contents] = true
		},
	}

	err := ParseFile("testfiles/fruit.xml", scb)
	if err != nil {
		panic(err)
	}

	chars := []string{"Coffee", "Ethiopia", "Brazil"}
	for _, c := range chars {
		if _, ok := m[c]; !ok {
			t.Errorf("expected to find %v characters", c)
		}
	}
}
