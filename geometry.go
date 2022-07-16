// Package geojson
// most code borrow from https://github.com/paulmach/go.geojson/blob/master/geometry.go, but this one is for mongodb BSON only, not JSON
package geojson

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// A GeometryType serves to enumerate the different GeoJSON geometry types.
type GeometryType string

// The geometry types supported by GeoJSON 1.0
const (
	GeometryPoint           GeometryType = "Point"
	GeometryMultiPoint      GeometryType = "MultiPoint"
	GeometryLineString      GeometryType = "LineString"
	GeometryMultiLineString GeometryType = "MultiLineString"
	GeometryPolygon         GeometryType = "Polygon"
	GeometryMultiPolygon    GeometryType = "MultiPolygon"
	GeometryCollection      GeometryType = "GeometryCollection"
)

// A Geometry correlates to a GeoJSON geometry object.
type Geometry struct {
	Type GeometryType `bson:"type" json:"type"`

	// Point []float64{longitude, latitude}
	// demo: { type: "Point", coordinates: [ 40, 5 ] }
	Point Point

	// { type: "LineString", coordinates: [ [ 40, 5 ], [ 41, 6 ] ] }
	LineString []Point

	// {
	// 	type: "Polygon",
	// 	coordinates: [ [ [ 0 , 0 ] , [ 3 , 6 ] , [ 6 , 1 ] , [ 0 , 0  ] ] ]
	// }
	Polygon [][]Point

	// {
	// 	type: "MultiPoint",
	// 	coordinates: [
	// [ -73.9580, 40.8003 ],
	// [ -73.9498, 40.7968 ],
	// [ -73.9737, 40.7648 ],
	// [ -73.9814, 40.7681 ]
	// ]
	// }
	MultiPoint []Point

	// {
	// 	type: "MultiLineString",
	// 	coordinates: [
	// [ [ -73.96943, 40.78519 ], [ -73.96082, 40.78095 ] ],
	// [ [ -73.96415, 40.79229 ], [ -73.95544, 40.78854 ] ],
	// [ [ -73.97162, 40.78205 ], [ -73.96374, 40.77715 ] ],
	// [ [ -73.97880, 40.77247 ], [ -73.97036, 40.76811 ] ]
	// ]
	// }
	MultiLineString [][]Point

	// {
	// 	type: "MultiPolygon",
	// 	coordinates: [
	// [ [ [ -73.958, 40.8003 ], [ -73.9498, 40.7968 ], [ -73.9737, 40.7648 ], [ -73.9814, 40.7681 ], [ -73.958, 40.8003 ] ] ],
	// [ [ [ -73.958, 40.8003 ], [ -73.9498, 40.7968 ], [ -73.9737, 40.7648 ], [ -73.958, 40.8003 ] ] ]
	// ]
	// }
	MultiPolygon [][][]Point

	Geometries []*Geometry
}

// Point presents a geometry point, must in format []float64{longitude, latitude}
// https://docs.mongodb.com/v4.2/reference/geojson/#geojson-point
// https://docs.mongodb.com/v4.2/geospatial-queries/#geo-overview-location-data
// warning: do not change any field of this struct
// To specify GeoJSON data, use an embedded document with:
// a field named type that specifies the GeoJSON object type and
// a field named coordinates that specifies the object’s coordinates.
// If specifying latitude and longitude coordinates, list the longitude first and then latitude:
// Valid longitude values are between -180 and 180, both inclusive.
// Valid latitude values are between -90 and 90, both inclusive.
type Point []float64 // []float64{longitude, latitude}

// NewPoint creates and initializes a point geometry with the give coordinate.
func NewPoint(coordinate Point) *Geometry {
	return &Geometry{
		Type:  GeometryPoint,
		Point: coordinate,
	}
}

// NewMultiPoint creates and initializes a multi-point geometry with the given coordinates.
func NewMultiPoint(coordinates ...Point) *Geometry {
	return &Geometry{
		Type:       GeometryMultiPoint,
		MultiPoint: coordinates,
	}
}

// NewLineString creates and initializes a line string geometry with the given coordinates.
func NewLineString(coordinates []Point) *Geometry {
	return &Geometry{
		Type:       GeometryLineString,
		LineString: coordinates,
	}
}

// NewMultiLineString creates and initializes a multi-line string geometry with the given lines.
func NewMultiLineString(lines ...[]Point) *Geometry {
	return &Geometry{
		Type:            GeometryMultiLineString,
		MultiLineString: lines,
	}
}

// NewPolygon creates and initializes a polygon geometry with the given polygon.
func NewPolygon(polygon [][]Point) *Geometry {
	return &Geometry{
		Type:    GeometryPolygon,
		Polygon: polygon,
	}
}

// NewMultiPolygon creates and initializes a multi-polygon geometry with the given polygons.
func NewMultiPolygon(polygons ...[][]Point) *Geometry {
	return &Geometry{
		Type:         GeometryMultiPolygon,
		MultiPolygon: polygons,
	}
}

// NewGeometryCollection creates and initializes a geometry collection geometry with the given geometries.
func NewGeometryCollection(geometries ...*Geometry) *Geometry {
	return &Geometry{
		Type:       GeometryCollection,
		Geometries: geometries,
	}
}

// defining a struct here lets us define the order of the BSON elements.
type geometry struct {
	Type        GeometryType `bson:"type" json:"type,omitempty"`
	Coordinates interface{}  `bson:"coordinates,omitempty" json:"coordinates,omitempty"`
	Geometries  interface{}  `bson:"geometries,omitempty" json:"geometries,omitempty"`
}

func (g *Geometry) toPureGeometry() *geometry {
	geo := &geometry{
		Type: g.Type,
	}

	switch g.Type {
	case GeometryPoint:
		geo.Coordinates = g.Point
	case GeometryMultiPoint:
		geo.Coordinates = g.MultiPoint
	case GeometryLineString:
		geo.Coordinates = g.LineString
	case GeometryMultiLineString:
		geo.Coordinates = g.MultiLineString
	case GeometryPolygon:
		geo.Coordinates = g.Polygon
	case GeometryMultiPolygon:
		geo.Coordinates = g.MultiPolygon
	case GeometryCollection:
		geo.Geometries = g.Geometries
	}
	return geo
}

// MarshalBSON converts the geometry object into the correct BSON.
// MarshalBSON implements bson.Marshaler
// nolint: gocritic
func (g Geometry) MarshalBSON() ([]byte, error) {
	geo := g.toPureGeometry()
	return bson.Marshal(geo)
}

// MarshalJSON for testing purpose
// nolint: gocritic
func (g Geometry) MarshalJSON() ([]byte, error) {
	geo := g.toPureGeometry()
	return bson.MarshalExtJSON(geo, false, false)
}

func UnmarshalGeometryRawJSON(data []byte) (*Geometry, error) {
	g := &Geometry{}
	err := bson.UnmarshalExtJSON(data, true, g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// UnmarshalGeometry decodes the binary BSON data into a GeoJSON geometry.
// Alternately one can call json.Unmarshal(g) directly for the same result.
func UnmarshalGeometry(data []byte) (*Geometry, error) {
	g := &Geometry{}
	err := bson.Unmarshal(data, g)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// UnmarshalBSON decodes the data into a GeoJSON geometry.
// This fulfills the bson.Unmarshaler interface.
func (g *Geometry) UnmarshalBSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	var object map[string]interface{}
	err := bson.Unmarshal(data, &object)
	if err != nil {
		return err
	}

	return decodeGeometry(g, object)
}

// UnmarshalJSON decodes the data into a GeoJSON geometry.
// This fulfills the json.Unmarshaler interface.
func (g *Geometry) UnmarshalJSON(data []byte) error {
	err := bson.UnmarshalExtJSON(data, true, g)
	if err != nil {
		return err
	}

	return nil
}

func decodeGeometry(g *Geometry, object map[string]interface{}) error {
	t, ok := object["type"]
	if !ok {
		return errors.New("type property not defined")
	}

	if s, ok := t.(string); ok {
		g.Type = GeometryType(s)
	} else {
		return errors.New("type property not string")
	}

	var err error

	switch g.Type {
	case GeometryPoint:
		g.Point, err = decodePosition(object["coordinates"])
	case GeometryMultiPoint:
		g.MultiPoint, err = decodePositionSet(object["coordinates"])
	case GeometryLineString:
		g.LineString, err = decodePositionSet(object["coordinates"])
	case GeometryMultiLineString:
		g.MultiLineString, err = decodePathSet(object["coordinates"])
	case GeometryPolygon:
		g.Polygon, err = decodePathSet(object["coordinates"])
	case GeometryMultiPolygon:
		g.MultiPolygon, err = decodePolygonSet(object["coordinates"])
	case GeometryCollection:
		g.Geometries, err = decodeGeometries(object["geometries"])
	}

	return err
}

func decodePosition(data interface{}) (Point, error) {
	coords, ok := data.(primitive.A)
	if !ok {
		return nil, fmt.Errorf("not a valid position, got %v", data)
	}

	result := make(Point, 0, len(coords))
	for _, coord := range coords {
		switch f := coord.(type) {
		case float64:
			result = append(result, f)
		case int:
			result = append(result, float64(f))
		case int32:
			result = append(result, float64(f))
		case int64:
			result = append(result, float64(f))
		default:
			return nil, fmt.Errorf("not a valid coordinate, got %v", coord)
		}
	}

	return result, nil
}

func decodePositionSet(data interface{}) ([]Point, error) {
	points, ok := data.(primitive.A)
	if !ok {
		return nil, fmt.Errorf("not a valid set of positions, got %v", data)
	}

	result := make([]Point, 0, len(points))
	for _, point := range points {
		if p, err := decodePosition(point); err == nil {
			result = append(result, p)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodePathSet(data interface{}) ([][]Point, error) {
	sets, ok := data.(primitive.A)
	if !ok {
		return nil, fmt.Errorf("not a valid path, got %v", data)
	}

	result := make([][]Point, 0, len(sets))

	for _, set := range sets {
		if s, err := decodePositionSet(set); err == nil {
			result = append(result, s)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodePolygonSet(data interface{}) ([][][]Point, error) {
	polygons, ok := data.(primitive.A)
	if !ok {
		return nil, fmt.Errorf("not a valid polygon, got %v", data)
	}

	result := make([][][]Point, 0, len(polygons))
	for _, polygon := range polygons {
		if p, err := decodePathSet(polygon); err == nil {
			result = append(result, p)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodeGeometries(data interface{}) ([]*Geometry, error) {
	if vs, ok := data.(primitive.A); ok {
		geometries := make([]*Geometry, 0, len(vs))
		for _, v := range vs {
			g := &Geometry{}

			vmap, ok := v.(map[string]interface{})
			if !ok {
				break
			}

			err := decodeGeometry(g, vmap)
			if err != nil {
				return nil, err
			}

			geometries = append(geometries, g)
		}

		if len(geometries) == len(vs) {
			return geometries, nil
		}
	}

	return nil, fmt.Errorf("not a valid set of geometries, got %v", data)
}

// IsPoint returns true with the geometry object is a Point type.
func (g *Geometry) IsPoint() bool {
	return g.Type == GeometryPoint
}

// IsMultiPoint returns true with the geometry object is a MultiPoint type.
func (g *Geometry) IsMultiPoint() bool {
	return g.Type == GeometryMultiPoint
}

// IsLineString returns true with the geometry object is a LineString type.
func (g *Geometry) IsLineString() bool {
	return g.Type == GeometryLineString
}

// IsMultiLineString returns true with the geometry object is a LineString type.
func (g *Geometry) IsMultiLineString() bool {
	return g.Type == GeometryMultiLineString
}

// IsPolygon returns true with the geometry object is a Polygon type.
func (g *Geometry) IsPolygon() bool {
	return g.Type == GeometryPolygon
}

// IsMultiPolygon returns true with the geometry object is a MultiPolygon type.
func (g *Geometry) IsMultiPolygon() bool {
	return g.Type == GeometryMultiPolygon
}

// IsCollection returns true with the geometry object is a GeometryCollection type.
func (g *Geometry) IsCollection() bool {
	return g.Type == GeometryCollection
}
