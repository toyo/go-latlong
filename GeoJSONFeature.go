package latlong

// GeoJSONFeature is Feature of GeoJSON
type GeoJSONFeature struct {
	Type     string          `json:"type"`
	Geometry GeoJSONGeometry `json:"geometry"`
	Property interface{}     `json:"properties"`
}
