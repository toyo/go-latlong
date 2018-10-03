package latlong_test

import (
	"testing"
	"time"

	latlong "github.com/toyo/go-latlong"
)

func TestSolarAngle(t *testing.T) {
	l := latlong.NewRectGridLocator("PM95UQ").Center()

	sa := l.SolarAngle(time.Unix(0, 0))
	correctResponsesa := 70.38745256769428
	if sa != correctResponsesa {
		t.Errorf("expected %+v, was %+v", correctResponsesa, sa)
	}

}
