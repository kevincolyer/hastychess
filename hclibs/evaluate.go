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

// Most Valuable Victim, Least Valuable Agressor
func MVVLVA(m Move, p *Pos) int {
	return csshash[p.Board[m.to]] + csshash[m.extra] - csshash[p.Board[m.from]] // If promotion then m.extra has value of piece we promote to, otherwise it is 0
	//     return csshash[p.Board[m.to]]-csshash[p.Board[m.from]] // If promotion then m.extra has value of piece we promote to, otherwise it is 0
}

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

// sub pst_score ( Int $gamestage, Position $p) { // actually pst and material score
//     ENTER { if $TUNING { my $n ='evaluate::pst_score' ;%tuning{$n}[2]= now; %tuning{$n}[0] //=0; %tuning{$n}[1] //=0 } }
//     LEAVE { if $TUNING { my $n ='evaluate::pst_score' ;my $dur = now - %tuning{$n}[2]; %tuning{$n}[0]+=$dur; %tuning{$n}[1]++; } }
// 
//     my Int $side=$p.side;
//     my Int $xside=$p.xside;
//     my Int $incheck=$p.in_check;
// #     my Int $enpassant=$p.en_passant;
//     my Int $score=0;
//     my Int $piece;
// #     my Int $gamestage=gamestage($p);
// #     say "css - entered css " ~ $p.board;
// #     my Int $flip = $side == WHITE ?? 1 !! -1; # when asses from other side
//     for @grid -> Int $i {
//         $piece = $p.board[ $i ];
//         if $piece != EMPTY {
//             if  ($piece +> 3) == $side {
//                 $score +=  %csshash{$piece} + @pst[ $gamestage ][ $piece ][ $i ]
//             } else {
//                 $score -=  %csshash{$piece} - @pst[ $gamestage ][ $piece ][ $i ]
//             }
//         }
//     }
//     return $score + ( ($incheck==$xside) - ($incheck==$side) ) * 20_000  # K-K' = ok if static position being analysed
// }
// 
// sub pst_score_delta (Int $from, Int $to, Int $gamestage, Position $p) returns Int {
//     ENTER { if $TUNING { my $n ='evaluate::pst_score_delta' ;%tuning{$n}[2]= now; %tuning{$n}[0] //=0; %tuning{$n}[1] //=0 } }
//     LEAVE { if $TUNING { my $n ='evaluate::pst_score_delta' ;my $dur = now - %tuning{$n}[2]; %tuning{$n}[0]+=$dur; %tuning{$n}[1]++; } }
//     #fake a move
//     my Int $score=0;
//     my Int $fp=$p.board[$from];
//     my Int $tp=$p.board[$to];
// 
//     $p.board[$from]=EMPTY;
//     $p.board[$to]=$fp;
//     $score += 20_000 * king_is_in_check( $p.king[$p.xside], $p );
// 
//     # restore
//     $p.board[$from]=$fp;
//     $p.board[$to]=$tp;
// 
//     # use pst: add the $to square moving score for the piece moving (for the gamestage)
//     $score += @pst[ $gamestage ][ $fp ][ $to ];
//     $score -= @pst[ $gamestage ][ $fp ][ $from ];
// 
// #     # if pawn promotes
// #     if $fp==P and $from +& 0x7 == 7  {
// #         ## argh want type here!
// #     }
// #     # if castling ...
// 
//     # add the piece score if we are capturing
//     $score += %csshash{ $tp } if $tp != EMPTY;
//     return $score;
// }
// 
// sub _load_pst (@i is copy) {
//     my Int @board = 0 xx 0x88;
//     for @grid -> Int $j {
//         @board[$j]= shift @i;
//     }
//     return @board;
// }
// 
// {
//     my Int @i;
// 
//     #// pawn w
//     @i=
//     0,  0,  0,  0,  0,  0,  0,  0,
//     50, 50, 50, 50, 50, 50, 50, 50,
//     10, 10, 20, 30, 30, 20, 10, 10,
//     5,  5, 10, 25, 25, 10,  5,  5,
//     0,  0,  0, 20, 20,  0,  0,  0,
//     5, -5,-10,  0,  0,-10, -5,  5,
//     5, 10, 10,-20,-20, 10, 10,  5,
//     0,  0,  0,  0,  0,  0,  0,  0
//     ;
//     @pst[$_][P]= [ _load_pst(@i)         ] for OPENING, MIDGAME, ENDGAME;
//     @pst[$_][p]= [ _load_pst(@i.reverse) ] for OPENING, MIDGAME, ENDGAME; # symetrical.
// 
//     #// knight w
//     @i=
//     -50,-40,-30,-30,-30,-30,-40,-50,
//     -40,-20,  0,  0,  0,  0,-20,-40,
//     -30,  0, 10, 15, 15, 10,  0,-30,
//     -30,  5, 15, 20, 20, 15,  5,-30,
//     -30,  0, 15, 20, 20, 15,  0,-30,
//     -30,  5, 10, 15, 15, 10,  5,-30,
//     -40,-20,  0,  5,  5,  0,-20,-40,
//     -50,-40,-30,-30,-30,-30,-40,-50,
//     ;
//     @pst[$_][N]= [ _load_pst(@i)         ] for OPENING, MIDGAME, ENDGAME;
//     @pst[$_][n]= [ _load_pst(@i.reverse) ] for OPENING, MIDGAME, ENDGAME; # symetrical.
// 
//     #// bishop w
//     @i=
//     -20,-10,-10,-10,-10,-10,-10,-20,
//     -10,  0,  0,  0,  0,  0,  0,-10,
//     -10,  0,  5, 10, 10,  5,  0,-10,
//     -10,  5,  5, 10, 10,  5,  5,-10,
//     -10,  0, 10, 10, 10, 10,  0,-10,
//     -10, 10, 10, 10, 10, 10, 10,-10,
//     -10,  5,  0,  0,  0,  0,  5,-10,
//     -20,-10,-10,-10,-10,-10,-10,-20,
//     ;
//     @pst[$_][B]= [ _load_pst(@i)         ] for OPENING, MIDGAME, ENDGAME;
//     @pst[$_][b]= [ _load_pst(@i.reverse) ] for OPENING, MIDGAME, ENDGAME; # symetrical.
// 
//     #//rook w
//     @i=
//     0,  0,  0,  0,  0,  0,  0,  0,
//     5, 10, 10, 10, 10, 10, 10,  5,
//     -5,  0,  0,  0,  0,  0,  0, -5,
//     -5,  0,  0,  0,  0,  0,  0, -5,
//     -5,  0,  0,  0,  0,  0,  0, -5,
//     -5,  0,  0,  0,  0,  0,  0, -5,
//     -5,  0,  0,  0,  0,  0,  0, -5,
//     0,  0,  0,  5,  5,  0,  0,  0
//     ;
//     @pst[$_][R]= [ _load_pst(@i)         ] for OPENING, MIDGAME, ENDGAME;
//     @pst[$_][r]= [ _load_pst(@i.reverse) ] for OPENING, MIDGAME, ENDGAME; # symetrical.
// 
//     #//queen w
//     @i=
//     -20,-10,-10, -5, -5,-10,-10,-20,
//     -10,  0,  0,  0,  0,  0,  0,-10,
//     -10,  0,  5,  5,  5,  5,  0,-10,
//     -5,  0,  5,  5,  5,  5,  0, -5,
//     0,  0,  5,  5,  5,  5,  0, -5,
//     -10,  5,  5,  5,  5,  5,  0,-10,
//     -10,  0,  5,  0,  0,  0,  0,-10,
//     -20,-10,-10, -5, -5,-10,-10,-20
//     ;
//     @pst[$_][Q]= [ _load_pst(@i)         ] for OPENING, MIDGAME, ENDGAME;
//     @pst[$_][q]= [ _load_pst(@i.reverse) ] for OPENING, MIDGAME, ENDGAME; # symetrical. ]=
//     #// kpc king w opening
//     @i=
//     0,0,0,0,0,0,0,0,
//     0,0,0,0,0,0,0,0,
//     0,0,0,0,0,0,0,0,
//     0,0,0,0,0,0,0,0,
//     0,0,0,0,0,0,0,0,
//     0,0,0,0,0,0,0,0,
//     20, 20,  0,  0,  0,  0, 20, 20,
//     20, 30, 10,  0,  0, 10, 30, 20
//     ;
//     @pst[OPENING][K]= [ _load_pst(@i)         ] ;
//     @pst[OPENING][k]= [ _load_pst(@i.reverse) ] ; # symetrical.
//     #// king w middle game
//     @i=
//     -30,-40,-40,-50,-50,-40,-40,-30,
//     -30,-40,-40,-50,-50,-40,-40,-30,
//     -30,-40,-40,-50,-50,-40,-40,-30,
//     -30,-40,-40,-50,-50,-40,-40,-30,
//     -20,-30,-30,-40,-40,-30,-30,-20,
//     -10,-20,-20,-20,-20,-20,-20,-10,
//     20, 20,  0,  0,  0,  0, 20, 20,
//     20, 30, 10,  0,  0, 10, 30, 20
//     ;
//     @pst[MIDGAME][K]= [ _load_pst(@i)         ] ;
//     @pst[MIDGAME][k]= [ _load_pst(@i.reverse) ] ; # symetrical.
// 
//     #// king w end game
//     @i=
//     -50,-40,-30,-20,-20,-30,-40,-50,
//     -30,-20,-10,  0,  0,-10,-20,-30,
//     -30,-10, 20, 30, 30, 20,-10,-30,
//     -30,-10, 30, 40, 40, 30,-10,-30,
//     -30,-10, 30, 40, 40, 30,-10,-30,
//     -30,-10, 20, 30, 30, 20,-10,-30,
//     -30,-30,  0,  0,  0,  0,-30,-30,
//     -50,-30,-30,-30,-30,-30,-30,-50
//     ;
//     @pst[ENDGAME][K]= [ _load_pst(@i)         ] ;
//     @pst[ENDGAME][k]= [ _load_pst(@i.reverse) ] ; # symetrical.#!/usr/bin/env perl6
// }