//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import "testing"
import "fmt"

func TestDecToAlg(t *testing.T) {
	if AlgToDec("h8") != H8 {
		t.Error("top square is not 0x77 != 119")
	}
	if DecToAlg(H8) != "h8" {
		t.Error("0x77 = top square != h8")
	}
	if AlgToDec("f5") != F5 {
		t.Error("square is 0xf5 != 69")
	}
	if DecToAlg(F5) != "f5" {
		t.Error("square is 69 != f5")
	}
	// if  AlgToDec("g9") {t.Error( "dies ok with invalid input")}
	// if  AlgToDec("8a") {t.Error( "Dies ok with invalid input")}
	// if  AlgToDec("8") {t.Error( "Dies ok with invalid input")}
	// if  AlgToDec("a") {t.Error( "Dies ok with invalid input")}
	if AlgToDec("a1") != A1 {
		t.Error("bottom square if a1 = 0")
	}
	if DecToAlg(A1) != "a1" {
		t.Error("0x77 = bottom square 0 is a1")
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

	if fmt.Sprintf("%v", FENToNewBoard(STARTFEN)) != `{ [6 2 5 7 3 5 2 6 0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 9 9 9 9 9 9 9 9 0 0 0 0 0 0 0 0 14 10 13 15 11 13 10 14 0 0 0 0 0 0 0 0] [0 0] [false false false false] [4 116] 0 -1 -1 0 1 0 0 0}` {
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
