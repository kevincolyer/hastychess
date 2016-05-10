package hclibs

import "strings"
import "fmt"
import "os"
import "strconv"
import "math/rand"

// import "math"

func Die(e string) {
	fmt.Println(e)
	os.Exit(1) // can't depend on panic
}

func Pick(i int) int {
	return rand.Intn(i)
}
func FENToNewBoard(f string) Pos {
	var p Pos
	FENToBoard(f, &p)
	return p
}

func FENToBoard(f string, p *Pos) *Pos {
	// parse fen
	fenfields := strings.Split(f, " ")
	// board =0 , side2move, castling, enpassant, halfmoveclock, fullmovecounter
	side2move := fenfields[1]
	castling := fenfields[2]
	//fmt.Println(castling)
	enpassant := fenfields[3]
	p.HalfMoveClock, _ = strconv.Atoi(fenfields[4])
	p.FullMoveClock, _ = strconv.Atoi(fenfields[5])
	// parse board
	ranks := strings.Split(fenfields[0], "/")
	if len(ranks) != 8 {
		Die(fmt.Sprintf("Invalid FEN %v - not enough ranks (found %v) ", f, len(ranks)))
	}
	for r, line := range ranks {
		if line == "" {
			line = "........"
		}
		line = strings.Replace(line, "1", ".", -1)
		line = strings.Replace(line, "2", "..", -1)
		line = strings.Replace(line, "3", "...", -1)
		line = strings.Replace(line, "4", "....", -1)
		line = strings.Replace(line, "5", ".....", -1)
		line = strings.Replace(line, "6", "......", -1)
		line = strings.Replace(line, "7", ".......", -1)
		line = strings.Replace(line, "8", "........", -1)

		rank := 7 - r // reverse loading the board as fen is top down
		if len(line) != 8 {
			Die(fmt.Sprintf("Invalid FEN line (%v) on rank %v - wrong length (found %v)", line, rank+1, len(line)))
		}
		files := strings.Split(line, "")
		for file, piece := range files {
			switch piece {
			case "P":
				p.Board[rank<<4+file] = PAWN
			case "Q":
				p.Board[rank<<4+file] = QUEEN
			case "K":
				p.Board[rank<<4+file] = KING
			case "B":
				p.Board[rank<<4+file] = BISHOP
			case "R":
				p.Board[rank<<4+file] = ROOK
			case "N":
				p.Board[rank<<4+file] = NIGHT

			case "p":
				p.Board[rank<<4+file] = pawn
			case "q":
				p.Board[rank<<4+file] = queen
			case "k":
				p.Board[rank<<4+file] = king
			case "b":
				p.Board[rank<<4+file] = bishop
			case "r":
				p.Board[rank<<4+file] = rook
			case "n":
				p.Board[rank<<4+file] = night
			case ".":
				p.Board[rank<<4+file] = EMPTY
			default:
				Die(fmt.Sprintf("unrecognised character in FEN definition [%v]", piece))
			}
		}
	}
	// continue parsing FEN
	// side to move
	if side2move == "" {
		side2move = "w"
	}
	if side2move == "w" {
		p.Side = WHITE
	} else {
		p.Side = BLACK
	}
	// castling rules
	p.Castled[0], p.Castled[1], p.Castled[2], p.Castled[3] = true, true, true, true
	if strings.Contains(castling, "Q") {
		p.Castled[WHITE*2+QS] = false
	}
	if strings.Contains(castling, "K") {
		p.Castled[WHITE*2+KS] = false
	}
	if strings.Contains(castling, "q") {
		p.Castled[BLACK*2+QS] = false
	}
	if strings.Contains(castling, "k") {
		p.Castled[BLACK*2+KS] = false
	}
	// enpassant
	if enpassant == "-" {
		p.EnPassant = -1
	} else {
		p.EnPassant = AlgToDec(enpassant)
	}
	// Find king and count taken pieces
	p.TakenPieces[0] = 16
	p.TakenPieces[1] = 16
	for _, sq := range GRID {
		pc := p.Board[sq]
		if PieceType(pc) == KING {
			p.King[PieceColour(pc)] = sq
		}
		if pc != EMPTY {
			p.TakenPieces[PieceColour(pc)]--
		} // exists so not taken.
	}
	p.InCheck = -1 // -1 == neither side
	if InCheck(p.King[WHITE], WHITE, p) {
		p.InCheck = WHITE
	}
	if InCheck(p.King[BLACK], BLACK, p) {
		p.InCheck = BLACK
	}
	return p
}

func BoardToStr(p *Pos) string {
	ptos := [...]string{".", "P", "N", "K", "-", "B", "R", "Q", "-", "p", "n", "k", "-", "b", "r", "q"}
	var s string
	for rank := 7; rank >= 0; rank-- { // reverse order
		s += " "
		for file := 0; file < 8; file++ {
			s += ptos[p.Board[rank<<4+file]]
			// 			fmt.Println(ptos[p.Board[rank<<4+file]])
		}
		s += fmt.Sprintf(" %v\n", rank+1)
	}
	s += " abcdefgh"
	return s
}

func BoardToStrWide(p *Pos) string {
	ptos := [...]string{".", "P", "N", "K", "-", "B", "R", "Q", "-", "p", "n", "k", "-", "b", "r", "q"}
	var s string
	for rank := 7; rank >= 0; rank-- { // reverse order
		s += fmt.Sprintf(" %v   ", rank+1)

		for file := 0; file < 8; file++ {
			s += ptos[p.Board[rank<<4+file]] + "  "
			// 			fmt.Println(ptos[p.Board[rank<<4+file]])
		}
		s += "\n"
	}
	s += "\n     A  B  C  D  E  F  G  H"
	return s
}

func (p *Pos) String() string {
	return fmt.Sprintf("%s\n", BoardToStrWide(p))
}

func PieceColour(piece int) int {
	if piece>>3 > 0 {
		return BLACK
	} else {
		return WHITE
	}
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func PieceType(p int) int {
	return p & 7 // mask only the piece type (colour masked out)
}

func AlgToDec(alg string) int {
	alg = strings.TrimSpace(strings.ToLower(alg))
	//die "Not proper algebraic notation" if $alg !~~ m/<[a..h]> <[1..8]>/;
	return int(alg[0]) - 97 + (int(alg[1])-49)*16
}

func DecToAlg(dec int) string {
	return string((dec&7)+97) + string((dec>>4)+48+1)
}

// func Xside(p Pos) int {
//     if p.Side == WHITE { return BLACK } else { return WHITE }
// }

func OtherSide(p Pos) int {
	if p.Side == WHITE {
		return BLACK
	} else {
		return WHITE
	}
}

func TtKey(p *Pos) string {
	return fmt.Sprintf("%v %v %v", p.Board, p.Castled, p.Side)
}

func Abs(n int) int {
	if n < 0 {
		n = -n
	}
	return n
}

func MhDistance(from, to int) int {
	return Abs(from>>4-to>>4) + Abs(from&7-to&7)
}

func MoveToAlg(m Move) (s string) {
	ptos := [...]string{".", "P", "N", "K", "-", "B", "R", "Q", "-", "p", "n", "k", "-", "b", "r", "q"}
	s = DecToAlg(m.from) + DecToAlg(m.to)
	if m.mtype&PROMOTE > 0 {
		s += ptos[m.extra] // this will show whites symbols, but that is ok. In MakeMove the correct piece is shown.
	}
	return
}

func AlgToMove(s string) (move Move) {
	r := []rune(s)
	//  fmt.Print(s+"----",len(r))
	//  fmt.Println( string(r[0])+string(r[1])+string(r[2])+string(r[3]))
	move.from = AlgToDec(string(r[0]) + string(r[1]))
	move.to = AlgToDec(string(r[2]) + string(r[3]))
	if len(r) == 5 {
		//                 e:=string(s)[4]
		e := r[4]
		move.mtype = PROMOTE
		if e == 'q' {
			move.extra = QUEEN
		}
		if e == 'b' {
			move.extra = BISHOP
		}
		if e == 'r' {
			move.extra = ROOK
		}
		if e == 'n' {
			move.extra = NIGHT
		}
	}
	return
}

/* TODO
func AlgMove(from, to, type, extra int) string {}
*/
