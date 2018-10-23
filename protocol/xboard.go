// Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package protocol

import (
	"bufio"
	"fmt"

	"github.com/kevincolyer/hastychess/hclibs"
	"io"
	"regexp"
	"strconv"
	"strings"
)



type xboard struct {
	In      io.Reader
	Out     io.Writer
	Options CLIOptions
}

func NewXboard(in io.Reader, out io.Writer, options CLIOptions) (*xboard, error) {
	p := xboard{Out: out, In: in, Options: options}
	_, e := io.WriteString(p.Out, "Welcome to "+options.NameVersion+"\n")
	return &p, e
}

func (p *xboard) Echo() (e error) {
	b := make([]byte, 256)
	if n, err := p.In.Read(b); err != io.EOF {
		io.WriteString(p.Out, string(b[:n]))
	} else {
		e = fmt.Errorf("EOF for input: none sent?")
	}
	return
}

func (p *xboard) Start() (e error) {
	io.WriteString(p.Out, "started ok\n")
	e = fmt.Errorf("Done!")
	return
}

//===========================================================================


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
	hclibs.GameDepthSearch = 8
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

