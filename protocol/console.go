//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer

// Package protocol ... Text ui
package protocol

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/kevincolyer/hastychess/engine"
	"github.com/kevincolyer/hastychess/hclibs"
)

// Console ... provides a text ui to communicate to game engine with
type Console struct {
	In      io.Reader
	Out     io.Writer
	Stderr  io.Writer
	Options CLIOptions
}

// EngineInfo ... datatype for messaging stats and game position info between componants
type EngineInfo struct {
	*hclibs.Statistics
	*hclibs.PV
}

func (proto *Console) o(s string) {
	io.WriteString(proto.Out, s)
}

func (proto *Console) oln(s string) {
	io.WriteString(proto.Out, s+"\n")
}

func (proto *Console) ofln(s interface{}) {
	io.WriteString(proto.Out, fmt.Sprintf("%v\n", s))
}

func (proto *Console) debug(s string) {
	io.WriteString(proto.Stderr, "Stderr: "+s)
}

func NewConsole(in io.Reader, out io.Writer, stderr io.Writer, options CLIOptions) (*Console, error) {
	p := Console{Out: out, In: in, Stderr: stderr, Options: options}
	p.debug("Started Console Protocol for " + options.NameVersion + "\n")
	return &p, nil
}

func (proto *Console) Echo() (e error) {
	b := make([]byte, 256)
	if n, err := proto.In.Read(b); err != io.EOF {
		io.WriteString(proto.Out, string(b[:n]))
		io.WriteString(proto.Stderr, "Stderr: "+string(b[:n]))
	} else {
		e = fmt.Errorf("EOF for input: none sent?")
	}
	return
}

func (proto *Console) Start() error {
	myEngine, err := engine.New(engine.Options{})
	if err != nil {
		return fmt.Errorf("Error creating engine: %v", err)
	}
	proto.MainLoop(myEngine)
	// 	if myEngine.Stop() {
	// 		io.WriteString(proto.Stderr, "Stderr: Engine stopped ok\n")
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
	proto       *Console
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

func (proto *Console) MainLoop(myEngine *engine.Engine) {
	hiwhite := color.New(color.FgHiWhite).SprintfFunc()
	ui := tui{Title: hiwhite("Hello and welcome to %v\n\n", proto.Options.NameVersion), Cls: "\033[H\033[2J",
		Pv: "pv", Stats: "stats", History: "History", Cmdline: "> ", Update: true,
	}
	pos := hclibs.FENToNewBoard(hclibs.STARTFEN)
	if proto.Options.RBC > 0 {
		pos = hclibs.FENToNewBoard(hclibs.NewRBCFEN(proto.Options.RBC))
	}

	ui.Board = hclibs.BoardToStrColour(&pos)
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
					ui.proto.oln("Info:\n" + ui.Info)
				}
				ui.proto.o(ui.Cmdline)
			}

			time.Sleep(time.Millisecond * 100)
			ui.spinner = spinners[spincounter]
			spincounter++
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

	// 	Channel to pass EngineInfo to ui.
	engineInfo := make(chan hclibs.EngineInfo, 1)
	go func(ei chan hclibs.EngineInfo) {
		// block until we have data...
		for {
			data := <-ei //:
			ui.Pv = fmt.Sprintf("%v", data.Pv)
			ui.Stats = fmt.Sprintf("%v", data.Stats)
			ui.Info = data.Info
		}
	}(engineInfo)

	// main input loop
	for quit == false {
		if pos.Side == hclibs.WHITE {
			ui.Cmdline = "WHITE >"
		} else {
			ui.Cmdline = "BLACK >"
		}
		ui.Update = false
		ui.UpdateOnce = true

		scanner.Scan()
		quit = (scanner.Err() == io.EOF)
		input := strings.TrimSpace(scanner.Text())

		ui.Info = ""
		ui.Result = ""
		ui.Update = true

		switch {
		case input == "quit" || input == "q":
			quit = true
			// 				break QUIT

		case strings.Contains(input, "move"):
			fields := strings.Fields(input)
			if len(fields) > 1 {
				move, err = hclibs.ParseUserMove(fields[1], &pos)
				if err != "" {
					ui.Result = err
					break //next
				}
			}
			result = hclibs.MakeUserMove(move, &pos)
			ui.Board = hclibs.BoardToStrColour(&pos)
			ui.Result = result

		case re.MatchString(input):
			move, err = hclibs.ParseUserMove(input, &pos)
			if err != "" {
				ui.Result = err
				break //next
			}
			result = hclibs.MakeUserMove(move, &pos)
			ui.Board = hclibs.BoardToStrColour(&pos)
			ui.Result = result

		case strings.Contains(input, "new"):
			pos = hclibs.FENToNewBoard(hclibs.STARTFEN)
			if proto.Options.RBC > 0 {
				pos = hclibs.FENToNewBoard(hclibs.NewRBCFEN(proto.Options.RBC))
			}
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&pos)

		// special commands that allow testing of certain positions
		case strings.Contains(input, "kiwipete"):
			pos = hclibs.FENToNewBoard("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&pos)

		case strings.Contains(input, "pos4"):
			pos = hclibs.FENToNewBoard("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&pos)

		case strings.Contains(input, "pos5"):
			pos = hclibs.FENToNewBoard("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&pos)

		case strings.Contains(input, "end1"):
			pos = hclibs.FENToNewBoard("8/8/4k3/7p/8/8/2K5/R6Q w - - 0 1")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&pos)

		case strings.Contains(input, "end2"):
			pos = hclibs.FENToNewBoard("8/8/8/1bn5/8/2k5/8/2K5 w - - 0 1")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&pos)

		case strings.Contains(input, "end3"):
			pos = hclibs.FENToNewBoard("8/8/8/1k6/8/7Q/3R4/2K5 w - - 0 1")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&pos)

		case strings.Contains(input, "end4"):
			pos = hclibs.FENToNewBoard("8/8/8/k7/8/7Q/1R6/2K5 w - - 0 1")
			// 				hclibs.GameOver = false
			ui.Board = hclibs.BoardToStrColour(&pos)

			// normal commands
		case strings.HasPrefix(input, "setboard"):
			fields := strings.Fields(input)
			if len(fields) == 1 {
				ui.Result = "Error (no Fen): " + input
				break //next
			}
			fen := strings.Join(fields[1:], " ")
			ui.Info = "Parsing fen [" + fen + "]"
			ui.Result = "Setting new position"
			pos = hclibs.FENToNewBoard(fen)
			ui.Board = hclibs.BoardToStrColour(&pos)

		case strings.HasPrefix(input, "rbc"):
			fields := strings.Fields(input)
			if len(fields) == 1 {
				ui.Result = "Error (missing number 1,2 or 3): " + input
				break //next
			}
			d, err := strconv.Atoi(fields[1])
			if err != nil {
				ui.Result = "Please specify a number"
				break //next
			}
			fen := hclibs.NewRBCFEN(d)
			ui.Result = "Setting new position"
			ui.Info = "RBC fen: " + fen
			pos = hclibs.FENToNewBoard(fen)
			ui.Board = hclibs.BoardToStrColour(&pos)

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
			_, s := hclibs.Divide(d, &pos)
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
			nodes := hclibs.Perft(d, &pos)
			elapsed := time.Since(start)
			ui.Update = false
			proto.o(fmt.Sprintf("Perft to depth %v gives %v nodes ", d, hclibs.Comma(nodes)))
			proto.o(fmt.Sprintf("(nps: %v)\n", hclibs.Comma(int(float64(nodes)/elapsed.Seconds()))))
			proto.o("Press return to continue")
			scanner.Scan()

		case input == "help":
			ui.Info = "Commands: move [a2a4],[a2a4], g[o], auto, quit, new, rbc #, ping, depth #, perft #, divide #,\n          setboard [fen], kiwipete, pos4, pos5, end1, end2, end3, end4\n"

			////////////////////////////////////////////////////////////////////////////////

		case strings.Contains(input, "go") || input == "g": // || hclibs.GameForce == true:
			ui.Cmdline = "Thinking..."
			ui.Status = "Thinking..."
			ui.ShowSpinner = true

			res, info, srch := hclibs.Go(&pos, engineInfo)

			ui.Board = hclibs.BoardToStrColour(&pos)
			ui.Info = ui.Info + info
			ui.Result = res
			ui.Stats = srch.Stats.String()
			ui.Pv = srch.PV.String()
			ui.Status = "Awaiting user input..."

			ui.ShowSpinner = false

			////////////////////////////////////////////////////////////////////////////////

		default:
			ui.Result = "[" + input + "] not understood"
		}

	}
	// tell ui go routine to stop
	uictl <- true
}
