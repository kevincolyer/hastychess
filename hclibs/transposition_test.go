//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import "github.com/dex4er/go-tap"
import "testing"

func TestRand64(t *testing.T) {
	// 		tap.Ok(true, "Ok")
	// 		tap.Is("Aaa", "Aaa", "Is")
	//	tap.Is(123, 123, "Is")
	//tap.DoneTesting()
	Rand64Reset()
	lim := 1000
	var a = make([]Hash, lim)
	var b = make([]Hash, lim)
	for i := 0; i < lim; i++ {
		a[i] = Rand64()
	}
	Rand64Reset()
	for i := 0; i < lim; i++ {
		b[i] = Rand64()
	}
	tap.Is(a[0], b[0], "testing random number generator 1")
	tap.Is(a[lim-1], b[lim-1], "testing random number generator 2")
	k := 0
	for _, i := range a {
		for _, j := range b {
			if i == j {
				k++
			}
		}
	}
	tap.Is(k, lim, "testing random number generator - all unique")
}

func TestInitHashSize(t *testing.T) {
	// 		tap.Ok(true, "Ok")
	// 		tap.Is("Aaa", "Aaa", "Is")
	//	tap.Is(123, 123, "Is")
	//tap.DoneTesting()
	size := 8
	e := InitHashSize(size)
	tap.Is(e, nil, "No error expected from function")
	tap.Is(len(tthash), size*1024*1024/8, "Is tthash the length we expected?")
	tap.Is(Zhash.mask, Hash(size*1024*1024/8-1), "Is Zhash.mask correct?")
}

func TestTTZKey(t *testing.T) {
	// put function test in here!!!!!
	p := FENToNewBoard(STARTFEN)
	key := TTZKey(&p)
	data, err := TTPeek(key, TTHASH)
	tap.Ok(err == false, "unitialised hash is empty")
	TTPoke(key, TTHASH, TtData{score: 1})
	data, err = TTPeek(key, TTHASH)
	tap.Is(data.score, 1, "Retrieved data from TT table OK")

}
