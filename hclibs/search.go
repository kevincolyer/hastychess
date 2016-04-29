package hclibs

import (
	"fmt"
	"sort"
)

func gamestage(p *Pos) int {
	if p.TakenPieces[p.Side] > 12 {
		return ENDGAME
	}
	if p.TakenPieces[p.Side] > 4 {
		return MIDGAME
	}
	return OPENING
}

func NegaMaxAB(p Pos, alpha int, beta int, depth int, pvmoves *[]Move) int {

	bestval := NEGINF
	side := p.Side
	var q Pos
	val := NEGINF
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
			if elem.nodetype==TTEXACT {
                            val = elem.score
                            if depth==0 {return val} // ???????????????????????????????????????????????????????????
                            // if at end of search and seen exact value before then return exact value...
                        }
			if elem.nodetype==TTLOWER {alpha = elem.score} // use previous deeper bounds to set this bound
			if elem.nodetype==TTUPPER {beta = elem.score} // use previous deeper bounds to set this bound
		}
	}

	// if we have not looked at this position before then get a value for it
	if val == NEGINF {

		val = ClaudeShannonScore(&p, lenmoves) // material score
	}
	
	// we are at a leaf node at the end of a search so...
	if depth == 0 {
                // NEED QUIESENCE SERACH ONE LEVEL DEEPER HERE...
                // put exact score in TT here!!!!! if we are not using it already
            if GameUseTt {// NEED TO MODIFY THIS SOME MORE
			
                        elem, ok = tt[ttkey]
			if ok {
				if elem.ply < q.Ply {
				} // if this search is deeper than prev then update
				StatTtUpdates++
				tt[ttkey] = TtData{val, q.Ply,TTEXACT,Move{}}
			} else {
				StatTtWrites++
				tt[ttkey] = TtData{val, q.Ply,TTEXACT,Move{}}
			}
		}
                return val
		
	} // at a leaf // note at a leaf we can't detect stalemate - need to look deeper for that...

        // checkmate and stalemate detection...
	if lenmoves == 0 || (p.TakenPieces[side] == 15 && p.TakenPieces[1-side] == 15) { // go till no moves or only kings
		// 		fmt.Printf("depth %v found 0 moves\n", depth)
		//fmt.Printf("%v\n", p)
		if p.InCheck == side {
			val = -CHECKMATE // checkmate to xside
		} else {
			val = -STALEMATE // stalemate
		}
		return val
	}

	// node stuff
	var consider []Movescore
	// search all moves one deeper with cs scores plus??? here

	for _, m := range moves {
		q = p
		MakeMove(m, &q)
		ttkey = TtKey(&q)
		if elem, ok := tt[ttkey]; ok {
			consider = append(consider, Movescore{m, elem.score, ttkey})

		} else {

			css := ClaudeShannonScore(&q, lenmoves) // ??????????????? should this not be negated???????
			consider = append(consider, Movescore{m, css, ttkey})
		}
	}
	sort.Sort(bymovescore(consider)) // sort descending ??????????????????
	// ALPHA == lower bound
	// BETA == upper bound

	for _, m := range consider {
		q = p
		MakeMove(m.move, &q)
		ttkey = m.ttkey

		val = -NegaMaxAB(q, -beta, -alpha, depth-1, pvmoves)
                // RE-SORT THE MOVES HERE... FROM REMAINING MOVES, IF ONE GREATER THAN BETA MOVE TO NEXT TO BE CONSIDERED
		if GameUseTt {// NEED TO MODIFY THIS SOME MORE
			var tttype int
			if val>beta {tttype=TTUPPER}
// 			if val<beta  && val>bestval && val<alpha {tttype=TTLOWER}
			if val<beta  && val>bestval && val>alpha {tttype=TTLOWER}
                        elem, ok = tt[ttkey]
			if ok {
				if elem.ply < q.Ply {
				} // if this search is deeper than prev then update
				StatTtUpdates++
				tt[ttkey] = TtData{val, q.Ply,tttype,m.move}
			} else {
				StatTtWrites++
				tt[ttkey] = TtData{val, q.Ply,tttype,m.move}
			}
		}
		StatNodes++

		// found a better upper bound
		if val >= beta {
			//             @pv[depth]=[ @m ]
			(*pvmoves)[depth] = m.move
			return val // best val above expexted upper bound
		}

		// found better value below upper bound
		if val > bestval {
			bestval = val
			(*pvmoves)[depth] = m.move
			if val > alpha { // is better than lower bound so move that up and make a new lower bound
				alpha = val
			}
		}
	}
	return bestval // found between lower and upper bounds
}

// search = generate all moves into PV struct
// do a depth evalute to 4 to score
// sort
// if no moves loose or draw
// if only one move - make it

// NEED QuiesenceSearch, SearchRoot, Nulls and PV's 
// SortMoves needs to be smarter by adding in captures and promtions etc.
// consider killer and history


// THIS SHOULD BE SEARCH ROOT OR A WAY TO ALLOW ME TO PLUG IN DIFFERENT SEARCHES
func Search(p Pos, initdepth, maxdepth int) PV {

	StatNodes = 0  // for internal statistics purpose...
	StatTtHits = 0 //
	StatTtUpdates = 0
	StatTtWrites = 0

	var pv []PV
	///////// Checkmate and stalemate detection
	consider := GenerateAllMoves(&p)
	num2consider := len(consider)
	// 	var score int
	// 	fmt.Printf("num2consider %v\n", num2consider)
	if num2consider == 0 {
		if p.InCheck > -1 {
			pv = []PV{PV{[]Move{Move{0, 0, 0, 0}}, -CHECKMATE, 0}}
		} else {
			pv = []PV{PV{[]Move{{0, 0, 0, 0}}, -STALEMATE, 0}}
		}
		return pv[0]
	}
	if num2consider == 1 {
                // if only one move to make, make it!
                StatNodes++
		pv = []PV{PV{[]Move{Move{consider[0].from, consider[0].to, consider[0].mtype, consider[0].extra}}, ClaudeShannonScore(&p, 1), 0}}
		return pv[0]
	}
	// for each possible move add a pv node, and look to depth initdepth to get an initial score to prune the tree with.

	/////////// Get initial sort so we can get an left hand set of nodes
	for i, _ := range consider {
		//fmt.Println(p)
		q := p // copy p
		//fmt.Println(q)
		//fmt.Println(p)

		MakeMove(consider[i], &q)
		pv = append(pv, PV{make([]Move, maxdepth+1), 0, 0})
		//pv = append(pv, PV{[]Move{Move{consider[i].from, consider[i].to, consider[i].mtype, consider[i].extra}}, NegaMaxAB(q, NEGINF, INF, initdepth), initdepth})
		pv[i].moves[0].from = consider[i].from
		pv[i].moves[0].to = consider[i].to
		pv[i].moves[0].mtype = consider[i].mtype
		pv[i].moves[0].extra = consider[i].extra
		pv[i].score = -NegaMaxAB(q, NEGINF, INF, initdepth, &pv[i].moves)
		pv[i].depth = initdepth
	}
	sort.Sort(bypv(pv))
	max := pv[0].score
	//fmt.Printf("Max=%v, Min=%v\n",max,pv[len(pv)-1].score)
	maxi := 0
	///////// Deeper sort (iterative deepning here?)
	for i, _ := range pv {
		if max > pv[i].score+25 {
			break
		} // if the best is > 1/4 a pawn the next choice then give up search - done

		// cut here somehow

		q := p //copy p
		MakeMove(pv[i].moves[0], &q)
		fmt.Printf("# Considering %v (initial score %v)... to depth %d... ", pv[i].moves[0], pv[i].score, maxdepth)
		pv[i].score = -NegaMaxAB(q, NEGINF, INF, maxdepth, &pv[i].moves)
		fmt.Printf("# %v\n", pv[i].score)
		if pv[i].score > max {
			max = pv[i].score
			maxi = i
			if GameUseStats {
				fmt.Printf("# Best line %v scored %v\n", pv[i].moves, pv[i].score)
			}
		} // set new max
		pv[i].depth = maxdepth
	}
	///////////////////// finished search report and cleanup
	fmt.Printf("# Chosen move %v score %v\n", pv[maxi].moves, pv[maxi].score)
	if GameUseStats {
		fmt.Printf("# StatNodes searched: %v\n", StatNodes)
		if GameUseTt {
			fmt.Printf("# TT table has %v entries\n", len(tt))
		}
	}
	// prune dead tt entries (from ply's in the past)
	if GameUseTt {
		if GameUseStats {
			fmt.Println("# Culling TT entries")
		}
		culled := 0
		for key, ttdata := range tt {
			if ttdata.ply < p.Ply {
				delete(tt, key)
				culled++

			}
		}
		if GameUseStats {
			fmt.Printf("# culled %v TT entries, %v remain\n", culled, len(tt))
		}
	}
	return pv[maxi]
}

// Iterative deepeing
// prune PV of all useless nodes
// do depth search to max depth or timeout
// resort pv tree each turn
