package gdal

/*
#include "go_gdal.h"
#include "gdal_version.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

/* --------------------------------------------- */
/* Misc functions                                */
/* --------------------------------------------- */

// ComputeMedianCutPCT computes an optimal pseudocolor table for an RGB image.
func ComputeMedianCutPCT(
	red, green, blue RasterBand,
	colors int,
	ct ColorTable,
	progress ProgressFunc,
	data interface{},
) int {
	callback := newGoGDALProgressCallback(progress, data)
	defer callback.close()

	err := C.GDALComputeMedianCutPCT(
		red.cval,
		green.cval,
		blue.cval,
		nil,
		C.int(colors),
		ct.cval,
		callback.fn,
		callback.arg,
	)
	return int(err)
}

// DitherRGB2PCT converts a 24-bit RGB image to an 8-bit pseudocolor image.
func DitherRGB2PCT(
	red, green, blue, target RasterBand,
	ct ColorTable,
	progress ProgressFunc,
	data interface{},
) int {
	callback := newGoGDALProgressCallback(progress, data)
	defer callback.close()

	err := C.GDALDitherRGB2PCT(
		red.cval,
		green.cval,
		blue.cval,
		target.cval,
		ct.cval,
		callback.fn,
		callback.arg,
	)
	return int(err)
}

// Checksum computes a checksum for the requested image region.
func (rb RasterBand) Checksum(xOff, yOff, xSize, ySize int) int {
	sum := C.GDALChecksumImage(rb.cval, C.int(xOff), C.int(yOff), C.int(xSize), C.int(ySize))
	return int(sum)
}

// ComputeProximity computes the proximity of all pixels in dest to selected
// pixels in src.
func (rb RasterBand) ComputeProximity(
	dest RasterBand,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	callback := newGoGDALProgressCallback(progress, data)
	defer callback.close()

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return ErrFromCPLErr(C.GDALComputeProximity(
		rb.cval,
		dest.cval,
		(**C.char)(unsafe.Pointer(&opts[0])),
		callback.fn,
		callback.arg,
	))
}

// FillNoData fills selected raster regions by interpolating from surrounding
// pixels.
func (rb RasterBand) FillNoData(
	mask RasterBand,
	distance float64,
	iterations int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	callback := newGoGDALProgressCallback(progress, data)
	defer callback.close()

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return ErrFromCPLErr(C.GDALFillNodata(
		rb.cval,
		mask.cval,
		C.double(distance),
		0,
		C.int(iterations),
		(**C.char)(unsafe.Pointer(&opts[0])),
		callback.fn,
		callback.arg,
	))
}

// Polygonize creates polygon coverage from raster data using an integer buffer.
func (rb RasterBand) Polygonize(
	mask RasterBand,
	layer Layer,
	fieldIndex int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	callback := newGoGDALProgressCallback(progress, data)
	defer callback.close()

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return ErrFromCPLErr(C.GDALPolygonize(
		rb.cval,
		mask.cval,
		layer.cval,
		C.int(fieldIndex),
		(**C.char)(unsafe.Pointer(&opts[0])),
		callback.fn,
		callback.arg,
	))
}

// FPolygonize creates polygon coverage from raster data using a floating point buffer.
func (rb RasterBand) FPolygonize(
	mask RasterBand,
	layer Layer,
	fieldIndex int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	callback := newGoGDALProgressCallback(progress, data)
	defer callback.close()

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return ErrFromCPLErr(C.GDALFPolygonize(
		rb.cval,
		mask.cval,
		layer.cval,
		C.int(fieldIndex),
		(**C.char)(unsafe.Pointer(&opts[0])),
		callback.fn,
		callback.arg,
	))
}

// SieveFilter wraps the corresponding GDAL/OGR operation.
func (rb RasterBand) SieveFilter(
	mask, dest RasterBand,
	threshold, connectedness int,
	options []string,
	progress ProgressFunc,
	data interface{},
) error {
	callback := newGoGDALProgressCallback(progress, data)
	defer callback.close()

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return ErrFromCPLErr(C.GDALSieveFilter(
		rb.cval,
		mask.cval,
		dest.cval,
		C.int(threshold),
		C.int(connectedness),
		(**C.char)(unsafe.Pointer(&opts[0])),
		callback.fn,
		callback.arg,
	))
}

/* --------------------------------------------- */
/* Warp functions                                */
/* --------------------------------------------- */

//Unimplemented: CreateGenImgProjTransformer
//Unimplemented: CreateGenImgProjTransformer2
//Unimplemented: CreateGenImgProjTransformer3
//Unimplemented: SetGenImgProjTransformerDstGeoTransform
//Unimplemented: DestroyGenImgProjTransformer
//Unimplemented: GenImgProjTransform

//Unimplemented: CreateReprojectionTransformer
//Unimplemented: DestroyReprojection
//Unimplemented: ReprojectionTransform
//Unimplemented: CreateGCPTransformer
//Unimplemented: CreateGCPRefineTransformer
//Unimplemented: DestroyGCPTransformer
//Unimplemented: GCPTransform

//Unimplemented: CreateTPSTransformer
//Unimplemented: DestroyTPSTransformer
//Unimplemented: TPSTransform

//Unimplemented: CreateRPCTransformer
//Unimplemented: DestroyRPCTransformer
//Unimplemented: RPCTransform

//Unimplemented: CreateGeoLocTransformer
//Unimplemented: DestroyGeoLocTransformer
//Unimplemented: GeoLocTransform

//Unimplemented: CreateApproxTransformer
//Unimplemented: DestroyApproxTransformer
//Unimplemented: ApproxTransform

//Unimplemented: SimpleImageWarp
//Unimplemented: SuggestedWarpOutput
//Unimplemented: SuggsetedWarpOutput2
//Unimplemented: SerializeTransformer
//Unimplemented: DeserializeTransformer

//Unimplemented: TransformGeolocations

/* --------------------------------------------- */
/* Contour line functions                        */
/* --------------------------------------------- */

//Unimplemented: CreateContourGenerator
//Unimplemented: FeedLine
//Unimplemented: Destroy
//Unimplemented: ContourWriter
//Unimplemented: ContourGenerate

/* --------------------------------------------- */
/* Rasterizer functions                          */
/* --------------------------------------------- */

// Burn geometries into raster
//Unimplmemented: RasterizeGeometries

// Burn geometries from the specified list of layers into the raster
//Unimplemented: RasterizeLayers

// Burn geometries from the specified list of layers into the raster
//Unimplemented: RasterizeLayersBuf

/* --------------------------------------------- */
/* Gridding functions                            */
/* --------------------------------------------- */

// GridAlgorithm represents Grid Algorithm code
type GridAlgorithm int

// GA_InverseDistancetoAPower and related constants are exported GDAL/OGR symbols.
const (
	GA_InverseDistancetoAPower                = GridAlgorithm(C.GGA_InverseDistanceToAPower)
	GA_MovingAverage                          = GridAlgorithm(C.GGA_MovingAverage)
	GA_NearestNeighbor                        = GridAlgorithm(C.GGA_NearestNeighbor)
	GA_MetricMinimum                          = GridAlgorithm(C.GGA_MetricMinimum)
	GA_MetricMaximum                          = GridAlgorithm(C.GGA_MetricMaximum)
	GA_MetricRange                            = GridAlgorithm(C.GGA_MetricRange)
	GA_MetricCount                            = GridAlgorithm(C.GGA_MetricCount)
	GA_MetricAverageDistance                  = GridAlgorithm(C.GGA_MetricAverageDistance)
	GA_MetricAverageDistancePts               = GridAlgorithm(C.GGA_MetricAverageDistancePts)
	GA_Linear                                 = GridAlgorithm(C.GGA_Linear)
	GA_InverseDistanceToAPowerNearestNeighbor = GridAlgorithm(C.GGA_InverseDistanceToAPowerNearestNeighbor)
)

// GridLinearOptions represents GridLinearOptions: Linear method control options.
type GridLinearOptions struct {
	// Radius: in case the point to be interpolated does not fit into a triangle of the Delaunay triangulation,
	// use that maximum distance to search a nearest neighbour, or use nodata otherwise. If set to -1, the search
	// distance is infinite. If set to 0, nodata value will be always used.
	Radius float64
	// NoDataValue: no data marker to fill empty points.
	NoDataValue float64
}

// GridInverseDistanceToAPowerOptions represents GridInverseDistanceToAPowerOptions: Inverse distance to a power method control options.
type GridInverseDistanceToAPowerOptions struct {
	// Power: Weighting power
	Power float64
	// Smoothing: Smoothing parameter
	Smoothing float64
	// AnisotropyRatio: Reserved for future use
	AnisotropyRatio float64
	// AnisotropyAngle: Reserved for future use
	AnisotropyAngle float64
	// Radius1: The first radius (X axis if rotation angle is 0) of search ellipse.
	Radius1 float64
	// Radius2: The second radius (Y axis if rotation angle is 0) of search ellipse.
	Radius2 float64
	// Angle: Angle of ellipse rotation in degrees. Ellipse rotated counter clockwise.
	Angle float64
	// MaxPoints: Maximum number of data points to use.
	// Do not search for more points than this number. If less amount of points found the grid node
	// considered empty and will be filled with NODATA marker.
	MaxPoints uint32
	// MinPoints: Minimum number of data points to use.
	// If less amount of points found the grid node considered empty and will be filled with NODATA marker.
	MinPoints uint32
	// NoDataValue: No data marker to fill empty points.
	NoDataValue float64
}

// GridInverseDistanceToAPowerNearestNeighborOptions represents GridInverseDistanceToAPowerNearestNeighborOptions: Inverse distance to a power, with nearest neighbour search, control options.
type GridInverseDistanceToAPowerNearestNeighborOptions struct {
	// Power: Weighting power
	Power float64
	// Radius: The radius of search circle
	Radius float64
	// Smoothing: Smoothing parameter
	Smoothing float64
	// MaxPoints: Maximum number of data points to use.
	// Do not search for more points than this number. If less amount of points found the grid node
	// considered empty and will be filled with NODATA marker.
	MaxPoints uint32
	// MinPoints: Minimum number of data points to use.
	// If less amount of points found the grid node considered empty and will be filled with NODATA marker.
	MinPoints uint32
	// NoDataValue: No data marker to fill empty points.
	NoDataValue float64
}

// GridMovingAverageOptions represents GridMovingAverageOptions: Moving average method control options.
type GridMovingAverageOptions struct {
	// Radius1: The first radius (X axis if rotation angle is 0) of search ellipse.
	Radius1 float64
	// Radius2: The second radius (Y axis if rotation angle is 0) of search ellipse.
	Radius2 float64
	// Angle: Angle of ellipse rotation in degrees. Ellipse rotated counter clockwise.
	Angle float64
	// MinPoints: Minimum number of data points to use.
	// If less amount of points found the grid node considered empty and will be filled with NODATA marker.
	MinPoints uint32
	// NoDataValue: No data marker to fill empty points.
	NoDataValue float64
}

// GridNearestNeighborOptions represents GridNearestNeighborOptions: Nearest neighbor method control options.
type GridNearestNeighborOptions struct {
	// Radius1: The first radius (X axis if rotation angle is 0) of search ellipse.
	Radius1 float64
	// Radius2: The second radius (Y axis if rotation angle is 0) of search ellipse.
	Radius2 float64
	// Angle: Angle of ellipse rotation in degrees. Ellipse rotated counter clockwise.
	Angle float64
	// NoDataValue: No data marker to fill empty points.
	NoDataValue float64
}

// GridDataMetricsOptions represents GridDataMetricsOptions: Data metrics method control options.
type GridDataMetricsOptions struct {
	// Radius1: The first radius (X axis if rotation angle is 0) of search ellipse.
	Radius1 float64
	// Radius2: The second radius (Y axis if rotation angle is 0) of search ellipse.
	Radius2 float64
	// Angle: Angle of ellipse rotation in degrees. Ellipse rotated counter clockwise.
	Angle float64
	// MinPoints: Minimum number of data points to use.
	// If less amount of points found the grid node considered empty and will be filled with NODATA marker.
	MinPoints uint32
	// NoDataValue: No data marker to fill empty points.
	NoDataValue float64
}

var errInvalidOptionsTypeWasPassed = errors.New("invalid options type was passed")

// GridCreate wraps the corresponding GDAL/OGR operation.
func GridCreate(
	algorithm GridAlgorithm,
	options interface{},
	x, y, z []float64,
	xMin, xMax, yMin, yMax float64,
	nX, nY uint,
	progress ProgressFunc,
	data interface{},
) ([]float64, error) {
	if len(x) != len(y) || len(x) != len(z) {
		return nil, errors.New("lengths of x, y, z should equal")
	}
	if len(x) == 0 {
		return nil, errors.New("x, y, z must not be empty")
	}
	if nX == 0 || nY == 0 {
		return nil, errors.New("nX and nY must be greater than zero")
	}

	poptions := unsafe.Pointer(nil)
	switch algorithm {
	case GA_InverseDistancetoAPower:
		soptions, ok := options.(GridInverseDistanceToAPowerOptions)
		if !ok {
			return nil, errInvalidOptionsTypeWasPassed
		}
		cOptions := C.goGDALGridInverseDistanceToAPowerOptionsInit()
		cOptions.dfPower = C.double(soptions.Power)
		cOptions.dfSmoothing = C.double(soptions.Smoothing)
		cOptions.dfAnisotropyRatio = C.double(soptions.AnisotropyRatio)
		cOptions.dfAnisotropyAngle = C.double(soptions.AnisotropyAngle)
		cOptions.dfRadius1 = C.double(soptions.Radius1)
		cOptions.dfRadius2 = C.double(soptions.Radius2)
		cOptions.dfAngle = C.double(soptions.Angle)
		cOptions.nMaxPoints = C.uint(soptions.MaxPoints)
		cOptions.nMinPoints = C.uint(soptions.MinPoints)
		cOptions.dfNoDataValue = C.double(soptions.NoDataValue)
		poptions = unsafe.Pointer(&cOptions)
	case GA_InverseDistanceToAPowerNearestNeighbor:
		soptions, ok := options.(GridInverseDistanceToAPowerNearestNeighborOptions)
		if !ok {
			return nil, errInvalidOptionsTypeWasPassed
		}
		cOptions := C.goGDALGridInverseDistanceToAPowerNearestNeighborOptionsInit()
		cOptions.dfPower = C.double(soptions.Power)
		cOptions.dfRadius = C.double(soptions.Radius)
		cOptions.dfSmoothing = C.double(soptions.Smoothing)
		cOptions.nMaxPoints = C.uint(soptions.MaxPoints)
		cOptions.nMinPoints = C.uint(soptions.MinPoints)
		cOptions.dfNoDataValue = C.double(soptions.NoDataValue)
		poptions = unsafe.Pointer(&cOptions)
	case GA_MovingAverage:
		soptions, ok := options.(GridMovingAverageOptions)
		if !ok {
			return nil, errInvalidOptionsTypeWasPassed
		}
		cOptions := C.goGDALGridMovingAverageOptionsInit()
		cOptions.dfRadius1 = C.double(soptions.Radius1)
		cOptions.dfRadius2 = C.double(soptions.Radius2)
		cOptions.dfAngle = C.double(soptions.Angle)
		cOptions.nMinPoints = C.uint(soptions.MinPoints)
		cOptions.dfNoDataValue = C.double(soptions.NoDataValue)
		poptions = unsafe.Pointer(&cOptions)
	case GA_NearestNeighbor:
		soptions, ok := options.(GridNearestNeighborOptions)
		if !ok {
			return nil, errInvalidOptionsTypeWasPassed
		}
		cOptions := C.goGDALGridNearestNeighborOptionsInit()
		cOptions.dfRadius1 = C.double(soptions.Radius1)
		cOptions.dfRadius2 = C.double(soptions.Radius2)
		cOptions.dfAngle = C.double(soptions.Angle)
		cOptions.dfNoDataValue = C.double(soptions.NoDataValue)
		poptions = unsafe.Pointer(&cOptions)
	case GA_MetricMinimum, GA_MetricMaximum, GA_MetricCount, GA_MetricRange,
		GA_MetricAverageDistance, GA_MetricAverageDistancePts:
		soptions, ok := options.(GridDataMetricsOptions)
		if !ok {
			return nil, errInvalidOptionsTypeWasPassed
		}
		cOptions := C.goGDALGridDataMetricsOptionsInit()
		cOptions.dfRadius1 = C.double(soptions.Radius1)
		cOptions.dfRadius2 = C.double(soptions.Radius2)
		cOptions.dfAngle = C.double(soptions.Angle)
		cOptions.nMinPoints = C.uint(soptions.MinPoints)
		cOptions.dfNoDataValue = C.double(soptions.NoDataValue)
		poptions = unsafe.Pointer(&cOptions)
	case GA_Linear:
		soptions, ok := options.(GridLinearOptions)
		if !ok {
			return nil, errInvalidOptionsTypeWasPassed
		}
		cOptions := C.goGDALGridLinearOptionsInit()
		cOptions.dfRadius = C.double(soptions.Radius)
		cOptions.dfNoDataValue = C.double(soptions.NoDataValue)
		poptions = unsafe.Pointer(&cOptions)
	}

	buffer := make([]float64, nX*nY)
	callback := newGoGDALProgressCallback(progress, data)
	defer callback.close()
	err := ErrFromCPLErr(C.GDALGridCreate(
		C.GDALGridAlgorithm(algorithm),
		poptions,
		C.uint(uint(len(x))),
		float64SlicePtr(x),
		float64SlicePtr(y),
		float64SlicePtr(z),
		C.double(xMin),
		C.double(xMax),
		C.double(yMin),
		C.double(yMax),
		C.uint(nX),
		C.uint(nY),
		C.GDALDataType(Float64),
		unsafe.Pointer(float64SlicePtr(buffer)),
		callback.fn,
		callback.arg,
	))
	return buffer, err
}

//Unimplemented: ComputeMatchingPoints
