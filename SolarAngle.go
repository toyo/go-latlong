package latlong

import (
	"time"

	"github.com/KlausBrunner/gosolarpos"
)

// SolarAngle returns an solar angle at time t.
func (latlong *Coordinate) SolarAngle(t time.Time) (zenithAngle float64) {
	azimuth, zenithAngle := gosolarpos.Grena3(t,
		latlong.Lat(),
		latlong.Lng(),
		gosolarpos.EstimateDeltaT(t), 1000, 20)
	_ = azimuth
	return
}
