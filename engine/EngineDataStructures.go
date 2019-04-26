// Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer

package engine

import "github.com/kevincolyer/hastychess/hclibs"

// fen
// move
// result
// stats
// pv
// options - tt, ttsize,

// Engine ... type for game engines
type Engine struct {
}

type Options struct {
}

type PV struct {
}

type Move struct {
}

type Stats struct {
}

type Result struct {
}

type Fen string

type Thinker interface {
	New(Options) (*Engine, error)
	Search(Fen, PV, Options) Move
	Ponder(Fen, PV, Options) PV
	MakeMove(Fen, Move) Result
	Stop() bool
	GetPV() PV
	GetStats() Stats
	ListMoves(Fen) []Move
	IsLegal(Fen, Move) bool
}

func New(e Options) (*Engine, error) {
	// load book here or it will not get loaded and then will not get used!
	hclibs.GameInit()
	return &Engine{}, nil
}

func (e *Engine) Stop() bool {
	return true
}
