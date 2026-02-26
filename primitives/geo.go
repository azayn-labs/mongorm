package primitives

import "go.mongodb.org/mongo-driver/v2/bson"

type GeoField struct {
	name string
}

func GeoType(name string) *GeoField {
	return &GeoField{name: name}
}

func (f *GeoField) BSONName() string {
	return f.name
}

func (f *GeoField) Eq(v any) bson.M {
	return bson.M{f.name: v}
}

func (f *GeoField) Ne(v any) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

func (f *GeoField) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

func (f *GeoField) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

func (f *GeoField) IsNull() bson.M {
	return bson.M{f.name: nil}
}

func (f *GeoField) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

func (f *GeoField) Near(geometry any) bson.M {
	return bson.M{
		f.name: bson.M{
			"$near": bson.M{
				"$geometry": geometry,
			},
		},
	}
}

func (f *GeoField) NearWithDistance(geometry any, minDistance *float64, maxDistance *float64) bson.M {
	near := bson.M{"$geometry": geometry}
	if minDistance != nil {
		near["$minDistance"] = *minDistance
	}
	if maxDistance != nil {
		near["$maxDistance"] = *maxDistance
	}

	return bson.M{f.name: bson.M{"$near": near}}
}

func (f *GeoField) NearSphere(geometry any) bson.M {
	return bson.M{
		f.name: bson.M{
			"$nearSphere": bson.M{
				"$geometry": geometry,
			},
		},
	}
}

func (f *GeoField) NearSphereWithDistance(geometry any, minDistance *float64, maxDistance *float64) bson.M {
	near := bson.M{"$geometry": geometry}
	if minDistance != nil {
		near["$minDistance"] = *minDistance
	}
	if maxDistance != nil {
		near["$maxDistance"] = *maxDistance
	}

	return bson.M{f.name: bson.M{"$nearSphere": near}}
}

func (f *GeoField) Within(geometry any) bson.M {
	return bson.M{
		f.name: bson.M{
			"$geoWithin": bson.M{
				"$geometry": geometry,
			},
		},
	}
}

func (f *GeoField) WithinBox(bottomLeft []float64, upperRight []float64) bson.M {
	return bson.M{
		f.name: bson.M{
			"$geoWithin": bson.M{
				"$box": bson.A{bottomLeft, upperRight},
			},
		},
	}
}

func (f *GeoField) WithinCenter(center []float64, radius float64) bson.M {
	return bson.M{
		f.name: bson.M{
			"$geoWithin": bson.M{
				"$center": bson.A{center, radius},
			},
		},
	}
}

func (f *GeoField) WithinCenterSphere(center []float64, radius float64) bson.M {
	return bson.M{
		f.name: bson.M{
			"$geoWithin": bson.M{
				"$centerSphere": bson.A{center, radius},
			},
		},
	}
}

func (f *GeoField) Intersects(geometry any) bson.M {
	return bson.M{
		f.name: bson.M{
			"$geoIntersects": bson.M{
				"$geometry": geometry,
			},
		},
	}
}
