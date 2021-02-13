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
func SearchRoot(p *Pos, srch *Search) (bestmove Move, bestscore int) {
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
		return // result routine works out what the win/lose/draw state is.
	}

	// 2. give a rough order
	srch.Info = fmt.Sprintf("moves to consider: %v\n", consider)
	OrderMoves(&consider, p, srch.PV)
// 	srch.Info += fmt.Sprintf("moves sorted     : %v\n", consider)
    srch.Info += fmt.Sprintf("moves sorted     : ")
    for i:=range consider {
     srch.Info += fmt.Sprintf("%v(%v) ", consider[i],consider[i].score)   
    }
    srch.Info += "\n"

	alpha := NEGINF
	beta := POSINF
	bestmove = consider[0]
	//bestscore = bestmove.score

	// reset pv
	srch.PV.count = 1
	srch.PV.moves[0] = bestmove

	// 3a. if iterative deepening loop from depth 2 to max depth in turn, sorting best score descending
	//depth := srch.MaxDepthToSearch
	srch.PV.ply = p.Ply // syncronise new PV
	searchdepth := 0
	// 	for depth := 2; depth < srch.MaxDepthToSearch+1; depth++ {
	for depth := 0; depth < srch.MaxDepthToSearch; depth++ {
		enterquiesce := (depth>2) //(depth == srch.MaxDepthToSearch)
		// create new child PV
		childpv := PV{ply: p.Ply + 1}
		count := 0

		for _, move := range consider {

			MakeMove(move, p)
			// need neg here as we switch sides in make move and evaluation happens relative to side
			val := -negamaxab(-beta, -alpha, depth, p, &childpv, enterquiesce, searchdepth+1, srch)
			//fmt.Printf("# move %v scored %v\n", move, val)
			UnMakeMove(move, p)
			// update for next round of sorting when iterative deepening. Do after unmakemove as the move score change is recorded in history array
			move.score = val
            

			if val > alpha {
				bestmove = move // (and hence score too)
                srch.Stats.AlphaRaised++
				alpha = val
				//update PV (stack based)
				copy(srch.PV.moves[1:], childpv.moves[:])
                srch.PV.count = childpv.count + 1
				srch.PV.moves[0] = bestmove
				srch.Stats.Score = bestscore
				// update PV with child PV


			}

			if srch.Stats.Nodes > srch.ExplosionLimit || srch.StopSearch() {
				return
			}

			count++
		}
		// re-sort for next loop when iterative deepening

		//sort.Slice(consider, func(i, j int) bool { return consider[i].score > consider[j].score }) // by score type descending
		//if bestscore <= alpha -50 || bestscore >= beta+50 {
		//alpha = NEGINF
		//beta = POSINF
		//	depth-- // search again but deeper
		//	fmt.Println("# Re-run this search but with a wide window")
		// 		} else {
		//
		// 			alpha = consider[0].score - 50
		// 			beta = consider[0].score + 50
		// 		}
		// 		elapsed := time.Since(srch.TimeStart)
		// 		if UCI() {
		// 			// upperbound or lowerbound or cp for exact
		// 		}
	}
	return
}

// negamaxab search: searches between window so prunes the search tree.
func negamaxab(alpha, beta, depth int, p *Pos, parentpv *PV, enterquiesce bool, searchdepth int, srch *Search) int {
	// Implimenting NegaMaxAB failsoft
	//https://www.chessprogramming.org/Alpha-Beta

    // create a new child pv to pass down
	childpv := PV{ply: p.Ply + 1}
	srch.Stats.Nodes++

	if srch.Stats.Nodes > srch.ExplosionLimit || srch.StopSearch() {
		return Eval(p, 0, Gamestage(p))
	}

	if depth == 0 {
		//parentpv.count = 0 // reset because we are at a leaf...

		// Quiese?
		if enterquiesce {
			// enter q search at this level - not deeper so no need to invert.
			return SearchQuiesce(p, alpha, beta, QUIESCEDEPTH, searchdepth, srch)
		}
		// return an exact score
		return Eval(p, 0, Gamestage(p))
	}

	consider := GenerateAllMoves(p)

	// need to know if we are in mate before we return an eval at a leaf as this is the only way we check for mate!!! Not done in eval!
	if len(consider) == 0 {
		if p.InCheck != -1 {
			//                         fmt.Println("found a checkmate")
			return -CHECKMATE + searchdepth
		}
		//                 fmt.Println("found a stalemate")
		return -STALEMATE + searchdepth // Unless we can't win because of lack of material then stalemate makes sense -- impliment this!
	}

	// give initial order for searching
	OrderMoves(&consider, p, parentpv)

	// choose the
	bestmove := consider[0]
	bestscore := NEGINF

	//reset the PV
	parentpv.moves[0] = bestmove // in case we don't find anything better set first move to return
	childpv.count = 0
	parentpv.count = 1
	count := 0

	// depth first search...
	// iterative deepening tries to convert to breadth first...
	for _, move := range consider {
		// prevent search explosion by reducing search width with increasing depth.
		// search entire space to depth of 2, then go deeper
		// May not be needed with iterative deepening...
// 		if searchdepth >= 2 && count > MAXSEARCHDEPTH-searchdepth {
// 			break
// 		}
		count++

		MakeMove(move, p)
		score := -negamaxab(-beta, -alpha, depth-1, p, &childpv, enterquiesce, searchdepth+1, srch)
		UnMakeMove(move, p)

		if score >= beta {
			// set killers here
			// this would be a killer move but the opponent wont let us reach it!
			// but lets keep it for the future.
			// history table here...
			srch.Stats.LowerCuts++
			srch.Stats.BetaRaised++
			return beta
		}

		if score > bestscore {
			bestscore = score
			bestmove = move

			// update PV on stack
			parentpv.moves[0] = bestmove
			srch.Stats.Score = bestscore
			copy(parentpv.moves[1:], childpv.moves[:])
			parentpv.count = childpv.count + 1
		}

		if bestscore > alpha {
			srch.Stats.AlphaRaised++
			alpha = bestscore
		}

		// Check if we have been asked to stop...
		if srch.StopSearch() {
			return alpha
		}
	}
	srch.Stats.UpperCuts++
	return alpha
}

func SearchQuiesce(p *Pos, alpha, beta int, qdepth int, searchdepth int, srch *Search) int {
	gamestage := Gamestage(p)
	srch.Stats.QNodes++

	// need a standpat score
	val := EvalQ(p, 0, gamestage) // custom evaluator here for QUIESENCE TODO
	//standpat := val

	// is so good our our opponent wont allow??
	if val >= beta {
		srch.Stats.LowerCuts++
		return beta
	}
	// is move better than previous best?
	if val >= alpha {
		srch.Stats.AlphaRaised++
		alpha = val
	}

	// finish searching:
	// to prevent search explosion while testing TODO remove!
	// when at end of search
	// someone signals we should stop
	if qdepth == 0 || srch.StopSearch() || srch.Stats.QNodes > srch.ExplosionLimit/4 {
		// 		fmt.Println("# Qnode explosion - bottling!")
		srch.Stats.UpperCuts++
		return alpha
	}

	// get moves - but only captures and promotions
	moves := GenerateMovesForQSearch(p)
	// nothing more to search...
	if len(moves) == 0 {
		srch.Stats.UpperCuts++
		return alpha
	}

	// score by BLIND - Better or Lower If Not Defended
	// 	for i := range moves {
	// 		if BLIND(moves[i], p) {
	// 			moves[i].score = 100
	//
	// 		}
	// 	}
    //squareControlledByOpponentPawnPenalty := 350;
	//capturedPieceValueMultiplier := 10;

	// score them by Most Valuable Victim - Least Valuable Aggressor
	for i := range moves {
		moves[i].score = MVVLVA(moves[i], p)
	}

	// And order descending to provoke cuts
	sort.Slice(moves, func(i, j int) bool { return moves[i].score > moves[j].score }) // by score type descending
    
//     for i:=range moves {
//      fmt.Printf("%vx%v(%v) ", p.Board[moves[i].from],p.Board[moves[i].to],moves[i].score  ) 
//     }
//     fmt.Println()
//    

	// loop over all moves, searching deeper until no moves left and all is "quiet" - return this score...)
	for _, m := range moves {
		//              No premature optimisation or cargo culting!

		// 		// adjust each score for delta cut offs and badmoves skipping to next each time
		// 		// delta - if not promotion and not endgame and is a low scoring capture then don't look deeper
		// 		// delta cut qnodes from 20M to 640,000 in one case!
		/*if m.mtype != PROMOTE && gamestage != ENDGAME && standpat+csshash[p.Board[m.to]]+200 < alpha {
			continue
		}*/
		//
		// 		// badmoves - cut qnodes from 640,000 to 64,000
		// 		// capture by pawn is ok so skip
 		if PieceType(p.Board[m.from]) == PAWN && m.mtype != PROMOTE {
 			continue
 		}

		// search deeper until quiet
		MakeMove(m, p)
		val = -SearchQuiesce(p, -beta, -alpha, qdepth-1, searchdepth+1, srch)
		UnMakeMove(m, p)

		// adjust window
        if val > beta {
				srch.Stats.BetaRaised++
				srch.Stats.LowerCuts++
				return beta
			}
		
		if val >= alpha {
			srch.Stats.AlphaRaised++
			alpha = val
		}
    }
	srch.Stats.UpperCuts++
	return alpha
}

// From https://www.chessprogramming.org/Move_Ordering
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

// use blind here to put bad captures (blind==0) back of the queue after ordinary captures

func OrderMoves(moves *[]Move, p *Pos, pv *PV) bool {
	// order by move type (pv, capture and promotion first down to quiet moves)
	for i := 0; i < len((*moves)); i++ {

		// boost good captures and punish bad captures
		if (*moves)[i].mtype == CAPTURE {
			mvvlva := MVVLVA((*moves)[i], p)
			if mvvlva > 0 {
				(*moves)[i].score = PieceType(p.Board[(*moves)[i].to])*2 + GOODCAPTURE
			}
			if mvvlva < 0 {
				(*moves)[i].score = PieceType(p.Board[(*moves)[i].from])*-2 + BADCAPTURE
			}
			if mvvlva == 0 {
				(*moves)[i].score = PieceType(p.Board[(*moves)[i].from])*2 + CAPTURE
			}
		}

		if (*moves)[i].mtype == EPCAPTURE || (*moves)[i].mtype == O_O_O || (*moves)[i].mtype == O_O {
			(*moves)[i].score = (*moves)[i].mtype + PieceType(p.Board[(*moves)[i].from])*2 // type of move + which piece is moving
		}

		if (*moves)[i].mtype == QUIET || (*moves)[i].mtype == ENPASSANT {
			//for now - TODO note weakness here is that is will rank Q moves in scan order! Tends to make rooks jitter about for no reason...
			// should rank by history... or hash table...
			(*moves)[i].score = QUIET
			// boost history
			// boost killers
			//TODO  boost check?????
		}

		// boost by pst
		(*moves)[i].score+=Pst[Gamestage(p)][(*moves)[i].piece][(*moves)[i].to]
		// cycle through pv to boost all moves in current move list to top
		for _, m := range pv.moves {
			if (*moves)[i].from == m.from && (*moves)[i].to == m.to { //&& (*moves)[i].extra == m.extra {
				(*moves)[i].score += PVBONUS
			}
		}
	}

	sort.Slice((*moves), func(i, j int) bool { return (*moves)[i].score > (*moves)[j].score }) // descending
	//assert
	l := len((*moves))
	if l > 0 && (*moves)[0].score < (*moves)[l-1].score {
		panic("sort is not descending")
	}
	return true
}

/*
QUIET     *   sorted by history not piecevalue
ENPASSANT *
KILLERS   c
O_O_O     *
O_O       *
EPCAPTURE *
BADCAPTURE c
CAPTURE   * (equal)
GOODCAPTURE c
PROMOTE   *
PVBONUS   c
CHECK     c
*/
