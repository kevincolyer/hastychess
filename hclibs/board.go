//Package hclibs ... Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import "strings"
import "fmt"
import "os"
import "strconv"
import "math/rand"
import "github.com/fatih/color"
import "time"

// import "math"

func init() {
	rand.Seed(time.Now().UTC().UnixNano()) // a really HOT cup of strong tea
}

func Die(e string) {
	fmt.Println(e)
	os.Exit(1) // can't depend on panic
}

func Pick(i int) int {
	return rand.Intn(i)
}
func (f Fen) NewBoard() Pos {
	return FENToNewBoard(string(f))
}
func FENToNewBoard(f string) Pos {
	var p Pos
	FENToBoard(f, &p)
	return p
}

func NewRBCFEN(diff int) string {
	//"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	b := []byte("rnbqkbnr/pppppppp")
	bp := []byte("rnbq")
	w := []byte("PPPPPPPP/RNBQKBNR")
	wp := []byte("RNBQ")
	if diff > 0 {
		if diff > 3 {
			diff = 3
		}
		for s := 0; s < diff*5; s++ {
			// swap black piece at index i with one from table at rndindex
			i := rbcrand(16)
			if b[i] != 'k' {
				b[i] = bp[rbcrand(4)]
			}
			// swap white piece at index i with one from table at rndindex
			i = rbcrand(16)
			if w[i] != 'K' {
				w[i] = wp[rbcrand(4)]
			}
		}
	}
	return string(b) + "/8/8/8/8/" + string(w) + " w KQkq - 0 1"
}

func rbcrand(n int) int {
	i := rand.Intn(n)
	if i == 8 {
		i++
	}
	return i
}
func BoardToFEN(p *Pos) string {

	ptos := [...]string{".", "P", "N", "K", "-", "B", "R", "Q", "-", "p", "n", "k", "-", "b", "r", "q"}
	var fen, castling, enpassant, side string

	for rank := 7; rank >= 0; rank-- {
		j := 0
		for file := 0; file < 8; file++ {
			if p.Board[rank<<4+file] == EMPTY {
				j++
				continue
			}
			if j > 0 {
				fen += fmt.Sprintf("%d", j)
				j = 0
			}
			fen += ptos[p.Board[rank<<4+file]]

		}
		if j == 8 {
			fen += "8"
		} // or ignore?
		if rank > 0 {
			fen += "/"
		}
	}

	if p.Side == 0 {
		side = "w"
	} else {
		side = "b"
	}

	if p.Castled[WHITE*2+KS] == false {
		castling += "K"
	}
	if p.Castled[WHITE*2+QS] == false {
		castling += "Q"
	}
	if p.Castled[BLACK*2+KS] == false {
		castling += "k"
	}
	if p.Castled[BLACK*2+QS] == false {
		castling += "q"
	}
	if castling == "" {
		castling = "-"
	}

	enpassant = DecToAlg(p.EnPassant)
	if p.EnPassant == -1 {
		enpassant = "-"
	}
	//fmt.Println(fen)
	return fmt.Sprintf("%s %s %s %s %d %d", fen, side, castling, enpassant, p.HalfMoveClock, p.FullMoveClock)
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
	p.Hash = TTZKey(p)
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

func BoardToStrColour(p *Pos) string {
	ptos := [...]string{" ", "P", "N", "K", "-", "B", "R", "Q", "-", "p", "n", "k", "-", "b", "r", "q"}
	var s string
	tog := false // true==black square
	whitepc := color.New(color.FgHiWhite).SprintFunc()
	whitesq := color.New(color.BgBlue).SprintFunc()
	blackpc := color.New(color.FgRed).SprintFunc()
	blacksq := color.New(color.BgBlack).SprintFunc()
	s += "     A  B  C  D  E  F  G  H\n\n"
	for rank := 7; rank >= 0; rank-- { // reverse order
		s += fmt.Sprintf(" %v  ", rank+1)

		for file := 0; file < 8; file++ {
			pc := ptos[p.Board[rank<<4+file]]
			if strings.ToUpper(pc) == pc {
				pc = whitepc(pc + " ")
			} else {
				pc = blackpc(pc + " ")
			}
			if tog {
				s += whitesq(" " + pc)
			} else {
				s += blacksq(" " + pc)
			}
			tog = !tog
		}
		s += fmt.Sprintf("  %v\n", rank+1)
		//s += "\n"
		tog = !tog
	}
	s += "\n     A  B  C  D  E  F  G  H\n"
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
	if color.NoColor == false {
		return BoardToStrColour(p)
	}
	return BoardToStrWide(p)
}

func Side(piece int) int {
	p := piece >> 3
	if p == BLACK {
		return BLACK
	}
	if p == WHITE {
		return WHITE
	}
	Die(fmt.Sprintf("Side was passed a piece value of %v which is not possible!", piece))
	return 0
}

func Xside(side int) int {
	return 1 - side
}
func Onboard(i int) bool {
	return i&0x88 == 0
}

func Offboard(i int) bool {
	return i&0x88 != 0
}

func PieceColour(piece int) int {
	if piece>>3 > 0 {
		return BLACK
	}
	return WHITE
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
	}
	return WHITE
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
	if m.mtype == PROMOTE {
		s += strings.ToLower(ptos[m.extra]) // force to lowe because xboard etc expect lower case in promotions?
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

/* TODO SAN notation so that...
 * https://en.wikipedia.org/wiki/Portable_Game_Notation#Example
 *
 * SANtoMove Move,error (if not parsable) - would only let you know if parsable not if legal in this position.
 * MovetoSAN(Move) string
 */

func MoveToSAN(move Move) (san string) {
	if move.mtype == O_O {
		return "O-O"
	}
	if move.mtype == O_O_O {
		return "O-O-O"
	}

	//prefix
	if move.piece != PAWN {
		san = pieceToInitial(move.piece)
	}

	// disambiguate here? (need to see either board or move list to determine - move list handier!)
	// TODO

	// if capture
	if move.mtype == CAPTURE || move.mtype == EPCAPTURE {
		san += "x"
	}

	// to square
	san += DecToAlg(move.to)

	// suffix
	if move.mtype == PROMOTE {
		san += "=" + pieceToInitial(move.extra)
	}
	if move.subtype == CHECK {
		san += "+"
	}
	// if move.subtype==CHECKMATE { san+="+" }  or RESULT
	return
}

func pieceToInitial(p int) (s string) {
	if p == PAWN {
		s = "P"
	}
	if p == KING {
		s = "K"
	}
	if p == QUEEN {
		s = "Q"
	}
	if p == BISHOP {
		s = "B"
	}
	if p == ROOK {
		s = "R"
	}
	if p == NIGHT {
		s = "N"
	}
	return
}

func initialToPiece(s string) (p int) {
	if s == "K" {
		p = KING
	}
	if s == "Q" {
		p = QUEEN
	}
	if s == "B" {
		p = BISHOP
	}
	if s == "R" {
		p = ROOK
	}
	if s == "N" {
		p = NIGHT
	}
	if s == "" || s == "P" {
		p = PAWN
	}
	return
}
