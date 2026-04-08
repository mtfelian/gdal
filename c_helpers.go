package gdal

/*
#include "go_gdal.h"
*/
import "C"
import "unsafe"

// goStringAndCPLFree copies a GDAL-allocated string into Go memory and then
// releases the original buffer with CPLFree.
func goStringAndCPLFree(cstr *C.char) string {
	if cstr == nil {
		return ""
	}
	result := C.GoString(cstr)
	C.CPLFree(unsafe.Pointer(cstr))
	return result
}

// cStringListToSlice converts a null-terminated GDAL string list into a Go
// slice without taking ownership of the individual strings.
func cStringListToSlice(list **C.char) []string {
	if list == nil {
		return nil
	}

	var result []string
	ptr := uintptr(unsafe.Pointer(list))
	for {
		item := *(**C.char)(unsafe.Pointer(ptr))
		if item == nil {
			break
		}
		result = append(result, C.GoString(item))
		ptr += unsafe.Sizeof(list)
	}

	return result
}

// cIntSlicePtr returns a pointer to the first element of data or nil for an
// empty slice.
func cIntSlicePtr(data []C.int) *C.int {
	if len(data) == 0 {
		return nil
	}
	return &data[0]
}

// cDatasetHandleSlicePtr returns a pointer to the first dataset handle or nil
// for an empty slice.
func cDatasetHandleSlicePtr(data []C.GDALDatasetH) *C.GDALDatasetH {
	if len(data) == 0 {
		return nil
	}
	return &data[0]
}

// cGIntBigSlicePtr returns a pointer to the first element of data or nil for
// an empty slice.
func cGIntBigSlicePtr(data []C.GIntBig) *C.GIntBig {
	if len(data) == 0 {
		return nil
	}
	return &data[0]
}

// float64SlicePtr returns a C-compatible pointer to the start of data or nil
// for an empty slice.
func float64SlicePtr(data []float64) *C.double {
	if len(data) == 0 {
		return nil
	}
	return (*C.double)(unsafe.Pointer(&data[0]))
}

// byteSlicePtr returns a raw pointer to the first byte of data or nil for an
// empty slice.
func byteSlicePtr(data []uint8) unsafe.Pointer {
	if len(data) == 0 {
		return nil
	}
	return unsafe.Pointer(&data[0])
}

// uCharSlicePtr returns a C-compatible pointer to the first byte of data or
// nil for an empty slice.
func uCharSlicePtr(data []uint8) *C.uchar {
	if len(data) == 0 {
		return nil
	}
	return (*C.uchar)(unsafe.Pointer(&data[0]))
}

// copyCIntArray copies a C array into Go-owned memory.
func copyCIntArray(data *C.int, count C.int) []int {
	if data == nil || count <= 0 {
		return nil
	}

	raw := unsafe.Slice(data, int(count))
	return CIntSliceToInt(raw)
}

// copyCGIntBigArray copies a C GIntBig array into Go-owned memory.
func copyCGIntBigArray(data *C.GIntBig, count C.int) []int64 {
	if data == nil || count <= 0 {
		return nil
	}

	raw := unsafe.Slice(data, int(count))
	result := make([]int64, len(raw))
	for i := range raw {
		result[i] = int64(raw[i])
	}
	return result
}

// copyCDoubleArray copies a C double array into Go-owned memory.
func copyCDoubleArray(data *C.double, count C.int) []float64 {
	if data == nil || count <= 0 {
		return nil
	}

	raw := unsafe.Slice(data, int(count))
	result := make([]float64, len(raw))
	for i := range raw {
		result[i] = float64(raw[i])
	}
	return result
}

// copyCUCharArray copies a C uchar array into Go-owned memory.
func copyCUCharArray(data *C.uchar, count C.int) []uint8 {
	if data == nil || count <= 0 {
		return nil
	}

	raw := unsafe.Slice(data, int(count))
	result := make([]uint8, len(raw))
	for i := range raw {
		result[i] = uint8(raw[i])
	}
	return result
}

// copyCUIntBigArray copies a C GUIntBig array into Go-owned memory.
func copyCUIntBigArray(data *C.GUIntBig, count C.int) []int {
	if data == nil || count <= 0 {
		return nil
	}

	raw := unsafe.Slice(data, int(count))
	return CUIntBigSliceToInt(raw)
}

// int64SliceToCGIntBig converts a Go slice into the C integer type expected by
// OGR list setters.
func int64SliceToCGIntBig(data []int64) []C.GIntBig {
	result := make([]C.GIntBig, len(data))
	for i := range data {
		result[i] = C.GIntBig(data[i])
	}
	return result
}

// majorObjectFromDataset exposes a dataset as a MajorObject for shared metadata
// helpers in tests.
func majorObjectFromDataset(dataset Dataset) MajorObject {
	return MajorObject{cval: C.GDALMajorObjectH(unsafe.Pointer(dataset.cval))}
}
