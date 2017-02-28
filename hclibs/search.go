//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import (
	//	"fmt"
	"sort"
)

// THIS SHOULD BE SEARCH ROOT OR A WAY TO ALLOW ME TO PLUG IN DIFFERENT SEARCHES
func SearchRoot(p Pos, initdepth, maxdepth int) (bestmove Move, bestscore int) {
	return
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
