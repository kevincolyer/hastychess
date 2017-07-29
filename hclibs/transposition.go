//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer

// Zorbist hashing functions (replaces a naive system of hashmaps and text string keys)
package hclibs

import "fmt"

type Hash uint64

type zhashstruct struct {
	psq    [128][20]Hash
	castle [4]Hash
	side   [2]Hash
	ep     [129]Hash // have to add 1 to ep value to as -1 indicates we are not in ep
	mask   Hash
}

var Zhash zhashstruct

// how do I declare an array but set its size later?
// don't bother! Slices are OK.
var tthash []TtData
var qstthash []TtData

type TT []TtData

var hashnext Hash = 1

func Rand64Reset() {
	hashnext = 1
}

// Rand64 generates uint64's and is pinched the one the cpw engine pinches from Sungorus!
func Rand64() Hash {
	hashnext = hashnext*1103515245 + 12345
	return hashnext
}

// in init create and fill the zorbist hash with the unique values for each piece and square and position
func init() {
	i := 0
	for i = 0; i < 128; i++ {
		Zhash.ep[i] = Rand64()
		for j := 0; j < 20; j++ {
			Zhash.psq[i][j] = Rand64()
		}
	}
	Zhash.ep[128] = Rand64()

	for i = 0; i < 4; i++ {
		Zhash.castle[i] = Rand64()
	}
	Zhash.side[BLACK] = Rand64()
	Zhash.side[WHITE] = Rand64()
}

//type TTZKey uint64

func (tt TT) Peek(key Hash) (data TtData) {
	key = key & Zhash.mask
	StatTtPeeks++
	return tt[key] // a nil entry is means never considered
}

// NOTE THIS IS NOT THREAD SAFE!!! THIS COULD BE A BIG DEAL!
func (tt TT) Poke(key Hash, data TtData) {
	StatTtWrites++
	key = key & Zhash.mask
	tt[key] = data
	return
}

// An attempt at a bit of a better thread safe model.
func (tt TT) SafePoke(key Hash, data TtData) {
	StatTtWrites++
	key = key & Zhash.mask
	tt[key] = data
	// try to clobber another thread. Probably doesn't matter so much.
	if tt[key] != data {
		tt[key] = data
	}
	return
}

func (tt TT) Clear() {
	for i := range tt {
		tt[i] = TtData{}
	}
	return
}

// make TtKey - scan board for pieces, xor in, xor in castling states, xor in side to move and EP
func TTZKey(p *Pos) (z Hash) {
	for _, square := range GRID {
		if p.Board[square] != EMPTY {
			z = z ^ Zhash.psq[square][p.Board[square]]
		}
	}

	if p.Castled[BLACK*2+QS] {
		z = z ^ Zhash.castle[BLACK*2+QS]
	}
	if p.Castled[BLACK*2+KS] {
		z = z ^ Zhash.castle[BLACK*2+KS]
	}
	if p.Castled[WHITE*2+QS] {
		z = z ^ Zhash.castle[WHITE*2+QS]
	}
	if p.Castled[WHITE*2+KS] {
		z = z ^ Zhash.castle[WHITE*2+KS]
	}

	z = z ^ Zhash.ep[p.EnPassant+1]

	z = z ^ Zhash.side[p.Side]
	return
}

// droid fish gives options to set size off tt in MB. Not sure how big an entry really is so lets suppose it is 8 bytes...

// initialises the hash to the size that the engine will set. size is given in human terms of number of entries e.g. 32M=32 million byte / 8
// this needs to always be done only once. ics engine sends a command to do this. Xboard also. Conoles we keep it set at
func (tt TT) InitHashSize(size uint64) uint64 {
	size = size * 1024 * 1024 / 8

	if size > (1<<TTMAXSIZE) || size <= 0 {
		panic(fmt.Errorf("size %d is larger than max allowd (or < 1)", size))
	}
	var power uint8 = 0
	for size > 0 {
		size = size >> 1
		power++
	}
	size = 1 << (power - 1)
	Zhash.mask = Hash(size - 1)
	return size
}
