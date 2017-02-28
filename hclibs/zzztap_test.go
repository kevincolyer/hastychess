//Hastychess, Copyright (C) GPLv3, 2016, Kevin Colyer
package hclibs

import "github.com/dex4er/go-tap"
import "testing"

func TestMain(t *testing.T) {
	// 		tap.Is("Aaa", "Aaa", "Is")
	//	tap.Is(123, 123, "Is")
	// This needs to be run last as Donetesting will fail and cause the go test to fail. But if tap is seen after donetesting it does not!
	tap.DoneTesting()

}
