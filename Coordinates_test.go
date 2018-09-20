package latlong

import (
	"testing"
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
