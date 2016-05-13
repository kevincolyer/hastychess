//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

func TtCull() {
	StatTtCulls = 0

}

// func TtPeek(TtKey string, ply int, val *int, alpha int, beta int) (m Move, bool {
//     if GameUseTt==false { return GameUseTt }
//     item,ok:=tt[Ttkey]
//     if ok == nil { return false }
//     // something found. Was this entry searched to same or deeper level that we are on now?
//     if item.ply >= ply {
//
//
//     }
//
// }
//
// func TtPoke(TtKey string, ply int, score int, tttype int ) bool {
//     if GameUseTt==false { return GameUseTt }
//
//
// }
