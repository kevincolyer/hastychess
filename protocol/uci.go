// Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer

package protocol

import (
	"fmt"
	"io"
	// 	"bufio"
	// 	"github.com/kevincolyer/hastychess/hclibs"
	// 	"regexp"
	// 	"strconv"
	// 	"strings"
	// 	"time"
)

type Uci struct {
	In      io.Reader
	Out     io.Writer
	Options CLIOptions
}

func NewUCI(in io.Reader, out io.Writer, options CLIOptions) (*Uci, error) {
	proto := Uci{Out: out, In: in, Options: options}
	_, e := io.WriteString(proto.Out, "Welcome to "+options.NameVersion+"\n")
	return &proto, e
}

func (proto *Uci) Echo() (e error) {
	b := make([]byte, 256)
	if n, err := proto.In.Read(b); err != io.EOF {
		io.WriteString(proto.Out, string(b[:n]))
	} else {
		e = fmt.Errorf("EOF for input: none sent?")
	}
	return
}

func (proto *Uci) Start() (e error) {
	io.WriteString(proto.Out, "started ok\n")
	e = fmt.Errorf("done")
	return
}

//uci feedback
// 			fmt.Printf("info depth %v score upperbound %v time %v nodes %v nps %v pv %v\n", depth, bestscore, Milliseconds(elapsed), srch.Stats.Nodes+srch.Stats.QNodes, int(float64(srch.Stats.Nodes+srch.Stats.QNodes)/elapsed.Seconds()), srch.PV)
// 		}

//===========================================================================

/***************************************************************
 * Use the UCI chess protocol
 *
 * lots of fun
 * See: http://wbec-ridderkerk.nl/html/UCIProtocol.html
 ****************************************************************/

// little helper function ("hack") needed to switch from console to uci if command line switches are ignored
// func ucihelper() {
// 	name := fmt.Sprintf("HastyChess v%v", hclibs.VERSION)
// 	fmt.Println("id name " + name + "\nid author Kevin Colyer 2016")
// 	// Send options to GUI here...
// 	//
// 	//
// 	fmt.Println("debug on")
// 	fmt.Println("uciok")
// }
//
// func mainIcs(scanner *bufio.Scanner) {
// 	// 	version := 0.99
// 	//	name := fmt.Sprintf("HastyChess v%v", hclibs.VERSION)
// 	// 	scanner := bufio.NewScanner(os.Stdin)
//
// 	p := hclibs.FENToNewBoard(hclibs.STARTFEN)
//
// 	//hclibs.GameUseBook = false // UCI gui does book - unless option below...
//
// 	hclibs.GameOver = false
// 	hclibs.GameDisplayOn = false
// 	hclibs.GameDepthSearch = 8
// 	hclibs.GameForce = false
// 	// main input loop
// 	for {
// 		for scanner.Scan() {
//
// 			//input := strings.ToLower(strings.TrimSpace(scanner.Text()))
// 			input := strings.TrimSpace(scanner.Text())
// 			fmt.Printf("info string echo server sent (%v)\n", input)
// 			// 			time.Sleep(time.Second)
// 			//fmt.Pri(ntf("You said [%v]\n", input)
// 			switch {
//
// 			case input == "uci":
// 				//         fmt.Println("id name " + name + "\nid author Kevin Colyer 2016")
// 				//         Send options to GUI here...
// 				//
// 				//
// 				//         fmt.Println("uciok")
// 				ucihelper()
// 				continue
//
// 			case strings.HasPrefix(input, "isready"):
// 				fmt.Println("readyok")
//
// 			case strings.HasPrefix(input, "quit"):
// 				return
//
// 			case strings.HasPrefix(input, "ucinewgame"):
// 				p = hclibs.FENToNewBoard(hclibs.STARTFEN)
// 				hclibs.GameOver = false
//
// 			case strings.HasPrefix(input, "go"):
// 				go uciGo(input, &p)
// 				//  lots of sub verbs here...
// 				continue
// 				// start searching now!
//
// 			case input == "ping":
// 				fmt.Println("pong")
// 				continue
// 				// start searching now!
// 			case strings.HasPrefix(input, "stop"):
// 				if hclibs.Control != nil {
// 					close(hclibs.Control)
// 				} // tells searches to finish
// 				continue
// 				// stop searching now!
//
// 			case strings.HasPrefix(input, "position"):
// 				uciPosition(input, &p)
// 				// position fen moves...
// 				// position startpos moves...
// 				continue
//
// 			}
// 		}
// 	}
//
// }
//
// func uciPosition(input string, p *hclibs.Pos) {
// 	fen, moves := ParsePositionInput(input)
// 	*p = hclibs.FENToNewBoard(fen)
//
// 	for _, m := range moves {
// 		/*if !hclibs.IsValidMove(m, p) {
// 					fmt.Printf("info currmove %v (Sent bad move by server %v in %v)\n", m, m, moves)
// 		        }*/
// 		hclibs.MakeMove(m, p) // note no error checking here!!!
// 	}
// }
//
// // parses a string "position <fen> <move>*" or "position startpos <move>*" returns
// // fen as string and slice of Move types.
// func ParsePositionInput(input string) (fen string, moves []hclibs.Move) {
//
// 	re := regexp.MustCompile("[a-h][1-8][a-h][1-8][qbnr]?")
//
// 	ef := 2 // end fen = start of moves (if any)
// 	f := strings.Split(input, " ")
// 	// can ignore first field == "position"
// 	if len(f) > 2 && f[1] == "fen" {
// 		fen = strings.Join(f[2:], " ")
// 		return
// 	}
// 	if f[1] == "startpos" {
// 		fen = hclibs.STARTFEN
// 	} else {
// 		// assume command is "moves"
// 		//look ahead to moves - what is skipped is a fen
// 		for ; ef < len(f); ef++ {
// 			if f[ef] == "moves" {
// 				break
// 			}
// 		}
//
// 		fen = strings.Join(f[1:ef], " ")
//
// 	}
//
// 	for ; ef < len(f); ef++ {
// 		if re.MatchString(f[ef]) {
// 			moves = append(moves, hclibs.AlgToMove(f[ef]))
// 		}
//
// 	}
// 	return
// }
//
// func uciGo(input string, p *hclibs.Pos) {
// 	// position fen moves...
// 	// position startpos moves...
//
// 	if hclibs.Control != nil && hclibs.StopSearch() == false { // if channel still open, close it
// 		close(hclibs.Control)
// 	}
// 	hclibs.Control = make(chan string)
// 	go func() {
// 		for m := range hclibs.Control {
// 			fmt.Println(m)
// 		} // send messages from search straight to console until channel closed
// 	}()
// 	// channel shut in "stop" command or at end of a search or at a time out...
//
// 	// expand imput string to parse sub verbs
// 	f := strings.Split(input, " ")
// 	if len(f) > 1 {
// 		f = f[1:]
// 	} // shift out "go"
// 	// chomp settings two at a time
// 	count := 0
// 	for len(f) > 1 {
// 		// double word commands
// 		count = 2
// 		// DEPTH
// 		if f[0] == "depth" && len(f) > 1 {
// 			d, err := strconv.Atoi(f[1])
// 			if err != nil {
// 				fmt.Println("info string Setting depth: Please specify a number")
// 				d = 0
// 			} else {
// 				hclibs.GameDepthSearch = d
// 				fmt.Println("info string  Set depth to ", d)
// 			}
//
// 		}
// 		// MOVETIME
// 		if f[0] == "movetime" && len(f) > 1 {
// 			d, err := strconv.Atoi(f[1])
// 			if err != nil {
// 				fmt.Println("info string Setting movetime: Please specify a number")
// 				d = 0
// 			} else {
//
// 				fmt.Println("info string Set movetime to ", d)
// 				hclibs.GameDurationToSearch = time.Duration(d * 1000 * 1000) // milliseconds to nanoseconds
// 			}
// 		}
// 		// Single verb commands
// 		// INFINITE (single)
// 		// PONDER etc.
// 		// 		fmt.Println("# comsuming tokens:",count)
// 		f = f[count:] // consume tokens
// 	}
//
// 	res, _ := hclibs.Go(p)
// 	fmt.Println(res)
// 	return
// }
