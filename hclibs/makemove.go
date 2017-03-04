//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import "fmt"

func MakeMove(m Move, p *Pos) {

	// make tentative move
	// is king in check? not legal - bail out
	//return False if not is_legal_move(from,to, side); // perhaps too much checking in this sub
	// yes legal
	tp := p.Board[m.to]
	fp := p.Board[m.from]
	from := m.from
	to := m.to
	extra := m.extra
	side := p.Side
	xside := 1 - side
	//fmt.Print(p)
	hply := p.Ply
	// copy current move into history array
	history[hply].move = m
	history[hply].TakenPieces = p.TakenPieces
	history[hply].Castled = p.Castled
	history[hply].King = p.King
	history[hply].Side = p.Side
	history[hply].InCheck = p.InCheck
	history[hply].EnPassant = p.EnPassant
	history[hply].Fifty = p.Fifty
	history[hply].FullMoveClock = p.FullMoveClock
	history[hply].HalfMoveClock = p.HalfMoveClock
	history[hply].Hash = p.Hash
	history[hply].JustTaken = -1 // so we don't overwrite with rubbish on future passes up and down the array and for Promote and capture!
	zhash := p.Hash              // for manipulating the hash with just the changes we need to make
	//
	if from == to && from == A1 {
		panic("I've been given a nonsence move")
	}
	// remove the FROM and TO squares from old hash ready to update with new values
	//	zhash = zhash ^ Zhash.psq[from][p.Board[from]]
	//	zhash = zhash ^ Zhash.psq[to][p.Board[to]]

	if m.mtype == O_O_O { //  update castled (left)
		fmt.Println("QS castle")
		if from-to != 2 {
			panic("castling error to and from")
		}
		//         debug ==2 && say "O-O-O";
		p.Castled[side*2+KS] = true
		zhash = zhash ^ Zhash.castle[side*2+KS]
		p.Castled[side*2+QS] = true // can't castle other side once castled!
		zhash = zhash ^ Zhash.castle[side*2+QS]

		rook := from - 4
		p.Board[from-1] = p.Board[rook] // move rook to right of king
		//		zhash = zhash ^ Zhash.psq[from-1][p.Board[from -1]]
		p.Board[rook] = EMPTY
		//		zhash = zhash ^ Zhash.psq[rook][p.Board[rook]]
	}
	if m.mtype == O_O { //  update castled (right)
		fmt.Println("KS castle")
		if to-from != 2 {
			panic("castling error to and from")
		}
		//         debug ==2 && say "O-O";
		p.Castled[side*2+KS] = true
		zhash = zhash ^ Zhash.castle[side*2+KS]
		p.Castled[side*2+QS] = true // can't castle other side once castled!
		zhash = zhash ^ Zhash.castle[side*2+QS]

		rook := from + 3
		p.Board[from+1] = p.Board[rook] // move rook to left of king
		//		zhash = zhash ^ Zhash.psq[from-1][p.Board[from -1]]
		p.Board[rook] = EMPTY
		//		zhash = zhash ^ Zhash.psq[rook][p.Board[rook]]
	}
	if m.mtype == PROMOTE {
		fmt.Println("Promote")
		fp = extra + (side << 3)
		if tp != EMPTY {
			// promote by capturing
			history[hply].JustTaken = tp
		}
		// let move below promote push it into right place
	}
	if m.mtype == EPCAPTURE {
		fmt.Println("EPCapture")
		epcapture := to
		// generate has already spotted this and created a taking move.
		//         debug ==2 && say "epcapture";

		if side == WHITE {
			epcapture += SOUTH
		} else {
			epcapture += NORTH
		} // Piece to take is vert above or below to
		//		zhash = zhash ^ Zhash.psq[epcapture][p.Board[epcapture]]
		p.Board[epcapture] = EMPTY // remove p from their list
		p.TakenPieces[side]++      // increase the count of pawns
		p.Fifty = -1
	}

	//	zhash = zhash ^ Zhash.ep[p.EnPassant+1]
	p.EnPassant = -1
	if m.mtype == ENPASSANT {
		fmt.Println("EP")
		p.EnPassant = extra
		//		zhash = zhash ^ Zhash.ep[extra+1]
	}

	// move piece
	p.Board[from] = EMPTY
	////	zhash = zhash ^ Zhash.psq[from][p.Board[from]] //new value // dont hash empty squares
	p.Board[to] = fp
	//	zhash = zhash ^ Zhash.psq[to][p.Board[to]] //new value
	//     p.pieces[side][ p.pieces[side].grep-index(from) ]=to;

	// CAPTURE! update pieces for from
	if tp != EMPTY { // Capturing
		fmt.Println("Capturing")
		history[hply].JustTaken = tp // record what was taken for unmakemove
		p.TakenPieces[side]++        // increase the count of pawns
		p.Fifty = -1
	} else {
		fmt.Println("quiet")
	}

	if fp == KING+(side<<3) {
		if p.Castled[side*2+QS] == false {
			zhash = zhash ^ Zhash.castle[side*2+QS]
			p.Castled[side*2+QS] = true
		}
		if p.Castled[side*2+KS] == false {
			zhash = zhash ^ Zhash.castle[side*2+KS]
			p.Castled[side*2+KS] = true // deny the opportunity to castle
		}
		p.King[side] = to // update king index;
	}

	// these checks disallow castling if rook has moved OR rook is taken! (to)
	if p.Castled[WHITE*2+QS] == false && (from == A1 || to == A1) { //a1
		p.Castled[WHITE*2+QS] = true // denies possibilitiy of castling to qs
		zhash = zhash ^ Zhash.castle[WHITE*2+QS]
	}
	if p.Castled[WHITE*2+KS] == false && (from == H1 || to == H1) { //h1
		p.Castled[WHITE*2+KS] = true // denies possibilitiy of castling to ks
		zhash = zhash ^ Zhash.castle[WHITE*2+KS]
	}
	if p.Castled[BLACK*2+QS] == false && (from == A8 || to == A8) { //a8
		p.Castled[BLACK*2+QS] = true // denies possibilitiy of castling to qs
		zhash = zhash ^ Zhash.castle[BLACK*2+QS]
	}
	if p.Castled[BLACK*2+KS] == false && (from == H8 || to == H8) { //h8
		p.Castled[BLACK*2+KS] = true // denies possibilitiy of castling to ks
		zhash = zhash ^ Zhash.castle[BLACK*2+KS]
	}

	if (fp & 0x7) == PAWN {
		p.Fifty = -1
	} // pawn move resets fifty move counter
	p.Fifty++ // increase the fifty move rule counter...

	// swap side
	zhash = zhash ^ Zhash.side[p.Side]
	p.Side = 1 - p.Side
	zhash = zhash ^ Zhash.side[p.Side]
	p.InCheck = -1
	// have to switch sides to get in_check to evaluate right
	if InCheck(p.King[xside], xside, p) {
		p.InCheck = xside
	} // other side now in check?

	if p.Side == WHITE {
		p.FullMoveClock++
	} // incremented after blacks turn
	// p.History[p.Ply] = History{move: Move{m.from, m.to, m.mtype, m.extra} }
	p.Ply++
	// update hash
	p.Hash = TTZKey(p)
	if p.Hash != zhash {
		panic("TTZKey hash does not match makemoves calculated hash!")
	}
	return
}

func UnMakeMove(m Move, p *Pos) bool {
	hply := p.Ply - 1
	// assert history holds correct move to UnMakeMove

	if m != history[hply].move {
		panic("UnMakeMove histry array stored move disagrees with given move to unmake")
	}
	//fmt.Printf("Unmaking %v which is %v in history\np.Ply is %v and hply is %v\n",m,history[hply].move,p.Ply,hply)
	// restore previous move from history array
	switch m.mtype {
	case QUIET, ENPASSANT:
		{
			p.Board[m.from] = p.Board[m.to]
			p.Board[m.to] = EMPTY
			//		fmt.Println("Quiet or Enpassant")
		}
	case CAPTURE:
		{
			p.Board[m.from] = p.Board[m.to]
			p.Board[m.to] = history[hply].JustTaken
			//		fmt.Println("Capture")
		}
	case EPCAPTURE:
		{
			p.Board[m.from] = p.Board[m.to]
			p.Board[m.to] = EMPTY
			// if white was capturing
			if history[hply].Side == WHITE {
				p.Board[m.to+SOUTH] = pawn
				//	fmt.Println("restoring black pawn")
			} else {
				p.Board[m.to+NORTH] = PAWN
				//	fmt.Println("restoring white pawn")
			}
		}
	case O_O_O:
		{
			rook := m.from - 4
			// move rook from left of king
			p.Board[rook] = p.Board[m.from-1]
			p.Board[m.from-1] = EMPTY
			p.Board[m.from] = p.Board[m.to]
			p.Board[m.to] = EMPTY
			//fmt.Println("Queen side C")
		}
	case O_O:
		{
			rook := m.from + 3
			// move rook from left of king
			p.Board[rook] = p.Board[m.from+1]
			p.Board[m.from+1] = EMPTY
			p.Board[m.from] = p.Board[m.to]
			p.Board[m.to] = EMPTY
			//fmt.Println("King side C")
		}
	case PROMOTE:
		{
			if history[hply].JustTaken == -1 {
				p.Board[m.to] = EMPTY
			} else {
				p.Board[m.to] = history[hply].JustTaken
			}
			if history[hply].Side == WHITE {
				p.Board[m.from] = PAWN
				// fmt.Println("promote white")
			} else {
				p.Board[m.from] = pawn
				// fmt.Println("promote black")
			}
		}
	default:
		panic("encountered an unknown move type!")
	}

	//
	p.TakenPieces = history[hply].TakenPieces
	p.Castled = history[hply].Castled
	p.King = history[hply].King
	p.Side = history[hply].Side
	p.InCheck = history[hply].InCheck
	p.EnPassant = history[hply].EnPassant
	p.Fifty = history[hply].Fifty
	p.FullMoveClock = history[hply].FullMoveClock
	p.HalfMoveClock = history[hply].HalfMoveClock
	//  assert move correctly reset
	p.Ply--
	//
	p.Hash = TTZKey(p)
	if p.Hash != history[hply].Hash {
		fmt.Println(m, "\n", p)
		panic("UnMakeMove did not give same position hash as stored in history array.")
	}
	return true
}
