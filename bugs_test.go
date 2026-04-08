package gdal

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func createMemoryRasterDataset(t *testing.T, xSize, ySize, bands int, dataType DataType) Dataset {
	t.Helper()

	driver, err := GetDriverByName(DriverNameMEM)
	if err != nil {
		t.Fatalf("GetDriverByName(MEM): %v", err)
	}

	ds := driver.Create("", xSize, ySize, bands, dataType, nil)
	if ds.cval == nil {
		t.Fatal("MEM driver returned nil dataset")
	}
	return ds
}

func createMemoryVectorLayer(t *testing.T, name string) (DataSource, Layer) {
	t.Helper()

	driver := OGRDriverByName(OGRDriverNameMemory)
	if driver.cval == nil {
		t.Fatal("OGR Memory driver not available")
	}

	ds, ok := driver.Create(name, nil)
	if !ok {
		t.Fatal("failed to create OGR memory datasource")
	}

	layer := ds.CreateLayer(name, SpatialReference{}, GT_Unknown, nil)
	if layer.cval == nil {
		ds.Destroy()
		t.Fatal("failed to create OGR memory layer")
	}

	return ds, layer
}

func createSpatialReferenceFromEPSG(t *testing.T, code int) SpatialReference {
	t.Helper()

	sr := CreateSpatialReference("")
	if err := sr.FromEPSG(code); err != nil {
		sr.Destroy()
		t.Fatalf("FromEPSG(%d): %v", code, err)
	}

	return sr
}

func addLayerField(t *testing.T, layer Layer, name string, fieldType FieldType) {
	t.Helper()

	fd := CreateFieldDefinition(name, fieldType)
	defer fd.Destroy()

	if err := layer.CreateField(fd, false); err != nil {
		t.Fatalf("CreateField(%s): %v", name, err)
	}
}

func TestFeatureListFieldsRoundTrip(t *testing.T) {
	ds, layer := createMemoryVectorLayer(t, "roundtrip")
	defer ds.Destroy()

	addLayerField(t, layer, "ints", FT_IntegerList)
	addLayerField(t, layer, "int64s", FT_Integer64List)
	addLayerField(t, layer, "reals", FT_RealList)
	addLayerField(t, layer, "strings", FT_StringList)
	addLayerField(t, layer, "blob", FT_Binary)

	feature := layer.Definition().Create()
	defer feature.Destroy()

	wantInts := []int{1, 2, 3, 4}
	wantInt64s := []int64{1 << 40, (1 << 40) + 7}
	wantReals := []float64{1.5, 2.5, 3.5}
	wantStrings := []string{"alpha", "beta"}
	wantBinary := []uint8{1, 2, 3, 4, 5}

	feature.SetFieldIntegerList(0, wantInts)
	feature.SetFieldInteger64List(1, wantInt64s)
	feature.SetFieldFloat64List(2, wantReals)
	feature.SetFieldStringList(3, wantStrings)
	feature.SetFieldBinary(4, wantBinary)

	if got := feature.FieldAsIntegerList(0); !reflect.DeepEqual(got, wantInts) {
		t.Fatalf("FieldAsIntegerList = %v, want %v", got, wantInts)
	}
	if got := feature.FieldAsInteger64List(1); !reflect.DeepEqual(got, wantInt64s) {
		t.Fatalf("FieldAsInteger64List = %v, want %v", got, wantInt64s)
	}
	if got := feature.FieldAsFloat64List(2); !reflect.DeepEqual(got, wantReals) {
		t.Fatalf("FieldAsFloat64List = %v, want %v", got, wantReals)
	}
	if got := feature.FieldAsStringList(3); !reflect.DeepEqual(got, wantStrings) {
		t.Fatalf("FieldAsStringList = %v, want %v", got, wantStrings)
	}
	if got := feature.FieldAsBinary(4); !reflect.DeepEqual(got, wantBinary) {
		t.Fatalf("FieldAsBinary = %v, want %v", got, wantBinary)
	}
}

func TestFeatureSetFromWithMapUsesFieldMap(t *testing.T) {
	ds, layer := createMemoryVectorLayer(t, "mapped")
	defer ds.Destroy()

	addLayerField(t, layer, "first", FT_Integer)
	addLayerField(t, layer, "second", FT_Integer)

	fd := layer.Definition()

	source := fd.Create()
	defer source.Destroy()
	source.SetFieldInteger(0, 10)
	source.SetFieldInteger(1, 20)

	target := fd.Create()
	defer target.Destroy()

	if err := target.SetFromWithMap(source, 0, []int{0, 1}); err != nil {
		t.Fatalf("SetFromWithMap: %v", err)
	}

	if got := target.FieldAsInteger(0); got != 10 {
		t.Fatalf("target field 0 = %d, want 10", got)
	}
	if got := target.FieldAsInteger(1); got != 20 {
		t.Fatalf("target field 1 = %d, want 20", got)
	}
}

func TestLayerReorderFieldsReordersDefinitions(t *testing.T) {
	ds, layer := createMemoryVectorLayer(t, "reorder")
	defer ds.Destroy()

	addLayerField(t, layer, "first", FT_Integer)
	addLayerField(t, layer, "second", FT_Integer)

	if err := layer.ReorderFields([]int{1, 0}); err != nil {
		t.Fatalf("ReorderFields: %v", err)
	}

	defn := layer.Definition()
	if got := defn.FieldDefinition(0).Name(); got != "second" {
		t.Fatalf("field 0 name = %q, want %q", got, "second")
	}
	if got := defn.FieldDefinition(1).Name(); got != "first" {
		t.Fatalf("field 1 name = %q, want %q", got, "first")
	}
}

func TestEmptyInputsReturnErrorsInsteadOfPanicking(t *testing.T) {
	if _, err := CreateFromWKB(nil, SpatialReference{}, 0); err == nil {
		t.Fatal("CreateFromWKB(nil) returned nil error")
	}

	var geom Geometry
	if err := geom.FromWKB(nil, 0); err == nil {
		t.Fatal("Geometry.FromWKB(nil) returned nil error")
	}

	if _, err := GridCreate(
		GA_Linear,
		GridLinearOptions{Radius: -1, NoDataValue: 0},
		nil, nil, nil,
		0, 0, 0, 0,
		10, 10,
		nil,
		nil,
	); err == nil {
		t.Fatal("GridCreate(nil, nil, nil) returned nil error")
	}

	if _, err := Warp("", nil, nil, nil); err == nil {
		t.Fatal("Warp with zero source datasets returned nil error")
	}

	if _, err := VectorTranslate("", nil, nil); err == nil {
		t.Fatal("VectorTranslate with zero source datasets returned nil error")
	}

	var ds Dataset
	if err := ds.BuildOverviews("NEAREST", 1, nil, 0, nil, nil, nil); err == nil {
		t.Fatal("BuildOverviews with missing overview list returned nil error")
	}

	var rb RasterBand
	if err := rb.IO(Read, 0, 0, 1, 1, []uint8{}, 1, 1, 0, 0); err == nil {
		t.Fatal("RasterBand.IO with empty buffer returned nil error")
	}
}

func TestCoordinateTransformHandlesZeroPointsAndNilZ(t *testing.T) {
	var zero CoordinateTransform
	if !zero.Transform(0, nil, nil, nil) {
		t.Fatal("Transform with zero points returned false")
	}

	src := CreateSpatialReference("")
	defer src.Destroy()
	if err := src.FromEPSG(4326); err != nil {
		t.Fatalf("src.FromEPSG: %v", err)
	}

	dst := CreateSpatialReference("")
	defer dst.Destroy()
	if err := dst.FromEPSG(3857); err != nil {
		t.Fatalf("dst.FromEPSG: %v", err)
	}

	ct := CreateCoordinateTransform(src, dst)
	defer ct.Destroy()

	x := []float64{37.6173}
	y := []float64{55.7558}
	if !ct.Transform(1, x, y, nil) {
		t.Fatal("Transform with nil z slice returned false")
	}

	if x[0] == 37.6173 || y[0] == 55.7558 {
		t.Fatal("Transform with nil z slice did not transform coordinates")
	}
}

func TestNilStringListsReturnEmptySlices(t *testing.T) {
	raster := createMemoryRasterDataset(t, 8, 8, 1, Byte)
	defer raster.Close()

	if got := raster.FileList(); len(got) != 0 {
		t.Fatalf("FileList length = %d, want 0", len(got))
	}

	if got := raster.RasterBand(1).CategoryNames(); len(got) != 0 {
		t.Fatalf("CategoryNames length = %d, want 0", len(got))
	}

	ds, layer := createMemoryVectorLayer(t, "strings")
	defer ds.Destroy()

	addLayerField(t, layer, "names", FT_StringList)
	feature := layer.Definition().Create()
	defer feature.Destroy()

	if got := feature.FieldAsStringList(0); len(got) != 0 {
		t.Fatalf("FieldAsStringList length = %d, want 0", len(got))
	}
}

func TestSpatialReferenceToPCIRoundTrip(t *testing.T) {
	sr := createSpatialReferenceFromEPSG(t, 32633)
	defer sr.Destroy()

	proj, units, params, err := sr.ToPCI()
	if err != nil {
		t.Fatalf("ToPCI: %v", err)
	}
	if proj == "" {
		t.Fatal("ToPCI returned an empty projection")
	}
	if units == "" {
		t.Fatal("ToPCI returned empty units")
	}
	if len(params) != 17 {
		t.Fatalf("ToPCI params length = %d, want 17", len(params))
	}

	paramsCopy := append([]float64(nil), params...)
	params[0]++
	runtime.GC()

	roundTrip := CreateSpatialReference("")
	defer roundTrip.Destroy()

	if err := roundTrip.FromPCI(proj, units, paramsCopy); err != nil {
		t.Fatalf("FromPCI: %v", err)
	}
	if !sr.IsSame(roundTrip) {
		t.Fatal("PCI round trip changed the spatial reference")
	}
}

func TestSpatialReferenceToUSGSRoundTrip(t *testing.T) {
	sr := createSpatialReferenceFromEPSG(t, 32633)
	defer sr.Destroy()

	proj, zone, params, datum, err := sr.ToUSGS()
	if err != nil {
		t.Fatalf("ToUSGS: %v", err)
	}
	if len(params) != 15 {
		t.Fatalf("ToUSGS params length = %d, want 15", len(params))
	}
	if zone <= 0 {
		t.Fatalf("ToUSGS zone = %d, want > 0", zone)
	}

	paramsCopy := append([]float64(nil), params...)
	params[0]++
	runtime.GC()

	proj2, zone2, params2, datum2, err := sr.ToUSGS()
	if err != nil {
		t.Fatalf("second ToUSGS: %v", err)
	}
	if proj2 != proj || zone2 != zone || datum2 != datum {
		t.Fatalf("second ToUSGS changed scalar output: got (%d, %d, %d), want (%d, %d, %d)", proj2, zone2, datum2, proj, zone, datum)
	}
	if !reflect.DeepEqual(params2, paramsCopy) {
		t.Fatalf("second ToUSGS params = %v, want %v", params2, paramsCopy)
	}

	roundTrip := CreateSpatialReference("")
	defer roundTrip.Destroy()

	if err := roundTrip.FromUSGS(proj, zone, paramsCopy, datum); err != nil {
		t.Fatalf("FromUSGS: %v", err)
	}
	if err := roundTrip.Validate(); err != nil {
		t.Fatalf("USGS round trip produced an invalid spatial reference: %v", err)
	}
	if !roundTrip.IsProjected() {
		t.Fatal("USGS round trip produced a non-projected spatial reference")
	}
}

func TestRasterOutputFormatOptionsRespectExplicitOfFlag(t *testing.T) {
	options := []string{"-of", "GTiff", "-tr", "10", "10"}
	got := ensureRasterOutputFormatOptions(options)
	if !reflect.DeepEqual(got, options) {
		t.Fatalf("ensureRasterOutputFormatOptions modified explicit -of options: got %v, want %v", got, options)
	}

	got = ensureRasterOutputFormatOptions([]string{"-a", "code"})
	want := []string{"-of", "MEM", "-a", "code"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ensureRasterOutputFormatOptions = %v, want %v", got, want)
	}
}

func TestGeometryExportsReturnContent(t *testing.T) {
	geom := Create(GT_Point)
	defer geom.Destroy()
	geom.AddPoint2D(1.25, 2.5)

	wkt, err := geom.ToWKT()
	if err != nil {
		t.Fatalf("ToWKT: %v", err)
	}
	if !strings.Contains(wkt, "POINT") {
		t.Fatalf("ToWKT = %q, want POINT content", wkt)
	}

	if got := geom.ToGML(); got == "" {
		t.Fatal("ToGML returned an empty string")
	}
	if got := geom.ToGML_Ex([]string{"FORMAT=GML3"}); got == "" {
		t.Fatal("ToGML_Ex returned an empty string")
	}
	if got := geom.ToKML(); got == "" {
		t.Fatal("ToKML returned an empty string")
	}
	if got := geom.ToJSON(); !strings.Contains(got, "Point") {
		t.Fatalf("ToJSON = %q, want Point content", got)
	}
	if got := geom.ToJSON_ex([]string{"COORDINATE_PRECISION=1"}); !strings.Contains(got, "Point") {
		t.Fatalf("ToJSON_ex = %q, want Point content", got)
	}
}

func TestDefaultHistogramReturnsUsableGoSlice(t *testing.T) {
	ds, err := Open("testdata/demproc.tif", ReadOnly)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer ds.Close()

	_, _, buckets, histogram, err := ds.RasterBand(1).DefaultHistogram(1, DummyProgress, nil)
	if err != nil {
		t.Fatalf("DefaultHistogram: %v", err)
	}
	if buckets <= 0 {
		t.Fatalf("DefaultHistogram buckets = %d, want > 0", buckets)
	}
	if len(histogram) != buckets {
		t.Fatalf("DefaultHistogram len = %d, want %d", len(histogram), buckets)
	}

	histogram[0]++
	runtime.GC()

	total := 0
	for _, value := range histogram {
		total += value
	}
	if total <= 0 {
		t.Fatal("DefaultHistogram returned an unusable slice")
	}
}

func TestMajorObjectMetadataMethods(t *testing.T) {
	ds := createMemoryRasterDataset(t, 4, 4, 1, Byte)
	defer ds.Close()

	object := majorObjectFromDataset(ds)

	object.SetMetadata([]string{"ALPHA=one"}, "")
	if got := object.MetadataItem("ALPHA", ""); got != "one" {
		t.Fatalf("MetadataItem(ALPHA) = %q, want %q", got, "one")
	}

	object.SetMetadataItem("BETA", "two", "")
	if got := object.MetadataItem("BETA", ""); got != "two" {
		t.Fatalf("MetadataItem(BETA) = %q, want %q", got, "two")
	}

	metadata := object.Metadata("")
	if len(metadata) < 2 {
		t.Fatalf("Metadata length = %d, want at least 2 items", len(metadata))
	}
}
