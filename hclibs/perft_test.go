//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

//
import "fmt"
import "github.com/dex4er/go-tap"
import "testing"
import "io/ioutil"
import "strings"
import "strconv"

//import "fmt"

func TestGenerateMoves(t *testing.T) {
	//	tap.Ok(true, "Ok")
	p := FENToNewBoard(STARTFEN)
	tap.Is(IsValidMove(Move{from: A2, to: A3, mtype: QUIET}, &p), true, "Testing IsValidMove - finds correct move")
	tap.Is(IsValidMove(Move{from: A2, to: B3}, &p), false, "Testing IsValidMove - finds incorrect move")
	//tap.Is("Aaa", "Aaa", "Is")
	tap.Is(len(GenerateAllMoves(&p)), 20, "20 moves counted on a new board")
	tap.Is(Perft(1, &p), 20, "first test of perft")
	tap.Is(Perft(2, &p), 400, "2nd test of perft")
	tap.Is(Perft(3, &p), 8902, "third test of perft")

	tap.Is(Divide(4, &p), 197281, "4th test of divide")
	// 	tap.Is(123, 123, "Is")

	dat, err := ioutil.ReadFile("perftsuite.epd")
	check(err)
	lines := strings.Split(string(dat), "\r\n")
	for l, i := range lines {
		//fmt.Println(i)
		if i == "" {
			break
		}

		///////// COMMENT OUT THIS TO GET GREATER TESTING
		if l > 1 {
			break
		}
		/////////////////////////////////////////////////

		items := strings.Split(i, ";")
		fen := items[0]
		j := len(items)
		//fmt.Printf("j=%d,items=%v\n",items)
		j = 4 // temp limit - if this commented out then run with go test -timeout 24h to avoid timeouts
		for k := 1; k < j; k++ {
			//fmt.Print(items[k])
			test := strings.Split(items[k], " ")
			//fmt.Printf("k=%d, test=%v\n",k,test)
			d, _ := strconv.Atoi(test[0][1:])
			comp, err := strconv.Atoi(test[1])
			if err != nil {
				fmt.Println("Error found ", err)
			}
			q := FENToNewBoard(fen)
			tap.Is(Perft(d, &q), comp, "line "+strconv.Itoa(l)+") "+fen+" depth "+test[0]+" is "+test[1])
		}
	}
	// 	p = FENToNewBoard("4k3/4p3/4K3/8/8/8/8/8 b - - 0 1")
	// 	Divide(1, p) // should be 2 but we get 4
	// 	fmt.Println(BoardToStr(&p))
	/*
			p = FENToNewBoard("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
			Divide(4, p) // should be 25 but we get 24
			fmt.Println(BoardToStrWide(&p))
			MakeMove(Move{B2, B3, 0, 0}, &p) //B2B3
			Divide(3, p) // should be 25 but we get 24
			fmt.Println(BoardToStrWide(&p))
			MakeMove(Move{A6,B5,0,0},&p) // A6B5
		  Divide(2, p) // should be 25 but we get 24
		  fmt.Println(BoardToStrWide(&p))
			MakeMove(Move{A2,A4,ENPASSANT,A4},&p) // A2a4
		  Divide(1, p) // should be 25 but we get 24
		  fmt.Println(BoardToStrWide(&p))
			MakeMove(Move{B4,A3,EPCAPTURE,0},&p) // b4a3 capture
			Divide(1, p) // should be 25 but we get 24
			fmt.Println(BoardToStrWide(&p))
	*/
	//DeepPerftTest(t)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func BenchmarkPerft(b *testing.B) {

	p := FENToNewBoard(STARTFEN)
	nodes := Perft(4, &p)
	fmt.Println("Nodes " + strconv.Itoa(nodes))
}
