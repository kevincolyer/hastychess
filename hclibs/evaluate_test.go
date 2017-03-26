//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import "fmt"
import "github.com/dex4er/go-tap"
import "testing"

func TestPstScore(t *testing.T) {

	p := FENToNewBoard(STARTFEN)
	//tap.Is("Aaa", "Aaa", "Is")
	// 	_, text := ParseUserMove("4a4a", &p)
	tap.Is(PstScore(&p, 0, Gamestage(&p)), 0, "New board stating evaluation is balanced, hence 0")
	p.Board[A2] = EMPTY
	tap.Is(PstScore(&p, 0, Gamestage(&p)), -55, "New board less one pawn (was on a +5 sq) = -55 (with dd/ss/ii)")
	// restore
	p.Side = 1 - p.Side
	tap.Is(PstScore(&p, 0, Gamestage(&p)), 55, "Test symetrical evaluation (as black now)")
	p.Side = 1 - p.Side

	// 	p.Board[A2] = PAWN
	// 	// test restored.
	// 	tap.Is(PstScore(&p, 0, Gamestage(&p)), 0, "New board stating evaluation is balanced, hence 0")
	//
	// 	MakeMove(Move{from: A2, to: A4}, &p)
	// 	p.Side = WHITE
	// 	tap.Is(PstScore(&p, 0, Gamestage(&p)), -5, "New board a2->a4 gives -5 for evaluation of white")
	// 	p.Side = BLACK
	// 	tap.Is(PstScore(&p, 0, Gamestage(&p)), 5, "New board a2->a4 gives -5 for evaluation of black")

	p = FENToNewBoard("8/5k2/8/8/8/8/5K2/8 w KkqQ - 0 1") // symetrical for test
	tap.Is(PstScore(&p, 0, Gamestage(&p)), 0, "Kings in end game zero sum")

	MakeMove(Move{from: F2, to: F1}, &p)
	p.Side = WHITE
	tap.Is(Pst[ENDGAME][KING][F1], -30, "is PSQ f1 for king end game what I expect?")
	tap.Is(Pst[ENDGAME][KING][F2], 0, "is PSQ f2 for king end game what I expect?")
	tap.Is(PstScore(&p, 0, Gamestage(&p)), -30, "Kings in end game = -30")

	// test check
	// 	p = FENToNewBoard("8/5k2/8/8/8/5Q2/5K2/8 w KkqQ - 0 1") // symetrical for test
	// 	fmt.Println(&p)
	// 	tap.Is(PstScore(&p, Gamestage(&p)) > CHECK, true, "Black is in check - true?")
	// 	// test as black that black is in check
	// 	p = FENToNewBoard("8/5k2/8/8/8/5Q2/5K2/8 b KkqQ - 0 1") // symetrical for test
	// 	fmt.Println(&p)
	// 	tap.Is(PstScore(&p, Gamestage(&p)) < -CHECK, true, "Black is in check - true?")
	// 	//fmt.Println(PstScore(&p, Gamestage(&p)))

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

func TestMVVLVA(t *testing.T) {
	var m Move
	p := FENToNewBoard("r7/1P3k1r/4p1P1/5b2/Pp2p1BN/5Q2/4pKr1/8 w KkqQ - 0 1")
	fmt.Println(&p)
	// Qxp
	m.from = F3
	m.to = E2
	val := MVVLVA(m, &p)
	tap.Is(val < 0, true, "Qxp F3E2 is net negative (not a worthwhile)")
	tap.Is(val, -900+100, "Qxp F3E2 is net negative")
	tap.Is(BLIND(m, &p), true, "QxP F3E2 is BLIND as p is lower but not defended")
	m.from = F3
	m.to = E4
	tap.Is(BLIND(m, &p), false, "QxP F3E4 is NOT BLIND")

	// pxQ
	m.from = E4
	m.to = F3
	val = MVVLVA(m, &p)
	tap.Is(val > 0, true, "pxQ is net positive")
	tap.Is(val, -100+900, "pxQ is net positive (a good capture)")
	tap.Is(BLIND(m, &p), true, "PxQ E4F3 is BLIND")

	// enpassant capture
	m.from = B4
	m.to = A3
	m.mtype = EPCAPTURE
	val = MVVLVA(m, &p)
	tap.Is(val, 0, "EP Capture ranks zero")
	// promote + capture
	m.from = B7
	m.to = A8
	m.mtype = PROMOTE
	m.extra = QUEEN
	val = MVVLVA(m, &p)
	tap.Is(val, 900-100+500, "P Promote capture rook to QUEEN (a very good capture)")
	// test Pxr > Bxp
	m.from = G6
	m.to = H7
	m.mtype = CAPTURE
	m.extra = 0
	val = MVVLVA(m, &p)
	// bxP
	m.from = F5
	m.to = G6
	val2 := MVVLVA(m, &p)
	tap.Is(val > val2, true, "Test Pxr > Bxp")
	// Rxp
	m.from = H7
	m.to = G6
	val = MVVLVA(m, &p)
	tap.Is(val < val2, true, "Test Rxp < Bxp")
	// g6xf5 == blind false
	m.from = F3
	m.to = E4
	tap.Is(BLIND(m, &p), false, "F3xE4 QxP guarded by bishop F5 is NOT BLIND")

}
