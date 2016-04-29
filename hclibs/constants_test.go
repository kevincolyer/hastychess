//  hasty chess libs
//  constants
package hclibs

import "fmt"
import "testing"

func TestGrid(t *testing.T) {
	// 	var GRID [64]int
	// 	GRID = MakeGrid()
	//     fmt.Println(GRID)
	if GRID[0] != A1 {
		t.Errorf("GRID 0 is not constant A1")
	}
	if GRID[7] != H1 {
		t.Errorf("GRID 7 is not constant H1")
	}
	if GRID[56] != A8 {
		t.Errorf("GRID 56 (%v) is not constant A8=%v", GRID[56], H1)
	}
	if GRID[63] != H8 {
		t.Errorf("GRID 63 is not constant H8=%v", H8)
	}
}

func TestArrays(t *testing.T) {
	if fmt.Sprintf("%v", NM) != "[31 33 -31 -33 14 18 -14 -18]" {
		t.Errorf("Array Error")
	}
	if fmt.Sprintf("%v", BM) != "[17 15 -15 -17]" {
		t.Errorf("Array Error")
	}
	if fmt.Sprintf("%v", KM) != "[16 1 -1 -16 17 15 -15 -17]" {
		t.Errorf("Array Error")
	}
	if fmt.Sprintf("%v", QM) !=
		"[16 1 -1 -16 17 15 -15 -17]" {
		t.Errorf("Array Error")
	}
	if fmt.Sprintf("%v", RM) !=
		"[16 1 -1 -16]" {
		t.Errorf("Array Error")
	}
	if fmt.Sprintf("%v", PM) !=
		"[[15 16 17 32] [-15 -16 -17 -32]]" {
		t.Errorf("Array Error")
	}
}

func TestAlg(t *testing.T) {
	if A1 != 0 {
		t.Errorf("Alg Error A1 !=0 (got %v)", A1)
	}
	if H1 != 7 {
		t.Errorf("Alg Error H1 !=7 (got %v)", H1)
	}
	if A8 != 112 {
		t.Errorf("Alg Error A8 !=112 (got %v)", A8)
	}
	if H8 != 119 {
		t.Errorf("Alg Error H8 !=119 (got %v)", H8)
	}
}
