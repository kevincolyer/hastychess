// Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package protocol



type CLIOptions struct {
	Engine         string
	Protocol       string
	Stats          bool
	Book           bool
	TT             bool
	ConsoleNoColor bool
	Fen            string
	Depth          int
	Difficulty     int    // several uses but for crazy chess 1-3
	GameType       string // normal eq empty string. or "crazy" etc
	NameVersion    string
}

