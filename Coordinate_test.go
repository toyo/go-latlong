package latlong

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestString(t *testing.T) {
	l := NewLatLongGridLocator("PM95UQ")

	Config.Lang = "ja"
	ls := l.String()
	correctResponsels := "北緯35.7度、東経139.7度"
	if ls != correctResponsels {
		t.Errorf("expected %+v, was %+v", correctResponsels, ls)
	}

	lp := l.PrecString()
	correctResponselp := "緯度誤差0.041667度、経度誤差0.083333度"
	if lp != correctResponselp {
		t.Errorf("expected %+v, was %+v", correctResponselp, lp)
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(&l)
	correctResponseJSON := "[139.708333,35.687500]\n"
	if err != nil {
		t.Error(err)
		return
	}
	JSON := b.String()

	if JSON != correctResponseJSON {
		t.Errorf("expected '%+v', was '%+v'", correctResponseJSON, JSON)
	}
}
