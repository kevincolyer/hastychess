//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

// import "fmt"

const BLACK = 1
const WHITE = 0
const STARTFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// board and piece numbering suggestions from https://cis.uab.edu/hyatt/boardrep.html
const MAXSEARCHDEPTH = 8  // for manual play - easier to know what you have done wrong
const MAXSEARCHDEPTHX = 8 // for xboard
const QUIESCEDEPTH = 10
const TTMAXSIZE = 20 // (is 2^20=1048576 entries when over this oldest tt nodes are culled back so tt is this size
const TTHASH = 1
const QSTTHASH = 2
const PTTHASH = 3
const ETTHASH = 4

// const USEBOOK=false
// const USETTABLE=true

// pawn=1, knight=2, king=3, bishop=5, rook=6 and queen=7
const EMPTY = 0
const PAWN = 1
const NIGHT = 2
const KING = 3
const BISHOP = 5
const ROOK = 6
const QUEEN = 7
const pawn = 1 + 8
const night = 2 + 8
const king = 3 + 8
const bishop = 5 + 8
const rook = 6 + 8
const queen = 7 + 8

// used for mtype in struct move
// const QUIET = 0
// const CAPTURE = 2
// const PROMOTE = 4
// const ENPASSANT = 8
// const O_O_O = 16
// const O_O = 32
// const EPCAPTURE = 64

// used for mtype in struct move
// this is useful for ordering moves - ascending order is interesting for us.
const (
	CAPTURE   = 1 << iota
	EPCAPTURE = 1 << iota
	PROMOTE   = 1 << iota
	O_O_O     = 1 << iota
	O_O       = 1 << iota
	ENPASSANT = 1 << iota
	QUIET     = 1 << iota
)

const OPENING = 0
const MIDGAME = 1
const ENDGAME = 2

const QS = 0
const KS = 1

const CHECKMATE = 200000
const STALEMATE = 50000
const CHECK = 20000
const NEGINF = -1000001
const INF = -NEGINF
const POSINF = INF

const TTEXACT = 1
const TTLOWER = 2
const TTUPPER = 3

// const A1  = 0
// const A8  = 112
// const H1  = 7
// const H8  = 119

const (
	A1, B1, C1, D1, E1, F1, G1, H1 = iota<<4 + 0, iota<<4 + 1, iota<<4 + 2, iota<<4 + 3, iota<<4 + 4, iota<<4 + 5, iota<<4 + 6, iota<<4 + 7
	A2, B2, C2, D2, E2, F2, G2, H2 = iota<<4 + 0, iota<<4 + 1, iota<<4 + 2, iota<<4 + 3, iota<<4 + 4, iota<<4 + 5, iota<<4 + 6, iota<<4 + 7
	A3, B3, C3, D3, E3, F3, G3, H3 = iota<<4 + 0, iota<<4 + 1, iota<<4 + 2, iota<<4 + 3, iota<<4 + 4, iota<<4 + 5, iota<<4 + 6, iota<<4 + 7
	A4, B4, C4, D4, E4, F4, G4, H4 = iota<<4 + 0, iota<<4 + 1, iota<<4 + 2, iota<<4 + 3, iota<<4 + 4, iota<<4 + 5, iota<<4 + 6, iota<<4 + 7
	A5, B5, C5, D5, E5, F5, G5, H5 = iota<<4 + 0, iota<<4 + 1, iota<<4 + 2, iota<<4 + 3, iota<<4 + 4, iota<<4 + 5, iota<<4 + 6, iota<<4 + 7
	A6, B6, C6, D6, E6, F6, G6, H6 = iota<<4 + 0, iota<<4 + 1, iota<<4 + 2, iota<<4 + 3, iota<<4 + 4, iota<<4 + 5, iota<<4 + 6, iota<<4 + 7
	A7, B7, C7, D7, E7, F7, G7, H7 = iota<<4 + 0, iota<<4 + 1, iota<<4 + 2, iota<<4 + 3, iota<<4 + 4, iota<<4 + 5, iota<<4 + 6, iota<<4 + 7
	A8, B8, C8, D8, E8, F8, G8, H8 = iota<<4 + 0, iota<<4 + 1, iota<<4 + 2, iota<<4 + 3, iota<<4 + 4, iota<<4 + 5, iota<<4 + 6, iota<<4 + 7
)

const (
	PROTOCONSOLE = iota + 1
	PROTOXBOARD
	PROTOUCI
)

const NORTH = 16
const NN = 32
const SOUTH = -16
const SS = -32
const EAST = 1
const WEST = -1
const NE = NORTH + EAST
const NW = NORTH + WEST
const SE = SOUTH + EAST
const SW = SOUTH + WEST

var RM = [4]int{NORTH, EAST, WEST, SOUTH}
var BM = [4]int{NE, NW, SE, SW}
var QM = [8]int{NORTH, EAST, WEST, SOUTH, NE, NW, SE, SW}

//  leapers
var KM = [8]int{NORTH, EAST, WEST, SOUTH, NE, NW, SE, SW}
var NM = [8]int{31, 33, -31, -33, 14, 18, -14, -18}

// # pwn
var PM = [2][4]int{
	[4]int{NW, NORTH, NE, NN},
	[4]int{SE, SOUTH, SW, SS},
}

// func MakeGrid() [64]int {
// 	var GRID [64]int
// 	var i int
// 	j := 0
// 	for ; i < 128; i++ {
// 		if i&0x88 == 0 {
// 			GRID[j] = i
// 			j++
// 		}
// 	}
// 	return GRID
// }

var GRID = [64]int{
	A1, B1, C1, D1, E1, F1, G1, H1,
	A2, B2, C2, D2, E2, F2, G2, H2,
	A3, B3, C3, D3, E3, F3, G3, H3,
	A4, B4, C4, D4, E4, F4, G4, H4,
	A5, B5, C5, D5, E5, F5, G5, H5,
	A6, B6, C6, D6, E6, F6, G6, H6,
	A7, B7, C7, D7, E7, F7, G7, H7,
	A8, B8, C8, D8, E8, F8, G8, H8,
}

var REVGRID = [64]int{ // this is actually how we display the board on a screen or in FEN - top down !
	A8, B8, C8, D8, E8, F8, G8, H8,
	A7, B7, C7, D7, E7, F7, G7, H7,
	A6, B6, C6, D6, E6, F6, G6, H6,
	A5, B5, C5, D5, E5, F5, G5, H5,
	A4, B4, C4, D4, E4, F4, G4, H4,
	A3, B3, C3, D3, E3, F3, G3, H3,
	A2, B2, C2, D2, E2, F2, G2, H2,
	A1, B1, C1, D1, E1, F1, G1, H1,
}
