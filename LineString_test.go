package latlong_test

import (
	"encoding/xml"
	"testing"

	latlong "github.com/toyo/go-latlong"
)

func TestLineString(t *testing.T) {

	const xmlstrings = `<ISO6709>+12+123/+12.3+123.4/+12.34+123.43/</ISO6709>`

	var ISO6709 latlong.LineString

	err := xml.Unmarshal([]byte(xmlstrings), &ISO6709)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	expct := *latlong.NewLatLongAlt(latlong.NewAngle(12, 1), latlong.NewAngle(123, 1), nil)
	if ISO6709.MultiPoint[0] != expct {
		t.Errorf("Not match got %#v expct %#v", ISO6709.MultiPoint[0], expct)
	}

	expct = *latlong.NewLatLongAlt(latlong.NewAngle(12.3, 0.1), latlong.NewAngle(123.4, 0.1), nil)
	if ISO6709.MultiPoint[1] != expct {
		t.Errorf("Not match got %#v expct %#v", ISO6709.MultiPoint[0], expct)
	}

	expct = *latlong.NewLatLongAlt(latlong.NewAngle(12.34, 0.01), latlong.NewAngle(123.43, 0.01), nil)
	if ISO6709.MultiPoint[2] != expct {
		t.Errorf("Not match got %#v expct %#v", ISO6709.MultiPoint[0], expct)
	}

}
