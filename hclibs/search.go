//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import (
	"fmt"
	"sort"
	"time"
)

func Milliseconds(d time.Duration) int {
	return int(d.Nanoseconds() / 1000000)
}

// THIS SHOULD BE SEARCH ROOT OR A WAY TO ALLOW ME TO PLUG IN DIFFERENT SEARCHES
func SearchRoot(p *Pos, maxdepth int, globalpv *PV, starttime time.Time) (bestmove Move, bestscore int) {
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
	// reset history table or age it?
	// reset killers

	consider := GenerateAllMoves(p)

	// 1a. check that we are not in checkmate or stalemate
	if len(consider) == 0 {
		// 		if p.InCheck == p.Side {
		// 			bestscore = -CHECKMATE
		// 		} else {
		// 			bestscore = -STALEMATE
		// 		}
		return // result routine works out what the win/lose/draw state is.
	}

	// 2. give a rough order
	// 	OrderMoves(consider, p)
	OrderMoves(&consider, p, globalpv)
	// 	fmt.Println(" ",consider[0])
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
	searchdepth := 0
	// 	for depth := 2; depth < maxdepth+1; depth++ {
	for depth := maxdepth; depth < maxdepth+1; depth++ {
		enterquiesce := (depth == maxdepth)
		childpv := PV{ply: p.Ply + 1}
		count := 0

		// from CPW: if only one move and searched 4 deep then move...
		if len(consider) == 1 && depth > 4 {
			break
		}

		if GameProtocol == PROTOCONSOLE {
			fmt.Printf("# Searching to depth %v\n", depth)
		}
		for _, move := range consider {
			//negamax sorts ENTIRE search space! With iterative deepening and some pruning we can cut the search space down.
			// so if done shallow search and looked at about 4 moves already and current move looks no better than best break and search deeper...
			// 			if depth > 2 && count > 3 && move.score < bestscore+25 {
			// 				break
			// 			}

			MakeMove(move, p)
			// need neg here as we switch sides in make move and evaluation happens relative to side
			val := -negamaxab(alpha, beta, depth, p, &childpv, enterquiesce, searchdepth+1)
			//fmt.Printf("# move %v scored %v\n", move, val)
			UnMakeMove(move, p)
			// update for next round of sorting when iterative deepening. Do after unmakemove as the move score change is recorded in history array
			move.score = val

			if StatNodes > PREVENTEXPLOSION || StopSearch() {
				return
			}

			if val >= bestscore {
				bestmove = move
				bestscore = val
				//update PV (stack based)
				globalpv.moves[0] = bestmove
				copy(globalpv.moves[1:], childpv.moves[:])
				globalpv.count = childpv.count + 1

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
		//if bestscore <= alpha -50 || bestscore >= beta+50 {
		alpha = NEGINF
		beta = POSINF
		//	depth-- // search again but deeper
		//	fmt.Println("# Re-run this search but with a wide window")
		// 		} else {
		//
		// 			alpha = consider[0].score - 50
		// 			beta = consider[0].score + 50
		// 		}
		elapsed := time.Since(starttime)
		if UCI() {
			// upperbound or lowerbound or cp for exact
			fmt.Printf("info depth %v score upperbound %v time %v nodes %v nps %v pv %v\n", depth, bestscore, Milliseconds(elapsed), StatNodes+StatQNodes, int(float64(StatNodes+StatQNodes)/elapsed.Seconds()), globalpv)
		}
		if GamePostStats == true && GameProtocol == PROTOXBOARD {
			// ply	Integer score Integer in centipawns.time in centiseconds (ex:1028 = 10.28 seconds). nodes Nodes searched. pv freeform
			fmt.Printf("%v %v %v %v %v\n", depth, bestscore, float64(Milliseconds(elapsed)/100), StatNodes+StatQNodes, globalpv)
		}
	}
	return
}

// negamaxab search: searches between window so prunes the search tree.
func negamaxab(alpha, beta, depth int, p *Pos, parentpv *PV, enterquiesce bool, searchdepth int) int {
	// 	var childpv PV
	childpv := PV{ply: p.Ply + 1}

	StatNodes++
	if depth == 0 || StatNodes > PREVENTEXPLOSION {
		parentpv.count = 0 // reset because we are at a leaf...

		// Quiese?
		if enterquiesce && StatNodes < PREVENTEXPLOSION {
			// enter q search at this level - not deeper so no need to invert.
			return SearchQuiesce(p, alpha, beta, QUIESCEDEPTH, searchdepth)
		}
		return Eval(p, 0, Gamestage(p))
	}
	// need to know if we are in mate before we return an eval at a leaf as this is the only way we check for mate!!! Not done in eval!

	consider := GenerateAllMoves(p)
	if len(consider) == 0 {
		if p.InCheck != -1 {
			//                         fmt.Println("found a checkmate")
			return -CHECKMATE + searchdepth
		}
		//                 fmt.Println("found a stalemate")
		return -STALEMATE + searchdepth // Unless we can't win because of lack of material then stalemate makes sense -- impliment this!
	}
	max := NEGINF
	OrderMoves(&consider, p, parentpv)

	// reset PV and choose the
	bestmove := consider[0]
	childpv.moves[0] = bestmove // in case we don't find anything better set first move to return
	childpv.count = 1
	count := 0
	for _, move := range consider {
		// prevent search explosion by reducing search width with increasing depth - probably doesn't work with iterative deepening...
		if count > MAXSEARCHDEPTH-searchdepth {
			break
		}
		MakeMove(move, p)
		score := -negamaxab(-beta, -alpha, depth-1, p, &childpv, enterquiesce, searchdepth+1)
		UnMakeMove(move, p)
		count++

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
			// set killers here
			// history table here...
			return beta
		}
		if StopSearch() {
			return alpha
		}
	}
	return alpha
}

func SearchQuiesce(p *Pos, alpha, beta int, qdepth int, searchdepth int) int {
	gamestage := Gamestage(p)
	StatQNodes++

	// need a standpat score
	val := EvalQ(p, 0, gamestage) // custom evaluator here for QUIESENCE TODO
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
	if qdepth == 0 || StopSearch() || StatQNodes > PREVENTEXPLOSION/4 {
		// 		fmt.Println("# Qnode explosion - bottling!")
		return alpha
	}

	// get moves - but only captures and promotions
	moves := GenerateMovesForQSearch(p)
	// nothing more to search...
	if len(moves) == 0 {
		return alpha
	}

	// score by BLIND - Better or Lower If Not Defended
	// 	for i := range moves {
	// 		if BLIND(moves[i], p) {
	// 			moves[i].score = 100
	//
	// 		}
	// 	}

	// score them by Most Valuable Victim - Least Valuable Aggressor
	for i := range moves {
		moves[i].score = MVVLVA(moves[i], p)
	}

	// And order descending to provoke cuts
	sort.Slice(moves, func(i, j int) bool { return moves[i].score > moves[j].score }) // by score type descending

	// loop over all moves, searching deeper until no moves left and all is "quiet" - return this score...)
	for _, m := range moves {
		//              No premature optimisation or cargo culting!

		// 		// adjust each score for delta cut offs and badmoves skipping to next each time
		// 		// delta - if not promotion and not endgame and is a low scoring capture then don't look deeper
		// 		// delta cut qnodes from 20M to 640,000 in one case!
		if m.mtype != PROMOTE && gamestage != ENDGAME && standpat+csshash[p.Board[m.to]]+200 < alpha {
			continue
		}
		//
		// 		// badmoves - cut qnodes from 640,000 to 64,000
		// 		// capture by pawn is ok so skip
		if PieceType(p.Board[m.from]) == PAWN && m.mtype != PROMOTE {
			continue
		}

		// search deeper until quiet
		MakeMove(m, p)
		val = -SearchQuiesce(p, -beta, -alpha, qdepth-1, searchdepth+1)
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

// From https://chessprogramming.wikispaces.com/Move+Ordering
// A typical move ordering consists as follows:
// 1.PV-move of the principal variation from the previous iteration of an iterative deepening framework for the leftmost path, often implicitly done by 2.
// 2.Hash move from hash tables
// 3.Winning captures/promotions
// 4.Equal captures/promotions
// 5.Killer moves (non capture), often with mate killers first
// 6.Non-captures sorted by history heuristic and that like
// 7.Losing captures (* but see below

/*
SORT_KING 400 // check
SORT_HASH 200 // hash move
SORT_CAPT 100 // capture
SORT_PROM  90 // promote
SORT_KILL  80*/ // killer move

// usr blind here to put bad captures (blind==0) back of the queue after ordinary captures

func OrderMoves(moves *[]Move, p *Pos, pv *PV) bool {
	// order by move type (capture and promotion first down to quiet moves)
	// 	plydelta := p.Ply - pv.ply
	//         if plydelta<0 {panic("this should not be!")}
	for i := 0; i < len(*moves); i++ {

		// boost or lower captures depending on good or bad
		// boost good captures and punish bad captures
		if (*moves)[i].mtype == CAPTURE {
			if BLIND((*moves)[i], p) {
				(*moves)[i].score = p.Board[(*moves)[i].from]*2 + GOODCAPTURE
			} else {
				(*moves)[i].score = p.Board[(*moves)[i].from]*2 + BADCAPTURE
			}
			if PieceType((*moves)[i].to) != PieceType((*moves)[i].from) {
				(*moves)[i].score = p.Board[(*moves)[i].from]*2 + CAPTURE
			}
		}

		if (*moves)[i].mtype == EPCAPTURE || (*moves)[i].mtype == O_O_O || (*moves)[i].mtype == O_O {
			(*moves)[i].score = (*moves)[i].mtype + p.Board[(*moves)[i].from]*2 // type of move + which piece is moving
		}

		if (*moves)[i].mtype == QUIET || (*moves)[i].mtype == ENPASSANT {
			//for now
			(*moves)[i].score = QUIET
			// boost history
			// boost killers
			// boost check?????
		}

		// boost PV to top here
		// cycle through pv to boost all moves in current move list to top
		for _, m := range pv.moves {

			if (*moves)[i].from == m.from && (*moves)[i].to == m.to && (*moves)[i].extra == m.extra {
				(*moves)[i].score += PVBONUS
			}
		}
	}

	sort.Slice((*moves), func(i, j int) bool { return (*moves)[i].score > (*moves)[j].score }) // descending
	//         fmt.Print((*moves)[0])
	return true
}

/*
BADCAPTURE c
QUIET     *   sorted by history not piecevalue
ENPASSANT *
KILLERS   c
O_O_O     *
O_O       *
EPCAPTURE *
CAPTURE   * (equal)
GOODCAPTURE c
PROMOTE   *
PVBONUS   c
CHECK     c
*/
