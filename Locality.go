package latlong

import (
	"context"

	"googlemaps.github.io/maps"
)

// Locality returns Japanese City, Town, Village name.
func (latlong *Point) Locality(ctx context.Context) (s string, err error) {

	r := &maps.GeocodingRequest{
		ResultType: []string{"locality"}, // "locality", "administrative_area_level_1", "country"
		Language:   Config.Lang,
		LatLng: &maps.LatLng{
			Lat: latlong.Lat.Degrees(),
			Lng: latlong.Lng.Degrees(),
		},
	}

	var c *maps.Client
	if Config.GoogleMapsAPIURL != "" {
		c, err = maps.NewClient(maps.WithAPIKey(Config.GoogleAPIKey), maps.WithBaseURL(Config.GoogleMapsAPIURL))
	} else {
		c, err = maps.NewClient(maps.WithAPIKey(Config.GoogleAPIKey), maps.WithHTTPClient(Config.HTTPClient))
	}
	if err == nil {
		var georesult []maps.GeocodingResult
		if georesult, err = c.ReverseGeocode(ctx, r); err == nil {
			for _, res := range georesult {
				for _, a := range res.AddressComponents {
					for _, t := range a.Types {
						if t == r.ResultType[0] {
							s = a.LongName
							return
						}
					}
				}
			}
		}
	}
	return
}
