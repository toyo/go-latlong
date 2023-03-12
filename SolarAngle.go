package latlong

import (
	"time"

	"github.com/klausbrunner/gosolarpos"
)

// SolarAngle returns an solar angle at time t.
func (latlong Point) SolarAngle(t time.Time) (zenithAngle float64) {
	azimuth, zenithAngle := gosolarpos.Grena3(t,
		latlong.Lat().Degrees(),
		latlong.Lng().Degrees(),
		gosolarpos.EstimateDeltaT(t), 1000, 20)
	_ = azimuth
	return
}
