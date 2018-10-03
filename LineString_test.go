package latlong_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/golang/geo/r3"
	"github.com/golang/geo/s2"
	latlong "github.com/toyo/go-latlong"
)

func TestLineString(t *testing.T) {

	const xmlstrings = `<ISO6709>+12+123/+12.3+123.4/+12.34+123.43/</ISO6709>`
	var answer = s2.Polyline{
		s2.Point{
			Vector: r3.Vector{
				X: -0.5327373653659239, Y: 0.8203436038418747, Z: 0.20791169081775934}},
		s2.Point{
			Vector: r3.Vector{
				X: -0.5378447709118935, Y: 0.8156844101282474, Z: 0.21303038627497659}},
		s2.Point{
			Vector: r3.Vector{
				X: -0.5381897230643339, Y: 0.8152783663496835, Z: 0.2137124407939944}}}
	var ISO6709 latlong.LineString

	err := xml.Unmarshal([]byte(xmlstrings), &ISO6709)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	if !answer.Equal(ISO6709.S2Region()) {
		t.Errorf("TestLineString %#v", ISO6709.S2Region())
	}
}

func TestNewLatLongAltsString(t *testing.T) {
	//str := `+32.9+130.9-10000/`
	//str := `+30.4402+130.2000/+30.4588+130.2200/+30.4545+130.2400/+30.4400+130.2545/+30.4200+130.2489/+30.4061+130.2400/+30.4037+130.2200/+30.4200+130.2056/+30.4398+130.2000/+30.4400+130.1998/+30.4402+130.2000/`
	ll := latlong.NewMultiPointISO6709("+12.34+123.45+3776/+0123.4-01234.5-3776/-001234+0012345-12345/")
	if ll == nil {
		t.Errorf("NewLatLongAltssISO6709 returned non nil error")
	}

	latlong.Config.Lang = "ja"
	lls := ll.String()
	correctResponsells := "北緯12.34度、東経123.45度、標高3776m,北緯1.390度、西経12.575度、ごく浅く,南緯0.2094度、東経1.3958度、深さ12km"
	if lls != correctResponsells {
		t.Errorf("expected %+v, was %+v", correctResponsells, lls)
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(ll.NewGeoJSONGeometry())
	correctResponseJSON := `{"type":"MultiPoint","coordinates":[[123.45,12.34,3776],[-12.575,1.390,-3776],[1.3958,-0.2094,-12345]]}
`
	if err != nil {
		t.Error(err)
		return
	}
	JSON := b.String()
	if JSON != correctResponseJSON {
		t.Errorf("expected '%+v', was '%+v'", correctResponseJSON, JSON)
	}

}

func TestLatLongstring(t *testing.T) {
	ll := latlong.NewMultiPointISO6709("+12+123/+12.3+123.4/+12.34+123.43/")
	if ll == nil {
		t.Errorf("NewLatLongsISO6709 returned nil error")
	}

	latlong.Config.Lang = "ja"
	lls := ll.String()
	correctResponsells := "北緯12度、東経123度,北緯12.3度、東経123.4度,北緯12.34度、東経123.43度"
	if lls != correctResponsells {
		t.Errorf("expected %+v, was %+v", correctResponsells, lls)
		for _, l := range *ll {
			t.Error(l.String())
			t.Error(l.PrecString())
		}
	}
}
