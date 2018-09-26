package latlong

import (
	"math"

	"github.com/golang/geo/s1"
)

const (
	circumferenceKmOfTheEarth = 40075 // km
	kmByRadian                = circumferenceKmOfTheEarth / (2 * math.Pi)
)

// EarthArcKmToAngle makes Distance to Angle.
func EarthArcKmToAngle(km float64) s1.Angle {
	return s1.Angle(km / kmByRadian)
}

// AngleToEarthArcKm makes Angle to Distance.
func AngleToEarthArcKm(angle s1.Angle) float64 {
	return angle.Radians() * kmByRadian
}
