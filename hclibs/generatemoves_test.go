//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import "fmt"
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

func TestInCheck(t *testing.T) {
	// 	tap.Fail("No test for this function!!!!! cost you dearly!!!!!")
	p := FENToNewBoard("8/5k2/8/8/8/5Q2/5K2/8 w KkqQ - 0 1") // symetrical for test
	fmt.Println(&p)
	tap.Is(p.Side == WHITE, true, "White to play")
	tap.Is(InCheck(p.King[BLACK], WHITE, &p), false, "Black King is in check (from white)")
	tap.Is(InCheck(p.King[BLACK], BLACK, &p), true, "I as black am in check on the black kings square") // I as black am in check on the black kings square
	tap.Is(InCheck(p.King[WHITE], WHITE, &p), false, "White King is NOT in check (from white)")
	tap.Is(InCheck(p.King[WHITE], BLACK, &p), true, "White King is in check (from black)") // not what you might expect!

	// checking for theoretical check on empty square
	tap.Is(InCheck(E7, BLACK, &p), false, "I as black am NOT in check on the square E7")
	tap.Is(InCheck(H5, BLACK, &p), true, "I as black am  in check on the square H5")
	fmt.Println(&p)
}

func TestIsAttacked(t *testing.T) {
	// 	tap.Fail("No test for this function!!!!! cost you dearly!!!!!")
	p := FENToNewBoard("8/5k2/8/8/8/5Q2/5K2/8 w KkqQ - 0 1") // symetrical for test
	fmt.Println(&p)
	tap.Is(p.Side == WHITE, true, "White to play")
	tap.Is(IsAttacked(p.King[BLACK], WHITE, &p), false, "Black King is in check (from white)")
	tap.Is(IsAttacked(p.King[BLACK], BLACK, &p), true, "I as black am in check on the black kings square") // I as black am in check on the black kings square
	tap.Is(IsAttacked(p.King[WHITE], WHITE, &p), false, "White King is NOT in check (from white)")
	tap.Is(IsAttacked(p.King[WHITE], BLACK, &p), true, "White King is in check (from black)") // not what you might expect!

	// checking for theoretical check on empty square
	tap.Is(IsAttacked(E7, BLACK, &p), false, "I as black am NOT in check on the square E7")
	tap.Is(IsAttacked(H5, BLACK, &p), true, "I as black am  in check on the square H5")
	fmt.Println(&p)
}
