//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

// generate moves
// validate moves? || in seperate file?

//import "fmt"

func GenerateAllMoves(p *Pos) (moves []Move) {
	for _, sq := range GRID {
		if p.Board[sq] != EMPTY && p.Board[sq]>>3 == p.Side {
			all := GenerateMoves(sq, p)
			//fmt.Println(all)
			for _, j := range all {
				// test legality here...
				legal, check := is_legal_move(j, p)
				if legal {
					if check { // TODO add check flag}
					}
					moves = append(moves, j)
				}
			}

		}
	}
	return
}

func GenerateAllPseudoMoves(p *Pos) (moves []Move) { // quicker than filtering all for check - can do that elsewhere...
	for _, sq := range GRID {
		if p.Board[sq] != EMPTY && p.Board[sq]>>3 == p.Side {
			all := GenerateMoves(sq, p)
			//fmt.Println(all)
			moves = append(moves, all...)
		}
	}
	return
}

func GenerateMovesForQSearch(p *Pos) (moves []Move) {
	for _, sq := range GRID {
		if p.Board[sq] != EMPTY && p.Board[sq]>>3 == p.Side {
			all := GenerateMoves(sq, p)
			//fmt.Println(all)
			// filter moves for QS search - just return the noisy ones
			for _, j := range all {
				// test legality here...
				if j.mtype == CAPTURE || j.mtype == PROMOTE { // || j.mtype==EPCAPTURE {
					legal, check := is_legal_move(j, p)
					if legal {
						if check { // TODO add check flag
						}
						moves = append(moves, j)
					}
				}
			}

		}
	}
	return
}

func GenerateMoves(from int, p *Pos) (moves []Move) {

	side := p.Side
	xside := 1 - side
	//xside_king:=KING + (xside << 3)
	piece := p.Board[from]
	if piece>>3 == xside {
		return
	} // can't move opponents pieces!
	var pc, to int // for holding pices. // for square to move to

	switch piece & 7 { // What kind of piece am I? mask off the side bit 4
	case NIGHT:
		for _, m := range NM {
			if (from+m)&0x88 == 0 { // not off board
				pc = p.Board[from+m]
				if pc == EMPTY || (((pc >> 3) == xside) && ((pc & 7) != KING)) {
					moves = append(moves, Move{from: from, to: from + m, mtype: move_type(pc, xside), extra: 0, piece: NIGHT}) // ! capture king
				}
			}
		}

	case KING:
		//             say "in generate_moves - castled is "~ p.castled.perl
		for _, m := range KM {
			if (from+m)&0x88 == 0 { // off board
				pc = p.Board[from+m]

				if !InCheck(from+m, side, p) && (pc == EMPTY || (pc>>3) == xside) && !king_is_near(from+m, p) {
					moves = append(moves, Move{from: from, to: from + m, mtype: move_type(pc, xside), extra: 0, piece: KING})
				}
			}
		}
		// bug below - as side=0 == black then if InCheck !=-1 but set to zero on initialisation then this is a problem!!!!! Solved in Board - set to -1
		if p.InCheck != side { //&& (p.Castled[side*2+KS] == false || p.Castled[side*2+QS] == false) { // CASTLING//

			if p.Castled[side*2+KS] == false { // kingsside
				if EMPTY == p.Board[from+1] && EMPTY == p.Board[from+2] && !InCheck(from+1, side, p) && !InCheck(from+2, side, p) && !king_is_near(from+2, p) {

					moves = append(moves, Move{from: from, to: from + 2, mtype: O_O, extra: 0, piece: KING})
				}
			}
			if p.Castled[side*2+QS] == false { // queens side
				if EMPTY == p.Board[from-1] && EMPTY == p.Board[from-2] && EMPTY == p.Board[from-3] && !InCheck(from-1, side, p) && !InCheck(from-2, side, p) && !king_is_near(from-2, p) {

					moves = append(moves, Move{from: from, to: from - 2, mtype: O_O_O, extra: 0, piece: KING})
				}
			}
		}

	case ROOK:
		moves = append(moves, slider_moves(from, p, RM[:], ROOK)...)

	case BISHOP:
		moves = append(moves, slider_moves(from, p, BM[:], BISHOP)...)

	case QUEEN:
		moves = append(moves, slider_moves(from, p, QM[:], QUEEN)...)

	case PAWN: // TODO capture to promote -- need a test && check for this...
		var promfile int
		if side == WHITE {
			promfile = 6
		} else {
			promfile = 1
		} // file that next move is promotion
		var m [4]int
		//
		if side == WHITE {
			m[0] = NW
			m[1] = NE
			m[2] = NORTH
			m[3] = NN
		} else {
			m[0] = SW
			m[1] = SE
			m[2] = SOUTH
			m[3] = SS
		}

		//here!
		i := 0
		for i < 2 { // look to LEFT && RIGHT - capture || check possible?
			to = from + m[i]
			if (to & 0x88) == 0 { // on board
				pc = p.Board[to]
				if pc != EMPTY && pc>>3 == xside && pc&7 != KING { // piece to capture - ! the king
					//                     next if pc +> 3 ==side   // blocked
					//                         my type=
					if (from >> 4) == promfile {
						// 4 possible promotions
						moves = append(moves, Move{from: from, to: to, mtype: PROMOTE, extra: QUEEN, piece: PAWN}, Move{from: from, to: to, mtype: PROMOTE, extra: BISHOP, piece: PAWN}, Move{from: from, to: to, mtype: PROMOTE, extra: ROOK, piece: PAWN}, Move{from: from, to: to, mtype: PROMOTE, extra: NIGHT, piece: PAWN})
					} else {
						moves = append(moves, Move{from: from, to: to, mtype: move_type(pc, xside), extra: 0, piece: PAWN})
					} // capture || check || promote
				}
			}
			i++
		}
		// forward moves
		file := from >> 4 // upper nybble is file
		var initfile int
		if side == WHITE {
			initfile = 1
		} else {
			initfile = 6
		} // counting from 0
		var epfile int
		if side == WHITE {
			epfile = 4 // file 5-1
		} else {
			epfile = 3 // file 4-1
		} // file to be on to do ep
		var j int
		if file == initfile {
			j = 3
		} else {
			j = 2
		} // last index of @m
		i = 2
		for i <= j { // loop vairies from look ahead 1 to look ahead 2
			to = from + m[i] // one || two spaces forward
			//                 if ! (to +& 0x88)  //  on board
			// en passant captures
			if p.EnPassant > -1 && file == epfile { // empty or -1?
				//                     debug ==2 && say "considering epcapture for p.en_passant"
				if from+m[0] == p.EnPassant { // ep to one side
					moves = append(moves, Move{from: from, to: from + m[0], mtype: EPCAPTURE, extra: 0, piece: PAWN})
				}
				if from+m[1] == p.EnPassant { // ep to other side
					moves = append(moves, Move{from: from, to: from + m[1], mtype: EPCAPTURE, extra: 0, piece: PAWN})
				}
				// || ! at all (too far away)
			}
			// enpassant capture must come before following check not after!
			pc = p.Board[to]
			if pc != EMPTY {
				break
			} // blocked by a piece

			// path only valid if checking for double move
			if i == 3 {

				if p.Board[from+m[2]] == EMPTY { // blocked in between inital && double move forward
					moves = append(moves, Move{from: from, to: to, mtype: ENPASSANT, extra: from + m[2], piece: PAWN})

					break
				}
			}

			// end of double move bit

			// promotions
			if file == promfile {
				// 4 possible promotions
				moves = append(moves, Move{from: from, to: to, mtype: PROMOTE, extra: QUEEN, piece: PAWN}, Move{from: from, to: to, mtype: PROMOTE, extra: BISHOP, piece: PAWN}, Move{from: from, to: to, mtype: PROMOTE, extra: ROOK, piece: PAWN}, Move{from: from, to: to, mtype: PROMOTE, extra: NIGHT, piece: PAWN})

				break
			}

			// just push on ahead - default
			moves = append(moves, Move{from: from, to: to, mtype: QUIET, extra: 0, piece: PAWN}) // default move

			i++
		}

	case EMPTY:
		panic("No piece at from {dec_to_alg(from)} // should never reach this...")
	}
	return
}

// Helper functions for generate
func move_type(piece, xside int) int {
	if piece == EMPTY {
		return QUIET
	}
	if (piece >> 3) != xside {
		return QUIET
	} // (our side) needed for sliders
	return CAPTURE
}

// Helper function for generate
func slider_moves(from int, p *Pos, dirs []int, mypiece int) (moves []Move) {

	side := p.Side
	xside := 1 - side
	n := 0
	piece := 0
	for _, m := range dirs {
		n = m + from
		// follow a ray
		for (n & 0x88) == 0 {
			piece = p.Board[n]
			// if king then too far
			if piece == KING+(xside<<3) {
				break
			}
			// if one of ours, too far!
			if piece != EMPTY && (piece>>3) == side {
				break
			}
			// capture
			if piece != EMPTY {
				moves = append(moves, Move{from: from, to: n, mtype: CAPTURE, extra: 0, piece: mypiece})
				break
			}
			// quiet
			moves = append(moves, Move{from: from, to: n, mtype: QUIET, extra: 0, piece: mypiece})
			n += m
		}
	}
	return
}

// Answers question: Am I in check? attacked? Side needs to be opponent for a true answer
// does not do a sanity check that the square king refers to is filled OR which colour king is there!!!!
// So looking at a white king from black's perspective will answer TRUE if a white piece "attacking"" is encountered!
// use with caution!

func InCheck(king, side int, p *Pos) bool {
	// ?????returns no, yes, stalemate, checkmate? looks from OTHER side perspective

	ssft := side << 3        //our sides shift
	xssft := (side ^ 1) << 3 //their sides shift
	sk := KING + ssft        // playing side's king
	var piece, ray int
	xB := BISHOP + xssft
	xQ := QUEEN + xssft
	xR := ROOK + xssft
	xN := NIGHT + xssft
	xP := PAWN + xssft

	for _, m := range BM { //// Queen and Bishops
		ray = king + m
		for ray&0x88 == 0 { // until off board
			piece = p.Board[ray]
			if piece != EMPTY && piece != sk { // side's king
				if piece == xB || piece == xQ {
					return true
				} // it's their queen or bishop
				break // it is something else...
			}
			ray += m
		}
	}
	for _, m := range RM { //// Queen and rooks
		ray = king + m
		for (ray & 0x88) == 0 { // until off board
			piece = p.Board[ray]
			if piece != EMPTY && piece != sk { // side's king
				if piece == xR || piece == xQ {
					return true
				} // their r or q
				break // something else
			}
			ray += m
		}
	}
	for _, m := range NM { //// Kights
		if (king+m)&0x88 == 0 { // off board check
			piece = p.Board[king+m]
			if piece != EMPTY && piece != sk {
				if piece == xN {
					return true
				} // yes a knight
			}
		}
	}
	// is pawn above or below? W is below attacking up, B is above attacking down
	// but we are searching from the king to the pawns...
	var pawns [2]int
	if side == WHITE {
		pawns[0] = 17
		pawns[1] = 15
	} else {
		pawns[0] = -17
		pawns[1] = -15
	}
	for _, m := range pawns {
		if (king+m)&0x88 == 0 { // on board
			if p.Board[king+m] == xP {
				return true
			}
		}
	}
	return false // nothing threatens me!!! Muhahahaha!
}

// Is attacked - answers question is SQUARE attacked by SIDE
func IsAttacked(square, side int, p *Pos) bool {
	// ?????returns no, yes, stalemate, checkmate? looks from OTHER side perspective

	// 	ssft := side << 3        //our sides shift
	xssft := side << 3 //their sides shift
	var piece, ray int
	xK := KING + xssft
	xB := BISHOP + xssft
	xQ := QUEEN + xssft
	xR := ROOK + xssft
	xN := NIGHT + xssft
	xP := PAWN + xssft

	//sliders
	for _, m := range BM { // Queen and Bishops and king
		ray = square + m
		if ray&0x88 == 0 && p.Board[ray] == xK {
			return true
		} // stop king sliding
		for ray&0x88 == 0 { // until off board
			piece = p.Board[ray]
			if piece != EMPTY {
				if piece == xB || piece == xQ {
					return true
				} // it's their queen or bishop
				break // it is something else...
			}
			ray += m
		}
	}
	for _, m := range RM { // Queen and rooks and king
		ray = square + m
		if ray&0x88 == 0 && p.Board[ray] == xK {
			return true
		} // stop king sliding
		for (ray & 0x88) == 0 { // until off board
			piece = p.Board[ray]
			if piece != EMPTY {
				if piece == xR || piece == xQ {
					return true
				} // their r or q
				break // something else
			}
			ray += m
		}
	}
	for _, m := range NM { // Kights
		if (square+m)&0x88 == 0 { // off board check
			piece = p.Board[square+m]
			if piece != EMPTY {
				if piece == xN {
					return true
				} // yes a knight
			}
		}
	}
	// is pawn above or below? W is below attacking up, B is above attacking down
	// but we are searching from the square to the pawns...
	var pawns [2]int
	if side == WHITE {
		pawns[0] = 17
		pawns[1] = 15
	} else {
		pawns[0] = -17
		pawns[1] = -15
	}
	for _, m := range pawns {
		if (square+m)&0x88 == 0 { // on board
			if p.Board[square+m] == xP {
				return true
			}
		}
	}
	return false // nothing threatens me!!! Muhahahaha!
}

func king_is_near(look int, p *Pos) bool {
	// search around king next move for other king, if one found then invalid
	k := p.King[1-p.Side]
	for _, i := range KM {
		if (look+i)&0x88 == 0 && look+i == k {
			return true
		}
	}
	return false
}

func is_legal_move(m Move, p *Pos) (retval, check bool) {

	side := p.Side
	king := p.King[side]
	// 	xside:=1-side
	// 	xking:=p.King[xside]
	var rook, epcap, eppawn int
	retval = true

	/////////////// make tentative move to check for legality
	// and store values to restore at end
	fp := p.Board[m.from]
	tp := p.Board[m.to]
	p.Board[m.from] = EMPTY
	p.Board[m.to] = fp
	// castle
	if m.mtype == O_O {
		rook = p.Board[m.from+3]
		p.Board[m.from+3] = EMPTY
		p.Board[m.from+1] = rook
	}
	if m.mtype == O_O_O {
		rook = p.Board[m.from-4]
		p.Board[m.from-4] = EMPTY
		p.Board[m.from-1] = rook
	}
	//     if type +& PROMOTE { // does not requre special reset code after // dont test promote as only want to know if a legal move (i.e. a pawn move)
	if m.mtype == EPCAPTURE {
		epcap = m.to
		if side == WHITE {
			epcap += -16
		} else {
			epcap += 16
		} // Piece to take is vert above or below to
		eppawn = p.Board[epcap]
		p.Board[epcap] = EMPTY
	}

	////// perform tests...
	if king == m.from { // if king, cant move king into check
		if InCheck(m.to, side, p) {
			retval = false
		}
	} else { // another piece cant leave king in check either...
		if InCheck(king, side, p) {
			retval = false
		} // cant move and allow king to be in check
	}

	// TODO test if move puts opponants king in check!!! return check=true
	// 	switch PieceType(fp) {
	//         case PAWN: if PieceColour(fp) == WHITE {
	//             if fp==xking+PM[side][0] || fp==xking+PM[side][2] { check = true ; break }
	//         }
	//         case KNIGHT:
	//             for i:=range NM
	//         case ROOK,QUEEN:
	//         case BISHOP,QUEEN:
	//     }
	// TODO other moves to check - 3 repetitions or fifty move?

	////// Completed test capture
	// make UNMOVE
	if m.mtype == EPCAPTURE {
		p.Board[epcap] = eppawn
	}
	// castle
	if m.mtype == O_O_O {
		p.Board[m.from-4] = rook
		p.Board[m.from-1] = EMPTY
	}
	if m.mtype == O_O {
		p.Board[m.from+3] = rook
		p.Board[m.from+1] = EMPTY
	}
	// reset board
	p.Board[m.from] = fp
	p.Board[m.to] = tp // reset board after move
	return
}

func IsValidMove(m Move, p *Pos) bool {
	moves := GenerateAllMoves(p)
	for _, j := range moves {

		if m == j {
			return true
		}
	}
	return false
}

func DoesMoveCheckKing(m, Move, p *Pos) (b bool) {
	return
}
