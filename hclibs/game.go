//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

// useful routines for playing games against people or computer
import "fmt"
import "strings"
import "regexp"
import "time"

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

var pv PV

type Search struct {
	Nodes     int
	QNodes    int
	TtHits    int
	TtWrites  int
	TtCulls   int
	TtUpdates int
	UpperCuts int
	LowerCuts int

	TimeStart           time.Time
	TimeElapsed         time.Duration
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

	UseTT   bool
	UseBook bool
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

func (stat Search) String() string {
	return fmt.Sprintf("STATS score %v | nodes %v | qnodes %v (%v%%)| nps %v | uppercuts %v | lowercuts %v |\n# STATS tt_hits %v (%v%%) | tt writes %v | tt updates %v | tt size %v | tt culls %v |\n", Comma(stat.Score), Comma(stat.Nodes), Comma(stat.QNodes), Comma(int((float64(stat.QNodes) / float64(stat.Nodes+stat.QNodes) * 100))), Comma(int(float64(stat.Nodes+stat.QNodes)/stat.TimeElapsed.Seconds())), Comma(stat.UpperCuts), Comma(stat.LowerCuts), Comma(stat.TtHits), Comma(int((float64(stat.TtHits) / float64(stat.Nodes) * 100))), Comma(stat.TtWrites), Comma(stat.TtUpdates), Comma(len(tt)), Comma(stat.TtCulls))
}

func Go(p *Pos) (res string, info string, srch Search) {
	// computer makes moves now!
	// 	var pv PV
	// 	var move Move
	// 	var score int
	var bookSuccess bool

	srch = Search{
		TimeStart:        time.Now(),
		Score:            NEGINF,
		ExplosionLimit:   2000000,
		MaxDepthToSearch: 6,
	}

	// 	StatTimeElapsed = 0
	// 	GameStopSearch = false

	srch.BestMove, bookSuccess = ChooseBookMove(p)
	if bookSuccess == false {
		// srch root
		// adjust pv if filled
		if pv.count > 0 {
			for pv.ply < p.Ply && pv.count > 0 {
				pv.ply++ // walking forward up plies
				pv.count--
				copy(pv.moves[0:], pv.moves[1:]) // shift movelist up by one
				if GameProtocol == PROTOCONSOLE {
					info += fmt.Sprintf("# pv chomp p.ply=%v, pv.ply=%v -- %v\n", p.Ply, pv.ply, pv)
				}
			}
		}
		if pv.ply > p.Ply {
			pv.ply = p.Ply
		} // in case reset game - pv is global (yuk) and not reset so far
		srch.TimeStart = time.Now()
		// some computation
		srch.BestMove, srch.Score = SearchRoot(p, srch.MaxDepthToSearch, &pv, &srch) // global variable for depth of search...
		srch.TimeElapsed = time.Since(srch.TimeStart)
		if GameProtocol == PROTOCONSOLE {
			info += "# fen: (" + BoardToFEN(p) + ")"
			info += fmt.Sprintf("\n# PV %v", pv)
			info += fmt.Sprintf("\n# %v", srch)
		}
	}

	info += result(p)

	res = fmt.Sprintf("move %v", MoveToAlg(srch.BestMove))

	//TODO if UCI() {
	//    res = "best" + res
	//} /*else {
	//		res += "#\n"
	//	}*/
	MakeMove(srch.BestMove, p)

	return
}

func result(p *Pos) (s string) {
	// ICS handles winning and losing. Plus sending these strings to KDE Knights crashes it!
	if GameProtocol == PROTOUCI {
		return
	}
	var win, lose string
	nummoves := len(GenerateAllMoves(p))

	if nummoves == 0 {
		GameOver = true
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
		GameOver = true
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

	if GameOver == true {
		s = "Game Over"
		return
	}
	MakeMove(m, p)
	s = result(p)
	return
}

func (s Search) StopSearch() bool {
	// are we stopping?
	if s.Stop {
		return true
	} // yes
	// otherwise only check every 1000 nodes
	if (s.Nodes+s.QNodes)%1000 != 0 {
		return false
	}
	// GameDurationToSearch ==0 means search forever
	if s.MaxDurationOfSearch == 0 {
		return false
	}
	// have we passed the time limit for searching?
	if time.Since(s.TimeStart) < s.MaxDurationOfSearch {
		return false
	}
	s.Stop = true
	return true
	// func StopSearch() bool {
	// 	// are we stopping?
	// 	if GameStopSearch {
	// 		return true
	// 	} // yes
	//
	// 	// otherwise only check every 1000 nodes
	// 	if (StatNodes+StatQNodes)%1000 != 0 {
	// 		return false
	// 	}
	// 	// GameDurationToSearch ==0 means search forever
	// 	if GameDurationToSearch == 0 {
	// 		return false
	// 	}
	// 	// have we passed the time limit for searching?
	// 	if time.Since(StatTimeStart) < GameDurationToSearch {
	// 		return false
	// 	}
	// 	//         fmt.Println(time.Since(StatTimeStart))
	// 	//         fmt.Println(StatTimeStart)
	// 	//         fmt.Println(GameDurationToSearch)
	// 	// yes, so halt now and forever
	// 	if !UCI() {
	// 		fmt.Print("# Out of time to search...\n")
	// 	}
	// 	GameStopSearch = true
	// 	return true

	// 	select {
	// 	case <-Control:
	// 		if !UCI() {
	// 			fmt.Print("# detected search stop\n")
	// 		} // open channel means we can keep searching
	// 		return true
	// 	default:
	//return false
	// 	}
}
