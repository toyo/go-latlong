package latlong

import (
	"encoding/json"
	"testing"
)

func TestGeoJSONGeometryLineString(t *testing.T) {
	jsonstring := `{ 
	"type": "LineString",
	"coordinates": [ [100.0, 0.0], [101.0, 1.0] ]
	}`

	var geom GeoJSONGeometry
	err := json.Unmarshal([]byte(jsonstring), &geom)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	lsp, err := geom.LineString()
	if err != nil {
		t.Errorf("LineString error: %v", err)
	}
	t.Logf("%s", lsp)
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

	var geom GeoJSONGeometry
	err := json.Unmarshal([]byte(jsonstring), &geom)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	lsp, err := geom.Circle()
	if err != nil {
		t.Errorf("Circle error: %v", err)
	}
	t.Logf("%s", lsp.String())
}

func TestGeoJSONGeometryPolygon(t *testing.T) {
	jsonstring := `{ "type": "Polygon",
	"coordinates": [
	  [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0], [100.0, 1.0], [100.0, 0.0] ] ]
   }`

	var geom GeoJSONGeometry
	err := json.Unmarshal([]byte(jsonstring), &geom)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	lsp, err := geom.Polygon()
	if err != nil {
		t.Errorf("Polygon error: %v", err)
	}
	t.Logf("%s", lsp)
}
