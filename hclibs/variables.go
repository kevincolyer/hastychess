//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import "fmt"
import "github.com/dustin/go-humanize"

//////////////////////////////////////////////////////////////////////////
type Move struct {
	from  int
	to    int
	mtype int // uses constants defined in constants.go
	extra int
}

func (mv Move) String() string {
	return fmt.Sprintf("%v", MoveToAlg(mv))
}

// type PVmoves []Move
//
// func (moves PVmoves) String() (s string) {
//     s = fmt.Sprintf("PV=$v  ", MoveToAlg(moves[0].from,moves[0].to,moves[0].extra))
//     for i:=len(moves); i==1;i-- {
//         s+= fmt.Sprintf("$v  ", MoveToAlg(moves[i].from,moves[i].to,moves[i].extra))
//     }
//     return
// }

//////////////////////////////////////////////////////////////////////////
type Pos struct {
	FEN           string
	Board         [128]int /* 0x88 board */
	TakenPieces   [2]int
	Castled       [4]bool
	King          [2]int
	Side          int
	InCheck       int // -1 == no side in check either 1 or 0
	EnPassant     int // -1 == not in enpassant otherwise square to check
	Fifty         int
	FullMoveClock int
	HalfMoveClock int
	Ply           int
	Hash          Hash
	//History       []Move
}

type History struct {
	move          Move
	TakenPieces   [2]int
	Castled       [4]bool
	King          [2]int
	Side          int
	InCheck       int // -1 == no side in check either 1 or 0
	EnPassant     int // -1 == not in enpassant otherwise square to check
	Fifty         int
	FullMoveClock int
	HalfMoveClock int
	//	Ply           int

}

//////////////////////////////////////////////////////////////////////////
//PV struct
type PV struct { // PV array
	moves [20]Move // too many perhaps depth x 2?
	count int      // for whole list of moves
	// 	depth int // depth searched to

}

/*
// used to provide a sort of PV struct by nodes
type bypv []PV

func (a bypv) Len() int           { return len(a) }
func (a bypv) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bypv) Less(i, j int) bool { return a[i].score > a[j].score } // the < means descending search*/

// pretty printer for PV struct
func (pv PV) String() string {
	return fmt.Sprintf("%v", pv.moves)
}

///////////////////////////////////////////////////////////////////////
type Movescore struct {
	move  Move
	score int
	ttkey string
}
type bymovescore []Movescore

func (a bymovescore) Len() int           { return len(a) }
func (a bymovescore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bymovescore) Less(i, j int) bool { return a[i].score > a[j].score } // > means descending.
///////////////////////////////////////////////////////////////////////
type TtData struct {
	score    int // score found
	ply      int // ply first discovered at (to avoid loops)
	nodetype int // TTEXACT TTUPPER OR TTLOWER
	move     Move
	age      int64
}

var tt map[string]TtData

////////////////////////////////////////////////////////////////////////
// type Book struct {
// 	moves []Move
// }

var book map[string][]Move

// for stat gathering
var StatNodes int
var StatQNodes int
var StatUpperCuts int
var StatLowerCuts int

var StatTimeStart int // not sure what type needed here
var StatTimeElapsed int

var StatTtHits int
var StatTtWrites int
var StatTtUpdates int
var StatTtCulls int
var TtAgeCounter int64

// for game
var GameOver bool
var GameDisplayOn bool
var GameDepthSearch int
var GameForce bool
var GameUseBook bool
var GameUseTt bool
var GameUseStats bool

type Proto int

var GameProtocol Proto

func UCI() bool {
	if GameProtocol == PROTOUCI {
		return true
	}
	return false
}

func Comma(i int) string {
	return humanize.Comma(int64(i))
}
func Commaf(i float64) string {
	return humanize.Comma(int64(i))
}

var Control chan string
