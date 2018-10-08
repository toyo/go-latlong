package latlong_test

import (
	"encoding/json"
	"testing"

	latlong "github.com/toyo/go-latlong"
)

func TestCircle(t *testing.T) {
	circle := latlong.NewCircle(
		latlong.NewPoint(latlong.NewAngle(35, 0), latlong.NewAngle(139, 0), nil),
		100)
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

	circle1 := geom.Geo()
	if circle1 == nil {
		t.Error("Unmarshal error")
	}
	t.Logf("%s", circle1.(latlong.Circle).String())
}
