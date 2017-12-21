package latlong

import (
	"math/rand"
	"testing"
	"time"
)

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

		l, _ := NewLatLongGeoHash(geohash)

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
	l, _ := NewLatLongGeoHash(geohash)

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
	l, _ := NewLatLongGeoHash(geohash)

	gh := l.GeoHash6()
	correctResponsegh := geohash
	if gh != correctResponsegh {
		t.Errorf("expected %+v, was %+v", correctResponsegh, gh)
	}
}
