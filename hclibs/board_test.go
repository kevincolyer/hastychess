//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import "testing"
import "fmt"

func TestDecToAlg(t *testing.T) {
	if AlgToDec("h8") != H8 {
		t.Error("top square is not 0x77 != 119")
	}
	if DecToAlg(H8) != "h8" {
		t.Errorf("0x77 = top square != h8 -> [%v]",DecToAlg(H8))
	}
	if AlgToDec("f5") != F5 {
		t.Error("square is 0xf5 != 69")
	}
	if DecToAlg(F5) != "f5" {
		t.Errorf("square is 69 != f5 -> [%v]",DecToAlg(F5))
	}

	if AlgToDec("a1") != A1 {
		t.Error("bottom square if a1 = 0")
	}
	if DecToAlg(A1) != "a1" {
		t.Errorf("0x77 = bottom square 0 is a1 -> [%v]",DecToAlg(A1))
	}
	if DecToAlg(AlgToDec("d4")) != "d4" {
		t.Error("d4 round trip ok")
	}
}

func TestBoardToStr(t *testing.T) {
	out :=
		` rnbqkbnr 8
 pppppppp 7
 ........ 6
 ........ 5
 ........ 4
 ........ 3
 PPPPPPPP 2
 RNBQKBNR 1
 abcdefgh`
	p := FENToNewBoard(STARTFEN)
	comp := BoardToStr(&p)

	if comp != out {
		t.Errorf("Board does not match STARTFEN\n%v", comp)
	}
}
func TestBoard(t *testing.T) {

	if fmt.Sprintf("%v", FENToNewBoard(STARTFEN)) != `{ [6 2 5 7 3 5 2 6 0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 9 9 9 9 9 9 9 9 0 0 0 0 0 0 0 0 14 10 13 15 11 13 10 14 0 0 0 0 0 0 0 0] [0 0] [false false false false] [4 116] 0 -1 -1 0 1 0 0 7424690183336351638}` {
		t.Errorf("Internal representation of Board failed - has something changed?\n%v", FENToNewBoard(STARTFEN))
	}
}

func TestBoardToFEN(t *testing.T) {
	p := FENToNewBoard(STARTFEN)
	comp := BoardToFEN(&p)

	if comp != STARTFEN {
		t.Errorf("Board does not match STARTFEN on roundtrip\n%v", comp)
	}
}

func TestNewRBCFEN(t *testing.T) {
	comp := NewRBCFEN(0)

	if comp != STARTFEN {
		t.Errorf("Board does not match STARTFEN on roundtrip\n[%v]", comp)
	}
	for i := 0; i < 100; i++ {
		p := FENToNewBoard(NewRBCFEN(i % 4))
		_ = p
	}
}

func TestSide(t *testing.T) {
	if Side(PAWN) != WHITE {
		t.Error("Side does not guess piece colour corretly - white")
	}
	if Side(pawn) != BLACK {
		t.Error("Side does not guess piece colour corretly - black")
	}

}

func TestMoveToAlg(t *testing.T) {
	if MoveToAlg(Move{from: A2, to: A3}) != "a2a3" {
		t.Error("quiet move alg incorrect (poor move struct)")
	}
	if MoveToAlg(Move{from: A2, to: A3, mtype: QUIET}) != "a2a3" {
		t.Error("quiet move alg incorrect")
	}
	if MoveToAlg(Move{from: A2, to: A3, mtype: CAPTURE}) != "a2a3" {
		t.Error("capture move alg incorrect")
	}
	if MoveToAlg(Move{from: A2, to: A3, mtype: EPCAPTURE}) != "a2a3" {
		t.Error("epcapture move alg incorrect")
	}
	if MoveToAlg(Move{from: A7, to: A8, mtype: PROMOTE, extra: QUEEN}) != "a7a8q" {
		t.Error("Promote to queen (white) incorrect")
	}

	if MoveToAlg(Move{from: A7, to: A8, mtype: PROMOTE, extra: queen}) != "a7a8q" {
		t.Errorf("Promote to queen (black) incorrect %v", MoveToAlg(Move{from: A7, to: A8, mtype: PROMOTE, extra: queen}))
	}

}
