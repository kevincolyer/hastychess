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

func Go(p *Pos) (res string, info string) {
	// computer makes moves now!
	// 	var pv PV
	var move Move
	var score int
	var success bool

	StatNodes = 0
	StatQNodes = 0
	StatTtHits = 0
	StatTtWrites = 0
	StatTtUpdates = 0
	StatUpperCuts = 0
	StatLowerCuts = 0

	StatTimeStart = 0 // not sure what type needed here
	StatTimeElapsed = 0

	//     my Int score = negamax(depth,p)
	//     say p.board
	move, success = ChooseBookMove(p)
	if success == true {
		info += fmt.Sprintf("# book move found")
	} else {
		// search root
		// adjust pv if filled
		if pv.count > 0 {
			for pv.ply < p.Ply && pv.count > 0 {
				pv.ply++ // walking forward up plies
				pv.count--
				copy(pv.moves[0:], pv.moves[1:]) // shift movelist up by one
				fmt.Printf("pv chomp p.ply=%v, pv.ply=%v -- %v\n", p.Ply, pv.ply, pv)
			}
		}
		start := time.Now()
		// some computation
		move, score = SearchRoot(p, GameDepthSearch, &pv) // global variable for depth of search...
		elapsed := time.Since(start)

		if GameUseStats {
			info += "# fen: (" + BoardToFEN(p) + ")"
			info += fmt.Sprintf("\n# STATS Score %v | nodes %v | qnodes %v (%v%%)| nps %v | uppercuts %v | lowercuts %v |\n# STATS tt_hits %v (%v%%) | tt writes %v | tt updates %v | tt size %v | tt culls %v |\n", Comma(score), Comma(StatNodes), Comma(StatQNodes), Comma(int((float64(StatQNodes) / float64(StatNodes+StatQNodes) * 100))), Comma(int(float64(StatNodes+StatQNodes)/elapsed.Seconds())), Comma(StatUpperCuts), Comma(StatLowerCuts), Comma(StatTtHits), Comma(int((float64(StatTtHits) / float64(StatNodes) * 100))), Comma(StatTtWrites), Comma(StatTtUpdates), Comma(len(tt)), Comma(StatTtCulls))
		}
	}

	info += result(p)
	if GameOver == false {
		res = fmt.Sprintf("move %v\n", MoveToAlg(move))
		if UCI() {
			res = "best" + res
		} else {
			res += "#\n"
		}
		MakeMove(move, p)
	}
	return
}

func result(p *Pos) (s string) {
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
	err = "# Not a valid move"
	input = strings.ToLower(strings.TrimSpace(input))
	re, e := regexp.Compile("[a-h][1-8][a-h][1-8][qbnr]?")
	if e != nil {
		panic("Regexp did not compile!")
	}
	if re.MatchString(input) == false {
		err = "# Unparseable"
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
			if m.mtype == PROMOTE {
				return
			}
			m = mv
			return
		}
	}
	return
}

func MakeUserMove(m Move, p *Pos) (s string) {
	s = ""
	if GameOver == true {
		s = "Game Over"
		return
	}
	MakeMove(m, p)
	s = result(p)
	return
}

func StopSearch() bool {
	select {
	case <-Control:
		fmt.Print("# detected search stop\n") // open channel means we can keep searching
		return true
	default:
		return false
	}
}
