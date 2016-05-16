//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package main

import "github.com/dex4er/go-tap"
import "testing"
import "hastychess/hclibs"
import "fmt"

func TestParsePositionInput(t *testing.T) {
	// 		tap.Ok(true, "Ok")
	var f string
	var m []hclibs.Move
	fen := hclibs.STARTFEN

	i := "position " + fen
	f, m = ParsePositionInput(i)
	tap.Is(len(m), 0, "Parses STARTFEN with 0 moves ok")
	i += " moves a2a4"
	f, m = ParsePositionInput(i)
	tap.Is(len(m), 1, "Parses STARTFEN with 1 move ok")
	i += " a2a4"
	f, m = ParsePositionInput(i)
	tap.Is(len(m), 2, "Parses STARTFEN with 2 move ok")

	j := "position startpos"
	f, m = ParsePositionInput(j)
	tap.Is(len(m), 0, "Parses startpos with 0 moves ok")
	tap.Is(f, fen, "Parses startpos with no moves ok")
	j += " moves a2a4"
	f, m = ParsePositionInput(j)
	tap.Is(len(m), 1, "Parses startpos with 1 moves ok")
	tap.Is(f, fen, "Parses startpos with 1 moves ok")
	j += " a2a4"
	f, m = ParsePositionInput(j)
	tap.Is(len(m), 2, "Parses startpos with 2 moves ok")
	tap.Is(f, fen, "Parses startpos with 2 moves ok")

	i = "position 8/1k6/8/5N2/8/4n3/8/2K5 b - - 0 1"
	j = "8/1k6/8/5N2/8/4n3/8/2K5 b - - 0 1"
	f, m = ParsePositionInput(i)
	tap.Is(len(m), 0, "Parses fen with 0 moves ok")
	tap.Is(f, j, "Parses fen ok")
	i += " moves a2a4 a2a3"
	f, m = ParsePositionInput(i)
	tap.Is(len(m), 2, "Parses fen with 2 moves ok")
	tap.Is(f, j, "Parses fen ok")

	fmt.Println(m)
	fmt.Println(m[1])
	// 		tap.Is("Aaa", "Aaa", "Is")
	//	tap.Is(123, 123, "Is")
	tap.DoneTesting()
}
