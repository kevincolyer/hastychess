//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import (
	"fmt"
	"sort"
)

// THIS SHOULD BE SEARCH ROOT OR A WAY TO ALLOW ME TO PLUG IN DIFFERENT SEARCHES
func SearchRoot(p *Pos, maxdepth int) (bestmove Move, bestscore int) {
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
	bestscore = NEGINF
	bestmove = Move{}
	//     var pv []PV

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
	fmt.Println(consider)
	OrderMoves(consider, p)
	fmt.Println(consider)

	// 3a. if iterative deepening loop from depth 2 to max depth in turn, sorting best score descending
	depth := maxdepth
	//         depth =1
	for _, move := range consider {
		MakeMove(move, p)
		val := -negamax(depth, p) // need neg here as we switch sides in make move and evaluation happens relative to side
		fmt.Printf("move %v scored %v\n", move, val)
		UnMakeMove(move, p)

		if val > bestscore {
			bestmove = move
			bestscore = val
			fmt.Printf("found bestscore %v move %v\n", bestscore, bestmove)
		}
	}
	return
}

// classical negamax search: negated minimax. This does no pruning!!!!
// It will search the entire search space until if finds the best move
func negamax(depth int, p *Pos) int {

	StatNodes++
	if depth == 0 || StatNodes > PREVENTEXPLOSION {
		return Eval(p, 1, Gamestage(p))
	}
	max := NEGINF

	consider := GenerateAllMoves(p)
	if len(consider) == 0 {
		if p.InCheck > 0 {
			return CHECKMATE
		}
		return STALEMATE
	}
	OrderMoves(consider, p)

	for _, move := range consider {
		MakeMove(move, p)
		score := -negamax(depth-1, p)
		UnMakeMove(move, p)
		if score > max {
			max = score
		}
	}
	return max
}

//func SearchQuiesce(p Pos, alpha, beta int, qdepth int) int {
// 	// need a standpat score
// 	var val int
// 	var q Pos
// 	gamestage := Gamestage(&p)
// 	val = EvalQ(&p, 1, gamestage) // custom evaluator here for QUIESENCE
// 	var mvscore []Movescore
// 	standpat := val
// 	StatQNodes++
//
// 	if val >= beta {
// 		return beta
// 	}
// 	if alpha <= val {
// 		alpha = val
// 	}
//
// 	// prevent search explosion while testing TODO remove!
//
// 	if StatQNodes > 1000000 {
// 		fmt.Println("Qnode explosion - bottling!")
// 		return alpha
//
// 	}
// 	if qdepth == 0 {
// 		return alpha
// 	} // cant search too deep
// 	// get moves - but only captures and promotions
// 	moves := GenerateMovesForQSearch(&p)
// 	if len(moves) == 0 {
// 		return alpha
// 	} // nothing more to search...
// 	// score them
// 	// 	fmt.Printf("# QS: looking at ply %v (%v moves to consider: %v)\n", p.Ply, len(moves), moves)
// 	// socre them by Most Valuable Victim - Least Valuable Aggressor
// 	for _, i := range moves {
// 		mvscore = append(mvscore, Movescore{move: i,
// 			score: MVVLVA(i, &p),
// 			ttkey: ""})
// 	}
// 	// And order descending to provoke cuts
// 	sort.Sort(bymovescore(mvscore))
// 	//         fmt.Printf("alpha= %v, beta=%v, val=%v, movescore=%v\n",alpha,beta,val,mvscore)
// 	// loop over all moves, searching deeper until no moves left and all is "quiet" - return this score...)
// 	for _, m := range mvscore {
// 		// adjust each score for delta cut offs and badmoves skipping to next each time
// 		// delta - if not promotion and not endgame and is a low scoring capture then continue
// 		if m.move.mtype != PROMOTE && Gamestage(&p) != ENDGAME && standpat+csshash[p.Board[m.move.to]]+200 < alpha {
// 			continue
// 		} // delta cut qnodes from 20M to 640,000 in one case!
//
// 		// 		// badmoves - cut qnodes from 640,000 to 64,000
// 		if p.Board[m.move.from]&7 == PAWN && m.move.mtype != PROMOTE {
// 			continue
// 		}
// 		// capture by pawn is ok
//
// 		// search deeper until quiet
// 		// 		fmt.Println("search one deeper")
// 		if StopSearch() {
// 			return alpha
// 		} // someone signals we should stop
// 		q = p
// 		MakeMove(m.move, &q)
// 		val = -SearchQuiesce(q, -beta, -alpha, qdepth-1)
// 		// 		fmt.Printf("returned from one deeper val = %v\n",val)
//
// 		// adjust window
// 		if val >= alpha {
// 			if val > beta {
// 				return beta
// 			}
// 			alpha = val
// 		}
// 	}
// 	return alpha // nothing better than this to return
// }

func OrderMoves(moves []Move, p *Pos) bool {
	// can add in pv at top and also any other things to help
	// order by move type (capture and promotion first down to quiet moves)
	// todo: sub sort by captured piece value // pv // any other factor!
	//sort.Slice(moves, func(i, j int) bool { return p.Board[moves[i].from] > p.Board[moves[j].from] }) // by piece type descending
	sort.Slice(moves, func(i, j int) bool { return moves[i].mtype < moves[j].mtype }) // by move type ascending
	return true
}
