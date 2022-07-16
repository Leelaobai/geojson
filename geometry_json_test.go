package geojson

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestGeometryMarshalJSONPoint_JSON(t *testing.T) {
	g := NewPoint([]float64{1, 2})
	blob, err := g.MarshalJSON()

	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"type":"Point"`)) {
		t.Errorf("json should have type Point")
	}

	if !bytes.Contains(blob, []byte(`"coordinates":[1.0,2.0]`)) {
		t.Errorf("json should marshal coordinates correctly, blob=%s", blob)
	}
}

func TestGeometryMarshalPoint_JSON(t *testing.T) {
	g := NewPoint([]float64{1, 2})
	blob, err := json.Marshal(g)

	if err != nil {
		t.Fatalf("should json.Marshal just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"type":"Point"`)) {
		t.Errorf("json should have type Point")
	}

	if !bytes.Contains(blob, []byte(`"coordinates":[1.0,2.0]`)) {
		t.Errorf("json should marshal coordinates correctly")
	}
}

func TestGeometryMarshalPointValue_JSON(t *testing.T) {
	g := NewPoint([]float64{1, 2})
	blob, err := json.Marshal(*g)

	if err != nil {
		t.Fatalf("should json.Marshal just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"type":"Point"`)) {
		t.Errorf("json should have type Point")
	}

	if !bytes.Contains(blob, []byte(`"coordinates":[1.0,2.0]`)) {
		t.Errorf("json should marshal coordinates correctly")
	}
}

func TestGeometryMarshalJSONMultiPoint_JSON(t *testing.T) {
	g := NewMultiPoint([]float64{1, 2}, []float64{3, 4})
	blob, err := g.MarshalJSON()

	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"type":"MultiPoint"`)) {
		t.Errorf("json should have type MultiPoint")
	}

	if !bytes.Contains(blob, []byte(`"coordinates":[[1.0,2.0],[3.0,4.0]]`)) {
		t.Errorf("json should marshal coordinates correctly")
	}
}

func TestGeometryMarshalJSONLineString_JSON(t *testing.T) {
	g := NewLineString([]Point{{1, 2}, {3, 4}})
	blob, err := g.MarshalJSON()

	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"type":"LineString"`)) {
		t.Errorf("json should have type LineString")
	}

	if !bytes.Contains(blob, []byte(`"coordinates":[[1.0,2.0],[3.0,4.0]]`)) {
		t.Errorf("json should marshal coordinates correctly")
	}
}

func TestGeometryMarshalJSONMultiLineString_JSON(t *testing.T) {
	g := NewMultiLineString(
		[]Point{{1, 2}, {3, 4}},
		[]Point{{5, 6}, {7, 8}},
	)
	blob, err := g.MarshalJSON()

	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"type":"MultiLineString"`)) {
		t.Errorf("json should have type MultiLineString")
	}

	if !bytes.Contains(blob, []byte(`"coordinates":[[[1.0,2.0],[3.0,4.0]],[[5.0,6.0],[7.0,8.0]]]`)) {
		t.Errorf("json should marshal coordinates correctly")
	}
}

func TestGeometryMarshalJSONPolygon_JSON(t *testing.T) {
	g := NewPolygon([][]Point{
		{{1, 2}, {3, 4}},
		{{5, 6}, {7, 8}},
	})
	blob, err := g.MarshalJSON()

	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"type":"Polygon"`)) {
		t.Errorf("json should have type Polygon")
	}

	if !bytes.Contains(blob, []byte(`"coordinates":[[[1.0,2.0],[3.0,4.0]],[[5.0,6.0],[7.0,8.0]]]`)) {
		t.Errorf("json should marshal coordinates correctly")
	}
}

func TestGeometryMarshalJSONMultiPolygon_JSON(t *testing.T) {
	g := NewMultiPolygon(
		[][]Point{
			{{1, 2}, {3, 4}},
			{{5, 6}, {7, 8}},
		},
		[][]Point{
			{{8, 7}, {6, 5}},
			{{4, 3}, {2, 1}},
		})
	blob, err := g.MarshalJSON()

	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"type":"MultiPolygon"`)) {
		t.Errorf("json should have type MultiPolygon")
	}

	if !bytes.Contains(blob, []byte(`"coordinates":[[[[1.0,2.0],[3.0,4.0]],[[5.0,6.0],[7.0,8.0]]],[[[8.0,7.0],[6.0,5.0]],[[4.0,3.0],[2.0,1.0]]]]`)) {
		t.Errorf("json should marshal coordinates correctly")
	}
}

func TestUnmarshalGeometryPoint_JSON(t *testing.T) {
	rawJSON := `{"type": "Point", "coordinates": [102.0, 0.5]}`

	g, err := UnmarshalGeometryRawJSON([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal geometry without issue, err %v", err)
	}

	if g.Type != "Point" {
		t.Errorf("incorrect type, got %v", g.Type)
	}

	if len(g.Point) != 2 {
		t.Errorf("should have 2 coordinate elements but got %d", len(g.Point))
	}
}

func TestUnmarshalGeometryMultiPoint_JSON(t *testing.T) {
	rawJSON := `{"type": "MultiPoint", "coordinates": [[1,2],[3,4]]}`

	g, err := UnmarshalGeometryRawJSON([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal geometry without issue, err %v", err)
	}

	if g.Type != "MultiPoint" {
		t.Errorf("incorrect type, got %v", g.Type)
	}

	if len(g.MultiPoint) != 2 {
		t.Errorf("should have 2 coordinate elements but got %d", len(g.MultiPoint))
	}
}

func TestUnmarshalGeometryLineString_JSON(t *testing.T) {
	rawJSON := `{"type": "LineString", "coordinates": [[1,2],[3,4]]}`

	g, err := UnmarshalGeometryRawJSON([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal geometry without issue, err %v", err)
	}

	if g.Type != "LineString" {
		t.Errorf("incorrect type, got %v", g.Type)
	}

	if len(g.LineString) != 2 {
		t.Errorf("should have 2 line string coordinates but got %d", len(g.LineString))
	}
}

func TestUnmarshalGeometryMultiLineString_JSON(t *testing.T) {
	rawJSON := `{"type": "MultiLineString", "coordinates": [[[1,2],[3,4]],[[5,6],[7,8]]]}`

	g, err := UnmarshalGeometryRawJSON([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal geometry without issue, err %v", err)
	}

	if g.Type != "MultiLineString" {
		t.Errorf("incorrect type, got %v", g.Type)
	}

	if len(g.MultiLineString) != 2 {
		t.Errorf("should have 2 line strings but got %d", len(g.MultiLineString))
	}
}

func TestUnmarshalGeometryPolygon_JSON(t *testing.T) {
	rawJSON := `{"type": "Polygon", "coordinates": [[[1,2],[3,4]],[[5,6],[7,8]]]}`

	g, err := UnmarshalGeometryRawJSON([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal geometry without issue, err %v", err)
	}

	if g.Type != "Polygon" {
		t.Errorf("incorrect type, got %v", g.Type)
	}

	if len(g.Polygon) != 2 {
		t.Errorf("should have 2 polygon paths but got %d", len(g.Polygon))
	}
}

func TestUnmarshalGeometryMultiPolygon_JSON(t *testing.T) {
	rawJSON := `{"type": "MultiPolygon", "coordinates": [[[[1,2],[3,4]],[[5,6],[7,8]]],[[[8,7],[6,5]],[[4,3],[2,1]]]]}`

	g, err := UnmarshalGeometryRawJSON([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal geometry without issue, err %v", err)
	}

	if g.Type != "MultiPolygon" {
		t.Errorf("incorrect type, got %v", g.Type)
	}

	if len(g.MultiPolygon) != 2 {
		t.Errorf("should have 2 polygons but got %d", len(g.MultiPolygon))
	}
}
