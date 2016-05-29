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
        lim:=1000
        var a=make([]uint64,lim)
        var b=make([]uint64,lim)
        for i:=0;i<lim;i++ {
            a[i]=Rand64()
        }
        Rand64Reset()
        for i:=0;i<lim;i++ {
            b[i]=Rand64()
        }
        tap.Is(a[0],b[0],"testing random number generator 1")
        tap.Is(a[lim-1],b[lim-1],"testing random number generator 2")
        k:=0
        for _,i:=range a {
            for _,j:=range b {
                if i==j { k++ }
            }
        }
        tap.Is(k,lim,"testing random number generator - all unique")
}

func TestInitHashSize(t *testing.T) {
	// 		tap.Ok(true, "Ok")
	// 		tap.Is("Aaa", "Aaa", "Is")
	//	tap.Is(123, 123, "Is")
	//tap.DoneTesting()
    size:=8
    e:=InitHashSize(size)
    tap.Is(e,nil,"No error expected from function")
    tap.Is(len(tthash),size*1024*1024/8,"Is tthash the length we expected?")
    tap.Is(Zhash.mask,uint64(size*1024*1024/8-1),"Is Zhash.mask correct?")
}