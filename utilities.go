package gdal

/*
#include "go_gdal.h"
#include "gdal_version.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

/* --------------------------------------------- */
/* Command line utility wrapper functions        */
/* --------------------------------------------- */

// stringArrayContains reports whether needle is present in array.
func stringArrayContains(array []string, needle string) bool {
	for _, s := range array {
		if s == needle {
			return true
		}
	}
	return false
}

// prependOptionIfMissing prepends a flag/value pair when the flag is absent.
func prependOptionIfMissing(options []string, flag, value string) []string {
	if stringArrayContains(options, flag) {
		return options
	}
	return append([]string{flag, value}, options...)
}

// ensureRasterOutputFormatOptions defaults command-style options to the MEM
// driver when no explicit output format was provided.
func ensureRasterOutputFormatOptions(options []string) []string {
	return prependOptionIfMissing(options, "-of", "MEM")
}

// ensureVectorOutputFormatOptions defaults command-style options to the Memory
// driver when no explicit output format was provided.
func ensureVectorOutputFormatOptions(options []string) []string {
	return prependOptionIfMissing(options, "-f", "MEM")
}

// Warp wraps the gdalwarp utility API.
//
// When dstDS and destDS are both empty, Warp creates an in-memory dataset and
// injects the MEM output format when the caller did not specify one.
func Warp(dstDS string, destDS *Dataset, sourceDS []Dataset, options []string) (Dataset, error) {
	if len(sourceDS) == 0 {
		return Dataset{}, fmt.Errorf("warp requires at least one source dataset")
	}
	if dstDS == "" && destDS == nil {
		dstDS = "MEM:::"
		options = ensureRasterOutputFormatOptions(options)
	}
	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))
	warpopts := C.GDALWarpAppOptionsNew(
		(**C.char)(unsafe.Pointer(&opts[0])),
		(*C.GDALWarpAppOptionsForBinary)(unsafe.Pointer(nil)))
	defer C.GDALWarpAppOptionsFree(warpopts)

	srcDS := make([]C.GDALDatasetH, len(sourceDS))
	for i, ds := range sourceDS {
		srcDS[i] = ds.cval
	}
	var cerr C.int
	cdstDS := C.CString(dstDS)
	defer C.free(unsafe.Pointer(cdstDS))
	var destDScval C.GDALDatasetH
	if destDS != nil {
		destDScval = destDS.cval
	}
	ds := C.GDALWarp(cdstDS, destDScval,
		C.int(len(sourceDS)),
		cDatasetHandleSlicePtr(srcDS),
		warpopts, &cerr)
	if cerr != 0 {
		return Dataset{}, fmt.Errorf("warp failed with code %d", cerr)
	}
	return Dataset{ds}, nil
}

// Translate wraps the gdal_translate utility API.
//
// When dstDS is empty, Translate creates an in-memory dataset and injects the
// MEM output format when the caller did not specify one.
func Translate(dstDS string, sourceDS Dataset, options []string) (Dataset, error) {
	if dstDS == "" {
		dstDS = "MEM:::"
		options = ensureRasterOutputFormatOptions(options)
	}
	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))
	translateopts := C.GDALTranslateOptionsNew(
		(**C.char)(unsafe.Pointer(&opts[0])),
		(*C.GDALTranslateOptionsForBinary)(unsafe.Pointer(nil)))
	defer C.GDALTranslateOptionsFree(translateopts)

	var cerr C.int
	cdstDS := C.CString(dstDS)
	defer C.free(unsafe.Pointer(cdstDS))
	ds := C.GDALTranslate(cdstDS,
		sourceDS.cval,
		translateopts, &cerr)
	if cerr != 0 {
		return Dataset{}, fmt.Errorf("translate failed with code %d", cerr)
	}
	return Dataset{ds}, nil
}

// VectorTranslate wraps the ogr2ogr-style vector translation API.
//
// When dstDS is empty, VectorTranslate creates an in-memory dataset and
// injects the Memory output format when the caller did not specify one.
func VectorTranslate(dstDS string, sourceDS []Dataset, options []string) (Dataset, error) {
	if len(sourceDS) == 0 {
		return Dataset{}, fmt.Errorf("vector translate requires at least one source dataset")
	}
	if dstDS == "" {
		dstDS = "MEM:::"
		options = ensureVectorOutputFormatOptions(options)
	}
	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))
	translateopts := C.GDALVectorTranslateOptionsNew(
		(**C.char)(unsafe.Pointer(&opts[0])),
		(*C.GDALVectorTranslateOptionsForBinary)(unsafe.Pointer(nil)))
	defer C.GDALVectorTranslateOptionsFree(translateopts)

	srcDS := make([]C.GDALDatasetH, len(sourceDS))
	for i, ds := range sourceDS {
		srcDS[i] = ds.cval
	}

	var cerr C.int
	cdstDS := C.CString(dstDS)
	defer C.free(unsafe.Pointer(cdstDS))
	ds := C.GDALVectorTranslate(cdstDS, nil,
		C.int(len(sourceDS)),
		cDatasetHandleSlicePtr(srcDS),
		translateopts, &cerr)
	if cerr != 0 {
		return Dataset{}, fmt.Errorf("vector translate failed with code %d", cerr)
	}
	return Dataset{ds}, nil
}

// Rasterize wraps the gdal_rasterize utility API.
//
// When dstDS is empty, Rasterize creates an in-memory dataset and injects the
// MEM output format when the caller did not specify one.
func Rasterize(dstDS string, sourceDS Dataset, options []string) (Dataset, error) {
	if dstDS == "" {
		dstDS = "MEM:::"
		options = ensureRasterOutputFormatOptions(options)
	}
	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))
	rasterizeopts := C.GDALRasterizeOptionsNew(
		(**C.char)(unsafe.Pointer(&opts[0])),
		(*C.GDALRasterizeOptionsForBinary)(unsafe.Pointer(nil)))
	defer C.GDALRasterizeOptionsFree(rasterizeopts)

	var cerr C.int
	cdstDS := C.CString(dstDS)
	defer C.free(unsafe.Pointer(cdstDS))
	ds := C.GDALRasterize(cdstDS, nil,
		sourceDS.cval,
		rasterizeopts, &cerr)
	if cerr != 0 {
		return Dataset{}, fmt.Errorf("rasterize failed with code %d", cerr)
	}
	return Dataset{ds}, nil
}

// DEMProcessing wraps the gdaldem utility API.
//
// When dstDS is empty, DEMProcessing creates an in-memory dataset and injects
// the MEM output format when the caller did not specify one.
func DEMProcessing(dstDS string, sourceDS Dataset, processing string, colorFileName string, options []string) (Dataset, error) {
	if dstDS == "" {
		dstDS = "MEM:::"
		options = ensureRasterOutputFormatOptions(options)
	}
	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))
	demprocessingopts := C.GDALDEMProcessingOptionsNew(
		(**C.char)(unsafe.Pointer(&opts[0])),
		(*C.GDALDEMProcessingOptionsForBinary)(unsafe.Pointer(nil)))
	defer C.GDALDEMProcessingOptionsFree(demprocessingopts)

	var cerr C.int
	cdstDS := C.CString(dstDS)
	defer C.free(unsafe.Pointer(cdstDS))

	cprocessing := C.CString(processing)
	defer C.free(unsafe.Pointer(cprocessing))
	ccolorFileName := C.CString(colorFileName)
	defer C.free(unsafe.Pointer(ccolorFileName))
	ds := C.GDALDEMProcessing(cdstDS,
		sourceDS.cval,
		cprocessing,
		ccolorFileName,
		demprocessingopts,
		&cerr)
	if cerr != 0 {
		return Dataset{}, fmt.Errorf("demprocessing failed with code %d", cerr)
	}
	return Dataset{ds}, nil
}

// ContourGenerate wraps GDALContourGenerateEx.
func ContourGenerate(band RasterBand, layer Layer, options []string, progress ProgressFunc, data interface{}) error {
	var opts **C.char
	if options != nil {
		o := make([]*C.char, len(options)+1)
		for i := 0; i < len(options); i++ {
			o[i] = C.CString(options[i])
			defer C.free(unsafe.Pointer(o[i]))
		}
		o[len(options)] = (*C.char)(unsafe.Pointer(nil))
		opts = (**C.char)(unsafe.Pointer(&o[0]))
	}

	callback := newGoGDALProgressCallback(progress, data)
	defer callback.close()

	return ErrFromCPLErr(
		C.GDALContourGenerateEx(band.cval,
			unsafe.Pointer(layer.cval),
			opts,
			callback.fn,
			callback.arg,
		))
}
