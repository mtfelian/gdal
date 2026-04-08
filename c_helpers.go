package gdal

/*
#include "go_gdal.h"
*/
import "C"
import "unsafe"

func goStringAndCPLFree(cstr *C.char) string {
	if cstr == nil {
		return ""
	}
	result := C.GoString(cstr)
	C.CPLFree(unsafe.Pointer(cstr))
	return result
}

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

func cIntSlicePtr(data []C.int) *C.int {
	if len(data) == 0 {
		return nil
	}
	return &data[0]
}

func cDatasetHandleSlicePtr(data []C.GDALDatasetH) *C.GDALDatasetH {
	if len(data) == 0 {
		return nil
	}
	return &data[0]
}

func cGIntBigSlicePtr(data []C.GIntBig) *C.GIntBig {
	if len(data) == 0 {
		return nil
	}
	return &data[0]
}

func float64SlicePtr(data []float64) *C.double {
	if len(data) == 0 {
		return nil
	}
	return (*C.double)(unsafe.Pointer(&data[0]))
}

func byteSlicePtr(data []uint8) unsafe.Pointer {
	if len(data) == 0 {
		return nil
	}
	return unsafe.Pointer(&data[0])
}

func uCharSlicePtr(data []uint8) *C.uchar {
	if len(data) == 0 {
		return nil
	}
	return (*C.uchar)(unsafe.Pointer(&data[0]))
}

func copyCIntArray(data *C.int, count C.int) []int {
	if data == nil || count <= 0 {
		return nil
	}

	raw := unsafe.Slice(data, int(count))
	return CIntSliceToInt(raw)
}

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

func copyCUIntBigArray(data *C.GUIntBig, count C.int) []int {
	if data == nil || count <= 0 {
		return nil
	}

	raw := unsafe.Slice(data, int(count))
	return CUIntBigSliceToInt(raw)
}

func int64SliceToCGIntBig(data []int64) []C.GIntBig {
	result := make([]C.GIntBig, len(data))
	for i := range data {
		result[i] = C.GIntBig(data[i])
	}
	return result
}

func majorObjectFromDataset(dataset Dataset) MajorObject {
	return MajorObject{cval: C.GDALMajorObjectH(unsafe.Pointer(dataset.cval))}
}
