// Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/kevincolyer/hastychess/hclibs"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	var flagXboard = flag.Bool("xboard", false, "Select xboard mode")
	var flagIcs = flag.Bool("ics", false, "Select ics server mode")
	var flagUci = flag.Bool("uci", false, "Select UCI (ics) server mode (same as -ics)")
	var flagConsole = flag.Bool("console", true, "Select console mode")
	var flagStats = flag.Bool("stats", true, "Enable printing of statistics")
	var flagUseBook = flag.Bool("book", true, "Enable the use of built in book moves")
	var flagUseTt = flag.Bool("tt", true, "Enable the use of Transposition Tables")
	var flagNoColor = flag.Bool("no-color", false, "Disable color output")
	flag.Parse()

	if *flagNoColor {
		color.NoColor = true // disables colorized output
	}

	hclibs.GameUseTt = *(flagUseTt)
	hclibs.GameUseStats = *(flagStats)
	hclibs.GameUseBook = *(flagUseBook)

	hclibs.GameInit()
	scanner := bufio.NewScanner(os.Stdin)
	switch {
	case *(flagXboard):
		hclibs.GameProtocol = hclibs.PROTOXBOARD
		color.NoColor = true
		mainXboard(scanner)
	case *(flagIcs):
		hclibs.GameProtocol = hclibs.PROTOUCI
		color.NoColor = true
		mainIcs(scanner)
	case *(flagUci):
		hclibs.GameProtocol = hclibs.PROTOUCI
		color.NoColor = true
		mainIcs(scanner)
	case *(flagConsole):
		hclibs.GameProtocol = hclibs.PROTOCONSOLE
		mainConsole(scanner)
	}
	fmt.Println("# Bye and thanks for playing!")

}

func mainConsole(scanner *bufio.Scanner) {

	var err string
	var result string
	var move hclibs.Move

	re, e := regexp.Compile("[a-h][1-8][a-h][1-8][qbnr]?")
	if e != nil {
		panic("Regexp did not compile!")
	}
	// 	version := 1.0
	hiwhite := color.New(color.FgHiWhite).PrintfFunc()
	hiwhite("Hello and welcome to HastyChess version %v\n\n", hclibs.VERSION)

	p := hclibs.FENToNewBoard(hclibs.STARTFEN)
	hclibs.GameOver = false
	hclibs.GameDisplayOn = true
	hclibs.GameDepthSearch = 4
	hclibs.GameForce = false
	if hclibs.GameDisplayOn {
		fmt.Println(&p)
	}
	quit := false
	fmt.Print("> ")
	hclibs.Control = make(chan string)
	// main input loop
QUIT:
	for quit == false {
	next:
		for scanner.Scan() {

			//input := strings.ToLower(strings.TrimSpace(scanner.Text()))
			input := strings.TrimSpace(scanner.Text())
			//fmt.Printf("You said [%v]\n", input)
			switch {
			case strings.Contains(input, "quit"):
				quit = true
				break QUIT

			// provide a way to change to xboard mode if use forgets to use commandline switch
			case strings.HasPrefix(input, "xboard"):
				color.NoColor = true
				hclibs.GameProtocol = hclibs.PROTOXBOARD
				mainXboard(scanner)
				return
			// provide a way to change to xboard mode if use forgets to use commandline switch
			case strings.HasPrefix(input, "uci"):
				color.NoColor = true
				hclibs.GameProtocol = hclibs.PROTOUCI
				ucihelper()
				mainIcs(scanner)
				return

			case strings.Contains(input, "move"):
				fields := strings.Fields(input)
				if len(fields) > 1 {
					move, err = hclibs.ParseUserMove(fields[1], &p)
					if err != "" {
						fmt.Println(err)
						break next
					}
				}
				result = hclibs.MakeUserMove(move, &p)
				fmt.Println(&p)
				fmt.Println(result)

			case re.MatchString(input):
				move, err = hclibs.ParseUserMove(input, &p)
				if err != "" {
					fmt.Println(err)
					break next
				}
				result = hclibs.MakeUserMove(move, &p)
				fmt.Println(&p)
				fmt.Println(result)

			case strings.Contains(input, "new"):
				p = hclibs.FENToNewBoard(hclibs.STARTFEN)
				hclibs.GameOver = false
				fmt.Println(&p)

			// special commands that allow testing of certain positions
			case strings.Contains(input, "kiwipete"):
				p = hclibs.FENToNewBoard("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
				hclibs.GameOver = false
				fmt.Println(&p)

			case strings.Contains(input, "pos4"):
				p = hclibs.FENToNewBoard("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
				hclibs.GameOver = false
				fmt.Println(&p)

			case strings.Contains(input, "pos5"):
				p = hclibs.FENToNewBoard("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
				hclibs.GameOver = false
				fmt.Println(&p)

			case strings.Contains(input, "end1"):
				p = hclibs.FENToNewBoard("8/8/4k3/7p/8/8/2K5/R6Q w - - 0 1")
				hclibs.GameOver = false
				fmt.Println(&p)

			case strings.Contains(input, "end2"):
				p = hclibs.FENToNewBoard("8/8/8/1bn5/8/2k5/8/2K5 w - - 0 1")
				hclibs.GameOver = false
				fmt.Println(&p)

			case strings.Contains(input, "end3"):
				p = hclibs.FENToNewBoard("8/8/8/1k6/8/7Q/3R4/2K5 w - - 0 1")
				hclibs.GameOver = false
				fmt.Println(&p)

			case strings.Contains(input, "end4"):
				p = hclibs.FENToNewBoard("8/8/8/k7/8/7Q/1R6/2K5 w - - 0 1")
				hclibs.GameOver = false
				fmt.Println(&p)

				// normal commands
			case strings.HasPrefix(input, "setboard"):
				fields := strings.Fields(input)
				if len(fields) == 1 {
					fmt.Println("Error (no Fen): " + input)
					break next
				}
				fen := strings.Join(fields[1:], " ")
				fmt.Println("# parsing fen [" + fen + "]")
				p = hclibs.FENToNewBoard(fen)
				// 				fmt.Println("Error (Not implemented yet!!!): " + input)
				break next
			case strings.Contains(input, "auto"):
				hclibs.GameForce = !hclibs.GameForce

			case strings.Contains(input, "go") || hclibs.GameForce == true:
				res, info := hclibs.Go(&p)
				fmt.Println(&p)
				fmt.Println(info)
				fmt.Println(res)

			case strings.Contains(input, "ping"):
				fmt.Println("pong")

			case strings.Contains(input, "divide"):
				fields := strings.Fields(input)
				if len(fields) == 1 {
					fmt.Println("Please specify a number to divide down to")
					break next
				}
				d, err := strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println("Please specify a number")
					break next
				}
				hclibs.Divide(d, &p)

			case strings.Contains(input, "perft"):
				fields := strings.Fields(input)
				if len(fields) == 1 {
					fmt.Println("Please specify a number to perft down to")
					break next
				}
				d, err := strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println("Please specify a number")
					break next
				}
				start := time.Now()
				nodes := hclibs.Perft(d, &p)
				elapsed := time.Since(start)
				fmt.Printf("\nPerft to depth %v gives %v nodes ", d, hclibs.Comma(nodes))
				fmt.Printf("(nps: %v)\n", hclibs.Comma(int(float64(nodes)/elapsed.Seconds())))

			case strings.Contains(input, "depth"):
				//  case strings.Contains(input,"fen") || strings.Contains(input,"setboard"):
				fields := strings.Fields(input)
				if len(fields) == 1 {
					fmt.Printf("Current depth is %d.\n", hclibs.GameDepthSearch)
					break next
				}
				d, err := strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println("Please specify a number")
					break next
				}
				if d > hclibs.MAXSEARCHDEPTH {
					d = hclibs.MAXSEARCHDEPTH
				}
				fmt.Printf("Current depth is %d. Setting depth to %d.\n", hclibs.GameDepthSearch, d)
				hclibs.GameDepthSearch = d
			}
			fmt.Print("> ")
		}
	}
}

func mainXboard(scanner *bufio.Scanner) {
	// see https://www.gnu.org/software/xboard/engine-intf.html
	// and http://hgm.nubati.net/newspecs.html for protocol info
	var err string
	var result string
	var move hclibs.Move

	re, e := regexp.Compile("^[a-h][1-8][a-h][1-8][qbnr]?")
	if e != nil {
		panic("Regexp did not compile!")
	}
	// 	version := 0.99
	name := fmt.Sprintf("HastyChess v%v", hclibs.VERSION)
	fmt.Printf("tellics say Hello and welcome to %v\n", name)

	fmt.Println("feature debug=1")
	fmt.Printf("feature myname=\"%v\"\n", name)

	// 	scanner := bufio.NewScanner(os.Stdin)
	p := hclibs.FENToNewBoard(hclibs.STARTFEN)
	hclibs.GameOver = false
	hclibs.GameDisplayOn = false
	// 	hclibs.GameDepthSearch = 4
	hclibs.GameForce = false
	if hclibs.GameDisplayOn {
		fmt.Println(&p)
	}
	quit := false
	// main input loop
QUIT:
	for quit == false {
	next:
		for scanner.Scan() {

			//input := strings.ToLower(strings.TrimSpace(scanner.Text()))
			input := strings.TrimSpace(scanner.Text())
			//fmt.Pri(ntf("You said [%v]\n", input)
			switch {

			case strings.HasPrefix(input, "accepted"):
				break next

			case strings.HasPrefix(input, "#"):
				break next

			case strings.HasPrefix(input, "xboard"):
				break next

			case strings.Contains(input, "protover 2"):
				fmt.Println("feature done=0")
				fmt.Printf("feature myname=\"%v\"", name)
				fmt.Println("feature usermove=1")
				fmt.Println("feature memory=0")
				fmt.Println("feature setboard=1")
				fmt.Println("feature ping=1")
				fmt.Println("feature sigint=0")  // SIGINT will halt a go program unless caught. Xboard sends SIGINT by default to let chess engine know it wants to talk!!! Was quite a problem!
				fmt.Println("feature sigterm=1") // We respond to SIGTERM - go stops!
				fmt.Println("feature variants=\"normal\"")
				fmt.Println("feature debug=1") // allows comments starting with hash symbols
				fmt.Println("feature done=1")

			case strings.Contains(input, "quit"):
				quit = true
				break QUIT

			case strings.HasPrefix(input, "result"):
				hclibs.GameOver = true
				break next

			case strings.HasPrefix(input, "usermove"), strings.HasPrefix(input, "move"):
				fields := strings.Fields(input)
				fmt.Println("# looking for move: fields ", len(fields))
				if len(fields) > 1 {
					move, err = hclibs.ParseUserMove(fields[1], &p)
					if err != "" {
						fmt.Println(err)
						break next
					}
					fmt.Println("# found a valid move")
				} else {
					fmt.Println("# not found a move on this line")
					break next
				}

				result = hclibs.MakeUserMove(move, &p)
				if result != "" {
					fmt.Println(result)
				}

				// make computer go if not in force mode
				fmt.Println("# gameforce=", hclibs.GameForce)
				if hclibs.GameForce == false {
					xboardGo(&p)
				}
				break next

				// Matches a2a3 type move
			case re.MatchString(input):
				move, err = hclibs.ParseUserMove(input, &p)
				if err != "" {
					fmt.Println(err)
					break next
				}
				result = hclibs.MakeUserMove(move, &p)
				if result != "" {
					fmt.Println(result)
				}

				// make computer go if not in force mode
				fmt.Println("# gameforce=", hclibs.GameForce)
				if hclibs.GameForce == false {
					xboardGo(&p)
				}
				break next

			case strings.Contains(input, "go"): // || hclibs.GameForce == true:
				hclibs.GameForce = false
				xboardGo(&p)
				break next

			case strings.Contains(input, "force"): // || hclibs.GameForce == true:
				hclibs.GameForce = true
				break next

			case strings.Contains(input, "new"):
				p = hclibs.FENToNewBoard(hclibs.STARTFEN)
				hclibs.GameOver = false
				break next

			case strings.Contains(input, "ping"):
				fields := strings.Fields(input)
				fmt.Print("pong")
				if len(fields) > 1 {
					fmt.Print(" " + fields[1])
				}
				fmt.Print("\n")
				break next

			case strings.HasPrefix(input, "sd"), strings.HasPrefix(input, "depth"): // not in spec but knights seems to send depth anyway.
				//  case strings.Contains(input,"fen") || strings.Contains(input,"setboard"):
				fields := strings.Fields(input)
				if len(fields) == 1 {
					fmt.Printf("# Current depth is %d.\n", hclibs.GameDepthSearch)
					break next
				}
				d, err := strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println("# Please specify a number")
					break next
				}
				if d > hclibs.MAXSEARCHDEPTHX {
					d = hclibs.MAXSEARCHDEPTHX
				}
				fmt.Printf("# Setting depth to %d.\n", d)
				hclibs.GameDepthSearch = d

			case strings.HasPrefix(input, "draw"):
				fmt.Println("offer draw")
				hclibs.GameOver = true
				break next

			case strings.HasPrefix(input, "setboard"):
				fields := strings.Fields(input)
				if len(fields) == 1 {
					fmt.Println("Error (no Fen): " + input)
					break next
				}
				fen := strings.Join(fields[1:], " ")
				fmt.Println("# parsing fen [" + fen + "]")
				p = hclibs.FENToNewBoard(fen)
				// 				fmt.Println("Error (Not implemented yet!!!): " + input)
				break next

			case strings.HasPrefix(input, "white"):
				p.Side = hclibs.WHITE
				break next

			case strings.HasPrefix(input, "black"):
				p.Side = hclibs.BLACK
				break next

			case strings.HasPrefix(input, "post"):
				hclibs.GamePostStats = true
				break next

			case strings.HasPrefix(input, "nopost"):
				hclibs.GamePostStats = false
				break next

			// no ops
			case strings.HasPrefix(input, "random"), strings.HasPrefix(input, "level"), strings.HasPrefix(input, "easy"), strings.HasPrefix(input, "hard"), strings.HasPrefix(input, "accepted"):
				break next

			// currently no ops - TODO
			// undo
			// 			case strings.HasPrefix(input, "time"), strings.HasPrefix(input, "otim"):
			// 				break next

			default:
				fmt.Printf("Error (unknown command): %v\n", input)

			}
		}
	}
}

func xboardGo(p *hclibs.Pos) {
	/*if hclibs.Control == nil || hclibs.StopSearch() == true {

		hclibs.Control = make(chan string)
	}*/
	res, info := hclibs.Go(p)
	// 	if hclibs.StopSearch() == false {
	// 		close(hclibs.Control)
	// 	}
	if len(info) > 0 {
		fmt.Println(info)
	}
	fmt.Println(res)
}

/***************************************************************
 * Use the UCI chess protocol
 *
 * lots of fun
 * See: http://wbec-ridderkerk.nl/html/UCIProtocol.html
 ****************************************************************/

// little helper function ("hack") needed to switch from console to uci if command line switches are ignored
func ucihelper() {
	name := fmt.Sprintf("HastyChess v%v", hclibs.VERSION)
	fmt.Println("id name " + name + "\nid author Kevin Colyer 2016")
	// Send options to GUI here...
	//
	//
	fmt.Println("debug on")
	fmt.Println("uciok")
}

func mainIcs(scanner *bufio.Scanner) {
	// 	version := 0.99
	//	name := fmt.Sprintf("HastyChess v%v", hclibs.VERSION)
	// 	scanner := bufio.NewScanner(os.Stdin)

	p := hclibs.FENToNewBoard(hclibs.STARTFEN)

	//hclibs.GameUseBook = false // UCI gui does book - unless option below...

	hclibs.GameOver = false
	hclibs.GameDisplayOn = false
	hclibs.GameDepthSearch = 8
	hclibs.GameForce = false
	// main input loop
	for {
		for scanner.Scan() {

			//input := strings.ToLower(strings.TrimSpace(scanner.Text()))
			input := strings.TrimSpace(scanner.Text())
			fmt.Printf("info string echo server sent (%v)\n", input)
			// 			time.Sleep(time.Second)
			//fmt.Pri(ntf("You said [%v]\n", input)
			switch {

			case input == "uci":
				//         fmt.Println("id name " + name + "\nid author Kevin Colyer 2016")
				//         Send options to GUI here...
				//
				//
				//         fmt.Println("uciok")
				ucihelper()
				continue

			case strings.HasPrefix(input, "isready"):
				fmt.Println("readyok")

			case strings.HasPrefix(input, "quit"):
				os.Exit(0)

			case strings.HasPrefix(input, "ucinewgame"):
				p = hclibs.FENToNewBoard(hclibs.STARTFEN)
				hclibs.GameOver = false

			case strings.HasPrefix(input, "go"):
				go uciGo(input, &p)
				//  lots of sub verbs here...
				continue
				// start searching now!

			case input == "ping":
				fmt.Println("pong")
				continue
				// start searching now!
			case strings.HasPrefix(input, "stop"):
				if hclibs.Control != nil {
					close(hclibs.Control)
				} // tells searches to finish
				continue
				// stop searching now!

			case strings.HasPrefix(input, "position"):
				uciPosition(input, &p)
				// position fen moves...
				// position startpos moves...
				continue

			}
		}
	}

}

func uciPosition(input string, p *hclibs.Pos) {
	fen, moves := ParsePositionInput(input)
	*p = hclibs.FENToNewBoard(fen)

	for _, m := range moves {
		/*if !hclibs.IsValidMove(m, p) {
					fmt.Printf("info currmove %v (Sent bad move by server %v in %v)\n", m, m, moves)
		        }*/
		hclibs.MakeMove(m, p) // note no error checking here!!!
	}
}

// parses a string "position <fen> <move>*" or "position startpos <move>*" returns
// fen as string and slice of Move types.
func ParsePositionInput(input string) (fen string, moves []hclibs.Move) {

	re := regexp.MustCompile("[a-h][1-8][a-h][1-8][qbnr]?")

	ef := 2 // end fen = start of moves (if any)
	f := strings.Split(input, " ")
	// can ignore first field == "position"
	if len(f) > 2 && f[1] == "fen" {
		fen = strings.Join(f[2:], " ")
		return
	}
	if f[1] == "startpos" {
		fen = hclibs.STARTFEN
	} else {
		// assume command is "moves"
		//look ahead to moves - what is skipped is a fen
		for ; ef < len(f); ef++ {
			if f[ef] == "moves" {
				break
			}
		}

		fen = strings.Join(f[1:ef], " ")

	}

	for ; ef < len(f); ef++ {
		if re.MatchString(f[ef]) {
			moves = append(moves, hclibs.AlgToMove(f[ef]))
		}

	}
	return
}

func uciGo(input string, p *hclibs.Pos) {
	// position fen moves...
	// position startpos moves...

	if hclibs.Control != nil && hclibs.StopSearch() == false { // if channel still open, close it
		close(hclibs.Control)
	}
	hclibs.Control = make(chan string)
	go func() {
		for m := range hclibs.Control {
			fmt.Println(m)
		} // send messages from search straight to console until channel closed
	}()
	// channel shut in "stop" command or at end of a search or at a time out...

	// expand imput string to parse sub verbs
	f := strings.Split(input, " ")
	if len(f) > 1 {
		f = f[1:]
	} // shift out "go"
	// chomp settings two at a time
	count := 0
	for len(f) > 1 {
		// double word commands
		count = 2
		// DEPTH
		if f[0] == "depth" && len(f) > 1 {
			d, err := strconv.Atoi(f[1])
			if err != nil {
				fmt.Println("info string Setting depth: Please specify a number")
				d = 0
			} else {
				hclibs.GameDepthSearch = d
				fmt.Println("info string  Set depth to ", d)
			}

		}
		// MOVETIME
		if f[0] == "movetime" && len(f) > 1 {
			d, err := strconv.Atoi(f[1])
			if err != nil {
				fmt.Println("info string Setting movetime: Please specify a number")
				d = 0
			} else {

				fmt.Println("info string Set movetime to ", d)
				hclibs.GameDurationToSearch = time.Duration(d * 1000 * 1000) // milliseconds to nanoseconds
			}
		}
		// Single verb commands
		// INFINITE (single)
		// PONDER etc.
		// 		fmt.Println("# comsuming tokens:",count)
		f = f[count:] // consume tokens
	}

	res, _ := hclibs.Go(p)
	fmt.Println(res)
	return
}
