package hclibs

import (
	"fmt"
	"github.com/dex4er/go-tap"
	"testing"
)

func TestSearch(t *testing.T) {
	p := FENToNewBoard(STARTFEN)
	tap.Is(len(GenerateAllMoves(&p)), 20, "20 moves counted on a new board")

	//tap.Ok(len(pv.moves) > 0, "test of searching") // initial depth 4 and 8 max depth
	// 	for j := 0; j < 2; j++ {
	// 		pv := Search(p, 2, 4)
	// 		MakeMove(pv.moves[0], &p)
	// 		fmt.Println(&p)
	// 	}
	// 	fmt.Printf("%v", pv.moves)

	p = FENToNewBoard("8/8/8/1k6/8/7Q/3R4/2K5 w - - 0 1") // for testing endings
	tap.Diag("Testing for check")
	for j := 0; j < 10; j++ {
		move, score := SearchRoot(p, 2, 4)
		MakeMove(move, &p)
		fmt.Println(&p)
		moves := GenerateAllMoves(&p)
		if score == CHECKMATE && len(moves) == 0 {
			tap.Ok(score == CHECKMATE, "Found checkmate")
			break
		}
		if score == STALEMATE && len(moves) == 0 {
			tap.Ok(score == STALEMATE, "Found stalement - this is a problem!")
			break
		}
		if len(moves) == 0 {
			tap.Isnt(len(moves), 0, "Should not have no moves and not be in checkmate or stalemate!")
			break
		}
	}
	tap.Diag("Testing for Stalemate")
	// 	p = FENToNewBoard("8/8/8/8/8/Ppk1b3/1P6/P1K5 w - - 0 1") // for testing endings
	// 	p = FENToNewBoard("5bnr/4p1pq/4Qpkr/7p/2P4P/8/PP1PPPP1/RNB1KBNR b KQ - 0 10") // for testing endings
	p = FENToNewBoard("5bnr/4p1pq/2Q1ppkr/7p/2P4P/8/PP1PPPP1/RNB1KBNR w KQ - 0 10") // for testing endings
	fmt.Println(&p)
	for j := 0; j < 10; j++ {
		move, score := SearchRoot(p, 2, 4)
		MakeMove(move, &p)
		fmt.Println(&p)
		moves := GenerateAllMoves(&p)
		if score == CHECKMATE && len(moves) == 0 {
			tap.Ok(score == CHECKMATE, "Found checkmate - this is a problem")
			break
		}
		if score == STALEMATE && len(moves) == 0 {
			tap.Ok(score == STALEMATE, "Found stalement - this is what we want")
			break
		}
		if len(moves) == 0 {
			tap.Isnt(len(moves), 0, "Should not have no moves and not be in checkmate or stalemate!")
			break
		}
	}
}
