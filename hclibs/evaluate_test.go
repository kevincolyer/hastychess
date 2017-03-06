//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import "fmt"
import "github.com/dex4er/go-tap"
import "testing"

// import "strings"

//import "fmt"

func TestPstScore(t *testing.T) {

	p := FENToNewBoard(STARTFEN)
	//tap.Is("Aaa", "Aaa", "Is")
	// 	_, text := ParseUserMove("4a4a", &p)
	tap.Is(PstScore(&p, Gamestage(&p)), 0, "New board stating evaluation is balanced, hence 0")
	p.Board[A2] = EMPTY
	tap.Is(PstScore(&p, Gamestage(&p)), -105, "New board less one pawn (was on a +5 sq) = -105")
	// restore
	p.Board[A2] = PAWN
	// test restored.
	tap.Is(PstScore(&p, Gamestage(&p)), 0, "New board stating evaluation is balanced, hence 0")

	MakeMove(Move{from: A2, to: A4}, &p)
	p.Side = WHITE
	tap.Is(PstScore(&p, Gamestage(&p)), -5, "New board a2->a4 gives -5 for evaluation of white")
	p.Side = BLACK
	tap.Is(PstScore(&p, Gamestage(&p)), 5, "New board a2->a4 gives -5 for evaluation of black")

	p = FENToNewBoard("8/5k2/8/8/8/8/5K2/8 w KkqQ - 0 1") // symetrical for test
	tap.Is(PstScore(&p, Gamestage(&p)), 0, "Kings in end game zero sum")
	fmt.Println(&p)

	MakeMove(Move{from: F2, to: F1}, &p)
	p.Side = WHITE
	tap.Is(Pst[ENDGAME][KING][F1], -30, "is PSQ f1 for king end game what I expect?")
	tap.Is(Pst[ENDGAME][KING][F2], 0, "is PSQ f2 for king end game what I expect?")
	tap.Is(PstScore(&p, Gamestage(&p)), -30, "Kings in end game = -30")
	fmt.Println(&p)

}

func TestGamestage(t *testing.T) {

	p := FENToNewBoard(STARTFEN)
	//tap.Is("Aaa", "Aaa", "Is")
	tap.Is(Gamestage(&p), OPENING, "Gamestage: test finds opening")
	p = FENToNewBoard("r3k2r/pppp4/8/8/8/8/PPPP4/1R2K2R w Kkq - 0 1")
	tap.Is(Gamestage(&p), MIDGAME, "Gamestage: test finds opening")
	p = FENToNewBoard("8/8/8/8/8/8/8/kK6 w KkqQ - 0 1")
	tap.Is(Gamestage(&p), ENDGAME, "Gamestage: test finds endgame")

}
