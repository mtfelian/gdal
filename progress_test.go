package gdal

import (
	"runtime"
	"testing"
)

func createFilledMemoryRasterDataset(t *testing.T, xSize, ySize int) Dataset {
	t.Helper()

	ds := createMemoryRasterDataset(t, xSize, ySize, 1, Byte)

	buffer := make([]uint8, xSize*ySize)
	for i := range buffer {
		buffer[i] = uint8(i % 251)
	}

	if err := ds.RasterBand(1).IO(Write, 0, 0, xSize, ySize, buffer, xSize, ySize, 0, 0); err != nil {
		ds.Close()
		t.Fatalf("RasterBand.IO(Write): %v", err)
	}

	return ds
}

func TestProgressCallbackBridgeNilHandle(t *testing.T) {
	if got := callGoGDALProgressFuncProxy(0, 0.5, "ignored"); got != 1 {
		t.Fatalf("callGoGDALProgressFuncProxy(nil) = %d, want 1", got)
	}
}

func TestProgressCallbackBridgeDispatches(t *testing.T) {
	marker := &struct{ name string }{name: "marker"}

	var (
		called      bool
		gotComplete float64
		gotMessage  string
		gotData     interface{}
	)

	callback := newGoGDALProgressCallback(func(complete float64, message string, progressArg interface{}) int {
		called = true
		gotComplete = complete
		gotMessage = message
		gotData = progressArg
		return 7
	}, marker)
	defer callback.close()

	runtime.GC()

	if got := callGoGDALProgressFuncProxy(uintptr(callback.handle), 0.25, "stage"); got != 7 {
		t.Fatalf("callGoGDALProgressFuncProxy(...) = %d, want 7", got)
	}
	if !called {
		t.Fatal("progress callback was not called")
	}
	if gotComplete != 0.25 {
		t.Fatalf("progress complete = %v, want 0.25", gotComplete)
	}
	if gotMessage != "stage" {
		t.Fatalf("progress message = %q, want %q", gotMessage, "stage")
	}
	if gotData != marker {
		t.Fatalf("progress data = %#v, want %#v", gotData, marker)
	}
}

func TestCopyWholeRasterNilProgress(t *testing.T) {
	src := createFilledMemoryRasterDataset(t, 512, 512)
	defer src.Close()

	dst := createMemoryRasterDataset(t, 512, 512, 1, Byte)
	defer dst.Close()

	runtime.GC()

	if err := src.CopyWholeRaster(dst, nil, nil, nil); err != nil {
		t.Fatalf("CopyWholeRaster(nil progress): %v", err)
	}

	if got, want := dst.RasterBand(1).Checksum(0, 0, 512, 512), src.RasterBand(1).Checksum(0, 0, 512, 512); got != want {
		t.Fatalf("destination checksum = %d, want %d", got, want)
	}
}

func TestCopyWholeRasterProgressCallback(t *testing.T) {
	src := createFilledMemoryRasterDataset(t, 512, 512)
	defer src.Close()

	dst := createMemoryRasterDataset(t, 512, 512, 1, Byte)
	defer dst.Close()

	marker := &struct{ name string }{name: "copy-whole-raster"}
	var (
		calls         int
		invalidData   bool
		invalidAmount bool
	)

	runtime.GC()

	if err := src.CopyWholeRaster(dst, nil, func(complete float64, message string, progressArg interface{}) int {
		calls++
		if progressArg != marker {
			invalidData = true
		}
		if complete < 0 || complete > 1 {
			invalidAmount = true
		}
		return 1
	}, marker); err != nil {
		t.Fatalf("CopyWholeRaster(progress): %v", err)
	}

	if calls == 0 {
		t.Fatal("progress callback was not called")
	}
	if invalidData {
		t.Fatal("progress callback received unexpected data")
	}
	if invalidAmount {
		t.Fatal("progress callback received invalid completion value")
	}
	if got, want := dst.RasterBand(1).Checksum(0, 0, 512, 512), src.RasterBand(1).Checksum(0, 0, 512, 512); got != want {
		t.Fatalf("destination checksum = %d, want %d", got, want)
	}
}
