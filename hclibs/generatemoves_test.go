//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

//import "fmt"
import "github.com/dex4er/go-tap"
import "testing"

//import "fmt"

func TestPerftDivide(t *testing.T) {
	//	tap.Ok(true, "Ok")
	p := FENToNewBoard(STARTFEN)
	tap.Is(len(GenerateMoves(A2, &p)), 2, "PAWN how many moves from A2 on new board")
	tap.Is(len(GenerateMoves(A1, &p)), 0, "ROOK how many moves from A1 on new board")
	tap.Is(len(GenerateMoves(A7, &p)), 0, "PAWN (wrong side) how many moves from A7 on new board")
	tap.Is(len(GenerateMoves(B1, &p)), 2, "Knight how many moves from B1 on new board")

	//tap.Is("Aaa", "Aaa", "Is")
	tap.Is(len(GenerateAllMoves(&p)), 20, "20 moves counterd on a new board")

	//tap.Is(123, 123, "Is")
}

func TestIncheck(t *testing.T) {
	tap.Fail("No test for this function!!!!! cost you dearly!!!!!")
}
