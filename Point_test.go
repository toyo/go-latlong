package latlong_test

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	latlong "github.com/toyo/go-latlong"
)

func TestString(t *testing.T) {
	ll := latlong.NewRectGridLocator("PM95UQ")
	//ll := NewRectGridLocator("JJ00AA")

	latlong.Config.Lang = "ja"
	ls := ll.Center().String()
	correctResponsels := "北緯35.69度、東経139.71度"
	if ls != correctResponsels {
		t.Errorf("expected %+v, was %+v", correctResponsels, ls)
		t.Error(ll.Center().String())
		t.Error(ll.Center().PrecString())
	}

	lp := ll.PrecString()
	correctResponselp := "緯度誤差0.041667度、経度誤差0.083333度"
	if lp != correctResponselp {
		t.Errorf("expected %+v, was %+v", correctResponselp, lp)
	}

}

func TestPointJSON(t *testing.T) {
	b := new(bytes.Buffer)
	centerll := latlong.NewLatLongAlt(latlong.NewAngle(35.69, 0.01), latlong.NewAngle(139.71, 0.01), nil)
	err := json.NewEncoder(b).Encode(centerll)
	correctResponseJSON := []byte(`[139.71,35.69]
`)
	if err != nil {
		t.Error(err)
		return
	}
	JSON := b.Bytes()
	if !bytes.Equal(JSON, correctResponseJSON) {
		t.Errorf("expected '%+v', was '%+v'", correctResponseJSON, JSON)
	}

	//t.Logf(`JSON=%s`, string(JSON))

	var ll latlong.Point
	json.NewDecoder(bytes.NewBuffer(JSON)).Decode(&ll)
	if ll != *centerll {
		t.Errorf("expected '%+v', was '%+v'", centerll, ll)
	}
}

//

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
		l := latlong.NewRectGridLocator(grid)
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

//

func randInit() {
	rand.Seed(time.Now().UnixNano())
}

func randomGeoHash1() (s rune) {
	return []rune("0123456789bcdefghjkmnpqrstuvwxyz")[int(rand.Intn(32))]
}

func randomGeoHash(n int) string {
	var r []rune
	for i := 0; i < n; i++ {
		r = append(r, randomGeoHash1())
	}
	return string(r)
}

func TestGeoHash(t *testing.T) {
	randInit()

	for n := 1; n < 12; n++ {
		geohash := randomGeoHash(n)
		//fmt.Println(geohash)

		l, _ := latlong.NewRectGeoHash(geohash)

		gh := l.GeoHash()
		correctResponsegh := geohash
		if gh != correctResponsegh {
			t.Errorf("expected %+v, was %+v", correctResponsegh, gh)
		}
	}
}

func TestGeoHash5(t *testing.T) {
	randInit()

	geohash := randomGeoHash(5)
	//fmt.Println(geohash)
	l, _ := latlong.NewRectGeoHash(geohash)

	gh := l.GeoHash5()
	correctResponsegh := geohash
	if gh != correctResponsegh {
		t.Errorf("expected %+v, was %+v", correctResponsegh, gh)
	}
}

func TestGeoHash6(t *testing.T) {
	randInit()

	geohash := randomGeoHash(6)
	//fmt.Println(geohash)
	l, _ := latlong.NewRectGeoHash(geohash)

	gh := l.GeoHash6()
	correctResponsegh := geohash
	if gh != correctResponsegh {
		t.Errorf("expected %+v, was %+v", correctResponsegh, gh)
	}
}
