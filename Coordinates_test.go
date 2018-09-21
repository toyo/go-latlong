package latlong

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/golang/geo/s2"
)

func TestNewLatLongAltsString(t *testing.T) {
	//str := `+32.9+130.9-10000/`
	//str := `+30.4402+130.2000/+30.4588+130.2200/+30.4545+130.2400/+30.4400+130.2545/+30.4200+130.2489/+30.4061+130.2400/+30.4037+130.2200/+30.4200+130.2056/+30.4398+130.2000/+30.4400+130.1998/+30.4402+130.2000/`
	ll, err := NewLatLongsISO6709("+12.34+123.45+3776/+0123.4-01234.5-3776/-001234+0012345-12345/")
	if err != nil {
		t.Errorf("NewLatLongAltssISO6709 returned non nil error: %v", err)
	}

	Config.Lang = "ja"
	lls := ll.String()
	correctResponsells := "北緯12.34度、東経123.45度、標高3776m,北緯1.39度、西経12.57度、ごく浅く,南緯0.209度、東経1.396度、深さ12km"
	if lls != correctResponsells {
		t.Errorf("expected %+v, was %+v", correctResponsells, lls)
	}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(&ll)
	correctResponseJSON := "[[123.45,12.34,3776],[-12.57,1.39,-3776],[1.396,-0.209,-12345]]\n"
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
	ll, err := NewLatLongsISO6709("+12+123/+12.3+123.4/+12.34+123.43/")
	if err != nil {
		t.Errorf("NewLatLongsISO6709 returned non nil error: %v", err)
	}

	Config.Lang = "ja"
	lls := ll.String()
	correctResponsells := "北緯12.34度、東経123.43度"
	if lls != correctResponsells {
		t.Errorf("expected %+v, was %+v", correctResponsells, lls)
		for _, l := range *ll {
			t.Error(l.String())
			t.Error(l.PrecString())
		}
	}
}

func TestPolygon(t *testing.T) {
	cs, err := NewLatLongsISO6709("+31.9388+130.8800/+31.9388+130.9000/+31.9372+130.9200/+31.9312+130.9400/+31.9265+130.9600/+31.9200+130.9783/+31.9000+130.9682/+31.8958+130.9600/+31.8909+130.9400/+31.8829+130.9200/+31.8812+130.9000/+31.8812+130.8800/+31.9000+130.8612/+31.9200+130.8612/+31.9388+130.8800/")
	if err != nil {
		t.Errorf("%v", err)
	}

	lo := cs.S2Loop()

	if lo.ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(31.93, 130.90))) != true {
		t.Error("Something wrong.")
	}

	if lo.ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(31.93, 140.90))) != false {
		t.Error("Something wrong.")
	}
}
