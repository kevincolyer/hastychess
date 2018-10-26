// Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package protocol

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/kevincolyer/hastychess/engine"
	"github.com/kevincolyer/hastychess/hclibs"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type console struct {
	In      io.Reader
	Out     io.Writer
	Stderr  io.Writer
	Options CLIOptions
}

func (p *console) o(s string) {
	io.WriteString(p.Out, "o "+s)
}

func (p *console) oln(s string) {
	io.WriteString(p.Out, "o "+s+"\n")
}

func (p *console) ofln(s interface{}) {
	io.WriteString(p.Out, fmt.Sprintf("o %v\n", s))
}

func (p *console) debug(s string) {
	io.WriteString(p.Stderr, "Stderr: "+s)
}

func NewConsole(in io.Reader, out io.Writer, stderr io.Writer, options CLIOptions) (*console, error) {
	p := console{Out: out, In: in, Stderr: stderr, Options: options}
	_, e := io.WriteString(p.Stderr, "Stderr: Started "+options.NameVersion+"\n")
	return &p, e
}

func (p *console) Echo() (e error) {
	b := make([]byte, 256)
	if n, err := p.In.Read(b); err != io.EOF {
		io.WriteString(p.Out, string(b[:n]))
		io.WriteString(p.Stderr, "Stderr: "+string(b[:n]))
	} else {
		e = fmt.Errorf("EOF for input: none sent?")
	}
	return
}

func (p *console) Start() error {
	myEngine, err := engine.New(engine.EngineOptions{})
	if err != nil {
		return fmt.Errorf("Error creating engine: %v", err)
	}
	p.debug("Engine started ok\n")
	p.o("hello world\n")
	p.MainLoop(myEngine)
	// 	if myEngine.Stop() {
	// 		io.WriteString(p.Stderr, "Stderr: Engine stopped ok\n")
	// 	}

	return nil
}

//===========================================================================

func (proto *console) MainLoop(myEngine *engine.Engine) {
	scanner := bufio.NewScanner(proto.In)
	var err string
	var result string
	var move hclibs.Move

	re := regexp.MustCompile("[a-h][1-8][a-h][1-8][qbnr]?")
	// 	version := 1.0
	hiwhite := color.New(color.FgHiWhite).SprintfFunc()
	proto.o(hiwhite("Hello and welcome to %v\n\n", proto.Options.NameVersion))

	p := hclibs.FENToNewBoard(hclibs.STARTFEN)
	hclibs.GameProtocol = hclibs.PROTOCONSOLE
	hclibs.GameOver = false
	hclibs.GameDisplayOn = true
	hclibs.GameDepthSearch = hclibs.MAXSEARCHDEPTH // 8 or 4 // don't delete this or search depth = 0!
	hclibs.GameForce = false
	if hclibs.GameDisplayOn {
		proto.ofln(&p)
	}
	quit := false
	proto.o("> ")
	hclibs.Control = make(chan string)
	// main input loop
QUIT:
	for quit == false {
	next:
		for scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())
			switch {
			case input == "quit" || input == "q":
				quit = true
				break QUIT

			// provide a way to change to xboard mode if user forgets to use commandline switch
			// 			case strings.HasPrefix(input, "xboard"):
			// 				color.NoColor = true
			// 				hclibs.GameProtocol = hclibs.PROTOXBOARD
			// 				mainXboard(scanner)
			// 				return
			// 			// provide a way to change to xboard mode if user forgets to use commandline switch
			// 			case strings.HasPrefix(input, "uci"):
			// 				color.NoColor = true
			// 				hclibs.GameProtocol = hclibs.PROTOUCI
			// 				ucihelper()
			// 				mainIcs(scanner)
			// 				return

			case strings.Contains(input, "move"):
				fields := strings.Fields(input)
				if len(fields) > 1 {
					move, err = hclibs.ParseUserMove(fields[1], &p)
					if err != "" {
						proto.oln(err)
						break next
					}
				}
				result = hclibs.MakeUserMove(move, &p)
				proto.ofln(&p)
				proto.oln(result)

			case re.MatchString(input):
				move, err = hclibs.ParseUserMove(input, &p)
				if err != "" {
					proto.ofln(err)
					break next
				}
				result = hclibs.MakeUserMove(move, &p)
				proto.ofln(&p)
				proto.oln(result)

			case strings.Contains(input, "new"):
				p = hclibs.FENToNewBoard(hclibs.STARTFEN)
				hclibs.GameOver = false
				proto.ofln(&p)

			// special commands that allow testing of certain positions
			case strings.Contains(input, "kiwipete"):
				p = hclibs.FENToNewBoard("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
				hclibs.GameOver = false
				proto.ofln(&p)

			case strings.Contains(input, "pos4"):
				p = hclibs.FENToNewBoard("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
				hclibs.GameOver = false
				proto.ofln(&p)

			case strings.Contains(input, "pos5"):
				p = hclibs.FENToNewBoard("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
				hclibs.GameOver = false
				proto.ofln(&p)

			case strings.Contains(input, "end1"):
				p = hclibs.FENToNewBoard("8/8/4k3/7p/8/8/2K5/R6Q w - - 0 1")
				hclibs.GameOver = false
				proto.ofln(&p)

			case strings.Contains(input, "end2"):
				p = hclibs.FENToNewBoard("8/8/8/1bn5/8/2k5/8/2K5 w - - 0 1")
				hclibs.GameOver = false
				proto.ofln(&p)

			case strings.Contains(input, "end3"):
				p = hclibs.FENToNewBoard("8/8/8/1k6/8/7Q/3R4/2K5 w - - 0 1")
				hclibs.GameOver = false
				proto.ofln(&p)

			case strings.Contains(input, "end4"):
				p = hclibs.FENToNewBoard("8/8/8/k7/8/7Q/1R6/2K5 w - - 0 1")
				hclibs.GameOver = false
				proto.ofln(&p)

				// normal commands
			case strings.HasPrefix(input, "setboard"):
				fields := strings.Fields(input)
				if len(fields) == 1 {
					proto.oln("Error (no Fen): " + input)
					break next
				}
				fen := strings.Join(fields[1:], " ")
				proto.oln("# parsing fen [" + fen + "]")
				p = hclibs.FENToNewBoard(fen)
				// 				proto.oln("Error (Not implemented yet!!!): " + input)
				break next

			case strings.Contains(input, "ping"):
				proto.oln("pong")

			case strings.Contains(input, "divide"):
				fields := strings.Fields(input)
				if len(fields) == 1 {
					proto.oln("Please specify a number to divide down to")
					break next
				}
				d, err := strconv.Atoi(fields[1])
				if err != nil {
					proto.oln("Please specify a number")
					break next
				}
				_, s := hclibs.Divide(d, &p)
				proto.ofln(s)

			case strings.Contains(input, "perft"):
				fields := strings.Fields(input)
				if len(fields) == 1 {
					proto.oln("Please specify a number to perft down to")
					break next
				}
				d, err := strconv.Atoi(fields[1])
				if err != nil {
					proto.oln("Please specify a number")
					break next
				}
				start := time.Now()
				nodes := hclibs.Perft(d, &p)
				elapsed := time.Since(start)
				proto.o(fmt.Sprintf("Perft to depth %v gives %v nodes ", d, hclibs.Comma(nodes)))
				proto.o(fmt.Sprintf("(nps: %v)\n", hclibs.Comma(int(float64(nodes)/elapsed.Seconds()))))

			case strings.Contains(input, "depth"):
				//  case strings.Contains(input,"fen") || strings.Contains(input,"setboard"):
				fields := strings.Fields(input)
				if len(fields) == 1 {
					proto.o(fmt.Sprintf("Current depth is %d.\n", hclibs.GameDepthSearch))
					break next
				}
				d, err := strconv.Atoi(fields[1])
				if err != nil {
					proto.oln("Please specify a number")
					break next
				}
				if d > hclibs.MAXSEARCHDEPTH {
					d = hclibs.MAXSEARCHDEPTH
				}
				proto.o(fmt.Sprintf("Current depth is %d. Setting depth to %d.\n", hclibs.GameDepthSearch, d))
				hclibs.GameDepthSearch = d

			case input == "help":
				proto.o("Commands: move [a2a4],[a2a4], g[o], auto, quit, new, ping, depth #, perft #, divide #,\n          setboard [fen], kiwipete, pos4, pos5, end1, end2, end3, end4\n")

			case strings.Contains(input, "auto"):
				hclibs.GameForce = !hclibs.GameForce

				////////////////////////////////////////////////////////////////////////////////
			case strings.Contains(input, "go") || input == "g" || hclibs.GameForce == true:
				res, info, _ := hclibs.Go(&p)
				proto.ofln(&p)
				proto.oln(info)
				proto.oln(res)

			}

			proto.o("> ")
		}
	}
}
