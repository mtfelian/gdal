package gdal

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := os.MkdirAll("./tmp", 0777); err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestVectorTranslate(t *testing.T) {
	srcDS, err := OpenEx("testdata/test.shp", OFReadOnly, nil, nil, nil)
	if err != nil {
		t.Errorf("Open: %v", err)
	}
	defer srcDS.Close()

	opts := []string{"-t_srs", "epsg:4326", "-f", "GeoJSON"}

	dstDS, err := VectorTranslate("./tmp/test4326.geojson", []Dataset{srcDS}, opts)
	if err != nil {
		t.Errorf("Warp: %v", err)
	}
	dstDS.Close()
	dstDS, err = OpenEx("./tmp/test4326.geojson", OFReadOnly|OFVector, []string{"geojson"}, nil, nil)
	if err != nil {
		t.Errorf("Open after translate: %v", err)
	}
	dstDS.Close()

}
func TestRasterize(t *testing.T) {
	srcDS, err := OpenEx("testdata/test.shp", OFReadOnly, nil, nil, nil)
	if err != nil {
		t.Errorf("Open: %v", err)
	}
	defer srcDS.Close()

	opts := []string{"-a", "code", "-tr", "10", "10"}

	dstDS, err := Rasterize("./tmp/rasterized.tif", srcDS, opts)
	if err != nil {
		t.Errorf("Warp: %v", err)
	}
	dstDS.Close()
	dstDS, err = Open("./tmp/rasterized.tif", ReadOnly)
	if err != nil {
		t.Errorf("Open after vector translate: %v", err)
	}
	dstDS.Close()
}

func TestWarp(t *testing.T) {
	srcDS, err := Open("testdata/tiles.gpkg", ReadOnly)
	if err != nil {
		t.Errorf("Open: %v", err)
	}
	defer srcDS.Close()

	opts := []string{"-t_srs", "epsg:3857", "-of", "GPKG"}

	dstDS, err := Warp("./tmp/tiles-3857.gpkg", nil, []Dataset{srcDS}, opts)
	if err != nil {
		t.Errorf("Warp: %v", err)
	}

	pngdriver, err := GetDriverByName("PNG")
	pngdriver.CreateCopy("./tmp/foo.png", dstDS, 0, nil, nil, nil)
	dstDS.Close()
}

func TestTranslate(t *testing.T) {
	srcDS, err := Open("testdata/tiles.gpkg", ReadOnly)
	if err != nil {
		t.Errorf("Open: %v", err)
	}
	defer srcDS.Close()

	opts := []string{"-of", "GTiff"}

	dstDS, err := Translate("./tmp/tiles.tif", srcDS, opts)
	if err != nil {
		t.Errorf("Warp: %v", err)
	}
	dstDS.Close()

	dstDS, err = Open("./tmp/tiles.tif", ReadOnly)
	if err != nil {
		t.Errorf("Open after raster translate: %v", err)
	}
	dstDS.Close()
}

func TestDEMProcessing(t *testing.T) {
	srcDS, err := Open("testdata/demproc.tif", ReadOnly)
	if err != nil {
		t.Errorf("Open: %v", err)
	}
	defer srcDS.Close()

	opts := []string{"-of", "GTiff"}

	dstDS, err := DEMProcessing("./tmp/demproc_output.tif", srcDS, "color-relief", "testdata/demproc_colors.txt", opts)
	if err != nil {
		t.Errorf("DEMProcessing: %v", err)
	}
	dstDS.Close()

	dstDS, err = Open("./tmp/demproc_output.tif", ReadOnly)
	if err != nil {
		t.Errorf("Open after raster DEM Processing: %v", err)
	}
	dstDS.Close()
}

func TestGenerateContours(t *testing.T) {
	srcDS, err := Open("testdata/demproc.tif", ReadOnly)
	if err != nil {
		t.Errorf("Open: %v", err)
	}
	defer srcDS.Close()

	band := srcDS.RasterBand(1)
	dstFileName := "./tmp/contours.json"
	ogrDriver := OGRDriverByName("GeoJSON")
	drv, ok := ogrDriver.Create(dstFileName, nil)
	if !ok {
		t.Errorf("ogrDriver.Create: %v", ok)
	}
	defer func() {
		if err := drv.Sync(); err != nil {
			t.Errorf("ds.Sync: %v", err)
		}
		drv.Destroy()
	}()
	sr := CreateSpatialReference("")
	if err := sr.FromEPSG(4326); err != nil {
		t.Errorf("sr.FromEPSG: %v", err)
	}
	layer := drv.CreateLayer("contour", sr, GT_Unknown, []string{})

	if err := ContourGenerate(band, layer, []string{"LEVEL_INTERVAL=20"}, DummyProgress, nil); err != nil {
		t.Errorf("ContourGenerate: %v", err)
	}
}
