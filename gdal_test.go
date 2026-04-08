// Copyright 2011 go-gdal. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdal

import (
	"testing"
)

func TestTiffDriver(t *testing.T) {
	_, err := GetDriverByName(DriverNameGTiff)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestMissingMetadata(t *testing.T) {
	ds, err := Open("testdata/tiles.gpkg", ReadOnly)
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer ds.Close()

	metadata := ds.Metadata("something-that-wont-exist")
	if len(metadata) != 0 {
		t.Errorf("got %d items, want 0", len(metadata))
	}
}

func TestMissingMetadataItem(t *testing.T) {
	ds, err := Open("testdata/demproc.tif", ReadOnly)
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer ds.Close()

	if got := ds.MetadataItem("missing-key", "missing-domain"); got != "" {
		t.Fatalf("dataset metadata item = %q, want empty string", got)
	}

	band := ds.RasterBand(1)
	if got := band.MetadataItem("missing-key", "missing-domain"); got != "" {
		t.Fatalf("raster band metadata item = %q, want empty string", got)
	}

	driver := ds.Driver()
	if got := driver.MetadataItem("missing-key", "missing-domain"); got != "" {
		t.Fatalf("driver metadata item = %q, want empty string", got)
	}
}
