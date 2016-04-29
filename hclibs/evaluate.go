package hclibs

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

var csshash = map[int]int{QUEEN: 900, queen: 900, //q =
	ROOK: 500, rook: 500,
	BISHOP: 300, bishop: 300,
	NIGHT: 300, night: 300,
	PAWN: 100, pawn: 100,
	KING: 0, king: 0} // ignore kings as evaluated elsewhere...

func ClaudeShannonScore(p *Pos, totalmoves int) int {
	side := p.Side
	xside := 1 - side
	incheck := p.InCheck
// 	enpassant := p.EnPassant
	score := 0
	var piece, up, dd, ss int
	if incheck == xside {
		score += CHECK
	}
	if incheck == side {
		score -= CHECK
	} // K-K'
	// Material valuation
	// could add bonuses at different game stages...
	for _, i := range GRID {
		piece = p.Board[i]
		if piece != EMPTY {
			if (piece >> 3) == side {
				score += csshash[piece]
			} else {
				score -= csshash[piece]
			}
			
			// Pawn mobility additions
			if piece&0x7 == PAWN {
				if (piece >> 3) == WHITE {
					up = NORTH
				} else {
					up = SOUTH
				}
				dd = i + up
				ss = dd
				// duubled
				for dd&0x88 == 0 {
					if p.Board[dd] == piece {
						score -= 50
					}
					dd += up
				}
				//ss - blocked TODO blocked by what?
				if ss&0x88 == 0 && EMPTY != p.Board[ss] && piece != p.Board[ss] {
					score -= 50
				} // blocked by opposite
				//ii - isolated
				for _, j := range QM {
					if (i+j)&0x88 == 0 && piece == p.Board[i+j] {
						score -= 50
						break
					}
				}
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
