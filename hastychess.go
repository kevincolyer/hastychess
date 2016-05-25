// Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/kevincolyer/hastychess/hclibs"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var flagXboard = flag.Bool("xboard", false, "Select xboard mode")
	var flagIcs = flag.Bool("ics", false, "Select ics server mode")
	var flagConsole = flag.Bool("console", true, "Select console mode")
	var flagStats = flag.Bool("stats", true, "Enable printing of statistics")
	var flagUseBook = flag.Bool("book", true, "Enable the use of built in book moves")
	var flagUseTt = flag.Bool("tt", true, "Enable the use of Transposition Tables")
	flag.Parse()

	hclibs.GameUseTt = *(flagUseTt)
	hclibs.GameUseStats = *(flagStats)
	hclibs.GameUseBook = *(flagUseBook)

	hclibs.GameInit()

	switch {
	case *(flagXboard):
		hclibs.GameProtocol = hclibs.PROTOXBOARD
		mainXboard()
	case *(flagIcs):
		hclibs.GameProtocol = hclibs.PROTOUCI
		mainIcs()
	case *(flagConsole):
		hclibs.GameProtocol = hclibs.PROTOCONSOLE
		mainConsole()
	}
	fmt.Println("Bye and thanks for playing!")

}

func mainConsole() {

	var err string
	var move hclibs.Move

	re, e := regexp.Compile("[a-h][1-8][a-h][1-8][qbnr]?")
	if e != nil {
		panic("Regexp did not compile!")
	}
	version := 0.99
	fmt.Printf("Hello and welcome to HastyChess version %v\n\n", version)

	scanner := bufio.NewScanner(os.Stdin)
	p := hclibs.FENToNewBoard(hclibs.STARTFEN)
	hclibs.GameOver = false
	hclibs.GameDisplayOn = true
	hclibs.GameDepthSearch = 4
	hclibs.GameForce = false
	if hclibs.GameDisplayOn {
		fmt.Println(hclibs.BoardToStrWide(&p))
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
				err = hclibs.MakeUserMove(move, &p)
				fmt.Println(hclibs.BoardToStrWide(&p))
				fmt.Println(err)

			case re.MatchString(input):
				move, err = hclibs.ParseUserMove(input, &p)
				if err != "" {
					fmt.Println(err)
					break next
				}
				err = hclibs.MakeUserMove(move, &p)
				fmt.Println(hclibs.BoardToStrWide(&p))
				fmt.Println(err)

			case strings.Contains(input, "new"):
				p = hclibs.FENToNewBoard(hclibs.STARTFEN)
				hclibs.GameOver = false
				fmt.Println(hclibs.BoardToStrWide(&p))

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
				hclibs.Divide(d, p)

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
				nodes := hclibs.Perft(d, p)
				fmt.Printf("\nPerft to depth %v gives %v nodes\n", d, nodes)

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
	var err string
	var move hclibs.Move

	re, e := regexp.Compile("[a-h][1-8][a-h][1-8][qbnr]?")
	if e != nil {
		panic("Regexp did not compile!")
	}
	version := 0.99
	name := fmt.Sprintf("HastyChess v%v", version)
	fmt.Printf("Hello and welcome to %v\n\n", name)

	fmt.Println("feature debug=1")

	scanner := bufio.NewScanner(os.Stdin)
	p := hclibs.FENToNewBoard(hclibs.STARTFEN)
	hclibs.GameOver = false
	hclibs.GameDisplayOn = false
	hclibs.GameDepthSearch = 4
	hclibs.GameForce = false
	if hclibs.GameDisplayOn {
		fmt.Println(hclibs.BoardToStrWide(&p))
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

			case strings.Contains(input, "move"):
				fields := strings.Fields(input)
				if len(fields) > 1 {
					move, err = hclibs.ParseUserMove(fields[1], &p)
					if err != "" {
						fmt.Println(err)
						break next
					}
				}
				err = hclibs.MakeUserMove(move, &p)
				//fmt.Println(hclibs.BoardToStrWide(&p))
				fmt.Println(err)

				// make computer go
				xboardGo(&p)

			case re.MatchString(input):
				move, err = hclibs.ParseUserMove(input, &p)
				if err != "" {
					fmt.Println(err)
					break next
				}
				err = hclibs.MakeUserMove(move, &p)
				//fmt.Println(hclibs.BoardToStrWide(&p))
				fmt.Println(err)

				// make computer go
				xboardGo(&p)

			case strings.Contains(input, "go") || hclibs.GameForce == true:
				xboardGo(&p)

			case strings.Contains(input, "new"):
				p = hclibs.FENToNewBoard(hclibs.STARTFEN)
				hclibs.GameOver = false
				/*
				   when "new" {
				   init_game;
				   $p=Position.new;
				   $time=1;
				   $otim=1;
				   $forced=False;
				   @stack=();
				*/

			case strings.Contains(input, "ping"):
				fmt.Println("pong")

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
			/*
			   when /$ setboard/ {
			       $_ ~~ /$ setboard \s+ (.*)/;
			       init_game;
			       $p=Position.new(FEN => $0);
			   }
			   when "force" {
			       $forced=True;                                            }*/

			case strings.HasPrefix(input, "white"):
				p.Side = hclibs.WHITE
				break next

			case strings.HasPrefix(input, "black"):
				p.Side = hclibs.BLACK
				break next
				/*
				   when /time/ { ... }
				   when /otim/ { ... }
				   when /post|random|hard|accepted|level/ {say "skip"; next; }
				*/
			default:
				fmt.Printf("# Error (unknown command): %v\n", input)

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

	version := 0.99
	name := fmt.Sprintf("HastyChess v%v", version)
	scanner := bufio.NewScanner(os.Stdin)
	p := hclibs.FENToNewBoard(hclibs.STARTFEN)

	hclibs.GameUseBook = false // UCI gui does book - unless option below...

	hclibs.GameOver = false
	hclibs.GameDisplayOn = false
	hclibs.GameDepthSearch = 4
	hclibs.GameForce = false
	/*if hclibs.GameDisplayOn {
			fmt.Println(hclibs.BoardToStrWide(&p))
	        }*/
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
