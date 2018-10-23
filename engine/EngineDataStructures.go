// Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer

package engine

// fen
// move
// result
// stats
// pv
// engineoptions - tt, ttsize,

type Engine struct {
}

type EngineOptions struct {
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

type EngineThinker interface {
	New(EngineOptions) (*Engine, error)
	Search(Fen, PV, EngineOptions) Move
	Ponder(Fen, PV, EngineOptions) PV
	MakeMove(Fen, Move) Result
	Stop() bool
	GetPV() PV
	GetStats() Stats
	ListMoves(Fen) []Move
	IsLegal(Fen, Move) bool
}

func New(e EngineOptions) (*Engine, error) {
	return &Engine{}, nil
}

func (e *Engine) Stop() bool {
	return true
}
