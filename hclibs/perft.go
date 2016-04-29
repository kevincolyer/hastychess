package hclibs

import (
	"fmt"
	"github.com/dex4er/go-tap"
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	"strconv"
)

// total count of nodes of a given depth
func Perft(depth int, p Pos) (nodes int) {

	var moves []Move
	if depth == 0 {
		return 1
	} // because b & w have turns...
	moves = append(moves, GenerateAllMoves(&p)...)

	for _, m := range moves {
		q := p
		MakeMove(m, &q)
		nodes += Perft(depth-1, q)
	}
	// returning from perft;
	return

}

// helper structs for divide
type divide struct {
	move  string
	nodes int
}

// pretty printer for divide struct
func (d divide) String() string {
	return fmt.Sprintf("%s: %d\n", d.move, d.nodes)
}

// used to provide a sort of divide struct by nodes
type by []divide

func (a by) Len() int           { return len(a) }
func (a by) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a by) Less(i, j int) bool { return a[i].move < a[j].move }

// the useful divide function shows counts of nodes for top level moves
func Divide(depth int, p Pos) (nodes_total int) {

	var moves []Move
	var result []divide
	if depth == 0 {
		return 1
	} // because b & w have turns...
	moves = append(moves, GenerateAllMoves(&p)...)
	nodes := 0

	for _, m := range moves {
		q := p
		MakeMove(m, &q)
		nodes = Perft(depth-1, q)

		result = append(result, divide{MoveToAlg(m), nodes})
		nodes_total += nodes
	}
	// returning from perft;
	sort.Sort(by(result))
	fmt.Println(result)
	fmt.Printf("Total moves: %v\n", len(result))
	fmt.Printf("Total nodes: %v\n", nodes_total)
	return
}

func DeepPerftTest(t *testing.T) {
	dat, err := ioutil.ReadFile("perftsuite.epd")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(dat), "\r\n")
	for l, i := range lines {
		//fmt.Println(i)
		if i == "" {
			break
		}
		items := strings.Split(i, ";")
		fen := items[0]
		j := len(items)
		//fmt.Printf("j=%d,items=%v\n",items)
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
			tap.Is(Perft(d, q), comp, "line "+strconv.Itoa(l)+") "+fen+" depth "+test[0]+" is "+test[1])
		}
	}
	tap.DoneTesting()
}
