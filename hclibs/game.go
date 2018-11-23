//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

// useful routines for playing games against people or computer
import "fmt"
import "strings"
import "regexp"
import "time"

// import "github.com/jinzhu/copier"

func GameInit() {
	//tt = make(map[string]TtData)
	if err := InitHashSize(8); err != nil {
		panic(err)
	}
	book = make(map[string][]Move)

	if GameUseBook {
		InitBook()
	}
	return
}

//var pv PV

type Statistics struct {
	Score       int
	Nodes       int
	QNodes      int
	TtHits      int
	TtWrites    int
	TtCulls     int
	TtUpdates   int
	UpperCuts   int
	LowerCuts   int
	TimeElapsed time.Duration
}

type Search struct {
	Stats *Statistics

	TimeStart           time.Time
	MaxDurationOfSearch time.Duration

	FEN    Fen
	P      *Pos
	NewFEN Fen

	Result           string
	Info             string
	PV               *PV
	ParentPV         *PV
	ChildPV          *PV
	Score            int
	BestMove         Move
	Stop             bool
	ExplosionLimit   int
	MaxDepthToSearch int

	EngineInfoChan chan EngineInfo
	UseTT          bool
	UseBook        bool
}

type EngineInfo struct {
	//Pv    *PV
	//Stats *Statistics
	Pv    PV
	Stats Statistics
	Info  string
}

func NewEngineInfo(srch Search) EngineInfo {
	return EngineInfo{
		Stats: *(srch.Stats),
		Pv:    *(srch.PV),
		Info:  srch.Info,
	}
}

func NewSearch(FEN Fen) *Search {
	if FEN == "" {
		FEN = Fen(STARTFEN)
	}
	srch := Search{
		Score:               NEGINF,
		ExplosionLimit:      2000000,
		MaxDurationOfSearch: time.Second * 30,
		MaxDepthToSearch:    8, // just a default
		FEN:                 FEN,
		UseTT:               true,
		UseBook:             true,
	}
	p := FEN.NewBoard()
	srch.P = &p
	s := Statistics{}
	srch.Stats = &s
	pv := PV{}
	srch.PV = &pv

	return &srch
}

func (srch Search) Search(depth int) (completed bool) {
	srch.MaxDepthToSearch = depth
	//srch.StartSearch()
	completed = !srch.Stop // flag is raised if we must stop or hit explosion limit
	if srch.BestMove.mtype == UNINITIALISED {
		panic("StartSearch has returned a nil best move")
	}
	return
}

func (stat Statistics) String() string {
	return fmt.Sprintf(
		"\nscore %v | nodes %v | qnodes %v (%v%%)| nps %v | uppercuts %v | lowercuts %v\ntt_hits %v (%v%%) | tt writes %v | tt updates %v | tt size %v | tt culls %v",
		Comma(stat.Score),
		Comma(stat.Nodes),
		Comma(stat.QNodes),
		Comma(int((float64(stat.QNodes) / float64(stat.Nodes+stat.QNodes) * 100))),
		Comma(int(float64(stat.Nodes+stat.QNodes)/stat.TimeElapsed.Seconds())),
		Comma(stat.UpperCuts),
		Comma(stat.LowerCuts),
		Comma(stat.TtHits),
		Comma(int((float64(stat.TtHits) / float64(stat.Nodes) * 100))),
		Comma(stat.TtWrites),
		Comma(stat.TtUpdates),
		Comma(len(tt)),
		Comma(stat.TtCulls),
	)
}

func Go(p *Pos, eiChan chan EngineInfo) (res string, info string, srch *Search) {
	// computer makes moves now!
	// 	var pv PV
	// 	var move Move
	// 	var score int
	var bookSuccess bool

	srch = NewSearch(Fen(""))
	srch.TimeStart = time.Now()
	srch.Score = NEGINF
	srch.ExplosionLimit = 2000000
	srch.MaxDepthToSearch = 8
	srch.EngineInfoChan = eiChan

	srch.BestMove, bookSuccess = ChooseBookMove(p)
	if bookSuccess == false {
		// srch root
		// adjust pv if filled
		if srch.PV.count > 0 {
			for srch.PV.ply < p.Ply && srch.PV.count > 0 {
				srch.PV.ply++ // walking forward up plies
				srch.PV.count--
				copy(srch.PV.moves[0:], srch.PV.moves[1:]) // shift movelist up by one
				info += fmt.Sprintf("pv chomp p.ply=%v, PV.ply=%v -- %v\n", p.Ply, srch.PV.ply, srch.PV)
			}
		}

		if srch.PV.ply > p.Ply {
			srch.PV.ply = p.Ply
		}

		srch.TimeStart = time.Now()
		srch.BestMove, srch.Stats.Score = SearchRoot(p, srch)
		srch.Stats.TimeElapsed = time.Since(srch.TimeStart)

	} else {
		info = "Book move found"
	}

	MakeMove(srch.BestMove, p)

	info += "fen: (" + BoardToFEN(p) + ")\n"
	info += result(p)
	res = fmt.Sprintf("move %v", MoveToAlg(srch.BestMove))
	return
}

func result(p *Pos) (s string) {
	// TODO ICS handles winning and losing. Plus sending these strings to KDE Knights crashes it!
	var win, lose string
	nummoves := len(GenerateAllMoves(p))

	if nummoves == 0 {
		// 		GameOver = true
		if p.InCheck == BLACK {
			win = "white"
			lose = "black"
		}
		if p.InCheck == WHITE {
			win = "black"
			lose = "white"
		}
		if p.InCheck == -1 {
			s += fmt.Sprintf("result 1/2 - 1/2 {draw - stalemate}\n")

		} else {
			s += fmt.Sprintf("result {%v}-{%v} {win mates}\n", win, lose)
		}
	}
	if p.Fifty >= 50 {
		// 		GameOver = true
		s += fmt.Sprintf("result 1/2 - 1/2 {draw - fifty move rule}\n")
	}
	return
}

func ParseUserMove(input string, p *Pos) (m Move, err string) {
	err = "Illegal move: " + input
	input = strings.ToLower(strings.TrimSpace(input))
	re, e := regexp.Compile("[a-h][1-8][a-h][1-8][qbnr]?")
	if e != nil {
		panic("Regexp did not compile!")
	}
	if re.MatchString(input) == false {
		err = "Error (unparseable as user move): " + input
		return
	}
	str := strings.Split(input, "")
	m.from = (AlgToDec(str[0] + str[1]))
	m.to = (AlgToDec(str[2] + str[3]))
	if len(str) > 4 {
		m.mtype = PROMOTE
		if str[4] == "q" {
			m.extra = QUEEN
		}
		if str[4] == "b" {
			m.extra = BISHOP
		}
		if str[4] == "r" {
			m.extra = ROOK
		}
		if str[4] == "n" {
			m.extra = NIGHT
		}
	}
	moves := GenerateAllMoves(p)
	for _, mv := range moves {
		if m.from == mv.from && m.to == mv.to {
			err = ""
			// if the move is found in the move list and it is a promote then use the one we created above
			if m.mtype == PROMOTE {
				return
			}
			// use the better defined matched move
			m = mv
			return
		}
	}
	return
}

func MakeUserMove(m Move, p *Pos) (s string) {

	// 	if GameOver == true {
	// 		s = "Game Over"
	// 		return
	// 	}
	MakeMove(m, p)
	s = result(p)
	return
}

func (srch Search) StopSearch() bool {
	// are we stopping?
	if srch.Stop {
		return true
	} // yes
	// otherwise only check every 8192 nodes
	if (srch.Stats.Nodes+srch.Stats.QNodes)&0x7fff != 0 {
		return false
	}
	// send statistics (note will block if no listener!)
	srch.Stats.TimeElapsed = time.Since(srch.TimeStart)
	srch.EngineInfoChan <- NewEngineInfo(srch)

	// GameDurationToSearch ==0 means search forever
	if srch.MaxDurationOfSearch == 0 {
		return false
	}
	// have we passed the time limit for searching?
	if time.Since(srch.TimeStart) < srch.MaxDurationOfSearch {
		return false
	}
	srch.Stop = true
	return true

}
