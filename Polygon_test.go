package latlong

import (
	"encoding/xml"
	"testing"

	"github.com/golang/geo/r3"
	"github.com/golang/geo/s2"
)

func TestPolygon1(t *testing.T) {

	const xmlstrings = `<ISO6709>+31.9388+130.8800/+31.9388+130.9000/+31.9372+130.9200/+31.9312+130.9400/+31.9265+130.9600/+31.9200+130.9783/+31.9000+130.9682/+31.8958+130.9600/+31.8909+130.9400/+31.8829+130.9200/+31.8812+130.9000/+31.8812+130.8800/+31.9000+130.8612/+31.9200+130.8612/+31.9388+130.8800/</ISO6709>`

	answer := s2.LoopFromPoints([]s2.Point{
		s2.Point{Vector: r3.Vector{X: -0.5553980515538867, Y: 0.6416214741327604, Z: 0.5290131267386939}},
		s2.Point{Vector: r3.Vector{X: -0.5553010232685868, Y: 0.6419349218777474, Z: 0.5287346495461317}},
		s2.Point{Vector: r3.Vector{X: -0.5554217360043626, Y: 0.6420744673087171, Z: 0.5284383347223471}},
		s2.Point{Vector: r3.Vector{X: -0.5557458360705103, Y: 0.6420232508649839, Z: 0.5281597400777982}},
		s2.Point{Vector: r3.Vector{X: -0.555969910600026, Y: 0.6418292198617572, Z: 0.5281597400777982}},
		s2.Point{Vector: r3.Vector{X: -0.5561836526853394, Y: 0.6416232691135639, Z: 0.5281849344856561}},
		s2.Point{Vector: r3.Vector{X: -0.5563592570682993, Y: 0.6413733696975146, Z: 0.5283034901619195}},
		s2.Point{Vector: r3.Vector{X: -0.5565534850090491, Y: 0.641145002960512, Z: 0.5283761004275717}},
		s2.Point{Vector: r3.Vector{X: -0.5566198423172453, Y: 0.6410360969045329, Z: 0.5284383347223471}},
		s2.Point{Vector: r3.Vector{X: -0.5566118367089491, Y: 0.6407986685417031, Z: 0.5287346495461317}},
		s2.Point{Vector: r3.Vector{X: -0.5563678158848891, Y: 0.6409311136679902, Z: 0.5288309379941646}},
		s2.Point{Vector: r3.Vector{X: -0.5561156272518964, Y: 0.6410925121862575, Z: 0.528900557709049}},
		s2.Point{Vector: r3.Vector{X: -0.5558555285186887, Y: 0.6412447392889905, Z: 0.5289894287690482}},
		s2.Point{Vector: r3.Vector{X: -0.5556219858580568, Y: 0.641427564553707, Z: 0.5290131267386939}},
		s2.Point{Vector: r3.Vector{X: -0.5553980515538867, Y: 0.6416214741327604, Z: 0.5290131267386939}}})

	var ISO6709 Polygon

	err := xml.Unmarshal([]byte(xmlstrings), &ISO6709)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	if !answer.Equal(ISO6709.S2Region()) {
		t.Errorf("TestLineString %#v", ISO6709.S2Region())
	}

	if ISO6709.S2Region().ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(31.93, 130.90))) != true {
		t.Error("Something wrong.")
	}

	if ISO6709.S2Region().ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(31.93, 140.90))) != false {
		t.Error("Something wrong.")
	}
}
