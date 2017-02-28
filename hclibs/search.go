//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import (
	"fmt"
	"sort"
)

// NEED QuiesenceSearch, SearchRoot, Nulls and PV's
// SortMoves needs to be smarter by adding in captures and promtions etc.
// consider killer and history

// THIS SHOULD BE SEARCH ROOT OR A WAY TO ALLOW ME TO PLUG IN DIFFERENT SEARCHES
func SearchRoot(p Pos, initdepth, maxdepth int) (bestmove Move, bestscore int) {
	var pv PV         // to collect our PV in
	GameUseTt = false // TODO fix the broken tt!
	enterQuiesce := false
	///////// Checkmate and stalemate detection

	consider := GenerateAllMoves(&p)
	considercount := len(consider)
	gamestage := Gamestage(&p)
	var mvscore []Movescore
	// 	var score int
	// 	fmt.Printf("considercount %v\n", considercount)
	if considercount == 0 {
		if p.InCheck > -1 {
			bestscore = -CHECKMATE

		} else {
			bestscore = -STALEMATE
		}
		bestmove = Move{} // nothing to od
		return
	}
	StatNodes++
	if considercount == 1 {
		// if only one move to make, make it!
		bestmove = consider[0]
		bestscore = Eval(&p, 1, gamestage)
		return
	}
	/////////// Get initial sort so we can get an left hand set of nodes
	for i, _ := range consider {
		mvscore = append(mvscore, Movescore{consider[i], NEGINF, ""})
	}
	// need PV in here to start first moves off
	OrderMoves(mvscore, &pv)
	//         sort.Slice(mvscore, func( i,j int) bool { return mvscore[i].move.mtype < mvscore[j].move.mtype} )
	bestscore = mvscore[0].score
	bestscore = NEGINF + 100
	bestmove = mvscore[0].move
	pv.moves[0] = bestmove

	//         fmt.Println(mvscore)

	for depth := 1; depth <= maxdepth; depth++ {
		var childpv PV
		if !UCI() {
			fmt.Printf("# searching to depth %d\n", depth)
		}
		//		if depth > 1 { enterQuiesce=true } // only enterQuiesce after first search to give an exact score for better ordering for later
		if depth == maxdepth {
			enterQuiesce = true
		} else {
			enterQuiesce = false
		} // only enterQuiesce after first search to give an exact score for better ordering for later
		///////// Deeper sort -- iterative deepening
		for i, _ := range mvscore {
			if StopSearch() {
				return
			}

			// prevent search explosion...
			// if the searching deeper and >4 searches and so far the best is < 1/4 a pawn above next choice then give up broader search

			if depth > 2 && i > 4 && bestscore > mvscore[i].score+26 {
				if !UCI() {
					fmt.Println("# nothing better, splitting at index ", i)
				}
				break
			}

			q := p //copy p
			MakeMove(mvscore[i].move, &q)

			score := NegaMaxAB(q, bestscore-50, bestscore+50, depth, enterQuiesce, &childpv)
			// TODO too much beardth at lower levels give qsearch explosion with poor move ordering...
			// need to improve q search move ordering...
			//   			if score <= bestscore-50 || score >= bestscore+50 {
			//    				score = NegaMaxAB(q, NEGINF, INF, depth, enterQuiesce) // broaden the search window - we found some values way beyond what was expected...
			//    			}
			mvscore[i].score = score

			// set new max
			if mvscore[i].score > bestscore {
				bestscore = mvscore[i].score
				bestmove = mvscore[i].move

				if GameUseStats && !UCI() {
					fmt.Printf("# good move %v scored %v found at index %v\n# PV=%v count %v\n", bestmove, bestscore, i, pv, pv.count)
				}
			}
			pv.moves[0] = bestmove
			copy(pv.moves[1:], childpv.moves[:])
			pv.count = childpv.count + 1
		}
		// now we have scores sort for the next round
		sort.Sort(bymovescore(mvscore))
		//         fmt.Println(mvscore)
	}

	///////////////////// finished search. Report and cleanup
	if !UCI() {
		fmt.Printf("# Chosen move %v score %v\n", bestmove, Comma(bestscore))
	}
	// prune dead tt entries (from ply's in the past)
	if GameUseTt {
		StatTtCulls = 0
		for key, ttdata := range tt {
			if ttdata.age+TTMAXSIZE < TtAgeCounter {
				delete(tt, key)
				StatTtCulls++

			}
		}
	}
	fmt.Printf("(%v) = %v (count %v)\n", bestscore, pv, pv.count)
	return
}

func NegaMaxAB(p Pos, alpha int, beta int, depth int, enterQuiesce bool, parentpv *PV) int {
	var childpv PV
	// i_alpha:= alpha // initial alpha
	bestval := NEGINF // the value we hope to find and return
	side := p.Side
	val := NEGINF      // a temporary value to compare with best val
	bestmove := Move{} // only for storing pv
	tttype := TTLOWER  // default node to write! lower=alpha / upper = beta

	var q Pos
	var elem TtData
	var ok bool
	var ttkey string

	moves := GenerateAllMoves(&p)
	lenmoves := len(moves)

	if GameUseTt {
		ttkey = TtKey(&p)
		elem, ok = tt[ttkey]
		if ok {
			StatTtHits++
			// if we have already searched deeper then use this value...
			if elem.ply >= p.Ply {

				if elem.nodetype == TTEXACT {
					bestval = elem.score
					if depth == 0 {
						return val
					} // ???????????????????????????????????????????????????????????
					// if at end of search and seen exact value before then return exact value...
				}
				if elem.nodetype == TTLOWER {
					alpha = elem.score
				} // use previous deeper bounds to set this bound
				if elem.nodetype == TTUPPER {
					beta = elem.score
				} // use previous deeper bounds to set this bound
			}
		}
	}
	gamestage := Gamestage(&p)

	// checkmate and stalemate detection...
	if lenmoves == 0 || (p.TakenPieces[side] == 15 && p.TakenPieces[1-side] == 15) { // go till no moves or only kings
		// 		fmt.Printf("depth %v found 0 moves\n", depth)
		//fmt.Printf("%v\n", p)
		if p.InCheck == side {
			bestval = -CHECKMATE // checkmate to xside
		} else {
			bestval = -STALEMATE // stalemate
		}
		if GameUseTt {
			elem, ok = tt[ttkey]
			if !ok || elem.ply <= p.Ply {
				StatTtWrites++
				tt[ttkey] = TtData{bestval, p.Ply, TTEXACT, Move{}, TtAgeCounter}
				TtAgeCounter++
			}
		}
		return bestval
	}
	if bestval == NEGINF {
		bestval = Eval(&p, lenmoves, gamestage) // material score
	}

	// LEAF NODE
	// we are at a leaf node at the end of a search so...
	if depth == 0 {
		// NEED QUIESCENCE SEARCH ONE LEVEL DEEPER HERE...
		// 		 		fmt.Printf("Entering q search\n")
		// if we have not looked at this position before then get a value for it
		if enterQuiesce {
			bestval = SearchQuiesce(p, alpha, beta, QUIESCEDEPTH) // why not neg???????
		}
		tttype := TTEXACT
		if bestval > beta {
			tttype = TTUPPER
		}
		if bestval <= alpha {
			tttype = TTLOWER
		}
		// 		fmt.Printf("leaving q search\n")
		//
		// put exact score in TT here!!!!! if we are not using it already
		if GameUseTt {
			// 			elem, ok = tt[ttkey]
			// 			if !ok || elem.ply <= p.Ply {
			// use always replace for q search (my choice)
			StatTtWrites++
			tt[ttkey] = TtData{bestval, p.Ply, tttype, Move{}, TtAgeCounter}
			TtAgeCounter++
			// 			}
		}
		parentpv.count = 0 // needed to prevent any extra copying onto stack
		return bestval

	} // at a leaf // note at a leaf we can't detect stalemate - need to look deeper for that...

	// Order moves (use function here...)
	sort.Slice(moves, func(i, j int) bool { return moves[i].mtype < moves[j].mtype })
	// ALPHA == lower bound
	// BETA == upper bound */
	bestmove = moves[0]
	parentpv.moves[0] = bestmove
	parentpv.count = 1

	if bestval > alpha {
		alpha = bestval
	} // adjust alpha for this level????

	for _, m := range moves {
		if StopSearch() {
			return alpha
		} // someone signals we should stop
		q = p
		MakeMove(m, &q)
		ttkey = TtKey(&q)

		val = -NegaMaxAB(q, -beta, -alpha, depth-1, enterQuiesce, &childpv)
		StatNodes++

		// found a better UPPER BOUND
		if val >= beta {
			// save in tt
			tttype = TTUPPER
			StatUpperCuts++
			if GameUseTt {
				elem, ok = tt[ttkey]
				if ok {
					if elem.ply <= q.Ply {
						StatTtUpdates++
					} else {
						StatTtWrites++
					}
				} // if this search is deeper than prev then update
				tt[ttkey] = TtData{val, q.Ply, tttype, m, TtAgeCounter}
				TtAgeCounter++
			}
			// end save in tt
			return val // best val above expexted upper bound
		}
		// found better value so reset LOWER BOUND
		if val > bestval {
			bestval = val
			bestmove = m
			if val > alpha { // is better than lower bound so move that up and make a new lower bound
				alpha = val
				tttype = TTEXACT
				StatLowerCuts++
			}
			// push move onto stack of parent pv and add child after. increase the total length counter
			parentpv.moves[0] = m

			copy(parentpv.moves[1:], childpv.moves[:])
			//  			copy(pv.moves[1:],childpv.moves)

			parentpv.count = childpv.count + 1
			//                         fmt.Println(pv)
		}
	} // END OF range over moves
	// save in tt
	if GameUseTt {
		elem, ok = tt[ttkey]
		if !ok || elem.ply < q.Ply {
			// if this search is deeper than prev then update
			//  ore not seen before so add
			tttype = TTLOWER
			StatTtWrites++
			tt[ttkey] = TtData{alpha, q.Ply, tttype, bestmove, TtAgeCounter}
			TtAgeCounter++
		}
	}
	return bestval // found between lower and upper bounds (which is alpha really)
}

func SearchQuiesce(p Pos, alpha, beta int, qdepth int) int {
	// need a standpat score
	var val int
	var q Pos
	gamestage := Gamestage(&p)
	val = EvalQ(&p, 1, gamestage) // custom evaluator here for QUIESENCE
	var mvscore []Movescore
	standpat := val
	StatQNodes++

	if val >= beta {
		return beta
	}
	if alpha <= val {
		alpha = val
	}

	// prevent search explosion while testing TODO remove!

	if StatQNodes > 1000000 {
		fmt.Println("Qnode explosion - bottling!")
		return alpha

	}
	if qdepth == 0 {
		return alpha
	} // cant search too deep
	// get moves - but only captures and promotions
	moves := GenerateMovesForQSearch(&p)
	if len(moves) == 0 {
		return alpha
	} // nothing more to search...
	// score them
	// 	fmt.Printf("# QS: looking at ply %v (%v moves to consider: %v)\n", p.Ply, len(moves), moves)
	// socre them by Most Valuable Victim - Least Valuable Aggressor
	for _, i := range moves {
		mvscore = append(mvscore, Movescore{move: i,
			score: MVVLVA(i, &p),
			ttkey: ""})
	}
	// And order descending to provoke cuts
	sort.Sort(bymovescore(mvscore))
	//         fmt.Printf("alpha= %v, beta=%v, val=%v, movescore=%v\n",alpha,beta,val,mvscore)
	// loop over all moves, searching deeper until no moves left and all is "quiet" - return this score...)
	for _, m := range mvscore {
		// adjust each score for delta cut offs and badmoves skipping to next each time
		// delta - if not promotion and not endgame and is a low scoring capture then continue
		if m.move.mtype != PROMOTE && Gamestage(&p) != ENDGAME && standpat+csshash[p.Board[m.move.to]]+200 < alpha {
			continue
		} // delta cut qnodes from 20M to 640,000 in one case!

		// 		// badmoves - cut qnodes from 640,000 to 64,000
		if p.Board[m.move.from]&7 == PAWN && m.move.mtype != PROMOTE {
			continue
		}
		// capture by pawn is ok

		// search deeper until quiet
		// 		fmt.Println("search one deeper")
		if StopSearch() {
			return alpha
		} // someone signals we should stop
		q = p
		MakeMove(m.move, &q)
		val = -SearchQuiesce(q, -beta, -alpha, qdepth-1)
		// 		fmt.Printf("returned from one deeper val = %v\n",val)

		// adjust window
		if val >= alpha {
			if val > beta {
				return beta
			}
			alpha = val
		}
	}
	return alpha // nothing better than this to return
}

func OrderMoves(mvscore []Movescore, pv *PV) bool {
	// can add in pv at top and also any other things to help
	sort.Slice(mvscore, func(i, j int) bool { return mvscore[i].move.mtype < mvscore[j].move.mtype })
	return true
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
