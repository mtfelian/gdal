/*
Package gdal wraps GDAL, the Geospatial Data Abstraction Library.

GDAL provides access to a large number of geospatial raster data formats.
This package also wraps the related OGR Simple Feature Library, which
provides similar functionality for vector formats.

# Limitations

Some less commonly used functions are not yet implemented. Most of the
missing pieces involve style tables, asynchronous I/O, and GCPs.

The documentation is limited, but the exposed functionality closely follows
the GDAL C API.

This wrapper was originally tested on Windows 7 with the MinGW32_x64
compiler and GDAL 1.11.

# Usage

A simple program that creates a georeferenced blank 256x256 GeoTIFF:

	package main

	import (
		"fmt"
		"flag"
		gdal "github.com/mtfelian/gdal"
	)

	func main() {
		flag.Parse()
		filename := flag.Arg(0)
		if filename == "" {
			fmt.Printf("Usage: tiff [filename]\n")
			return
		}
		buffer := make([]uint8, 256 * 256)

		driver, err := gdal.GetDriverByName(DriverNameGTiff)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		dataset := driver.Create(filename, 256, 256, 1, gdal.Byte, nil)
		defer dataset.Close()

		spatialRef := gdal.CreateSpatialReference("")
		spatialRef.FromEPSG(3857)
		srString, err := spatialRef.ToWKT()
		dataset.SetProjection(srString)
		dataset.SetGeoTransform([]float64{444720, 30, 0, 3751320, 0, -30})
		raster := dataset.RasterBand(1)
		raster.IO(gdal.Write, 0, 0, 256, 256, buffer, 256, 256, 0, 0)
	}

More examples are available in the ./examples subdirectory.
*/
package gdal
