//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

// import "fmt"

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
	if from == to && from == A1 {
		panic("I've been given a nonsence move")
	}
	if m.mtype == O_O_O { //  update castled (left)
		if from-to != 2 {
			panic("castling error to and from")
		}
		//         debug ==2 && say "O-O-O";
		p.Castled[side*2+KS] = true
		p.Castled[side*2+QS] = true // can't castle other side once castled!
		rook := from - 4
		//         die "unexpectedly the rook is not where I think it should be" if p.Board[rook]==EMPTY;
		//         p.pieces[side][ p.pieces[side].grep-index(rook) ]=from-1;
		p.Board[from-1] = p.Board[rook] // move rook to right of king
		p.Board[rook] = EMPTY
	}
	if m.mtype == O_O { //  update castled (right)
		if to-from != 2 {
			panic("castling error to and from")
		}
		//         debug ==2 && say "O-O";
		p.Castled[side*2+KS] = true
		p.Castled[side*2+QS] = true // can't castle other side once castled!
		rook := from + 3

		p.Board[from+1] = p.Board[rook] // move rook to left of king
		p.Board[rook] = EMPTY
	}
	if m.mtype == PROMOTE {
		fp = extra + (side << 3)
		// let move below promote push it into right place
	}
	if m.mtype == EPCAPTURE {
		epcapture := to
		// generate has already spotted this and created a taking move.
		//         debug ==2 && say "epcapture";

		if side == WHITE {
			epcapture += SOUTH
		} else {
			epcapture += NORTH
		} // Piece to take is vert above or below to
		p.Board[epcapture] = EMPTY // remove p from their list
		p.TakenPieces[side]++      // increase the count of pawns
		p.Fifty = -1
	}

	p.EnPassant = -1
	if m.mtype == ENPASSANT {
		p.EnPassant = extra
		//         debug ==2 && say "setting en_passant to {p.en_passant}";
	}

	// move piece
	p.Board[from] = EMPTY
	p.Board[to] = fp
	//     p.pieces[side][ p.pieces[side].grep-index(from) ]=to;

	// update pieces for from
	if tp != EMPTY { // Capturing
		// update pieces for to noting captures
		//         debug ==2 && say "capturing";
		//         die "epcaptures should not come here..." if m.mtype +& EPCAPTURE;
		//         my Int idx=-1;
		//         for ^p.pieces[xside] {
		//             if p.pieces[xside][_]==to {
		//                 idx = _;
		//                 last;
		//             };
		//         };
		//         die "to index not found" if idx == -1;
		//         p.pieces[xside].splice( idx ,1) ;
		//         debug ==2 && say p.pieces[xside].perl;
		p.TakenPieces[side]++ // increase the count of pawns
		p.Fifty = -1
	}

	if fp == KING+(side<<3) {
		p.Castled[side*2+QS] = true
		p.Castled[side*2+KS] = true // deny the opportunity to castle
		p.King[side] = to           // update king index;
	}

	// these checks disallow castling if rook has moved OR rook is taken! (to)
	if from == A1 || to == A1 { //a1
		p.Castled[WHITE*2+QS] = true // denies possibilitiy of castling to qs
	}
	if from == H1 || to == H1 { //h1
		p.Castled[WHITE*2+KS] = true // denies possibilitiy of castling to ks
	}
	if from == A8 || to == A8 { //a8
		p.Castled[BLACK*2+QS] = true // denies possibilitiy of castling to qs
	}
	if from == H8 || to == H8 { //h8
		p.Castled[BLACK*2+KS] = true // denies possibilitiy of castling to ks
	}

	if (fp & 0x7) == PAWN {
		p.Fifty = -1
	} // pawn move resets fifty move counter
	p.Fifty++ // increase the fifty move rule counter...

	// swap side
	p.Side = 1 - p.Side
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
	return
}

func UnMakeMove(p *Pos, from, to, mtype, extra int) {
	return
}
