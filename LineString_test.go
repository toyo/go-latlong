package latlong

import (
	"encoding/xml"
	"testing"

	"github.com/golang/geo/r3"
	"github.com/golang/geo/s2"
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
	var ISO6709 LineString

	err := xml.Unmarshal([]byte(xmlstrings), &ISO6709)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	if !answer.Equal(ISO6709.S2Region()) {
		t.Errorf("TestLineString %#v", ISO6709.S2Region())
	}
}
