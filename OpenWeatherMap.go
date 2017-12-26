package latlong

import (
	owm "github.com/briandowns/openweathermap"
)

// CurrentWeatherData return a pointer of struct for CurrentWeatherData
func (ll *LatLong) CurrentWeatherData(unit, lang, apikey string) *owm.CurrentWeatherData {
	w, err := owm.NewCurrent(unit, lang, apikey)
	if err != nil {
		panic(err)
	}

	w.CurrentByCoordinates(
		&owm.Coordinates{
			Longitude: ll.Lng(),
			Latitude:  ll.Lat(),
		},
	)

	return w
}

// CurrentPressure return a GrndLevel Pressure of the point.
func (ll *LatLong) CurrentPressure(unit, lang, apikey string) float64 {
	w := ll.CurrentWeatherData(unit, lang, apikey)

	return w.Main.GrndLevel
}
