package latlong_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/golang/geo/s2"
	latlong "github.com/toyo/go-latlong"
)

func TestCircle(t *testing.T) {
	circle := latlong.Circle{Point: latlong.Point{
		LatLng: s2.LatLng{Lat: math.Pi * 35 / 180, Lng: math.Pi * 139 / 180}},
		ChordAngle: latlong.Km(100).EarthChordAngle()}

	b, err := json.Marshal(circle.NewGeoJSONGeometry())
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}
	t.Logf("%v", string(b))

	var geom latlong.GeoJSONGeometry
	err = json.Unmarshal(b, &geom)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	circle1, err := geom.Circle()
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}
	t.Logf("%s", circle1.String())
}
