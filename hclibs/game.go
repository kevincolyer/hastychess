package hclibs

// useful routines for playing games against people or computer
import "fmt"
import "strings"
import "regexp"


func GameInit() {
	tt = make(map[string]TtData)
	return
}


func Go(p *Pos) (res string) {
	// computer makes moves now!
	var pv PV
	var move Move
	var success bool

	StatNodes = 0
	StatTtHits = 0

	//     my Int score = negamax(depth,p)
	//     say p.board
	move, success = ChooseBookMove(p)
	if success == false {
		//          for 1..depth -> Int i {
		pv = Search(*(p), 2, GameDepthSearch) // global variable for depth of search...
		res += fmt.Sprintf("# (score) PV=%v", pv.moves)
		move = pv.moves[0]
	} else {
		res += fmt.Sprintf("# book move found")
	}
	if GameUseStats {
		res += fmt.Sprintf("\n# STATS (Score %v | nodes %v | tt_hits %v (%v%%)  | tt writes %v | tt updates %v | tt size %v)\n", Comma(pv.score), Comma(StatNodes), Comma(StatTtHits), Commaf(float64(StatTtHits)/float64(StatNodes)*100) ,Comma(StatTtWrites), Comma(StatTtUpdates), Comma(len(tt)))
	}
	res += fmt.Sprintf("move %v\n", MoveToAlg(move))
	MakeMove(move, p)
	return res + result(p)
}

func result(p *Pos) (s string) {
	var win, lose string
	nummoves := len(GenerateAllMoves(p))

	if nummoves == 0 {
		GameOver = true
		if p.InCheck == 0 {
			win = "white"
			lose = "black"
		}
		if p.InCheck == -1 {
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
