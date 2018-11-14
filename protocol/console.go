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

type EngineInfo struct {
	*hclibs.Statistics
	*hclibs.PV
}

func (p *console) o(s string) {
	io.WriteString(p.Out, s)
}

func (p *console) oln(s string) {
	io.WriteString(p.Out, s+"\n")
}

func (p *console) ofln(s interface{}) {
	io.WriteString(p.Out, fmt.Sprintf("%v\n", s))
}

func (p *console) debug(s string) {
	io.WriteString(p.Stderr, "Stderr: "+s)
}

func NewConsole(in io.Reader, out io.Writer, stderr io.Writer, options CLIOptions) (*console, error) {
	p := console{Out: out, In: in, Stderr: stderr, Options: options}
	p.debug("Started Console Protocol for " + options.NameVersion + "\n")
	return &p, nil
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
	p.MainLoop(myEngine)
	// 	if myEngine.Stop() {
	// 		io.WriteString(p.Stderr, "Stderr: Engine stopped ok\n")
	// 	}

	return nil
}

//===========================================================================

type tui struct {
	Cls     string
	Title   string
	Result  string
	Board   string
	Status  string
	Pv      string
	Stats   string
	History string
	Cmdline string
	Info    string

	spinner     string
	ShowSpinner bool
	Update      bool
	UpdateOnce  bool
	proto       *console
}

// title        status spinner
// board        pv
//              stats
// result       history
// cmdline

func splitIt(s string) []string {
	return strings.Split(s, "\n")
}

func splitItN(s string) []string {
	return splitIt(s + "\n")
}

func pad(s string, size int) string {
	reps := size - len(s)
	if reps <= 0 {
		return s
	}
	return s + (strings.Repeat(" ", reps))
}

func (proto *console) MainLoop(myEngine *engine.Engine) {
	hiwhite := color.New(color.FgHiWhite).SprintfFunc()
	ui := tui{Title: hiwhite("Hello and welcome to %v\n\n", proto.Options.NameVersion), Cls: "\033[H\033[2J",
		Pv: "pv", Stats: "stats", History: "History", Info: "Info", Cmdline: "> ", Update: true,
	}
	p := hclibs.FENToNewBoard(hclibs.STARTFEN)
	ui.Board = hclibs.BoardToStrColour(&p)
	ui.Status = "awaiting user input"
	ui.proto = proto
	uictl := make(chan bool)
	// closure over ui to provide the new sparkly text ui system!!! woot!
	go func(done chan bool) {
		var lc []string
		var rc []string
		spinners := [4]string{"|", "/", "-", "\\"}
		spincounter := 0
		ui.spinner = spinners[0]
		for {

			select {
			case <-done:
				return
			default:
			}
			if ui.Update || ui.UpdateOnce {
				ui.UpdateOnce = false

				ui.proto.o(ui.Cls)
				ui.proto.o(ui.Title)
				lc = nil
				lc = append(lc, splitIt(ui.Board)...)
				lc = append(lc, splitItN("Response: "+ui.Result)...)
				// 			lc = append(lc, splitItN(ui.Cmdline)...)
				// 			lc = append(lc, splitItN(ui.Info)...)

				rc = nil
				sp := "Status: " + ui.Status
				if ui.ShowSpinner {
					sp += " " + ui.spinner
				}
				rc = append(rc, splitItN(sp)...)
				rc = append(rc, splitItN("PV: "+ui.Pv)...)
				rc = append(rc, splitItN("Stats: "+ui.Stats)...)
				rc = append(rc, splitItN("History: "+ui.History)...)

				lencols := hclibs.Max(len(rc), len(lc))
				for i := 0; i < lencols; i++ {
					l := ""
					r := ""

					if i < len(lc) {
						l = lc[i]
					}
					if i < len(rc) {
						r = rc[i]
					}
					// magic numbers from guesswork
					ui.proto.oln(pad(l, 31) + " | " + pad(r, 60))
				}
				if ui.Info != "" {
					ui.proto.oln("Info: " + ui.Info)
				}
				ui.proto.o(ui.Cmdline)
			}

			time.Sleep(time.Millisecond * 100)
			ui.spinner = spinners[spincounter]
			spincounter += 1
			spincounter %= 4

		}
	}(uictl)

	scanner := bufio.NewScanner(proto.In)
	var err string
	var result string
	var move hclibs.Move

	re := regexp.MustCompile("[a-h][1-8][a-h][1-8][qbnr]?")
	// 	proto.o(ui.cls)
	quit := false

	// 	hclibs.Control = make(chan string)
	engineInfo := make(chan hclibs.EngineInfo)
	go func(ei chan hclibs.EngineInfo) {
		//             var data EngineInfo
		i := 0
		for {
			select {
			case data := <-ei:
				ui.Pv = fmt.Sprintf("%v", data.Pv)
				ui.Stats = fmt.Sprintf("%v", data.Stats)
				ui.History = fmt.Sprintf("%v", i)
				i++
			default:
			}
		}
	}(engineInfo)
	// main input loop

	// QUIT:
	for quit == false {
		// 	net:
		//                 time.Sleep(time.Millisecond * 100)
		ui.Update = false
		ui.UpdateOnce = true
		// 		for scanner.Scan() {
		scanner.Scan()
		quit = (scanner.Err() == io.EOF)
		//                 }

		input := strings.TrimSpace(scanner.Text())
		ui.Update = true
		ui.Info = ""
		ui.Result = ""
		switch {
		case input == "quit" || input == "q":
			quit = true
			// 				break QUIT

		case strings.Contains(input, "move"):
			fields := strings.Fields(input)
			if len(fields) > 1 {
				move, err = hclibs.ParseUserMove(fields[1], &p)
				if err != "" {
					ui.Result = err
					break //next
				}
			}
			result = hclibs.MakeUserMove(move, &p)
			ui.Board = hclibs.BoardToStrColour(&p)
			ui.Result = result

		case re.MatchString(input):
			move, err = hclibs.ParseUserMove(input, &p)
			if err != "" {
				ui.Result = err
				break //next
			}
			result = hclibs.MakeUserMove(move, &p)
			ui.Board = hclibs.BoardToStrColour(&p)
			ui.Result = result

		case strings.Contains(input, "new"):
			p = hclibs.FENToNewBoard(hclibs.STARTFEN)
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&p)

		// special commands that allow testing of certain positions
		case strings.Contains(input, "kiwipete"):
			p = hclibs.FENToNewBoard("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&p)

		case strings.Contains(input, "pos4"):
			p = hclibs.FENToNewBoard("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&p)

		case strings.Contains(input, "pos5"):
			p = hclibs.FENToNewBoard("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&p)

		case strings.Contains(input, "end1"):
			p = hclibs.FENToNewBoard("8/8/4k3/7p/8/8/2K5/R6Q w - - 0 1")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&p)

		case strings.Contains(input, "end2"):
			p = hclibs.FENToNewBoard("8/8/8/1bn5/8/2k5/8/2K5 w - - 0 1")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&p)

		case strings.Contains(input, "end3"):
			p = hclibs.FENToNewBoard("8/8/8/1k6/8/7Q/3R4/2K5 w - - 0 1")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&p)

		case strings.Contains(input, "end4"):
			p = hclibs.FENToNewBoard("8/8/8/k7/8/7Q/1R6/2K5 w - - 0 1")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&p)

			// normal commands
		case strings.HasPrefix(input, "setboard"):
			fields := strings.Fields(input)
			if len(fields) == 1 {
				ui.Result = "Error (no Fen): " + input
				break //next
			}
			fen := strings.Join(fields[1:], " ")
			ui.Result = "Parsing fen [" + fen + "]"
			p = hclibs.FENToNewBoard(fen)
			// 				proto.oln("Error (Not implemented yet!!!): " + input)
			// 				break next

		case strings.Contains(input, "ping"):
			ui.Status = "pong"

		case strings.Contains(input, "divide"):
			fields := strings.Fields(input)
			if len(fields) == 1 {
				ui.Result = "Please specify a number to divide down to"
				break //next
			}
			d, err := strconv.Atoi(fields[1])
			if err != nil {
				ui.Result = "Please specify a number"
				break //next
			}
			_, s := hclibs.Divide(d, &p)
			ui.Update = false
			proto.ofln(s)
			proto.o("Press return to continue")
			scanner.Scan()

		case strings.Contains(input, "perft"):
			fields := strings.Fields(input)
			if len(fields) == 1 {
				ui.Result = "Please specify a number to perft down to"
				break //next
			}
			d, err := strconv.Atoi(fields[1])
			if err != nil {
				ui.Result = "Please specify a number"
				break //next
			}
			start := time.Now()
			nodes := hclibs.Perft(d, &p)
			elapsed := time.Since(start)
			ui.Update = false
			proto.o(fmt.Sprintf("Perft to depth %v gives %v nodes ", d, hclibs.Comma(nodes)))
			proto.o(fmt.Sprintf("(nps: %v)\n", hclibs.Comma(int(float64(nodes)/elapsed.Seconds()))))
			proto.o("Press return to continue")
			scanner.Scan()
			// 			case strings.Contains(input, "depth"):
			// 				//  case strings.Contains(input,"fen") || strings.Contains(input,"setboard"):
			// 				fields := strings.Fields(input)
			// 				if len(fields) == 1 {
			// 					proto.o(fmt.Sprintf("Current depth is %d.\n", hclibs.GameDepthSearch))
			// 					break next
			// 				}
			// 				d, err := strconv.Atoi(fields[1])
			// 				if err != nil {
			// 					proto.oln("Please specify a number")
			// 					break next
			// 				}
			// 				if d > hclibs.MAXSEARCHDEPTH {
			// 					d = hclibs.MAXSEARCHDEPTH
			// 				}
			// 				proto.o(fmt.Sprintf("Current depth is %d. Setting depth to %d.\n", hclibs.GameDepthSearch, d))
			// 				hclibs.GameDepthSearch = d

		case input == "help":
			ui.Info = "Commands: move [a2a4],[a2a4], g[o], auto, quit, new, ping, depth #, perft #, divide #,\n          setboard [fen], kiwipete, pos4, pos5, end1, end2, end3, end4\n"

			// 			case strings.Contains(input, "auto"):
			// 				hclibs.GameForce = !hclibs.GameForce

			////////////////////////////////////////////////////////////////////////////////
		case strings.Contains(input, "go") || input == "g": // || hclibs.GameForce == true:
			ui.Status = "Thinking..."
			ui.ShowSpinner = true
			res, info, srch := hclibs.Go(&p, engineInfo)
			// proto.o(cls)
			ui.Board = hclibs.BoardToStrColour(&p)
			ui.Info = info
			ui.Result = res
			ui.Stats = srch.Stats.String()
			ui.Pv = srch.PV.String()
			ui.Status = "Awaiting user input..."
			ui.ShowSpinner = false

		default:
			ui.Result = "[" + input + "] not understood"
		}
		//

	}
	// tell ui go routine to stop
	uictl <- true
}
