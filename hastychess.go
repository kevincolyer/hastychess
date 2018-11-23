// Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package main

import (
	"flag"
	"fmt"
	"github.com/kevincolyer/hastychess/protocol"
	"os"
)

type ProtocolSpeaker interface {
	Start() error
	//     Echo() error
}

func main() {
	o := protocol.CLIOptions{NameVersion: "HastyChess 1.91 (Bonehead)"}

	var flagXboard = flag.Bool("xboard", false, "Select xboard mode")
	var flagIcs = flag.Bool("ics", false, "Select ics server mode")
	var flagUci = flag.Bool("uci", false, "Select UCI (ics) server mode (same as -ics)")
	var flagConsole = flag.Bool("console", true, "Select console mode")
	var flagPerft = flag.Bool("perft", false, "Select perft mode")
	var flagDivide = flag.Bool("divide", false, "Select divide mode")
	var flagTestprotocol = flag.Bool("testprotocols", false, "Select test engine for protocol debugging")
	var flagVersion = flag.Bool("version", false, "Give version info")

	var flagEngine = flag.String("engine", "", "Select engine (mode) one of chess,test,perft or divide")
	var flagProtocol = flag.String("protocol", "", "Select protocol to use: console,ics or xboard")
	var flagStats = flag.Bool("stats", true, "Enable printing of statistics")
	var flagBook = flag.Bool("book", true, "Enable the use of built in book moves")
	var flagTT = flag.Bool("tt", true, "Enable the use of Transposition Tables")
	var flagConsoleNoColor = flag.Bool("no-color", false, "Disable color output")
	var flagFen = flag.String("fen", "", "Send fen to engine (quotes needed)")
	var flagDepth = flag.Int("depth", 0, "Depth for search - used for perft and divide")
	var flagDebug = flag.Bool("debug", false, "Print extra debug information")
	var flagRBC = flag.Int("rbc", 0, "Really Bad Chess mode - 0 off, 1 easy, 2 hard, 3+ really hard")

	flag.Parse()

	if *(flagVersion) {
		fmt.Println(o.NameVersion)
		os.Exit(0)
	}

	o.Engine = "chess"     // default
	o.Protocol = "console" // default

	o.Engine = *(flagEngine)
	o.Protocol = *(flagProtocol)
	o.Stats = *(flagStats)
	o.Book = *(flagBook)
	o.TT = *(flagTT)
	o.ConsoleNoColor = *(flagConsoleNoColor)
	o.Fen = *(flagFen)
	o.Depth = *(flagDepth)
	o.Debug = *(flagDebug)
	o.RBC = *(flagRBC)

	switch {
	case *(flagXboard):
		o.Protocol = "xboard"
		o.ConsoleNoColor = true
	case *(flagIcs) || *(flagUci):
		o.Protocol = "ics"
		o.ConsoleNoColor = true
	case *(flagConsole):
		o.Protocol = "console"
	case *(flagPerft):
		o.Engine = "perft"
	case *(flagDivide):
		o.Engine = "divide"
	case *(flagTestprotocol):
		o.Protocol = "test"
	}

	// start protocol to speak
	var err error
	var myProtocol ProtocolSpeaker

	switch o.Protocol {
	case "console":
		myProtocol, err = protocol.NewConsole(os.Stdout, os.Stdin, os.Stderr, o)
		//             case "xboard":
		//                 myProtocol, err = protocol.NewXboard(os.Stdin, os.Stdout,o)
		//             case "ics":
		//                 myProtocol, err = protocol.NewIcs(os.Stdin, os.Stdout,o)
	}
	if err != nil {
		panic(err)
	}

	// All OK so start talking...
	err = myProtocol.Start()
	if err != nil {
		fmt.Println(err)
	}

	// protocol exited to here.

	//So clean up and...

	// exit
	os.Exit(0)
}
