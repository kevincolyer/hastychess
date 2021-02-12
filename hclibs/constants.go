//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer

package hclibs

// Versioning

const VERSION = "1.11 'Blockhead'"

const BLACK = 1
const WHITE = 0
const STARTFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// board and piece numbering suggestions from https://cis.uab.edu/hyatt/boardrep.html

const PREVENTEXPLOSION = 2000000 // 2mil nodes to made debugging easier!!! TODO remove!
const MAXSEARCHDEPTH = 10        // for manual play - easier to know what you have done wrong
const MAXSEARCHDEPTHX = 8        // for xboard
const QUIESCEDEPTH = 4
const TTMAXSIZE = 20 // (is 2^20=1048576 entries when over this oldest tt nodes are culled back so tt is this size
const TTHASH = 1
const QSTTHASH = 2
const PTTHASH = 3
const ETTHASH = 4

// const USEBOOK=false
// const USETTABLE=true

// pawn=1, knight=2, king=3, bishop=5, rook=6 and queen=7 ---? WHY?

//WHITE
const EMPTY = 0
const PAWN = 1
const NIGHT = 3
const ROOK = 4
const BISHOP = 5
const KING = 6
const QUEEN = 7

//black
const pawn = PAWN + 8
const night = NIGHT + 8
const king = KING + 8
const bishop = BISHOP + 8
const rook = ROOK + 8
const queen = QUEEN + 8

// used for mtype in struct move
// const QUIET = 0
// const CAPTURE = 2
// const PROMOTE = 4
// const ENPASSANT = 8
// const O_O_O = 16
// const O_O = 32
// const EPCAPTURE = 64

// Values used for mtype in struct move
// this is useful for ordering moves - ascending order is interesting for us.

const (
	UNINITIALISED = 0
	QUIET         = 50 // sorted by history
	ENPASSANT     = 51
	KILLERS       = 100
	O_O_O         = 200
	O_O           = 201
	BADCAPTURE    = 300
	EPCAPTURE     = 301
	CAPTURE       = 400
	GOODCAPTURE   = 500
	PROMOTE       = 600
	PVBONUS       = 700
	INCHECK       = 800
)

const OPENING = 0
const MIDGAME = 1
const ENDGAME = 2

const QS = 0
const KS = 1

const CHECKMATE = 200000 // 200_000
const STALEMATE = 0  // 50000
const CHECK = 20000      // sign dictates who we award it to
const NEGINF = -1000001  // -1_000_001
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

//  pawn
var PM = [2][4]int{
	[4]int{NW, NORTH, NE, NN},
	[4]int{SE, SOUTH, SW, SS},
}

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

// REVGRID ... this is actually how we display the board on a screen or in FEN - top down !
var REVGRID = [64]int{
	A8, B8, C8, D8, E8, F8, G8, H8,
	A7, B7, C7, D7, E7, F7, G7, H7,
	A6, B6, C6, D6, E6, F6, G6, H6,
	A5, B5, C5, D5, E5, F5, G5, H5,
	A4, B4, C4, D4, E4, F4, G4, H4,
	A3, B3, C3, D3, E3, F3, G3, H3,
	A2, B2, C2, D2, E2, F2, G2, H2,
	A1, B1, C1, D1, E1, F1, G1, H1,
}
