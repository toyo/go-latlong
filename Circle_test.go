package latlong

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/golang/geo/s2"
)

func TestCircle(t *testing.T) {
	circle := Circle{Point: Point{LatLng: s2.LatLng{Lat: math.Pi * 35 / 180, Lng: math.Pi * 139 / 180},
		latprec: 2, lngprec: 3}, ChordAngle: Km(100).EarthChordAngle()}

	b, err := json.Marshal(circle)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}
	t.Logf("%v", string(b))

	var geom Circle
	err = json.Unmarshal(b, &geom)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	t.Logf("%s", geom.String())
}
