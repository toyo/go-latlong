package latlong

import (
	"math/rand"
	"testing"
)

func randomGridLocator(n int) (bs []byte) {

	bs = make([]byte, 0, n)
	for i := 0; i < n; i++ {
		switch i % 4 {
		case 0:
			fallthrough
		case 1:
			if i < 4 {
				bs = append(bs, byte('A'+rand.Intn(18)))
			} else {
				bs = append(bs, byte('A'+rand.Intn(24)))
			}
		case 2:
			fallthrough
		case 3:
			bs = append(bs, byte('0'+rand.Intn(10)))
		}
	}
	return
}

func TestGridLocator(t *testing.T) {
	randInit()

	for n := 0; n <= 12; n++ {
		grid := string(randomGridLocator(n))
		//fmt.Println(n, grid)
		l := NewLatLongGridLocator(grid)
		gl := l.GridLocator()
		le := len(grid)
		var correctResponsegl string
		if le%2 == 1 {
			correctResponsegl = grid[:le-1]
		} else {
			correctResponsegl = grid
		}

		if gl != correctResponsegl {
			t.Errorf("expected %+v, was %+v", correctResponsegl, gl)
		}

	}
}
