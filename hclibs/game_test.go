//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

//
// import "fmt"
import "github.com/dex4er/go-tap"
import "testing"

// import "strings"

//import "fmt"

func TestParseUserMove(t *testing.T) {

	p := FENToNewBoard(STARTFEN)
	//tap.Is("Aaa", "Aaa", "Is")
	_, text := ParseUserMove("4a4a", &p)
	tap.Is(text, "# Unparseable", "Unparseable input")
	_, text = ParseUserMove("4a4az", &p)
	tap.Is(text, "# Unparseable", "Unparseable input")
	_, text = ParseUserMove("a4a4", &p)
	tap.Is(text, "# Not a valid move", "not a valid move")
	_, text = ParseUserMove("a4a4q", &p)
	tap.Is(text, "# Not a valid move", "not a valid move")
	_, text = ParseUserMove("a1a2", &p)
	tap.Is(text, "# Not a valid move", "not a valid move")
	_, text = ParseUserMove("a2a3", &p)
	tap.Is(text, "", "Valid move")
	_, text = ParseUserMove("a2a3q", &p)
	// 	tap.Is(text, "# Not a valid move", "not a valid move") // parseable
	_, text = ParseUserMove("A2A3", &p)
	tap.Is(text, "", "Valid move")

}
