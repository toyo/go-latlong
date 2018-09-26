package latlong

import (
	"github.com/golang/geo/s1"
)

const (
	circumferenceKmOfTheEarth = 40075 // km
	kmByChordAngle            = circumferenceKmOfTheEarth / 8
)

// Km is kilo-meter.
type Km float64

// EarthChordAngle makes Distance to ChordAngle.
func (km Km) EarthChordAngle() s1.ChordAngle {
	return s1.ChordAngleFromSquaredLength(float64(km) / kmByChordAngle)
}

// EarthAngle makes Distance to Angle.
func (km Km) EarthAngle() s1.Angle {
	return km.EarthChordAngle().Angle()
}

// EarthArcFromChordAngle makes ChordAngle to Distance.
func EarthArcFromChordAngle(chordangle s1.ChordAngle) Km {
	return Km(chordangle * kmByChordAngle)
}

// EarthArcFromAngle makes Angle to Distance.
func EarthArcFromAngle(angle s1.Angle) Km {
	return EarthArcFromChordAngle(s1.ChordAngleFromAngle(angle))
}
