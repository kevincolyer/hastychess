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
	return PstScore(p, nummoves, gamestage)
	//		return ClaudeShannonScore(p, nummoves)
}

func EvalQ(p *Pos, nummoves, gamestage int) int {
	return PstScore(p, nummoves, gamestage) - 1 // -1 is so we don't trip the q search beta giving same score as normal eval.
	// 	return ClaudeShannonScore(p, nummoves)
}

// see https://chessprogramming.wikispaces.com/Simplified+evaluation+function for these values
// Check = king? == 20000 "which sometimes might be useful for discovering whether the king was taken"
var csshash = map[int]int{
	QUEEN: 900, queen: 900,
	ROOK: 500, rook: 500,
	BISHOP: 330, bishop: 330,
	NIGHT: 320, night: 320,
	PAWN: 100, pawn: 100,
	KING: 0, king: 0} // ignore kings as evaluated elsewhere...

// Most Valuable Victim, Least Valuable Agressor
func MVVLVA(m Move, p *Pos) int {
	if m.mtype == EPCAPTURE {
		return 0
	}
	if m.mtype == PROMOTE {
		return -csshash[p.Board[m.from]] + csshash[m.extra] + csshash[p.Board[m.to]] // If promotion then m.extra has value of piece we promote to, otherwise it is 0
	}
	// standard capture
	return -csshash[p.Board[m.from]] + csshash[p.Board[m.to]] // pxQ == 800 Qxp==-800
	// 	return -csshash[p.Board[m.to]] + csshash[m.extra] + csshash[p.Board[m.from]] // If promotion then m.extra has value of piece we promote to, otherwise it is 0
	//     return csshash[p.Board[m.to]]-csshash[p.Board[m.from]] // If promotion then m.extra has value of piece we promote to, otherwise it is 0
}

// Cargo culted from see.cpp of CPW engine
/******************************************************************************
*  This is not yet proper static exchange evaluation, but an approximation    *
*  proposed by Harm Geert Mueller under the acronym BLIND (better, or lower   *
*  if not defended. As the name indicates, it detects only obviously good     *
*  captures, but it seems enough to improve move ordering.                    *
******************************************************************************/
func BLIND(m Move, p *Pos) bool {

	/* captures by pawn do not lose material */
	if PieceType(p.Board[m.from]) == PAWN {
		return true
	}

	// BETTER
	/* Captures "lower takes higher" (as well as BxN) are good by definition. */
	if csshash[p.Board[m.to]] >= csshash[p.Board[m.from]]-50 {
		return true
	}

	// LOWER if not guarded ie. QxP guarded by p
	// 	/* Make the first capture, so that X-ray defender show up*/
	from := p.Board[m.from]
	p.Board[m.from] = EMPTY
	/* Captures of undefended pieces are good by definition */
	if !IsAttacked(m.to, 1-p.Side, p) { // need a better IsAttacked function.
		p.Board[m.from] = from
		return true
	}
	p.Board[m.from] = from
	return false // of other captures we know nothing, Jon Snow!
}

// TODO test coverage here!!!!
func PstScore(p *Pos, nummoves, gamestage int) (score int) { // actually Pst and material score
	var piece, up, dd, ss, pawnscore int
	//         var piece int
	for _, i := range GRID {
		piece = p.Board[i]
		if piece != EMPTY {
			if Side(piece) == p.Side {
				score += (csshash[piece] + Pst[gamestage][piece][i])
			} else {
				score -= (csshash[piece] + Pst[gamestage][piece][i])
			}
		}

		// Pawn mobility additions
		if piece&0x7 != PAWN {
			continue
		}

		// pawn mobility or blockage
		pawnscore = 0
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
			// one of our pawns
			if p.Board[dd] == piece {
				pawnscore -= 50
			}
			dd += up
		}

		//ss - blocked PAWNS TODO blocked by what?
		if Onboard(ss) && p.Board[ss] != EMPTY && p.Board[ss] != piece {
			pawnscore -= 50
		}
		//ii - isolated
		for _, j := range QM {
			if Onboard(i+j) && p.Board[i+j] == piece {
				// found a fellow pawn so not isolated
				break
			}
			// isolated penalty
			pawnscore -= 50
		}
		if Side(piece) == p.Side {
			score += pawnscore
		} else {
			score -= pawnscore
		}
	}

	// M-M'
	// just a rough estimate of how many moves...
	/*	if nummoves>0 {
		    incheck:=p.InCheck
		    enpassant:=p.EnPassant
		    p.InCheck=-1
		    p.EnPassant = 0
		    p.Side=1-p.Side

		    score += (nummoves - len(GenerateAllMoves(p))) * 10

		    p.Side=1-p.Side
		    p.EnPassant=enpassant
		    p.InCheck=incheck
		}
	*/
	//      Consider check else where - the King score is used for pl move generation. giving this as the final score means it is check or nothing!
	//      use killer or killer-mate?
	// opponent is in check
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

	////rook w - opening
	// Attempt to get the rook to stay put on it's square and not shift about before castling
	i = []int{
		0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, 10, 10, 10, 10, 5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		60, 0, 0, 5, 5, 0, 0, 60,
	}
	Pst[OPENING][ROOK] = loadPst(i)
	Pst[OPENING][rook] = loadPstRev(i) // rook b

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

	Pst[MIDGAME][ROOK] = loadPst(i)
	Pst[ENDGAME][ROOK] = loadPst(i)
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
	// what about castling?
	if p.TakenPieces[p.Side] > 10 {
		return ENDGAME
	}
	if p.Castled[p.Side] || p.TakenPieces[p.Side] > 4 {
		return MIDGAME
	}
	return OPENING
}
