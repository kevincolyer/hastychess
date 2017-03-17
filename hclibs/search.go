//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import (
	"fmt"
	"sort"
        "time"
)

func  Milliseconds(d time.Duration) int {
    return int(d.Nanoseconds()/1000000)
}

// THIS SHOULD BE SEARCH ROOT OR A WAY TO ALLOW ME TO PLUG IN DIFFERENT SEARCHES
func SearchRoot(p *Pos, maxdepth int, globalpv *PV,starttime time.Time) (bestmove Move, bestscore int) {
	/*
	   // 1. get all moves to consider
	   // 1a. check that we are not in checkmate or stalemate
	   // 2. give a rough order
	   // 3a. if iterative deepening loop from depth 2 to max depth in turn, sorting best score descending
	   // 4. loop all moves to consider
	   // 5. make move
	   // 6. negamax search the moves (or negamaxab or negascout) to depth
	   // 7. if found best move record score and move to return
	   // 8. loop back to 4 or 3a.
	   // 9. return best move and score
	*/
	consider := GenerateAllMoves(p)

	// 1a. check that we are not in checkmate or stalemate
	if len(consider) == 0 {
		if p.InCheck == p.Side {
			bestscore = CHECKMATE
		} else {
			bestscore = STALEMATE
		}
		return
	}

	// 2. give a rough order
	// 	OrderMoves(consider, p)
	OrderMoves(consider, p, globalpv)
	if GameProtocol == PROTOCONSOLE {
		fmt.Printf("# moves to consider: %v\n", consider)
	}
	alpha := NEGINF
	beta := POSINF
	bestscore = NEGINF
	bestmove = consider[0]
	// reset pv
	globalpv.count = 1
	globalpv.moves[0] = bestmove

	// 3a. if iterative deepening loop from depth 2 to max depth in turn, sorting best score descending
	//depth := maxdepth
	globalpv.ply = p.Ply // syncronise new PV

	for depth := 2; depth < maxdepth+1; depth++ {
		enterquiesce := (depth == maxdepth)
		childpv := PV{ply: p.Ply + 1}
		count := 0
		
		if GameProtocol == PROTOCONSOLE {
			fmt.Printf("# Searching to depth %v\n", depth)
		}
		for _, move := range consider {
			//negamax sorts ENTIRE search space! With iterative deepening and some pruning we can cut the search space down.
			// so if done shallow search and looked at about 4 moves already and current move looks no better than best break and search deeper...
			if depth > 2 && count > 4 && move.score < bestscore+25 {
				break
			}
			MakeMove(move, p)
			val := -negamaxab(alpha, beta, depth, p, &childpv, enterquiesce) // need neg here as we switch sides in make move and evaluation happens relative to side
			//fmt.Printf("# move %v scored %v\n", move, val)
			UnMakeMove(move, p)
			move.score = val // update for next round of sorting when iterative deepening. Do after unmakemove as the move score change is recorded in history array
                        
			if StatNodes > PREVENTEXPLOSION {
				return
			}

			if val >= bestscore {
				bestmove = move
				bestscore = val
				//fmt.Printf("# found bestscore %v move %v\n", bestscore, bestmove)
				//update PV (stack based)
				globalpv.moves[0] = bestmove
				copy(globalpv.moves[1:], childpv.moves[:])
				globalpv.count = childpv.count + 1
				//                         if val==bestscore {fmt.Printf("# depth: %v score: %v pv: %v\n", depth,bestscore,globalpv)}
				if GameProtocol == PROTOCONSOLE {
					fmt.Printf("# depth: %v score: %v pv: %v\n", depth, bestscore, globalpv)
				}
				
			}
			if val > alpha {
				alpha = val
			}
			// stop search as found better
			if val > beta {
				break
			}

			count++
		}
		// re-sort for next loop when iterative deepening

		sort.Slice(consider, func(i, j int) bool { return consider[i].score > consider[j].score }) // by score type descending
		if bestscore < alpha || bestscore > beta {
			alpha = NEGINF
			beta = POSINF
			depth-- // search again but deeper
		} else {

			alpha = consider[0].score - 50
			beta = consider[0].score + 50
		}
                elapsed:=time.Since(starttime)
		if UCI()  {
                                    fmt.Printf("info depth %v score cp %v time %v nodes %v nps %v pv %v\n",depth,bestscore, Milliseconds(elapsed), StatNodes+StatQNodes, int(float64(StatNodes+StatQNodes)/elapsed.Seconds()),globalpv)
                                }
	}
	return
}

// negamaxab search: searches between window so prunes the search tree.
func negamaxab(alpha, beta, depth int, p *Pos, parentpv *PV, enterquiesce bool) int {
	// 	var childpv PV
	childpv := PV{ply: p.Ply + 1}

	StatNodes++
	if depth == 0 || StatNodes > PREVENTEXPLOSION {
		parentpv.count = 0 // reset because we are at a leaf...

		// Quiese?
		if enterquiesce {
			return SearchQuiesce(p, alpha, beta, QUIESCEDEPTH)
		} else {
			return Eval(p, 1, Gamestage(p))
		}
	}
	max := NEGINF

	consider := GenerateAllMoves(p)
	if len(consider) == 0 {
		if p.InCheck > 0 {
			return CHECKMATE
		}
		return STALEMATE
	}

	OrderMoves(consider, p, parentpv)

	// reset PV and choose the
	bestmove := consider[0]
	parentpv.moves[0] = bestmove // in case we don't find anything better set first move to return
	parentpv.count = 1

	for _, move := range consider {

		MakeMove(move, p)
		score := -negamaxab(-beta, -alpha, depth-1, p, &childpv, enterquiesce)
		UnMakeMove(move, p)

		if score > max {
			max = score
			bestmove = move

			// update PV on stack
			parentpv.moves[0] = bestmove
			copy(parentpv.moves[1:], childpv.moves[:])
			parentpv.count = childpv.count + 1
		}
		if max > alpha {
			alpha = max
		}
		if alpha >= beta {
			return beta
		}
	}
	return max
}

func SearchQuiesce(p *Pos, alpha, beta int, qdepth int) int {
	gamestage := Gamestage(p)
	StatQNodes++

	// need a standpat score
	val := EvalQ(p, 1, gamestage) // custom evaluator here for QUIESENCE TODO
	standpat := val

	// is move worse than previous worst?
	if val >= beta {
		return beta
	}
	// is move less good than previous best?
	if alpha <= val {
		alpha = val
	}

	// finish searching:
	// to prevent search explosion while testing TODO remove!
	// when at end of search
	// someone signals we should stop
	if StatQNodes > PREVENTEXPLOSION/4 || qdepth == 0 || StopSearch() {
		// 		fmt.Println("# Qnode explosion - bottling!")
		return alpha
	}

	// get moves - but only captures and promotions
	moves := GenerateMovesForQSearch(p)
	// nothing more to search...
	if len(moves) == 0 {
		return alpha
	}

	// score them by Most Valuable Victim - Least Valuable Aggressor
	for i := range moves {
		moves[i].score = MVVLVA(moves[i], p)
	}
	// And order descending to provoke cuts
	sort.Slice(moves, func(i, j int) bool { return moves[i].score > moves[j].score }) // by score type descending

	// loop over all moves, searching deeper until no moves left and all is "quiet" - return this score...)
	for _, m := range moves {
		// adjust each score for delta cut offs and badmoves skipping to next each time
		// delta - if not promotion and not endgame and is a low scoring capture then continue
		// delta cut qnodes from 20M to 640,000 in one case!
		if m.mtype != PROMOTE && gamestage != ENDGAME && standpat+csshash[p.Board[m.to]]+200 < alpha {
			continue
		}

		// badmoves - cut qnodes from 640,000 to 64,000
		// capture by pawn is ok so skip
		if p.Board[m.from]&7 == PAWN && m.mtype != PROMOTE {
			continue
		}

		// search deeper until quiet
		MakeMove(m, p)
		val = -SearchQuiesce(p, -beta, -alpha, qdepth-1)
		UnMakeMove(m, p)

		// adjust window
		if val >= alpha {
			if val > beta {
				return beta
			}
			alpha = val
		}
	}
	return alpha
}

func OrderMoves(moves []Move, p *Pos, pv *PV) bool {
	// order by move type (capture and promotion first down to quiet moves)
	// 	plydelta := p.Ply - pv.ply
	//         if plydelta<0 {panic("this should not be!")}
	for i := range moves {
		moves[i].score = moves[i].mtype + p.Board[moves[i].from]*2 // type of move + which piece is moving

		// boost PV to top here
		// if moves[i].from == pv.moves[plydelta].from && moves[i].to == pv.moves[plydelta].to && moves[i].extra == pv.moves[plydelta].extra {

		// cycle through pv to boost all moves in current move list to top
		for _, m := range pv.moves {

			if moves[i].from == m.from && moves[i].to == m.to && moves[i].extra == m.extra {
				moves[i].score += PVBONUS
			}
		}
	}

	sort.Slice(moves, func(i, j int) bool { return moves[i].score > moves[j].score }) // descending
	return true
}
