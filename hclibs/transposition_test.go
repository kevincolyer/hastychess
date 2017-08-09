//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import "github.com/dex4er/go-tap"
import "testing"
import "fmt"

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
	var size = 8
	l := ttable.InitHashSize(size)
	ttable = make([]TtData, l)
	tap.Is(len(ttable), size*1024*1024/8, "Is tthash the length we expected?")
	tap.Is(Zhash.mask, Hash(size*1024*1024/8-1), "Is Zhash.mask correct?")

}

func TestTTZKey(t *testing.T) {
	// put function test in here!!!!!
	p := FENToNewBoard(STARTFEN)
	key := TTZKey(&p)
	ttable = make([]TtData, ttable.InitHashSize(8))

	data := ttable.Peek(key)
	tap.Ok(data.IsInUse() == false, "unitialised hash is empty")

	ttable.SafePoke(key, TtData{score: 1})
	ttable.Clear()
	tap.Ok(data.IsInUse() == false, "emptied hash is empty")

	// test insertion and retrevial
	ttable.SafePoke(key, TtData{score: 1})
	data = ttable.Peek(key)
	tap.Is(data.score, 1, "Retrieved data from TT table OK")

	tap.Is(fmt.Sprintf("%x", Zhash.mask), "fffff", "hash mask calculated correctly for 1024x1024-1")

	// test manual masking of hash and retrevial
	ttable.SafePoke(key, TtData{score: 2})
	tap.Is(ttable[key&Zhash.mask], ttable.Peek(key), "Test manual masking of hash and retreval")
	tap.Is(ttable[key&Zhash.mask].score, 2, "Test manual masking of hash and retreval (value)")

	// test one positions hash is different from anothers
	fmt.Println(&p)
	MakeMove(Move{from: A2, to: A4, mtype: QUIET}, &p)
	fmt.Println(&p)
	tap.Ok(key != p.Hash, "Different positions have different hashes")
	fmt.Printf("key=%v p.Hash=%v\n", key, p.Hash)
	ttable.SafePoke(p.Hash, TtData{score: 3})
	tap.Is(ttable.Peek(key).score, 2, "Check hashes are not clobbering")
	tap.Is(ttable.Peek(p.Hash).score, 3, "Check hashes are not clobbering")
	UnMakeMove(Move{from: A2, to: A4, mtype: QUIET}, &p)
	tap.Is(p.Hash, key, "UnMake Move resets hash OK")

}
