package latlong

import (
	"math"
	"strconv"

	"github.com/golang/geo/s1"
)

const (
	radiusKmOfTheEarth        = 6371.01                          // radius km of the Earth.
	circumferenceKmOfTheEarth = radiusKmOfTheEarth * 2 * math.Pi // circumference km of the Earth.
	kmByChordAngle            = circumferenceKmOfTheEarth / 8    //
)

// Km is kilo-meter.
type Km float64

// EarthAngle makes Distance to Angle.
func (km Km) EarthAngle() s1.Angle {
	return s1.Angle(km / radiusKmOfTheEarth)
}

// EarthArcFromChordAngle makes ChordAngle to Distance.
func EarthArcFromChordAngle(chordangle s1.ChordAngle) Km {
	return Km(chordangle * kmByChordAngle)
}

// EarthArcFromAngle makes Angle to Distance.
func EarthArcFromAngle(angle s1.Angle) Km {
	return Km(angle * radiusKmOfTheEarth)
}

func (km Km) String() string {
	return strconv.FormatFloat(float64(km), 'f', 0, 64) + "km"
}
