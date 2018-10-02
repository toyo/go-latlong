package latlong

import (
	"strconv"

	"github.com/golang/geo/s1"
)

const (
	radiusKmOfTheEarth = 6371.01 // radius km of the Earth.
)

// Km is kilo-meter.
type Km float64

// EarthAngle makes Distance to Angle.
func (km Km) EarthAngle() s1.Angle {
	return s1.Angle(km / radiusKmOfTheEarth)
}

// EarthArcFromAngle makes Angle to Distance.
func EarthArcFromAngle(angle s1.Angle) Km {
	return Km(angle * radiusKmOfTheEarth)
}

// EarthChordAngle makes Distance to ChordAngle.
func (km Km) EarthChordAngle() s1.ChordAngle {
	return s1.ChordAngleFromAngle(km.EarthAngle())
}

// EarthArcFromChordAngle makes ChordAngle to Distance.
func EarthArcFromChordAngle(chordangle s1.ChordAngle) Km {
	return EarthArcFromAngle(chordangle.Angle())
}

func (km Km) String() string {
	return strconv.FormatFloat(float64(km), 'f', 1, 64) + "km"
}
