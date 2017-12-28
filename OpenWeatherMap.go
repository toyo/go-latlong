package latlong

import (
	owm "github.com/briandowns/openweathermap"
)

// CurrentWeatherData return a pointer of struct for CurrentWeatherData
func (ll *LatLong) CurrentWeatherData(unit, lang, apikey string) (w *owm.CurrentWeatherData, err error) {
	if w, err = owm.NewCurrent(unit, lang, apikey); err == nil {
		if err = w.CurrentByCoordinates(&owm.Coordinates{Longitude: ll.Lng(), Latitude: ll.Lat()}); err == nil {
			return
		}
	}
	return
}

// CurrentPressure return a GrndLevel Pressure of the point.
func (ll *LatLong) CurrentPressure(unit, lang, apikey string) float64 {
	w, _ := ll.CurrentWeatherData(unit, lang, apikey)

	return w.Main.GrndLevel
}
