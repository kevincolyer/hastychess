package main

import (
	"bufio"
	"flag"
	"fmt"
	"hastychess/hclibs"
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
		mainXboard()
	case *(flagIcs):
		mainIcs()
	case *(flagConsole):
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
				err = hclibs.Go(&p)
				fmt.Println(&p)
				fmt.Println(err)

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
				err = hclibs.Go(&p)
				fmt.Println(err)

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
				err = hclibs.Go(&p)
				fmt.Println(err)

			case strings.Contains(input, "go") || hclibs.GameForce == true:
				err = hclibs.Go(&p)
				//fmt.Println(&p)
				fmt.Println(err)

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

func mainIcs() {
	fmt.Println("ICS mode not yet implimented")
	return
}
