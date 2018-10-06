package latlong_test

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/golang/geo/s2"
	latlong "github.com/toyo/go-latlong"
)

func TestPolygon1(t *testing.T) {

	const xmlstrings = `<ISO6709>+31.9388+130.8800/+31.9388+130.9000/+31.9372+130.9200/+31.9312+130.9400/+31.9265+130.9600/+31.9200+130.9783/+31.9000+130.9682/+31.8958+130.9600/+31.8909+130.9400/+31.8829+130.9200/+31.8812+130.9000/+31.8812+130.8800/+31.9000+130.8612/+31.9200+130.8612/+31.9388+130.8800/</ISO6709>`

	var ISO6709 latlong.Polygon

	err := xml.Unmarshal([]byte(xmlstrings), &ISO6709)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	b, err := json.Marshal(&ISO6709)
	if err != nil {
		t.Error("Error")
	}

	expct := `{"MultiPoint":[[130.8800,31.9388],[130.9000,31.9388],[130.9200,31.9372],[130.9400,31.9312],[130.9600,31.9265],[130.9783,31.9200],[130.9682,31.9000],[130.9600,31.8958],[130.9400,31.8909],[130.9200,31.8829],[130.9000,31.8812],[130.8800,31.8812],[130.8612,31.9000],[130.8612,31.9200],[130.8800,31.9388]]}`
	if string(b) != expct {
		t.Errorf("Wrong got %s expct %s", string(b), expct)
	}

	if ISO6709.S2Region().ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(31.93, 130.90))) != true {
		t.Error("Something wrong.")
	}

	if ISO6709.S2Region().ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(31.93, 140.90))) != false {
		t.Error("Something wrong.")
	}
}
