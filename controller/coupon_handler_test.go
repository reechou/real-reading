package controller

import (
	//"fmt"
	"testing"

	"github.com/dchest/uniuri"
)

func TestCoupon(t *testing.T) {
	uris := make([]string, 100000)
	for i := range uris {
		uris[i] = uniuri.New()
		//fmt.Printf("%d : %s\n", i+1, uris[i])
	}
	for i, u := range uris {
		for j, u2 := range uris {
			if i != j && u == u2 {
				t.Fatalf("not unique: %d:%q and %d:%q", i, u, j, u2)
			}
		}
	}
}
