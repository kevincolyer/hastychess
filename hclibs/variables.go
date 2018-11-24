//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import "fmt"
import "github.com/dustin/go-humanize"

// import "time"

//////////////////////////////////////////////////////////////////////////
type Move struct {
	from    int
	to      int
	mtype   int // uses constants defined in constants.go
	piece   int // color masked
	extra   int // promotion
	subtype int // check
	score   int
}

func (mv Move) String() string {
	// 	return fmt.Sprintf("%v", MoveToAlg(mv))
	return fmt.Sprintf("%v", MoveToSAN(mv))
}

type Fen string

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
	JustTaken     int
	Castled       [4]bool
	King          [2]int
	Side          int
	InCheck       int // -1 == no side in check either 1 or 0
	EnPassant     int // -1 == not in enpassant otherwise square to check
	Fifty         int
	FullMoveClock int
	HalfMoveClock int
	Hash          Hash
	//	Ply           int

}

var history [1024]History

//////////////////////////////////////////////////////////////////////////
//PV struct
type PV struct { // PV array
	moves [20]Move // too many perhaps depth x 2?
	count int      // for whole list of moves
	ply   int      // to syncronise the ply level

}

// pretty printer for PV struct
func (pv PV) String() (res string) {
	if pv.count < 1 {
		res = fmt.Sprintf("%v", pv.moves[:pv.count])
	} else {
		res = fmt.Sprintf("%v", pv.moves[:pv.count-1])
	}
	res = res[1 : len(res)-1]
	return
}

///////////////////////////////////////////////////////////////////////
type Movescore struct {
	move  Move
	score int
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
// var StatNodes int
// var StatQNodes int
// var StatUpperCuts int
// var StatLowerCuts int
//
// var StatTimeStart time.Time
// var StatTimeElapsed int
//
// var StatTtHits int
// var StatTtWrites int
// var StatTtUpdates int
// var StatTtCulls int
// var TtAgeCounter int64

// for game
// var GameOver bool
// var GmeDisplayOn bool
// var GameDepthSearch int
// var GameForce bool
// var GameUseBook bool
// var GameUseTt bool
// var GameUseStats bool
// var GamePostStats bool
// var GameStopSearch bool
// var GameDurationToSearch time.Duration

type Proto int

// var GameProtocol Proto

// func UCI() bool {
// 	if GameProtocol == PROTOUCI {
// 		return true
// 	}
// 	return false
// }

func Comma(i int) string {
	return humanize.Comma(int64(i))
}
func Commaf(i float64) string {
	return humanize.Comma(int64(i))
}

// var Control chan string
