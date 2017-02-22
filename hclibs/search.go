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
func SearchRoot(p Pos, initdepth, maxdepth int) (move Move, score int) {

	///////// Checkmate and stalemate detection
	consider := GenerateAllMoves(&p)
	num2consider := len(consider)
	gamestage := Gamestage(&p)
	var mvscore []Movescore
	// 	var score int
	// 	fmt.Printf("num2consider %v\n", num2consider)
	if num2consider == 0 {
		if p.InCheck > -1 {
			score = -CHECKMATE

		} else {
			score = -STALEMATE
		}
		move = Move{}
		return
	}
	StatNodes++
	if num2consider == 1 {
		// if only one move to make, make it!
		move = consider[0]
		score = Eval(&p, 1, gamestage)
		return
	}
	/////////// Get initial sort so we can get an left hand set of nodes
	for i, _ := range consider {
		q := p // copy p

		MakeMove(consider[i], &q)
		mvscore = append(mvscore, Movescore{consider[i], Eval(&q, num2consider, gamestage), ""})
		//pv = append(pv, PV{[]Move{Move{consider[i].from, consider[i].to, consider[i].mtype/*/*, consider[i].extra}}, NegaMaxAB(q, NEGINF, INF, initdepth), initdepth})
	}
	// 	sort.Sort(bymovescore(mvscore))

	//fmt.Printf("Max=%v, Min=%v\n",max,pv[len(pv)-1].score)

	for depth := 2; depth <= maxdepth; depth++ {
		sort.Sort(bymovescore(mvscore))
		max := mvscore[0].score
		score = max
		move = mvscore[0].move
		if !UCI() {
			fmt.Printf("# searching to depth %d\n", depth)
		}
		///////// Deeper sort -- iterative deepening
		for i, _ := range mvscore {
			// if the best is > 1/4 a pawn above next choice then give up search - done
			// and we have not found it in 4 searches...
			if depth > 2 && max > mvscore[i].score+25 && i < 4 {
				if !UCI() {
					fmt.Println("# nothing better, splitting at index ", i)
				}
				break
			}
			if StopSearch() {
				return
			} // someone signals we should stop
			q := p //copy p
			MakeMove(mvscore[i].move, &q)
			temp := NegaMaxAB(q, max-50, max+50, depth, true)
			if temp <= max-50 || temp >= max+50 {

				temp = NegaMaxAB(q, NEGINF, INF, depth, true) // broaden the search window - we found some values way beyond what was expected...
			}
			mvscore[i].score = temp
			//                         mvscore[i].score = NegaMaxAB(q, NEGINF, INF, depth,depth==maxdepth) // only enter q search at deepest leafnodes...
			// 			mvscore[i].score = NegaMaxAB(q, NEGINF, INF, depth, false)
			// fmt.Printf("# %v\n", mvscore[i].score)
			if mvscore[i].score > max {
				max = mvscore[i].score
				score = max
				move = mvscore[i].move
				if GameUseStats && !UCI() {
					fmt.Printf("# good move %v scored %v found at index %v\n", move, score, i)

				}
			} // set new max
		}
	}
	///////////////////// finished search report and cleanup
	if !UCI() {
		fmt.Printf("# Chosen move %v score %v\n", move, Comma(score))
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
	return
}

func NegaMaxAB(p Pos, alpha int, beta int, depth int, enterQuiesce bool) int {

	bestval := NEGINF
	side := p.Side
	val := NEGINF
	bestmove := Move{}
	tttype := TTLOWER // default node to write! lower=alpha / upper = beta

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
					val = elem.score
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
	// if we have not looked at this position before then get a value for it
	if val == NEGINF {
		val = Eval(&p, lenmoves, gamestage) // material score
	}

	// checkmate and stalemate detection...
	if lenmoves == 0 || (p.TakenPieces[side] == 15 && p.TakenPieces[1-side] == 15) { // go till no moves or only kings
		// 		fmt.Printf("depth %v found 0 moves\n", depth)
		//fmt.Printf("%v\n", p)
		if p.InCheck == side {
			val = -CHECKMATE // checkmate to xside
		} else {
			val = -STALEMATE // stalemate
		}
		if GameUseTt {
			elem, ok = tt[ttkey]
			if !ok || elem.ply <= p.Ply {
				StatTtWrites++
				tt[ttkey] = TtData{val, p.Ply, TTEXACT, Move{}, TtAgeCounter}
				TtAgeCounter++
			}
		}
		return val
	}

	// LEAF NODE
	// we are at a leaf node at the end of a search so...
	if depth == 0 {
		// NEED QUIESCENCE SEARCH ONE LEVEL DEEPER HERE...
		// 		fmt.Printf("Entering q search\n")
		if enterQuiesce {
			val = SearchQuiesce(p, alpha, beta, 12)
		}
		tttype := TTEXACT
		if val > beta {
			tttype = TTUPPER
		}
		if val <= alpha {
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
			tt[ttkey] = TtData{val, p.Ply, tttype, Move{}, TtAgeCounter}
			TtAgeCounter++
			// 			}
		}
		return val

	} // at a leaf // note at a leaf we can't detect stalemate - need to look deeper for that...

	// EVALUATION
	/* Dont think we need this - all this does is get a material score and sorts it...????
	var consider []Movescore

	for _, m := range moves {
		q = p
		MakeMove(m, &q)
		ttkey = TtKey(&q)
		css := Eval(&q, lenmoves, gamestage) // ??????????????? should this not be negated??????? No - these are my moves...
		consider = append(consider, Movescore{m, css, ttkey})
	}
	sort.Sort(bymovescore(consider)) // sort descending - bst moves first

	// ALPHA == lower bound
	// BETA == upper bound */
	bestmove = moves[0]
	for _, m := range moves {
		if StopSearch() {
			return alpha
		} // someone signals we should stop
		q = p
		MakeMove(m, &q)
		ttkey = TtKey( &q )

		val = -NegaMaxAB(q, -beta, -alpha, depth-1, enterQuiesce)
		// re-sort the moves here... from remaining moves, if one greater than beta move to next to be considered????????
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
		}
	}
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
	return bestval // found between lower and upper bounds
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

func Gamestage(p *Pos) int {
	if p.TakenPieces[p.Side] > 12 {
		return ENDGAME
	}
	if p.TakenPieces[p.Side] > 4 {
		return MIDGAME
	}
	return OPENING
}
