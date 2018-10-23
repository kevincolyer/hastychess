// Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package protocol

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/kevincolyer/hastychess/hclibs"
        	"github.com/fatih/color"
)



type console struct {
	In      io.Reader
	Out     io.Writer
	Options CLIOptions
}

func NewConsole(in io.Reader, out io.Writer, options CLIOptions) (*console, error) {
	p := console{Out: out, In: in, Options: options}
	_, e := io.WriteString(p.Out, "Welcome to "+options.NameVersion+"\n")
	return &p, e
}

func (p *console) Echo() (e error) {
	b := make([]byte, 256)
	if n, err := p.In.Read(b); err != io.EOF {
		io.WriteString(p.Out, string(b[:n]))
	} else {
		e = fmt.Errorf("EOF for input: none sent?")
	}
	return
}

func (p *console) Start() (e error) {
	io.WriteString(p.Out, "started ok\n")
	e = fmt.Errorf("Done!")
	return
}

//===========================================================================


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
	hclibs.GameDepthSearch = hclibs.MAXSEARCHDEPTH // 8 or 4 // don't delete this or search depth = 0!
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

