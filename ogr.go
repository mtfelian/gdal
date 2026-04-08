package gdal

/*
#include "go_gdal.h"
#include "go_ogr_wkb.h"
#include "gdal_version.h"
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

func init() {
	C.OGRRegisterAll()
}

/* -------------------------------------------------------------------- */
/*      Significant constants.                                          */
/* -------------------------------------------------------------------- */

// GeometryType enumerates well known binary geometry types.
type GeometryType uint32

// GT_Unknown and related constants are exported GDAL/OGR symbols.
const (
	GT_Unknown               = GeometryType(C.wkbUnknown)
	GT_Point                 = GeometryType(C.wkbPoint)
	GT_LineString            = GeometryType(C.wkbLineString)
	GT_Polygon               = GeometryType(C.wkbPolygon)
	GT_MultiPoint            = GeometryType(C.wkbMultiPoint)
	GT_MultiLineString       = GeometryType(C.wkbMultiLineString)
	GT_MultiPolygon          = GeometryType(C.wkbMultiPolygon)
	GT_GeometryCollection    = GeometryType(C.wkbGeometryCollection)
	GT_None                  = GeometryType(C.wkbNone)
	GT_LinearRing            = GeometryType(C.wkbLinearRing)
	GT_Point25D              = GeometryType(C.wkbPoint25D)
	GT_LineString25D         = GeometryType(C.wkbLineString25D)
	GT_Polygon25D            = GeometryType(C.wkbPolygon25D)
	GT_MultiPoint25D         = GeometryType(C.wkbMultiPoint25D)
	GT_MultiLineString25D    = GeometryType(C.wkbMultiLineString25D)
	GT_MultiPolygon25D       = GeometryType(C.wkbMultiPolygon25D)
	GT_GeometryCollection25D = GeometryType(C.wkbGeometryCollection25D)
)

/* -------------------------------------------------------------------- */
/*      Envelope functions                                              */
/* -------------------------------------------------------------------- */

// Envelope is an exported GDAL/OGR type.
type Envelope struct {
	cval C.OGREnvelope
}

// MinX wraps the corresponding GDAL/OGR operation.
func (env Envelope) MinX() float64 {
	return float64(env.cval.MinX)
}

// MaxX wraps the corresponding GDAL/OGR operation.
func (env Envelope) MaxX() float64 {
	return float64(env.cval.MaxX)
}

// MinY wraps the corresponding GDAL/OGR operation.
func (env Envelope) MinY() float64 {
	return float64(env.cval.MinY)
}

// MaxY wraps the corresponding GDAL/OGR operation.
func (env Envelope) MaxY() float64 {
	return float64(env.cval.MaxY)
}

// SetMinX wraps the corresponding GDAL/OGR operation.
func (env *Envelope) SetMinX(val float64) {
	env.cval.MinX = C.double(val)
}

// SetMaxX wraps the corresponding GDAL/OGR operation.
func (env *Envelope) SetMaxX(val float64) {
	env.cval.MaxX = C.double(val)
}

// SetMinY wraps the corresponding GDAL/OGR operation.
func (env *Envelope) SetMinY(val float64) {
	env.cval.MinY = C.double(val)
}

// SetMaxY wraps the corresponding GDAL/OGR operation.
func (env *Envelope) SetMaxY(val float64) {
	env.cval.MaxY = C.double(val)
}

// IsInit wraps the corresponding GDAL/OGR operation.
func (env Envelope) IsInit() bool {
	return env.cval.MinX != 0 || env.cval.MinY != 0 || env.cval.MaxX != 0 || env.cval.MaxY != 0
}

func min(a, b C.double) C.double {
	if a < b {
		return a
	}
	return b
}

func max(a, b C.double) C.double {
	if a > b {
		return a
	}
	return b
}

// Union returns the union of this envelope with another one.
func (env Envelope) Union(other Envelope) Envelope {
	if env.IsInit() {
		env.cval.MinX = min(env.cval.MinX, other.cval.MinX)
		env.cval.MinY = min(env.cval.MinY, other.cval.MinY)
		env.cval.MaxX = max(env.cval.MaxX, other.cval.MaxX)
		env.cval.MaxY = max(env.cval.MaxY, other.cval.MaxY)
	} else {
		env.cval.MinX = other.cval.MinX
		env.cval.MinY = other.cval.MinY
		env.cval.MaxX = other.cval.MaxX
		env.cval.MaxY = other.cval.MaxY
	}
	return env
}

// Intersect returns the intersection of this envelope with another.
func (env Envelope) Intersect(other Envelope) Envelope {
	if env.Intersects(other) {
		if env.IsInit() {
			env.cval.MinX = max(env.cval.MinX, other.cval.MinX)
			env.cval.MinY = max(env.cval.MinY, other.cval.MinY)
			env.cval.MaxX = min(env.cval.MaxX, other.cval.MaxX)
			env.cval.MaxY = min(env.cval.MaxY, other.cval.MaxY)
		} else {
			env.cval.MinX = other.cval.MinX
			env.cval.MinY = other.cval.MinY
			env.cval.MaxX = other.cval.MaxX
			env.cval.MaxY = other.cval.MaxY
		}
	} else {
		env.cval.MinX = 0
		env.cval.MinY = 0
		env.cval.MaxX = 0
		env.cval.MaxY = 0
	}
	return env
}

// Intersects reports whether one envelope intersects another.
func (env Envelope) Intersects(other Envelope) bool {
	return env.cval.MinX <= other.cval.MaxX &&
		env.cval.MaxX >= other.cval.MinX &&
		env.cval.MinY <= other.cval.MaxY &&
		env.cval.MaxY >= other.cval.MinY
}

// Contains reports whether one envelope completely contains another.
func (env Envelope) Contains(other Envelope) bool {
	return env.cval.MinX <= other.cval.MinX &&
		env.cval.MaxX >= other.cval.MaxX &&
		env.cval.MinY <= other.cval.MinY &&
		env.cval.MaxY >= other.cval.MaxY
}

/* -------------------------------------------------------------------- */
/*      Misc functions                                                  */
/* -------------------------------------------------------------------- */

// CleanupOGR cleans up all OGR related resources.
func CleanupOGR() {
	C.OGRCleanupAll()
}

// BoolToCInt converts a go bool to a C int.
func BoolToCInt(in bool) (out C.int) {
	if in {
		out = 1
	} else {
		out = 0
	}
	return
}

/* -------------------------------------------------------------------- */
/*      Geometry functions                                              */
/* -------------------------------------------------------------------- */

// Geometry is an exported GDAL/OGR type.
type Geometry struct {
	cval C.OGRGeometryH
}

// CreateFromWKB creates a geometry object from its well known binary representation.
func CreateFromWKB(wkb []uint8, srs SpatialReference, bytes int) (Geometry, error) {
	var newGeom Geometry
	if len(wkb) == 0 {
		return newGeom, fmt.Errorf("wkb must not be empty")
	}
	if bytes <= 0 || bytes > len(wkb) {
		return newGeom, fmt.Errorf("wkb byte count %d is out of range for buffer length %d", bytes, len(wkb))
	}
	cString := unsafe.Pointer(&wkb[0])
	return newGeom, ErrFromOGRErr(C.go_CreateFromWkb(
		cString, srs.cval, &newGeom.cval, C.int(bytes),
	))
}

// CreateFromWKT creates a geometry object from its well known text representation.
func CreateFromWKT(wkt string, srs SpatialReference) (Geometry, error) {
	cString := C.CString(wkt)
	defer C.free(unsafe.Pointer(cString))
	var newGeom Geometry
	return newGeom, ErrFromOGRErr(C.OGR_G_CreateFromWkt(
		&cString, srs.cval, &newGeom.cval,
	))
}

// CreateFromJson creates a geometry object from its GeoJSON representation.
func CreateFromJson(_json string) Geometry {
	cString := C.CString(_json)
	defer C.free(unsafe.Pointer(cString))
	var newGeom Geometry
	newGeom.cval = C.OGR_G_CreateGeometryFromJson(cString)
	return newGeom
}

// Destroy geometry object
func (geometry Geometry) Destroy() {
	C.OGR_G_DestroyGeometry(geometry.cval)
}

// Create an empty geometry of the desired type
func Create(geomType GeometryType) Geometry {
	geom := C.OGR_G_CreateGeometry(C.OGRwkbGeometryType(geomType))
	return Geometry{geom}
}

// ApproximateArcAngles strokes arc to linestring.
func ApproximateArcAngles(
	x, y, z,
	primaryRadius,
	secondaryRadius,
	rotation,
	startAngle,
	endAngle,
	stepSizeDegrees float64,
) Geometry {
	geom := C.OGR_G_ApproximateArcAngles(
		C.double(x),
		C.double(y),
		C.double(z),
		C.double(primaryRadius),
		C.double(secondaryRadius),
		C.double(rotation),
		C.double(startAngle),
		C.double(endAngle),
		C.double(stepSizeDegrees))
	return Geometry{geom}
}

// ForceToPolygon converts to polygon.
func (geometry Geometry) ForceToPolygon() Geometry {
	newGeom := C.OGR_G_ForceToPolygon(geometry.cval)
	return Geometry{newGeom}
}

// ForceToMultiPolygon converts to multipolygon.
func (geometry Geometry) ForceToMultiPolygon() Geometry {
	newGeom := C.OGR_G_ForceToMultiPolygon(geometry.cval)
	return Geometry{newGeom}
}

// ForceToMultiPoint converts to multipoint.
func (geometry Geometry) ForceToMultiPoint() Geometry {
	newGeom := C.OGR_G_ForceToMultiPoint(geometry.cval)
	return Geometry{newGeom}
}

// ForceToMultiLineString converts to multilinestring.
func (geometry Geometry) ForceToMultiLineString() Geometry {
	newGeom := C.OGR_G_ForceToMultiLineString(geometry.cval)
	return Geometry{newGeom}
}

// Dimension returns the dimension of this geometry.
func (geometry Geometry) Dimension() int {
	dim := C.OGR_G_GetDimension(geometry.cval)
	return int(dim)
}

// CoordinateDimension returns the dimension of the coordinates in this geometry.
func (geometry Geometry) CoordinateDimension() int {
	dim := C.OGR_G_GetCoordinateDimension(geometry.cval)
	return int(dim)
}

// SetCoordinateDimension sets the dimension of the coordinates in this geometry.
func (geometry Geometry) SetCoordinateDimension(dim int) {
	C.OGR_G_SetCoordinateDimension(geometry.cval, C.int(dim))
}

// Clone creates a copy of this geometry.
func (geometry Geometry) Clone() Geometry {
	newGeom := C.OGR_G_Clone(geometry.cval)
	return Geometry{newGeom}
}

// Envelope computes and returns the bounding envelope for this geometry.
func (geometry Geometry) Envelope() Envelope {
	var env Envelope
	C.OGR_G_GetEnvelope(geometry.cval, &env.cval)
	return env
}

// Unimplemented: GetEnvelope3D

// FromWKB assigns a geometry from well known binary data.
func (geometry Geometry) FromWKB(wkb []uint8, bytes int) error {
	if len(wkb) == 0 {
		return fmt.Errorf("wkb must not be empty")
	}
	if bytes <= 0 || bytes > len(wkb) {
		return fmt.Errorf("wkb byte count %d is out of range for buffer length %d", bytes, len(wkb))
	}
	cString := unsafe.Pointer(&wkb[0])
	return ErrFromOGRErr(C.go_ImportFromWkb(geometry.cval, cString, C.int(bytes)))
}

// ToWKB converts a geometry to well known binary data.
func (geometry Geometry) ToWKB() ([]uint8, error) {
	b := make([]uint8, geometry.WKBSize())
	cString := (*C.uchar)(unsafe.Pointer(&b[0]))
	err := ErrFromOGRErr(C.go_ExportToWkb(geometry.cval, C.OGRwkbByteOrder(C.wkbNDR), cString))
	return b, err
}

// WKBSize returns size of related binary representation.
func (geometry Geometry) WKBSize() int {
	size := C.OGR_G_WkbSize(geometry.cval)
	return int(size)
}

// FromWKT assigns geometry object from its well known text representation.
func (geometry Geometry) FromWKT(wkt string) error {
	cString := C.CString(wkt)
	defer C.free(unsafe.Pointer(cString))
	return ErrFromOGRErr(C.OGR_G_ImportFromWkt(geometry.cval, &cString))
}

// ToWKT returns geometry as WKT.
func (geometry Geometry) ToWKT() (string, error) {
	var p *C.char
	err := ErrFromOGRErr(C.OGR_G_ExportToWkt(geometry.cval, &p))
	return goStringAndCPLFree(p), err
}

// Type returns geometry type.
func (geometry Geometry) Type() GeometryType {
	gt := C.OGR_G_GetGeometryType(geometry.cval)
	return GeometryType(gt)
}

// Name returns geometry name.
func (geometry Geometry) Name() string {
	name := C.OGR_G_GetGeometryName(geometry.cval)
	return C.GoString(name)
}

// Unimplemented: DumpReadable

// FlattenTo2D converts geometry to strictly 2D.
func (geometry Geometry) FlattenTo2D() {
	C.OGR_G_FlattenTo2D(geometry.cval)
}

// CloseRings wraps the corresponding GDAL/OGR operation.
func (geometry Geometry) CloseRings() {
	C.OGR_G_CloseRings(geometry.cval)
}

// CreateFromGML creates a geometry from its GML representation.
func CreateFromGML(gml string) Geometry {
	cString := C.CString(gml)
	defer C.free(unsafe.Pointer(cString))
	geom := C.OGR_G_CreateFromGML(cString)
	return Geometry{geom}
}

// ToGML converts a geometry to GML format.
func (geometry Geometry) ToGML() string {
	return goStringAndCPLFree(C.OGR_G_ExportToGML(geometry.cval))
}

// ToGML_Ex converts a geometry to GML format with options.
func (geometry Geometry) ToGML_Ex(options []string) string {
	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return goStringAndCPLFree(C.OGR_G_ExportToGMLEx(geometry.cval, (**C.char)(unsafe.Pointer(&opts[0]))))
}

// ToKML converts a geometry to KML format.
func (geometry Geometry) ToKML() string {
	return goStringAndCPLFree(C.OGR_G_ExportToKML(geometry.cval, nil))
}

// ToJSON converts a geometry to JSON format.
func (geometry Geometry) ToJSON() string {
	return goStringAndCPLFree(C.OGR_G_ExportToJson(geometry.cval))
}

// ToJSON_ex converts a geometry to JSON format with options.
func (geometry Geometry) ToJSON_ex(options []string) string {
	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	return goStringAndCPLFree(C.OGR_G_ExportToJsonEx(geometry.cval, (**C.char)(unsafe.Pointer(&opts[0]))))
}

// SpatialReference returns the spatial reference associated with this geometry.
func (geometry Geometry) SpatialReference() SpatialReference {
	spatialRef := C.OGR_G_GetSpatialReference(geometry.cval)
	return SpatialReference{spatialRef}
}

// SetSpatialReference assigns a spatial reference to this geometry.
func (geometry Geometry) SetSpatialReference(spatialRef SpatialReference) {
	C.OGR_G_AssignSpatialReference(geometry.cval, spatialRef.cval)
}

// Transform applies coordinate transformation to geometry.
func (geometry Geometry) Transform(ct CoordinateTransform) error {
	return ErrFromOGRErr(C.OGR_G_Transform(geometry.cval, ct.cval))
}

// TransformTo wraps the corresponding GDAL/OGR operation.
func (geometry Geometry) TransformTo(sr SpatialReference) error {
	return ErrFromOGRErr(C.OGR_G_TransformTo(geometry.cval, sr.cval))
}

// Simplify the geometry
func (geometry Geometry) Simplify(tolerance float64) Geometry {
	newGeom := C.OGR_G_Simplify(geometry.cval, C.double(tolerance))
	return Geometry{newGeom}
}

// SimplifyPreservingTopology simplifies the geometry while preserving topology.
func (geometry Geometry) SimplifyPreservingTopology(tolerance float64) Geometry {
	newGeom := C.OGR_G_SimplifyPreserveTopology(geometry.cval, C.double(tolerance))
	return Geometry{newGeom}
}

// Segmentize modifies the geometry such that it has no line segment longer than the given distance.
func (geometry Geometry) Segmentize(distance float64) {
	C.OGR_G_Segmentize(geometry.cval, C.double(distance))
}

// Intersects reports whether these features intersect.
func (geometry Geometry) Intersects(other Geometry) bool {
	val := C.OGR_G_Intersects(geometry.cval, other.cval)
	return val != 0
}

// Equals reports whether these features are equal.
func (geometry Geometry) Equals(other Geometry) bool {
	val := C.OGR_G_Equals(geometry.cval, other.cval)
	return val != 0
}

// Disjoint reports whether the features are disjoint.
func (geometry Geometry) Disjoint(other Geometry) bool {
	val := C.OGR_G_Disjoint(geometry.cval, other.cval)
	return val != 0
}

// Touches reports whether this feature touches the other.
func (geometry Geometry) Touches(other Geometry) bool {
	val := C.OGR_G_Touches(geometry.cval, other.cval)
	return val != 0
}

// Crosses reports whether this feature crosses the other.
func (geometry Geometry) Crosses(other Geometry) bool {
	val := C.OGR_G_Crosses(geometry.cval, other.cval)
	return val != 0
}

// Within reports whether this geometry is within the other.
func (geometry Geometry) Within(other Geometry) bool {
	val := C.OGR_G_Within(geometry.cval, other.cval)
	return val != 0
}

// Contains reports whether this geometry contains the other.
func (geometry Geometry) Contains(other Geometry) bool {
	val := C.OGR_G_Contains(geometry.cval, other.cval)
	return val != 0
}

// Overlaps reports whether this geometry overlaps the other.
func (geometry Geometry) Overlaps(other Geometry) bool {
	val := C.OGR_G_Overlaps(geometry.cval, other.cval)
	return val != 0
}

// Boundary computes boundary for the geometry.
func (geometry Geometry) Boundary() Geometry {
	newGeom := C.OGR_G_Boundary(geometry.cval)
	return Geometry{newGeom}
}

// ConvexHull computes convex hull for the geometry.
func (geometry Geometry) ConvexHull() Geometry {
	newGeom := C.OGR_G_ConvexHull(geometry.cval)
	return Geometry{newGeom}
}

// Buffer computes buffer of the geometry.
func (geometry Geometry) Buffer(distance float64, segments int) Geometry {
	newGeom := C.OGR_G_Buffer(geometry.cval, C.double(distance), C.int(segments))
	return Geometry{newGeom}
}

// Intersection computes intersection of this geometry with the other.
func (geometry Geometry) Intersection(other Geometry) Geometry {
	newGeom := C.OGR_G_Intersection(geometry.cval, other.cval)
	return Geometry{newGeom}
}

// Union computes union of this geometry with the other.
func (geometry Geometry) Union(other Geometry) Geometry {
	newGeom := C.OGR_G_Union(geometry.cval, other.cval)
	return Geometry{newGeom}
}

// UnionCascaded wraps the corresponding GDAL/OGR operation.
func (geometry Geometry) UnionCascaded() Geometry {
	newGeom := C.OGR_G_UnionCascaded(geometry.cval)
	return Geometry{newGeom}
}

// Unimplemented: PointOn Surface (until 2.0)
// Return a point guaranteed to lie on the surface
// func (geom Geometry) PointOnSurface() Geometry {
//  newGeom := C.OGR_G_PointOnSurface(geom.cval)
//  return Geometry{newGeom}
// }

// Difference computes difference between this geometry and the other.
func (geometry Geometry) Difference(other Geometry) Geometry {
	newGeom := C.OGR_G_Difference(geometry.cval, other.cval)
	return Geometry{newGeom}
}

// SymmetricDifference computes symmetric difference between this geometry and the other.
func (geometry Geometry) SymmetricDifference(other Geometry) Geometry {
	newGeom := C.OGR_G_SymDifference(geometry.cval, other.cval)
	return Geometry{newGeom}
}

// Distance computes distance between thie geometry and the other.
func (geometry Geometry) Distance(other Geometry) float64 {
	dist := C.OGR_G_Distance(geometry.cval, other.cval)
	return float64(dist)
}

// Distance3D computes 3D distance between thie geometry and the other. This method is built on the SFCGAL library, check it for the definition of the geometry operation. If OGR is built without the SFCGAL library, this method will always return -1.0.
func (geometry Geometry) Distance3D(other Geometry) float64 {
	dist := C.OGR_G_Distance3D(geometry.cval, other.cval)
	return float64(dist)
}

// Length computes length of geometry.
func (geometry Geometry) Length() float64 {
	length := C.OGR_G_Length(geometry.cval)
	return float64(length)
}

// Area computes area of geometry.
func (geometry Geometry) Area() float64 {
	area := C.OGR_G_Area(geometry.cval)
	return float64(area)
}

// Centroid computes centroid of geometry.
func (geometry Geometry) Centroid() Geometry {
	var centroid Geometry
	C.OGR_G_Centroid(geometry.cval, centroid.cval)
	return centroid
}

// Empty wraps the corresponding GDAL/OGR operation.
func (geometry Geometry) Empty() {
	C.OGR_G_Empty(geometry.cval)
}

// IsEmpty reports whether the geometry is empty.
func (geometry Geometry) IsEmpty() bool {
	val := C.OGR_G_IsEmpty(geometry.cval)
	return val != 0
}

// IsNull reports whether the geometry is null.
func (geometry Geometry) IsNull() bool {
	return geometry.cval == nil
}

// IsValid reports whether the geometry is valid.
func (geometry Geometry) IsValid() bool {
	val := C.OGR_G_IsValid(geometry.cval)
	return val != 0
}

// IsSimple reports whether the geometry is simple.
func (geometry Geometry) IsSimple() bool {
	val := C.OGR_G_IsSimple(geometry.cval)
	return val != 0
}

// IsRing reports whether the geometry is a ring.
func (geometry Geometry) IsRing() bool {
	val := C.OGR_G_IsRing(geometry.cval)
	return val != 0
}

// Polygonize a set of sparse edges
func (geometry Geometry) Polygonize() Geometry {
	newGeom := C.OGR_G_Polygonize(geometry.cval)
	return Geometry{newGeom}
}

// PointCount returns number of points in the geometry.
func (geometry Geometry) PointCount() int {
	count := C.OGR_G_GetPointCount(geometry.cval)
	return int(count)
}

// Unimplemented: Points

// X returns the X coordinate of a point in the geometry.
func (geometry Geometry) X(index int) float64 {
	x := C.OGR_G_GetX(geometry.cval, C.int(index))
	return float64(x)
}

// Y returns the Y coordinate of a point in the geometry.
func (geometry Geometry) Y(index int) float64 {
	y := C.OGR_G_GetY(geometry.cval, C.int(index))
	return float64(y)
}

// Z returns the Z coordinate of a point in the geometry.
func (geometry Geometry) Z(index int) float64 {
	z := C.OGR_G_GetZ(geometry.cval, C.int(index))
	return float64(z)
}

// Point returns the coordinates of a point in the geometry.
func (geometry Geometry) Point(index int) (x, y, z float64) {
	C.OGR_G_GetPoint(
		geometry.cval,
		C.int(index),
		(*C.double)(&x),
		(*C.double)(&y),
		(*C.double)(&z))
	return
}

// SetPoint sets the coordinates of a point in the geometry.
func (geometry Geometry) SetPoint(index int, x, y, z float64) {
	C.OGR_G_SetPoint(
		geometry.cval,
		C.int(index),
		C.double(x),
		C.double(y),
		C.double(z))
}

// SetPoint2D sets the coordinates of a point in the geometry, ignoring the 3rd dimension.
func (geometry Geometry) SetPoint2D(index int, x, y float64) {
	C.OGR_G_SetPoint_2D(geometry.cval, C.int(index), C.double(x), C.double(y))
}

// AddPoint adds a new point to the geometry (line string or polygon only).
func (geometry Geometry) AddPoint(x, y, z float64) {
	C.OGR_G_AddPoint(geometry.cval, C.double(x), C.double(y), C.double(z))
}

// AddPoint2D adds a new point to the geometry (line string or polygon only), ignoring the 3rd dimension.
func (geometry Geometry) AddPoint2D(x, y float64) {
	C.OGR_G_AddPoint_2D(geometry.cval, C.double(x), C.double(y))
}

// GeometryCount returns the number of elements in the geometry, or number of geometries in the container.
func (geometry Geometry) GeometryCount() int {
	count := C.OGR_G_GetGeometryCount(geometry.cval)
	return int(count)
}

// Geometry returns geometry from a geometry container.
func (geometry Geometry) Geometry(index int) Geometry {
	newGeom := C.OGR_G_GetGeometryRef(geometry.cval, C.int(index))
	return Geometry{newGeom}
}

// AddGeometry adds a geometry to a geometry container.
func (geometry Geometry) AddGeometry(other Geometry) error {
	return ErrFromOGRErr(C.OGR_G_AddGeometry(geometry.cval, other.cval))
}

// AddGeometryDirectly adds a geometry to a geometry container and assign ownership to that container.
func (geometry Geometry) AddGeometryDirectly(other Geometry) error {
	return ErrFromOGRErr(C.OGR_G_AddGeometryDirectly(geometry.cval, other.cval))
}

// RemoveGeometry removes a geometry from the geometry container.
func (geometry Geometry) RemoveGeometry(index int, delete bool) error {
	return ErrFromOGRErr(C.OGR_G_RemoveGeometry(geometry.cval, C.int(index), BoolToCInt(delete)))
}

// BuildPolygonFromEdges builds a polygon / ring from a set of lines.
func (geometry Geometry) BuildPolygonFromEdges(autoClose bool, tolerance float64) (Geometry, error) {
	var cErr C.OGRErr
	newGeom := C.OGRBuildPolygonFromEdges(
		geometry.cval,
		0,
		BoolToCInt(autoClose),
		C.double(tolerance),
		&cErr,
	)
	return Geometry{newGeom}, ErrFromOGRErr(cErr)
}

/* -------------------------------------------------------------------- */
/*      Field definition functions                                      */
/* -------------------------------------------------------------------- */

// FieldType enumerates well known binary geometry types.
type FieldType int

// FT_Integer and related constants are exported GDAL/OGR symbols.
const (
	FT_Integer       = FieldType(C.OFTInteger)
	FT_IntegerList   = FieldType(C.OFTIntegerList)
	FT_Real          = FieldType(C.OFTReal)
	FT_RealList      = FieldType(C.OFTRealList)
	FT_String        = FieldType(C.OFTString)
	FT_StringList    = FieldType(C.OFTStringList)
	FT_Binary        = FieldType(C.OFTBinary)
	FT_Date          = FieldType(C.OFTDate)
	FT_Time          = FieldType(C.OFTTime)
	FT_DateTime      = FieldType(C.OFTDateTime)
	FT_Integer64     = FieldType(C.OFTInteger64)
	FT_Integer64List = FieldType(C.OFTInteger64List)
)

// Justification is an exported GDAL/OGR type.
type Justification int

// J_Undefined and related constants are exported GDAL/OGR symbols.
const (
	J_Undefined = Justification(C.OJUndefined)
	J_Left      = Justification(C.OJLeft)
	J_Right     = Justification(C.OJRight)
)

// FieldDefinition is an exported GDAL/OGR type.
type FieldDefinition struct {
	cval C.OGRFieldDefnH
}

// Field is an exported GDAL/OGR type.
type Field struct {
	cval *C.OGRField
}

// CreateFieldDefinition creates a new field definition.
func CreateFieldDefinition(name string, fieldType FieldType) FieldDefinition {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	fieldDef := C.OGR_Fld_Create(cName, C.OGRFieldType(fieldType))
	return FieldDefinition{fieldDef}
}

// Destroy the field definition
func (fd FieldDefinition) Destroy() {
	C.OGR_Fld_Destroy(fd.cval)
}

// Name returns the name of the field.
func (fd FieldDefinition) Name() string {
	name := C.OGR_Fld_GetNameRef(fd.cval)
	return C.GoString(name)
}

// SetName sets the name of the field.
func (fd FieldDefinition) SetName(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.OGR_Fld_SetName(fd.cval, cName)
}

// Type returns the type of this field.
func (fd FieldDefinition) Type() FieldType {
	fType := C.OGR_Fld_GetType(fd.cval)
	return FieldType(fType)
}

// SetType sets the type of this field.
func (fd FieldDefinition) SetType(fType FieldType) {
	C.OGR_Fld_SetType(fd.cval, C.OGRFieldType(fType))
}

// Justification returns the justification for this field.
func (fd FieldDefinition) Justification() Justification {
	justify := C.OGR_Fld_GetJustify(fd.cval)
	return Justification(justify)
}

// SetJustification sets the justification for this field.
func (fd FieldDefinition) SetJustification(justify Justification) {
	C.OGR_Fld_SetJustify(fd.cval, C.OGRJustification(justify))
}

// Width returns the formatting width for this field.
func (fd FieldDefinition) Width() int {
	width := C.OGR_Fld_GetWidth(fd.cval)
	return int(width)
}

// SetWidth sets the formatting width for this field.
func (fd FieldDefinition) SetWidth(width int) {
	C.OGR_Fld_SetWidth(fd.cval, C.int(width))
}

// Precision returns the precision for this field.
func (fd FieldDefinition) Precision() int {
	precision := C.OGR_Fld_GetPrecision(fd.cval)
	return int(precision)
}

// SetPrecision sets the precision for this field.
func (fd FieldDefinition) SetPrecision(precision int) {
	C.OGR_Fld_SetPrecision(fd.cval, C.int(precision))
}

// Set defining parameters of field in a single call
func (fd FieldDefinition) Set(
	name string,
	fType FieldType,
	width, precision int,
	justify Justification,
) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	C.OGR_Fld_Set(
		fd.cval,
		cName,
		C.OGRFieldType(fType),
		C.int(width),
		C.int(precision),
		C.OGRJustification(justify),
	)
}

// IsIgnored returns whether this field should be ignored when fetching features.
func (fd FieldDefinition) IsIgnored() bool {
	ignore := C.OGR_Fld_IsIgnored(fd.cval)
	return ignore != 0
}

// SetIgnored sets whether this field should be ignored when fetching features.
func (fd FieldDefinition) SetIgnored(ignore bool) {
	C.OGR_Fld_SetIgnored(fd.cval, BoolToCInt(ignore))
}

// Name returns human readable name for the field type.
func (ft FieldType) Name() string {
	name := C.OGR_GetFieldTypeName(C.OGRFieldType(ft))
	return C.GoString(name)
}

/* -------------------------------------------------------------------- */
/*      Feature definition functions                                    */
/* -------------------------------------------------------------------- */

// FeatureDefinition is an exported GDAL/OGR type.
type FeatureDefinition struct {
	cval C.OGRFeatureDefnH
}

// CreateFeatureDefinition creates a new feature definition object.
func CreateFeatureDefinition(name string) FeatureDefinition {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	fd := C.OGR_FD_Create(cName)
	return FeatureDefinition{fd}
}

// Destroy a feature definition object
func (fd FeatureDefinition) Destroy() {
	C.OGR_FD_Destroy(fd.cval)
}

// Release drops a reference, and delete object if no references remain.
func (fd FeatureDefinition) Release() {
	C.OGR_FD_Release(fd.cval)
}

// Name returns the name of this feature definition.
func (fd FeatureDefinition) Name() string {
	name := C.OGR_FD_GetName(fd.cval)
	return C.GoString(name)
}

// FieldCount returns the number of fields in the feature definition.
func (fd FeatureDefinition) FieldCount() int {
	count := C.OGR_FD_GetFieldCount(fd.cval)
	return int(count)
}

// FieldDefinition returns the definition of the indicated field.
func (fd FeatureDefinition) FieldDefinition(index int) FieldDefinition {
	fieldDefn := C.OGR_FD_GetFieldDefn(fd.cval, C.int(index))
	return FieldDefinition{fieldDefn}
}

// FieldIndex returns the index of the named field.
func (fd FeatureDefinition) FieldIndex(name string) int {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	index := C.OGR_FD_GetFieldIndex(fd.cval, cName)
	return int(index)
}

// AddFieldDefinition adds a new field definition to this feature definition.
func (fd FeatureDefinition) AddFieldDefinition(fieldDefn FieldDefinition) {
	C.OGR_FD_AddFieldDefn(fd.cval, fieldDefn.cval)
}

// DeleteFieldDefinition deletes a field definition from this feature definition.
func (fd FeatureDefinition) DeleteFieldDefinition(index int) error {
	return ErrFromOGRErr(C.OGR_FD_DeleteFieldDefn(fd.cval, C.int(index)))
}

// GeometryType returns the geometry base type of this feature definition.
func (fd FeatureDefinition) GeometryType() GeometryType {
	gt := C.OGR_FD_GetGeomType(fd.cval)
	return GeometryType(gt)
}

// SetGeometryType sets the geometry base type for this feature definition.
func (fd FeatureDefinition) SetGeometryType(geomType GeometryType) {
	C.OGR_FD_SetGeomType(fd.cval, C.OGRwkbGeometryType(geomType))
}

// IsGeometryIgnored returns if the geometry can be ignored when fetching features.
func (fd FeatureDefinition) IsGeometryIgnored() bool {
	isIgnored := C.OGR_FD_IsGeometryIgnored(fd.cval)
	return isIgnored != 0
}

// SetGeometryIgnored sets whether the geometry can be ignored when fetching features.
func (fd FeatureDefinition) SetGeometryIgnored(val bool) {
	C.OGR_FD_SetGeometryIgnored(fd.cval, BoolToCInt(val))
}

// IsStyleIgnored returns if the style can be ignored when fetching features.
func (fd FeatureDefinition) IsStyleIgnored() bool {
	isIgnored := C.OGR_FD_IsStyleIgnored(fd.cval)
	return isIgnored != 0
}

// SetStyleIgnored sets whether the style can be ignored when fetching features.
func (fd FeatureDefinition) SetStyleIgnored(val bool) {
	C.OGR_FD_SetStyleIgnored(fd.cval, BoolToCInt(val))
}

// Reference increments the reference count by one.
func (fd FeatureDefinition) Reference() int {
	count := C.OGR_FD_Reference(fd.cval)
	return int(count)
}

// Dereference decrements the reference count by one.
func (fd FeatureDefinition) Dereference() int {
	count := C.OGR_FD_Dereference(fd.cval)
	return int(count)
}

// ReferenceCount returns the current reference count.
func (fd FeatureDefinition) ReferenceCount() int {
	count := C.OGR_FD_GetReferenceCount(fd.cval)
	return int(count)
}

/* -------------------------------------------------------------------- */
/*      Feature functions                                               */
/* -------------------------------------------------------------------- */

// Feature is an exported GDAL/OGR type.
type Feature struct {
	cval C.OGRFeatureH
}

// Create a feature from this feature definition
func (fd FeatureDefinition) Create() Feature {
	feature := C.OGR_F_Create(fd.cval)
	return Feature{feature}
}

// Destroy this feature
func (feature Feature) Destroy() {
	C.OGR_F_Destroy(feature.cval)
}

// Definition returns feature definition.
func (feature Feature) Definition() FeatureDefinition {
	fd := C.OGR_F_GetDefnRef(feature.cval)
	return FeatureDefinition{fd}
}

// SetGeometry sets feature geometry.
func (feature Feature) SetGeometry(geom Geometry) error {
	return ErrFromOGRErr(C.OGR_F_SetGeometry(feature.cval, geom.cval))
}

// SetGeometryDirectly sets feature geometry, passing ownership to the feature.
func (feature Feature) SetGeometryDirectly(geom Geometry) error {
	return ErrFromOGRErr(C.OGR_F_SetGeometryDirectly(feature.cval, geom.cval))
}

// Geometry returns geometry of this feature.
func (feature Feature) Geometry() Geometry {
	geom := C.OGR_F_GetGeometryRef(feature.cval)
	return Geometry{geom}
}

// StealGeometry returns geometry of this feature and assume ownership.
func (feature Feature) StealGeometry() Geometry {
	geom := C.OGR_F_StealGeometry(feature.cval)
	return Geometry{geom}
}

// Clone duplicates feature.
func (feature Feature) Clone() Feature {
	newFeature := C.OGR_F_Clone(feature.cval)
	return Feature{newFeature}
}

// Equal reports whether two features are the same.
func (feature Feature) Equal(f2 Feature) bool {
	equal := C.OGR_F_Equal(feature.cval, f2.cval)
	return equal != 0
}

// FieldCount returns number of fields on this feature.
func (feature Feature) FieldCount() int {
	count := C.OGR_F_GetFieldCount(feature.cval)
	return int(count)
}

// FieldDefinition returns definition for the indicated field.
func (feature Feature) FieldDefinition(index int) FieldDefinition {
	defn := C.OGR_F_GetFieldDefnRef(feature.cval, C.int(index))
	return FieldDefinition{defn}
}

// FieldIndex returns the field index for the given field name.
func (feature Feature) FieldIndex(name string) int {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	index := C.OGR_F_GetFieldIndex(feature.cval, cName)
	return int(index)
}

// IsFieldSet returns if a field has ever been assigned a value.
func (feature Feature) IsFieldSet(index int) bool {
	set := C.OGR_F_IsFieldSet(feature.cval, C.int(index))
	return set != 0
}

// UnnsetField wraps the corresponding GDAL/OGR operation.
func (feature Feature) UnnsetField(index int) {
	C.OGR_F_UnsetField(feature.cval, C.int(index))
}

// RawField returns a reference to the internal field value.
func (feature Feature) RawField(index int) Field {
	field := C.OGR_F_GetRawFieldRef(feature.cval, C.int(index))
	return Field{field}
}

// FieldAsInteger returns field value as integer.
func (feature Feature) FieldAsInteger(index int) int {
	val := C.OGR_F_GetFieldAsInteger(feature.cval, C.int(index))
	return int(val)
}

// FieldAsInteger64 returns field value as 64-bit integer.
func (feature Feature) FieldAsInteger64(index int) int64 {
	val := C.OGR_F_GetFieldAsInteger64(feature.cval, C.int(index))
	return int64(val)
}

// FieldAsFloat64 returns field value as float64.
func (feature Feature) FieldAsFloat64(index int) float64 {
	val := C.OGR_F_GetFieldAsDouble(feature.cval, C.int(index))
	return float64(val)
}

// FieldAsString returns field value as string.
func (feature Feature) FieldAsString(index int) string {
	val := C.OGR_F_GetFieldAsString(feature.cval, C.int(index))
	return C.GoString(val)
}

// FieldAsIntegerList returns field as list of integers.
func (feature Feature) FieldAsIntegerList(index int) []int {
	var count C.int
	cArray := C.OGR_F_GetFieldAsIntegerList(feature.cval, C.int(index), &count)
	return copyCIntArray(cArray, count)
}

// FieldAsInteger64List returns field as list of 64-bit integers.
func (feature Feature) FieldAsInteger64List(index int) []int64 {
	var count C.int
	cArray := C.OGR_F_GetFieldAsInteger64List(feature.cval, C.int(index), &count)
	return copyCGIntBigArray(cArray, count)
}

// FieldAsFloat64List returns field as list of float64.
func (feature Feature) FieldAsFloat64List(index int) []float64 {
	var count C.int
	cArray := C.OGR_F_GetFieldAsDoubleList(feature.cval, C.int(index), &count)
	return copyCDoubleArray(cArray, count)
}

// FieldAsStringList returns field as list of strings.
func (feature Feature) FieldAsStringList(index int) []string {
	return cStringListToSlice(C.OGR_F_GetFieldAsStringList(feature.cval, C.int(index)))
}

// FieldAsBinary returns field as binary data.
func (feature Feature) FieldAsBinary(index int) []uint8 {
	var count C.int
	cArray := C.OGR_F_GetFieldAsBinary(feature.cval, C.int(index), &count)
	return copyCUCharArray(cArray, count)
}

// FieldAsDateTime returns field as date and time.
func (feature Feature) FieldAsDateTime(index int) (time.Time, bool) {
	var year, month, day, hour, minute, second, tzFlag C.int
	success := C.OGR_F_GetFieldAsDateTime(
		feature.cval,
		C.int(index),
		&year,
		&month,
		&day,
		&hour,
		&minute,
		&second,
		&tzFlag,
	)
	t := time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, time.UTC)
	return t, success != 0
}

// SetFieldInteger sets field to integer value.
func (feature Feature) SetFieldInteger(index, value int) {
	C.OGR_F_SetFieldInteger(feature.cval, C.int(index), C.int(value))
}

// SetFieldInteger64 sets field to 64-bit integer value.
func (feature Feature) SetFieldInteger64(index int, value int64) {
	C.OGR_F_SetFieldInteger64(feature.cval, C.int(index), C.GIntBig(value))
}

// SetFieldFloat64 sets field to float64 value.
func (feature Feature) SetFieldFloat64(index int, value float64) {
	C.OGR_F_SetFieldDouble(feature.cval, C.int(index), C.double(value))
}

// SetFieldString sets field to string value.
func (feature Feature) SetFieldString(index int, value string) {
	cVal := C.CString(value)
	defer C.free(unsafe.Pointer(cVal))
	C.OGR_F_SetFieldString(feature.cval, C.int(index), cVal)
}

// SetFieldIntegerList sets field to list of integers.
func (feature Feature) SetFieldIntegerList(index int, value []int) {
	cValue := IntSliceToCInt(value)
	C.OGR_F_SetFieldIntegerList(
		feature.cval,
		C.int(index),
		C.int(len(value)),
		cIntSlicePtr(cValue),
	)
}

// SetFieldInteger64List sets field to list of 64-bit integers.
func (feature Feature) SetFieldInteger64List(index int, value []int64) {
	cValue := int64SliceToCGIntBig(value)
	C.OGR_F_SetFieldInteger64List(
		feature.cval,
		C.int(index),
		C.int(len(value)),
		cGIntBigSlicePtr(cValue),
	)
}

// SetFieldFloat64List sets field to list of float64.
func (feature Feature) SetFieldFloat64List(index int, value []float64) {
	C.OGR_F_SetFieldDoubleList(
		feature.cval,
		C.int(index),
		C.int(len(value)),
		float64SlicePtr(value),
	)
}

// SetFieldStringList sets field to list of strings.
func (feature Feature) SetFieldStringList(index int, value []string) {
	length := len(value)
	cValue := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cValue[i] = C.CString(value[i])
		defer C.free(unsafe.Pointer(cValue[i]))
	}
	cValue[length] = (*C.char)(unsafe.Pointer(nil))

	C.OGR_F_SetFieldStringList(
		feature.cval,
		C.int(index),
		(**C.char)(unsafe.Pointer(&cValue[0])),
	)
}

// SetFieldRaw sets field from the raw field pointer.
func (feature Feature) SetFieldRaw(index int, field Field) {
	C.OGR_F_SetFieldRaw(feature.cval, C.int(index), field.cval)
}

// SetFieldBinary sets field as binary data.
func (feature Feature) SetFieldBinary(index int, value []uint8) {
	C.OGR_F_SetFieldBinary(
		feature.cval,
		C.int(index),
		C.int(len(value)),
		byteSlicePtr(value),
	)
}

// SetFieldDateTime sets field as date / time.
func (feature Feature) SetFieldDateTime(index int, dt time.Time) {
	C.OGR_F_SetFieldDateTime(
		feature.cval,
		C.int(index),
		C.int(dt.Year()),
		C.int(dt.Month()),
		C.int(dt.Day()),
		C.int(dt.Hour()),
		C.int(dt.Minute()),
		C.int(dt.Second()),
		C.int(1),
	)
}

// FID returns feature indentifier.
func (feature Feature) FID() int64 {
	fid := C.OGR_F_GetFID(feature.cval)
	return int64(fid)
}

// SetFID sets feature identifier.
func (feature Feature) SetFID(fid int64) error {
	return ErrFromOGRErr(C.OGR_F_SetFID(feature.cval, C.GIntBig(fid)))
}

// Unimplemented: DumpReadable

// SetFrom sets one feature from another.
func (feature Feature) SetFrom(other Feature, forgiving int) error {
	return ErrFromOGRErr(C.OGR_F_SetFrom(feature.cval, other.cval, C.int(forgiving)))
}

// SetFromWithMap sets one feature from another, using field map.
func (feature Feature) SetFromWithMap(other Feature, forgiving int, fieldMap []int) error {
	if len(fieldMap) == 0 {
		return fmt.Errorf("fieldMap must not be empty")
	}
	cFieldMap := IntSliceToCInt(fieldMap)
	return ErrFromOGRErr(C.OGR_F_SetFromWithMap(
		feature.cval,
		other.cval,
		C.int(forgiving),
		cIntSlicePtr(cFieldMap),
	))
}

// StlyeString returns style string for this feature.
func (feature Feature) StlyeString() string {
	style := C.OGR_F_GetStyleString(feature.cval)
	return C.GoString(style)
}

// SetStyleString sets style string for this feature.
func (feature Feature) SetStyleString(style string) {
	cStyle := C.CString(style)
	C.OGR_F_SetStyleStringDirectly(feature.cval, cStyle)
}

// IsNull returns true if this contains a null pointer.
func (feature Feature) IsNull() bool {
	return feature.cval == nil
}

/* -------------------------------------------------------------------- */
/*      Layer functions                                                 */
/* -------------------------------------------------------------------- */

// Layer is an exported GDAL/OGR type.
type Layer struct {
	cval C.OGRLayerH
}

// Name returns the layer name.
func (layer Layer) Name() string {
	name := C.OGR_L_GetName(layer.cval)
	return C.GoString(name)
}

// Type returns the layer geometry type.
func (layer Layer) Type() GeometryType {
	gt := C.OGR_L_GetGeomType(layer.cval)
	return GeometryType(gt)
}

// SpatialFilter returns the current spatial filter for this layer.
func (layer Layer) SpatialFilter() Geometry {
	geom := C.OGR_L_GetSpatialFilter(layer.cval)
	return Geometry{geom}
}

// SetSpatialFilter sets a new spatial filter for this layer.
func (layer Layer) SetSpatialFilter(filter Geometry) {
	C.OGR_L_SetSpatialFilter(layer.cval, filter.cval)
}

// SetSpatialFilterRect sets a new rectangular spatial filter for this layer.
func (layer Layer) SetSpatialFilterRect(minX, minY, maxX, maxY float64) {
	C.OGR_L_SetSpatialFilterRect(
		layer.cval,
		C.double(minX), C.double(minY), C.double(maxX), C.double(maxY),
	)
}

// SetAttributeFilter sets a new attribute query filter.
func (layer Layer) SetAttributeFilter(filter string) error {
	cFilter := C.CString(filter)
	defer C.free(unsafe.Pointer(cFilter))
	return ErrFromOGRErr(C.OGR_L_SetAttributeFilter(layer.cval, cFilter))
}

// ResetReading resets reading to start on the first featre.
func (layer Layer) ResetReading() {
	C.OGR_L_ResetReading(layer.cval)
}

// NextFeature returns the next available feature from this layer.
func (layer Layer) NextFeature() *Feature {
	feature := C.OGR_L_GetNextFeature(layer.cval)
	if feature == nil {
		return nil
	}
	return &Feature{feature}
}

// SetNextByIndex moves read cursor to the provided index.
func (layer Layer) SetNextByIndex(index int64) error {
	return ErrFromOGRErr(C.OGR_L_SetNextByIndex(layer.cval, C.GIntBig(index)))
}

// Feature returns a feature by its index.
func (layer Layer) Feature(index int64) Feature {
	feature := C.OGR_L_GetFeature(layer.cval, C.GIntBig(index))
	return Feature{feature}
}

// SetFeature rewrites the provided feature.
func (layer Layer) SetFeature(feature Feature) error {
	return ErrFromOGRErr(C.OGR_L_SetFeature(layer.cval, feature.cval))
}

// Create and write a new feature within a layer
func (layer Layer) Create(feature Feature) error {
	return ErrFromOGRErr(C.OGR_L_CreateFeature(layer.cval, feature.cval))
}

// Delete indicated feature from layer
func (layer Layer) Delete(index int64) error {
	return ErrFromOGRErr(C.OGR_L_DeleteFeature(layer.cval, C.GIntBig(index)))
}

// Definition returns the schema information for this layer.
func (layer Layer) Definition() FeatureDefinition {
	defn := C.OGR_L_GetLayerDefn(layer.cval)
	return FeatureDefinition{defn}
}

// SpatialReference returns the spatial reference system for this layer.
func (layer Layer) SpatialReference() SpatialReference {
	sr := C.OGR_L_GetSpatialRef(layer.cval)
	return SpatialReference{sr}
}

// FeatureCount returns the feature count for this layer.
func (layer Layer) FeatureCount(force bool) (count int, ok bool) {
	count = int(C.OGR_L_GetFeatureCount(layer.cval, BoolToCInt(force)))
	return count, count != -1
}

// Extent returns the extent of this layer.
func (layer Layer) Extent(force bool) (env Envelope, err error) {
	err = ErrFromOGRErr(C.OGR_L_GetExtent(layer.cval, &env.cval, BoolToCInt(force)))
	return
}

// TestCapability reports whether this layer supports the named capability.
func (layer Layer) TestCapability(capability string) bool {
	cString := C.CString(capability)
	defer C.free(unsafe.Pointer(cString))
	val := C.OGR_L_TestCapability(layer.cval, cString)
	return val != 0
}

// CreateField creates a new field on a layer.
func (layer Layer) CreateField(fd FieldDefinition, approxOK bool) error {
	return ErrFromOGRErr(C.OGR_L_CreateField(layer.cval, fd.cval, BoolToCInt(approxOK)))
}

// DeleteField deletes a field from the layer.
func (layer Layer) DeleteField(index int) error {
	return ErrFromOGRErr(C.OGR_L_DeleteField(layer.cval, C.int(index)))
}

// ReorderFields wraps the corresponding GDAL/OGR operation.
func (layer Layer) ReorderFields(layerMap []int) error {
	if len(layerMap) == 0 {
		return fmt.Errorf("layerMap must not be empty")
	}
	cLayerMap := IntSliceToCInt(layerMap)
	return ErrFromOGRErr(C.OGR_L_ReorderFields(layer.cval, cIntSlicePtr(cLayerMap)))
}

// ReorderField wraps the corresponding GDAL/OGR operation.
func (layer Layer) ReorderField(oldIndex, newIndex int) error {
	return ErrFromOGRErr(C.OGR_L_ReorderField(layer.cval, C.int(oldIndex), C.int(newIndex)))
}

// AlterFieldDefn wraps the corresponding GDAL/OGR operation.
func (layer Layer) AlterFieldDefn(index int, newDefn FieldDefinition, flags int) error {
	return ErrFromOGRErr(C.OGR_L_AlterFieldDefn(layer.cval, C.int(index), newDefn.cval, C.int(flags)))
}

// StartTransaction begins a transation on data sources which support it.
func (layer Layer) StartTransaction() error {
	return ErrFromOGRErr(C.OGR_L_StartTransaction(layer.cval))
}

// CommitTransaction commits a transaction on data sources which support it.
func (layer Layer) CommitTransaction() error {
	return ErrFromOGRErr(C.OGR_L_CommitTransaction(layer.cval))
}

// RollbackTransaction rolls back the current transaction on data sources which support it.
func (layer Layer) RollbackTransaction() error {
	return ErrFromOGRErr(C.OGR_L_RollbackTransaction(layer.cval))
}

// Sync flushes pending changes to the layer.
func (layer Layer) Sync() error {
	return ErrFromOGRErr(C.OGR_L_SyncToDisk(layer.cval))
}

// FIDColumn returns the name of the FID column.
func (layer Layer) FIDColumn() string {
	name := C.OGR_L_GetFIDColumn(layer.cval)
	return C.GoString(name)
}

// GeometryColumn returns the name of the geometry column.
func (layer Layer) GeometryColumn() string {
	name := C.OGR_L_GetGeometryColumn(layer.cval)
	return C.GoString(name)
}

// SetIgnoredFields sets which fields can be ignored when retrieving features from the layer.
func (layer Layer) SetIgnoredFields(names []string) error {
	length := len(names)
	cNames := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cNames[i] = C.CString(names[i])
		defer C.free(unsafe.Pointer(cNames[i]))
	}
	cNames[length] = (*C.char)(unsafe.Pointer(nil))

	return ErrFromOGRErr(C.OGR_L_SetIgnoredFields(layer.cval, (**C.char)(unsafe.Pointer(&cNames[0]))))
}

// Return the intersection of two layers
// Unimplemented: Intersection
// Will be new in 2.0

// Return the union of two layers
// Unimplemented: Union
// Will be new in 2.0

// Return the symmetric difference of two layers
// Unimplemented: SymDifference
// Will be new in 2.0

// Identify features in this layer with ones from the provided layer
// Unimplemented: Identity
// Will be new in 2.0

// Update this layer with features from the provided layer
// Unimplemented: Update
// Will be new in 2.0

// Clip off areas that are not covered by the provided layer
// Unimplemented: Clip
// Will be new in 2.0

// Remove areas that are covered by the provided layer
// Unimplemented: Erase
// Will be new in 2.0

/* -------------------------------------------------------------------- */
/*      Data source functions                                           */
/* -------------------------------------------------------------------- */

// DataSource is an exported GDAL/OGR type.
type DataSource struct {
	cval C.OGRDataSourceH
}

// OpenDataSource opens a file / data source with one of the registered drivers.
func OpenDataSource(name string, update int) DataSource {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ds := C.OGROpen(cName, C.int(update), nil)
	return DataSource{ds}
}

// OpenSharedDataSource opens a shared file / data source with one of the registered drivers.
func OpenSharedDataSource(name string, update int) DataSource {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ds := C.OGROpenShared(cName, C.int(update), nil)
	return DataSource{ds}
}

// Release drops a reference to this datasource and destroy if reference is zero.
func (ds DataSource) Release() error {
	return ErrFromOGRErr(C.OGRReleaseDataSource(ds.cval))
}

// OpenDataSourceCount returns the number of opened data sources.
func OpenDataSourceCount() int {
	count := C.OGRGetOpenDSCount()
	return int(count)
}

// OpenDataSourceByIndex returns the i'th datasource opened.
func OpenDataSourceByIndex(index int) DataSource {
	ds := C.OGRGetOpenDS(C.int(index))
	return DataSource{ds}
}

// Destroy closes datasource and releases resources.
func (ds DataSource) Destroy() {
	C.OGR_DS_Destroy(ds.cval)
}

// Name returns the name of the data source.
func (ds DataSource) Name() string {
	name := C.OGR_DS_GetName(ds.cval)
	return C.GoString(name)
}

// LayerCount returns the number of layers in this data source.
func (ds DataSource) LayerCount() int {
	count := C.OGR_DS_GetLayerCount(ds.cval)
	return int(count)
}

// LayerByIndex returns a layer of this data source by index.
func (ds DataSource) LayerByIndex(index int) Layer {
	layer := C.OGR_DS_GetLayer(ds.cval, C.int(index))
	return Layer{layer}
}

// LayerByName returns a layer of this data source by name.
func (ds DataSource) LayerByName(name string) Layer {
	cString := C.CString(name)
	defer C.free(unsafe.Pointer(cString))
	layer := C.OGR_DS_GetLayerByName(ds.cval, cString)
	return Layer{layer}
}

// Delete the layer from the data source
func (ds DataSource) Delete(index int) error {
	return ErrFromOGRErr(C.OGR_DS_DeleteLayer(ds.cval, C.int(index)))
}

// Driver returns the driver that the data source was opened with.
func (ds DataSource) Driver() OGRDriver {
	driver := C.OGR_DS_GetDriver(ds.cval)
	return OGRDriver{driver}
}

// CreateLayer creates a new layer on the data source.
func (ds DataSource) CreateLayer(
	name string,
	sr SpatialReference,
	geomType GeometryType,
	options []string,
) Layer {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	layer := C.OGR_DS_CreateLayer(
		ds.cval,
		cName,
		sr.cval,
		C.OGRwkbGeometryType(geomType),
		(**C.char)(unsafe.Pointer(&opts[0])),
	)
	return Layer{layer}
}

// CopyLayer duplicates an existing layer.
func (ds DataSource) CopyLayer(
	source Layer,
	name string,
	options []string,
) Layer {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	layer := C.OGR_DS_CopyLayer(
		ds.cval,
		source.cval,
		cName,
		(**C.char)(unsafe.Pointer(&opts[0])),
	)
	return Layer{layer}
}

// TestCapability reports whether the data source has the indicated capability.
func (ds DataSource) TestCapability(capability string) bool {
	cString := C.CString(capability)
	defer C.free(unsafe.Pointer(cString))
	val := C.OGR_DS_TestCapability(ds.cval, cString)
	return val != 0
}

// ExecuteSQL wraps the corresponding GDAL/OGR operation.
func (ds DataSource) ExecuteSQL(sql string, filter Geometry, dialect string) Layer {
	cSQL := C.CString(sql)
	defer C.free(unsafe.Pointer(cSQL))
	cDialect := C.CString(dialect)
	defer C.free(unsafe.Pointer(cDialect))

	layer := C.OGR_DS_ExecuteSQL(ds.cval, cSQL, filter.cval, cDialect)
	return Layer{layer}
}

// ReleaseResultSet wraps the corresponding GDAL/OGR operation.
func (ds DataSource) ReleaseResultSet(layer Layer) {
	C.OGR_DS_ReleaseResultSet(ds.cval, layer.cval)
}

// Sync flushes pending changes to the data source.
func (ds DataSource) Sync() error {
	return ErrFromOGRErr(C.OGR_DS_SyncToDisk(ds.cval))
}

/* -------------------------------------------------------------------- */
/*      Driver functions                                                */
/* -------------------------------------------------------------------- */

// OGRDriver is an exported GDAL/OGR type.
type OGRDriver struct {
	cval C.OGRSFDriverH
}

// Name returns name of driver (file format).
func (driver OGRDriver) Name() string {
	name := C.OGR_Dr_GetName(driver.cval)
	return C.GoString(name)
}

// Open attempts to open file with this driver.
func (driver OGRDriver) Open(filename string, update int) (newDS DataSource, ok bool) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	ds := C.OGR_Dr_Open(driver.cval, cFilename, C.int(update))
	return DataSource{ds}, ds != nil
}

// TestCapability reports whether this driver supports the named capability.
func (driver OGRDriver) TestCapability(capability string) bool {
	cString := C.CString(capability)
	defer C.free(unsafe.Pointer(cString))
	val := C.OGR_Dr_TestCapability(driver.cval, cString)
	return val != 0
}

// Create a new data source based on this driver
func (driver OGRDriver) Create(name string, options []string) (newDS DataSource, ok bool) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	ds := C.OGR_Dr_CreateDataSource(driver.cval, cName, (**C.char)(unsafe.Pointer(&opts[0])))
	return DataSource{ds}, ds != nil
}

// Copy creates a new datasource with this driver by copying all layers of the existing datasource.
func (driver OGRDriver) Copy(source DataSource, name string, options []string) (newDS DataSource, ok bool) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	ds := C.OGR_Dr_CopyDataSource(driver.cval, source.cval, cName, (**C.char)(unsafe.Pointer(&opts[0])))
	return DataSource{ds}, ds != nil
}

// Delete a data source
func (driver OGRDriver) Delete(filename string) error {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	return ErrFromOGRErr(C.OGR_Dr_DeleteDataSource(driver.cval, cFilename))
}

// Register adds a driver to the list of registered drivers.
func (driver OGRDriver) Register() {
	C.OGRRegisterDriver(driver.cval)
}

// Deregister removes a driver from the list of registered drivers.
func (driver OGRDriver) Deregister() {
	C.OGRDeregisterDriver(driver.cval)
}

// OGRDriverCount returns the number of registered drivers.
func OGRDriverCount() int {
	count := C.OGRGetDriverCount()
	return int(count)
}

// OGRDriverByIndex returns the indicated driver by index.
func OGRDriverByIndex(index int) OGRDriver {
	driver := C.OGRGetDriver(C.int(index))
	return OGRDriver{driver}
}

// OGRDriverByName returns the indicated driver by name.
func OGRDriverByName(name string) OGRDriver {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	driver := C.OGRGetDriverByName(cName)
	return OGRDriver{driver}
}

/* -------------------------------------------------------------------- */
/*      Style manager functions                                         */
/* -------------------------------------------------------------------- */

// StyleMgr is an exported GDAL/OGR type.
type StyleMgr struct {
	cval C.OGRStyleMgrH
}

// StyleTool is an exported GDAL/OGR type.
type StyleTool struct {
	cval C.OGRStyleToolH
}

// StyleTable is an exported GDAL/OGR type.
type StyleTable struct {
	cval C.OGRStyleTableH
}

// Unimplemented: CreateStyleManager

// Unimplemented: Destroy

// Unimplemented: InitFromFeature

// Unimplemented: InitStyleString

// Unimplemented: PartCount

// Unimplemented: PartCount

// Unimplemented: AddPart

// Unimplemented: AddStyle

// Unimplemented: CreateStyleTool

// Unimplemented: Destroy

// Unimplemented: Type

// Unimplemented: Unit

// Unimplemented: SetUnit

// Unimplemented: ParamStr

// Unimplemented: ParamNum

// Unimplemented: ParamDbl

// Unimplemented: SetParamStr

// Unimplemented: SetParamNum

// Unimplemented: SetParamDbl

// Unimplemented: StyleString

// Unimplemented: RGBFromString

// Unimplemented: CreateStyleTable

// Unimplemented: Destroy

// Unimplemented: Save

// Unimplemented: Load

// Unimplemented: Find

// Unimplemented: ResetStyleStringReading

// Unimplemented: NextStyle

// Unimplemented: LastStyleName
