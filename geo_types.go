package mongorm

// GeoPoint represents a GeoJSON Point.
type GeoPoint struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

// GeoLineString represents a GeoJSON LineString.
type GeoLineString struct {
	Type        string      `bson:"type" json:"type"`
	Coordinates [][]float64 `bson:"coordinates" json:"coordinates"`
}

// GeoPolygon represents a GeoJSON Polygon.
type GeoPolygon struct {
	Type        string        `bson:"type" json:"type"`
	Coordinates [][][]float64 `bson:"coordinates" json:"coordinates"`
}

func NewGeoPoint(longitude float64, latitude float64) *GeoPoint {
	return &GeoPoint{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
	}
}

func NewGeoLineString(coordinates ...[]float64) *GeoLineString {
	return &GeoLineString{
		Type:        "LineString",
		Coordinates: coordinates,
	}
}

func NewGeoPolygon(rings ...[][]float64) *GeoPolygon {
	return &GeoPolygon{
		Type:        "Polygon",
		Coordinates: rings,
	}
}
