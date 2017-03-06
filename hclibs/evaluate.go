//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

// import "fmt"

// f(p) = 200(K-K')
//        + 9(Q-Q')
//        + 5(R-R')
//        + 3(B-B' + N-N')
//        + 1(P-P')
//        - 0.5(D-D' + S-S' + I-I')
//        + 0.1(M-M') + ...
//
// KQRBNP = number of kings, queens, rooks, bishops, knights and pawns
// D,S,I = doubled, blocked and isolated pawns
// M = Mobility (the number of legal moves)

func Eval(p *Pos, nummoves, gamestage int) int {
	return PstScore(p, gamestage)
	//		return ClaudeShannonScore(p, nummoves)
}

func EvalQ(p *Pos, nummoves, gamestage int) int {
	return PstScore(p, gamestage) - 1 // -1 is so we don't trip the q search beta giving same score as normal eval.
	// 	return ClaudeShannonScore(p, nummoves)
}

// see https://chessprogramming.wikispaces.com/Simplified+evaluation+function for these values
var csshash = map[int]int{
	QUEEN: 900, queen: 900,
	ROOK: 500, rook: 500,
	BISHOP: 330, bishop: 330,
	NIGHT: 320, night: 320,
	PAWN: 100, pawn: 100,
	KING: 0, king: 0} // ignore kings as evaluated elsewhere...

// Most Valuable Victim, Least Valuable Agressor
func MVVLVA(m Move, p *Pos) int {
	return csshash[p.Board[m.to]] + csshash[m.extra] - csshash[p.Board[m.from]] // If promotion then m.extra has value of piece we promote to, otherwise it is 0
	// 	return -csshash[p.Board[m.to]] + csshash[m.extra] + csshash[p.Board[m.from]] // If promotion then m.extra has value of piece we promote to, otherwise it is 0
	//     return csshash[p.Board[m.to]]-csshash[p.Board[m.from]] // If promotion then m.extra has value of piece we promote to, otherwise it is 0
}

func ClaudeShannonScore(p *Pos, totalmoves int) int {
	side := p.Side
	xside := Xside(side)
	incheck := p.InCheck
	// 	enpassant := p.EnPassant
	score := 0
	var piece, up, dd, ss int
	// K-K'
	if incheck == xside {
		score += CHECK
	}
	if incheck == side {
		score -= CHECK
	}
	// Material valuation
	// could add bonuses at different game stages...
	for _, i := range GRID {
		piece = p.Board[i]
		if piece == EMPTY {
			continue
		}

		// add the piece value if ours - subtract if theirs
		if Side(piece) == side {
			score += csshash[piece]
		} else {
			score -= csshash[piece]
		}

		if piece&0x7 != PAWN {
			continue
		}

		// Pawn mobility additions
		if Side(piece) == WHITE {
			up = NORTH
		} else {
			up = SOUTH
		}

		// doubled pawns
		dd = i + up
		ss = dd
		// start from this pawn and count upwards
		// if still on board...
		for Onboard(dd) {
			if p.Board[dd] == piece {
				score -= 50
			} // one of our pawns
			dd += up
		}

		//ss - blocked PAWNS TODO blocked by what?
		if Onboard(ss) && EMPTY != p.Board[ss] && piece != p.Board[ss] {
			score -= 50
		} // blocked by opposite
		//ii - isolated
		for _, j := range QM {
			if Onboard(i+j) && piece == p.Board[i+j] {
				score -= 50
				break
			}
		}
	}

	// REMOVING MOBILITY FROM EVALUATION AS IT DOES NOT OFFER ENOUGH COMPARED TO THE ABOVE. THERE ARE OTHER FACTORS TOO THAT SHOULD BE ADDED FOR A GOOD EVALUATOR. COME BACK TO THIS AT SOME POINT
	// 	p.Side = xside
	// 	p.InCheck = -1
	// 	p.EnPassant = 0
	// 	score += (totalmoves - len(GenerateAllMoves(p))) * 10 // M-M' // just a rough estimate of how many moves...
	// 	p.EnPassant = enpassant
	// 	p.InCheck = incheck // restore
	// 	p.Side = side
	return score
}

// TODO test coverage here!!!!
func PstScore(p *Pos, gamestage int) (score int) { // actually Pst and material score

	piece := 0
	for _, i := range GRID {
		piece = p.Board[i]
		if piece != EMPTY {
			if Side(piece) == p.Side {
				score += (csshash[piece] + Pst[gamestage][piece][i])
			} else {
				score -= (csshash[piece] + Pst[gamestage][piece][i])
			}
		}
	}
	// opponant is in check
	if p.InCheck == Xside(p.Side) {
		score += CHECK
	}
	// i am in check :-(
	if p.InCheck == p.Side {
		score -= CHECK
	}
	return score
}

func loadPst(i []int) (board [128]int) {
	k := 0
	for _, j := range REVGRID { // beacuse pst below is reversed
		board[j] = i[k]
		k++
	}
	return
}

func loadPstRev(i []int) (board [128]int) {
	k := 0
	for _, j := range GRID { // beacuse pst below is reversed
		board[j] = i[k]
		k++
	}
	return
}

// piece square table is for 3 phases of game, 12 different pieces and the game board squares for each
var Pst [3][16][128]int

// set up piece square table: Note this is like display and not how internal rep is!
func init() {
	var i []int

	// pawn w
	i = []int{
		0, 0, 0, 0, 0, 0, 0, 0,
		50, 50, 50, 50, 50, 50, 50, 50,
		10, 10, 20, 30, 30, 20, 10, 10,
		5, 5, 10, 25, 25, 10, 5, 5,
		0, 0, 0, 20, 20, 0, 0, 0,
		5, -5, -10, 0, 0, -10, -5, 5,
		5, 10, 10, -20, -20, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	Pst[OPENING][PAWN] = loadPst(i)
	Pst[MIDGAME][PAWN] = loadPst(i)
	Pst[ENDGAME][PAWN] = loadPst(i)
	Pst[OPENING][pawn] = loadPstRev(i) // pawn b
	Pst[MIDGAME][pawn] = loadPstRev(i)
	Pst[ENDGAME][pawn] = loadPstRev(i)

	//// knight w
	i = []int{
		-50, -40, -30, -30, -30, -30, -40, -50,
		-40, -20, 0, 0, 0, 0, -20, -40,
		-30, 0, 15, 15, 15, 15, 0, -30,
		-30, 5, 15, 20, 20, 15, 5, -30,
		-30, 0, 15, 20, 20, 15, 0, -30,
		-30, 5, 15, 15, 15, 15, 5, -30,
		-40, -20, 0, 5, 5, 0, -20, -40,
		-50, -40, -30, -30, -30, -30, -40, -50,
	}
	Pst[OPENING][NIGHT] = loadPst(i)
	Pst[MIDGAME][NIGHT] = loadPst(i)
	Pst[ENDGAME][NIGHT] = loadPst(i)
	Pst[OPENING][night] = loadPstRev(i) // knight b
	Pst[MIDGAME][night] = loadPstRev(i)
	Pst[ENDGAME][night] = loadPstRev(i)

	//// bishop w
	i = []int{
		-20, -10, -10, -10, -10, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 10, 10, 5, 0, -10,
		-10, 5, 5, 10, 10, 5, 5, -10,
		-10, 0, 10, 10, 10, 10, 0, -10,
		-10, 10, 10, 10, 10, 10, 10, -10,
		-10, 5, 0, 0, 0, 0, 5, -10,
		-20, -10, -10, -10, -10, -10, -10, -20,
	}

	Pst[OPENING][BISHOP] = loadPst(i)
	Pst[MIDGAME][BISHOP] = loadPst(i)
	Pst[ENDGAME][BISHOP] = loadPst(i)
	Pst[OPENING][bishop] = loadPstRev(i) // bishop b
	Pst[MIDGAME][bishop] = loadPstRev(i)
	Pst[ENDGAME][bishop] = loadPstRev(i)

	////rook w
	i = []int{
		0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, 10, 10, 10, 10, 5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		0, 0, 0, 5, 5, 0, 0, 0,
	}

	Pst[OPENING][ROOK] = loadPst(i)
	Pst[MIDGAME][ROOK] = loadPst(i)
	Pst[ENDGAME][ROOK] = loadPst(i)
	Pst[OPENING][rook] = loadPstRev(i) // rook b
	Pst[MIDGAME][rook] = loadPstRev(i)
	Pst[ENDGAME][rook] = loadPstRev(i)

	////queen w
	i = []int{
		-20, -10, -10, -5, -5, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 5, 5, 5, 0, -10,
		-5, 0, 5, 5, 5, 5, 0, -5,
		0, 0, 5, 5, 5, 5, 0, -5,
		-10, 5, 5, 5, 5, 5, 0, -10,
		-10, 0, 5, 0, 0, 0, 0, -10,
		-20, -10, -10, -5, -5, -10, -10, -20,
	}

	Pst[OPENING][QUEEN] = loadPst(i)
	Pst[MIDGAME][QUEEN] = loadPst(i)
	Pst[ENDGAME][QUEEN] = loadPst(i)
	Pst[OPENING][queen] = loadPstRev(i) // queen b
	Pst[MIDGAME][queen] = loadPstRev(i)
	Pst[ENDGAME][queen] = loadPstRev(i)

	//// kpc king w opening
	i = []int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		20, 50, 10, 0, 0, 10, 50, 20,
	}
	Pst[OPENING][KING] = loadPst(i)
	Pst[OPENING][king] = loadPstRev(i)

	//// king w middle game
	i = []int{
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-20, -30, -30, -40, -40, -30, -30, -20,
		-10, -20, -20, -20, -20, -20, -20, -10,
		20, 20, 0, 0, 0, 0, 20, 20,
		20, 30, 10, 0, 0, 10, 30, 20,
	}
	Pst[MIDGAME][KING] = loadPst(i)
	Pst[MIDGAME][king] = loadPstRev(i)

	//// king w end game
	i = []int{
		-50, -40, -30, -20, -20, -30, -40, -50,
		-30, -20, -10, 0, 0, -10, -20, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -30, 0, 0, 0, 0, -30, -30,
		-50, -30, -30, -30, -30, -30, -30, -50,
	}
	Pst[ENDGAME][KING] = loadPst(i)
	Pst[ENDGAME][king] = loadPstRev(i)
	//         fmt.Println(Pst)
}

func Gamestage(p *Pos) int {
	if p.TakenPieces[p.Side] > 12 {
		return ENDGAME
	}
	if p.TakenPieces[p.Side] > 4 {
		return MIDGAME
	}
	return OPENING
}
