package latlong

// GeoJSONFeatureCollection is FeatureCollection of GeoJSON
type GeoJSONFeatureCollection struct {
	Type     string           `json:"type"`
	Features []GeoJSONFeature `json:"features"`
}

// NewGeoJSONFeatureCollection creates GeoJSONFeatureCollection
func NewGeoJSONFeatureCollection() *GeoJSONFeatureCollection {
	var g GeoJSONFeatureCollection
	g.Type = "FeatureCollection"
	return &g
}

// AddFeature adds GeoJSONFeature
func (g *GeoJSONFeatureCollection) AddFeature(f GeoJSONFeature) *GeoJSONFeatureCollection {
	gg := NewGeoJSONFeatureCollection()
	gg.Features = append(gg.Features, f)
	return gg
}
