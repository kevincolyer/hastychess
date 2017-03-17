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

	switch {
	case *(flagXboard):
		hclibs.GameProtocol = hclibs.PROTOXBOARD
		color.NoColor = true
		mainXboard()
	case *(flagIcs):
		hclibs.GameProtocol = hclibs.PROTOUCI
		color.NoColor = true
		mainIcs()
	case *(flagUci):
		hclibs.GameProtocol = hclibs.PROTOUCI
		color.NoColor = true
		mainIcs()
	case *(flagConsole):
		hclibs.GameProtocol = hclibs.PROTOCONSOLE
		mainConsole()
	}
	fmt.Println("Bye and thanks for playing!")

}

func mainConsole() {

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

	scanner := bufio.NewScanner(os.Stdin)
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

			input := strings.ToLower(strings.TrimSpace(scanner.Text()))
			//fmt.Printf("You said [%v]\n", input)
			switch {
			case strings.Contains(input, "quit"):
				quit = true
				break QUIT
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

func mainXboard() {
	// see https://www.gnu.org/software/xboard/engine-intf.html and XXX for protocol info
	var err string
	var result string
	var move hclibs.Move

	re, e := regexp.Compile("[a-h][1-8][a-h][1-8][qbnr]?")
	if e != nil {
		panic("Regexp did not compile!")
	}
	// 	version := 0.99
	name := fmt.Sprintf("HastyChess v%v", hclibs.VERSION)
	fmt.Printf("tellics Hello and welcome to %v\n\n", name)

	fmt.Println("feature debug=1")

	scanner := bufio.NewScanner(os.Stdin)
	p := hclibs.FENToNewBoard(hclibs.STARTFEN)
	hclibs.GameOver = false
	hclibs.GameDisplayOn = false
	hclibs.GameDepthSearch = 4
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

			input := strings.ToLower(strings.TrimSpace(scanner.Text()))
			//fmt.Pri(ntf("You said [%v]\n", input)
			switch {

			case strings.HasPrefix(input, "accepted"):
				break next

			case strings.HasPrefix(input, "#"):
				fmt.Println()
				break next

				// ignore when xboard tells me I have lost
			case strings.HasPrefix(input, "result"):
				fmt.Println()
				break next

			case strings.Contains(input, "xboard"):
				fmt.Println()

			case strings.Contains(input, "protover 2"):
				fmt.Println("feature done=0")
				fmt.Printf("feature myname=\"%v\"\n", name)
				fmt.Println("feature usermove=1")
				fmt.Println("feature setboard=1")
				fmt.Println("feature ping=1")
				fmt.Println("feature sigint=0")
				fmt.Println("feature variants=\"normal\"")
				fmt.Println("feature debug=1") // allows comments starting with hash symbols
				fmt.Println("feature done=1")

			case strings.Contains(input, "quit"):
				quit = true
				break QUIT

			case strings.Contains(input, "move"), strings.Contains(input, "move"):
				fields := strings.Fields(input)
				if len(fields) > 1 {
					move, err = hclibs.ParseUserMove(fields[1], &p)
					if err != "" {
						fmt.Println(err)
						break next
					}
				}
				result = hclibs.MakeUserMove(move, &p)
				fmt.Println(result)

				// make computer go if not in force mode
				if hclibs.GameForce == false {
					xboardGo(&p)
				}

				// Matches a2a3 type move
			case re.MatchString(input):
				move, err = hclibs.ParseUserMove(input, &p)
				if err != "" {
					fmt.Println(err)
					break next
				}
				result = hclibs.MakeUserMove(move, &p)
				fmt.Println(result)

				// make computer go if not in force mode
				if hclibs.GameForce == false {
					xboardGo(&p)
				}

			case strings.Contains(input, "go"): // || hclibs.GameForce == true:
				hclibs.GameForce = false
				xboardGo(&p)

			case strings.Contains(input, "force"): // || hclibs.GameForce == true:
				hclibs.GameForce = true

			case strings.Contains(input, "new"):
				p = hclibs.FENToNewBoard(hclibs.STARTFEN)
				hclibs.GameOver = false

			case strings.Contains(input, "ping"):
				fields := strings.Fields(input)
				fmt.Print("pong")
				if len(fields) > 1 {
					fmt.Print(" " + fields[1])
				}
				fmt.Print("\n")

			case strings.Contains(input, "depth"):
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
				fmt.Printf("# Current depth is %d. Setting depth to %d.\n", hclibs.GameDepthSearch, d)
				hclibs.GameDepthSearch = d

			case strings.HasPrefix(input, "draw"):
				fmt.Println("offer draw")
				hclibs.GameOver = true
				break next

			case strings.HasPrefix(input, "setboard"):
				fmt.Println("Error (Not implemented yet!!!): " + input)
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
			case strings.HasPrefix(input, "random"), strings.HasPrefix(input, "level"), strings.HasPrefix(input, "hard"), strings.HasPrefix(input, "accepted"):
				break next

			// currently no ops - TODO
			case strings.HasPrefix(input, "time"), strings.HasPrefix(input, "otim"):
				break next

			default:
				fmt.Printf("Error (unknown command): %v\n", input)

			}
		}
	}
}

func xboardGo(p *hclibs.Pos) {
	if hclibs.Control == nil || hclibs.StopSearch() == true {

		hclibs.Control = make(chan string)
	}
	res, info := hclibs.Go(p)
	if hclibs.StopSearch() == false {
		close(hclibs.Control)
	}
	fmt.Println(info)
	fmt.Println(res)
}

/***************************************************************
 * Use the UCI chess protocol
 *
 * lots of fun
 * See: http://wbec-ridderkerk.nl/html/UCIProtocol.html
 *
 ***************************************************************/

func mainIcs() {

	// 	version := 0.99
	name := fmt.Sprintf("HastyChess v%v", hclibs.VERSION)
	scanner := bufio.NewScanner(os.Stdin)
	p := hclibs.FENToNewBoard(hclibs.STARTFEN)

	hclibs.GameUseBook = false // UCI gui does book - unless option below...

	hclibs.GameOver = false
	hclibs.GameDisplayOn = false
	hclibs.GameDepthSearch = 4
	hclibs.GameForce = false
	// main input loop
	for {
		for scanner.Scan() {

			input := strings.ToLower(strings.TrimSpace(scanner.Text()))
			fmt.Printf("# echo server sent (%v)\n", input)
			// 			time.Sleep(time.Second)
			//fmt.Pri(ntf("You said [%v]\n", input)
			switch {

			case input == "uci":
				fmt.Println("id name " + name + "\nid author Kevin Colyer 2016")
				// Send options to GUI here...
				//
				//
				fmt.Println("uciok")
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

	if f[1] == "startpos" {
		fen = hclibs.STARTFEN
	} else { //look ahead to moves - what is skipped is a fen
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

	// 	res,info:=hclibs.Go(p)
	res, _ := hclibs.Go(p)
	// we need to massage the bestmove
	// not sure what to do with all our stats info. Could try to send it too.
	// 	fmt.Println(info+"best"+res)
	fmt.Println(res)
	return
}
