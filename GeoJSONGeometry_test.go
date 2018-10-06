package latlong_test

import (
	"encoding/json"
	"testing"

	latlong "github.com/toyo/go-latlong"
)

func TestGeoJSONGeometryLineString(t *testing.T) {
	jsonstring := `{ 
	"type": "LineString",
	"coordinates": [ [100.0, 0.0], [101.0, 1.0] ]
	}`

	var geom latlong.GeoJSONGeometry
	err := json.Unmarshal([]byte(jsonstring), &geom)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	lsp := geom.Geo()
	if lsp == nil {
		t.Errorf("LineString error: %v", err)
	}
	t.Logf("%s", lsp.(latlong.LineString))
}

func TestGeoJSONGeometryCircle(t *testing.T) {
	jsonstring := `{
	"type": "Circle",
	"coordinates": [ 100.0, 0.0 ],
	"radius": 0.5,
	"properties": {
	  "radius_units": "km"
	}
  }`

	var geom latlong.GeoJSONGeometry
	err := json.Unmarshal([]byte(jsonstring), &geom)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	lsp := geom.Geo()
	if lsp == nil {
		t.Error("Circle error")
	}
	t.Logf("%s", lsp.(latlong.Circle).String())
}

func TestGeoJSONGeometryPolygon(t *testing.T) {
	jsonstring := `{ "type": "Polygon",
	"coordinates": [
	  [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0], [100.0, 1.0], [100.0, 0.0] ] ]
   }`

	var geom latlong.GeoJSONGeometry
	err := json.Unmarshal([]byte(jsonstring), &geom)
	if err != nil {
		t.Error("Unmarshal error")
	}

	lsp := geom.Geo()
	if lsp == nil {
		t.Error("Polygon error")
	}
	t.Logf("%s", lsp.(latlong.Polygon))
}

func TestGeoJSONGeometryPoint(t *testing.T) {
	p := latlong.NewPointISO6709([]byte(`+35.000+139.0000`))
	llj := latlong.NewGeoJSONGeometry(p)
	b, err := json.Marshal(llj)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}
	t.Logf("%v", string(b))

	var ll1j latlong.GeoJSONGeometry
	err = json.Unmarshal(b, &ll1j)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	if !llj.Equal(ll1j) {
		t.Errorf("Unmatched expct %#v got %#v", llj, ll1j)
	}
}
