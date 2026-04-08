package gdal

/*
#include "go_gdal.h"
#include "gdal_version.h"
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

/* -------------------------------------------------------------------- */
/*      Spatial reference functions.                                    */
/* -------------------------------------------------------------------- */

// SpatialReference wraps OGRSpatialReferenceH.
type SpatialReference struct {
	cval C.OGRSpatialReferenceH
}

// CreateSpatialReference creates a new spatial reference, optionally seeded
// from WKT.
func CreateSpatialReference(wkt string) SpatialReference {
	cString := C.CString(wkt)
	defer C.free(unsafe.Pointer(cString))
	sr := C.OSRNewSpatialReference(cString)
	return SpatialReference{sr}
}

// FromWKT initializes sr from WKT.
func (sr SpatialReference) FromWKT(wkt string) error {
	cString := C.CString(wkt)
	defer C.free(unsafe.Pointer(cString))
	return ErrFromOGRErr(C.OSRImportFromWkt(sr.cval, &cString))
}

// ToWKT exports sr as WKT.
func (sr SpatialReference) ToWKT() (string, error) {
	var p *C.char
	err := ErrFromOGRErr(C.OSRExportToWkt(sr.cval, &p))
	return goStringAndCPLFree(p), err
}

// ToPrettyWKT exports sr as formatted WKT.
func (sr SpatialReference) ToPrettyWKT(simplify bool) (string, error) {
	var p *C.char
	err := ErrFromOGRErr(C.OSRExportToPrettyWkt(
		sr.cval, &p, BoolToCInt(simplify),
	))
	return goStringAndCPLFree(p), err
}

// FromEPSG initializes sr from an EPSG code.
func (sr SpatialReference) FromEPSG(code int) error {
	return ErrFromOGRErr(C.OSRImportFromEPSG(sr.cval, C.int(code)))
}

// FromEPSGA initializes sr from an EPSG code using EPSG axis ordering.
func (sr SpatialReference) FromEPSGA(code int) error {
	return ErrFromOGRErr(C.OSRImportFromEPSGA(sr.cval, C.int(code)))
}

// Destroy releases sr.
func (sr SpatialReference) Destroy() {
	C.OSRDestroySpatialReference(sr.cval)
}

// Clone returns a copy of sr.
func (sr SpatialReference) Clone() SpatialReference {
	newSR := C.OSRClone(sr.cval)
	return SpatialReference{newSR}
}

// CloneGeogCS returns a copy of the geographic coordinate system in sr.
func (sr SpatialReference) CloneGeogCS() SpatialReference {
	newSR := C.OSRCloneGeogCS(sr.cval)
	return SpatialReference{newSR}
}

// Reference increments and returns sr's reference count.
func (sr SpatialReference) Reference() int {
	count := C.OSRReference(sr.cval)
	return int(count)
}

// Dereference decrements and returns sr's reference count.
func (sr SpatialReference) Dereference() int {
	count := C.OSRDereference(sr.cval)
	return int(count)
}

// Release decrements sr's reference count and destroys it when it reaches zero.
func (sr SpatialReference) Release() {
	C.OSRRelease(sr.cval)
}

// Validate reports whether sr contains a valid spatial reference definition.
func (sr SpatialReference) Validate() error {
	return ErrFromOGRErr(C.OSRValidate(sr.cval))
}

// FromProj4 initializes sr from a PROJ string.
func (sr SpatialReference) FromProj4(input string) error {
	cString := C.CString(input)
	defer C.free(unsafe.Pointer(cString))
	return ErrFromOGRErr(C.OSRImportFromProj4(sr.cval, cString))
}

// ToProj4 exports sr as a PROJ string.
func (sr SpatialReference) ToProj4() (string, error) {
	var p *C.char
	err := ErrFromOGRErr(C.OSRExportToProj4(sr.cval, &p))
	return goStringAndCPLFree(p), err
}

// FromESRI initializes sr from an ESRI projection string.
func (sr SpatialReference) FromESRI(input string) error {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	cLines := make([]*C.char, len(lines)+1)
	for i := range lines {
		cLines[i] = C.CString(lines[i])
		defer C.free(unsafe.Pointer(cLines[i]))
	}

	return ErrFromOGRErr(C.OSRImportFromESRI(sr.cval, &cLines[0]))
}

// FromPCI initializes sr from a PCI projection definition.
func (sr SpatialReference) FromPCI(proj, units string, params []float64) error {
	if len(params) < 17 {
		return fmt.Errorf("pci projection definition requires 17 parameters")
	}

	cProj := C.CString(proj)
	defer C.free(unsafe.Pointer(cProj))
	cUnits := C.CString(units)
	defer C.free(unsafe.Pointer(cUnits))

	return ErrFromOGRErr(C.OSRImportFromPCI(
		sr.cval,
		cProj,
		cUnits,
		(*C.double)(unsafe.Pointer(&params[0])),
	))
}

// FromUSGS initializes sr from a USGS projection definition.
func (sr SpatialReference) FromUSGS(projsys, zone int, params []float64, datum int) error {
	if len(params) < 15 {
		return fmt.Errorf("usgs projection definition requires 15 parameters")
	}

	return ErrFromOGRErr(C.OSRImportFromUSGS(
		sr.cval,
		C.long(projsys),
		C.long(zone),
		(*C.double)(unsafe.Pointer(&params[0])),
		C.long(datum),
	))
}

// FromXML initializes sr from XML, currently limited to GML.
func (sr SpatialReference) FromXML(xml string) error {
	cXML := C.CString(xml)
	defer C.free(unsafe.Pointer(cXML))
	return ErrFromOGRErr(C.OSRImportFromXML(sr.cval, cXML))
}

// FromERM initializes sr from an ERMapper projection definition.
func (sr SpatialReference) FromERM(proj, datum, units string) error {
	cProj := C.CString(proj)
	defer C.free(unsafe.Pointer(cProj))
	cDatum := C.CString(datum)
	defer C.free(unsafe.Pointer(cDatum))
	cUnits := C.CString(units)
	defer C.free(unsafe.Pointer(cUnits))

	return ErrFromOGRErr(C.OSRImportFromERM(sr.cval, cProj, cDatum, cUnits))
}

// FromURL initializes sr from a URL-backed definition.
func (sr SpatialReference) FromURL(url string) error {
	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	return ErrFromOGRErr(C.OSRImportFromUrl(sr.cval, cURL))
}

// ToPCI exports sr as a PCI projection definition.
func (sr SpatialReference) ToPCI() (proj, units string, params []float64, errVal error) {
	var p, u *C.char
	var cParams *C.double
	err := ErrFromOGRErr(C.OSRExportToPCI(
		sr.cval, &p, &u, &cParams,
	))
	// GDAL allocates cParams for the caller, so copy it into Go memory before
	// releasing the native buffer.
	params = copyCDoubleArray(cParams, 17)
	C.CPLFree(unsafe.Pointer(cParams))
	return goStringAndCPLFree(p), goStringAndCPLFree(u), params, err
}

// ToUSGS exports sr as a USGS GCTP projection definition.
func (sr SpatialReference) ToUSGS() (proj, zone int, params []float64, datum int, errVal error) {
	var cProj, cZone, cDatum C.long
	var cParams *C.double
	err := ErrFromOGRErr(C.OSRExportToUSGS(
		sr.cval,
		&cProj,
		&cZone,
		&cParams,
		&cDatum,
	))
	// GDAL allocates cParams for the caller, so copy it into Go memory before
	// releasing the native buffer.
	params = copyCDoubleArray(cParams, 15)
	C.CPLFree(unsafe.Pointer(cParams))
	proj = int(cProj)
	zone = int(cZone)
	datum = int(cDatum)

	return proj, zone, params, datum, err
}

// ToXML exports sr as XML.
func (sr SpatialReference) ToXML() (xml string, errVal error) {
	var x *C.char
	err := ErrFromOGRErr(C.OSRExportToXML(sr.cval, &x, nil))
	return goStringAndCPLFree(x), err
}

// ToMICoordSys exports sr in MapInfo CoordSys format.
func (sr SpatialReference) ToMICoordSys() (output string, errVal error) {
	var x *C.char
	err := ErrFromOGRErr(C.OSRExportToMICoordSys(sr.cval, &x))
	return goStringAndCPLFree(x), err
}

// Export coordinate system in ERMapper format
// Unimplemented: ToERM

// MorphToESRI converts in place to ESRI WKT format.
func (sr SpatialReference) MorphToESRI() error {
	return ErrFromOGRErr(C.OSRMorphToESRI(sr.cval))
}

// MorphFromESRI converts in place from ESRI WKT format.
func (sr SpatialReference) MorphFromESRI() error {
	return ErrFromOGRErr(C.OSRMorphFromESRI(sr.cval))
}

// AttrValue returns indicated attribute of named node.
func (sr SpatialReference) AttrValue(key string, child int) (value string, ok bool) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	val := C.OSRGetAttrValue(sr.cval, cKey, C.int(child))
	return C.GoString(val), val != nil
}

// SetAttrValue sets attribute value in spatial reference.
func (sr SpatialReference) SetAttrValue(path, value string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	return ErrFromOGRErr(C.OSRSetAttrValue(sr.cval, cPath, cValue))
}

// SetAngularUnits sets the angular units for the geographic coordinate system.
func (sr SpatialReference) SetAngularUnits(units string, radians float64) error {
	cUnits := C.CString(units)
	defer C.free(unsafe.Pointer(cUnits))
	return ErrFromOGRErr(C.OSRSetAngularUnits(sr.cval, cUnits, C.double(radians)))
}

// AngularUnits returns the angular units for the geographic coordinate system.
func (sr SpatialReference) AngularUnits() (string, float64) {
	var x *C.char
	factor := C.OSRGetAngularUnits(sr.cval, &x)
	defer C.free(unsafe.Pointer(x))
	return C.GoString(x), float64(factor)
}

// SetLinearUnits sets the linear units for the projection.
func (sr SpatialReference) SetLinearUnits(name string, toMeters float64) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ErrFromOGRErr(C.OSRSetLinearUnits(sr.cval, cName, C.double(toMeters)))
}

// SetTargetLinearUnits sets the linear units for the target node.
func (sr SpatialReference) SetTargetLinearUnits(target, units string, toMeters float64) error {
	cTarget := C.CString(target)
	defer C.free(unsafe.Pointer(cTarget))
	cUnits := C.CString(units)
	defer C.free(unsafe.Pointer(cUnits))
	return ErrFromOGRErr(C.OSRSetTargetLinearUnits(sr.cval, cTarget, cUnits, C.double(toMeters)))
}

// SetLinearUnitsAndUpdateParameters sets the linear units for the target node and update all existing linear parameters.
func (sr SpatialReference) SetLinearUnitsAndUpdateParameters(name string, toMeters float64) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ErrFromOGRErr(C.OSRSetLinearUnitsAndUpdateParameters(sr.cval, cName, C.double(toMeters)))
}

// LinearUnits returns linear projection units.
func (sr SpatialReference) LinearUnits() (string, float64) {
	var x *C.char
	factor := C.OSRGetLinearUnits(sr.cval, &x)
	defer C.free(unsafe.Pointer(x))
	return C.GoString(x), float64(factor)
}

// TargetLinearUnits returns linear units for target.
func (sr SpatialReference) TargetLinearUnits(target string) (string, float64) {
	cTarget := C.CString(target)
	defer C.free(unsafe.Pointer(cTarget))
	var x *C.char
	factor := C.OSRGetTargetLinearUnits(sr.cval, cTarget, &x)
	defer C.free(unsafe.Pointer(x))
	return C.GoString(x), float64(factor)
}

// PrimeMeridian returns prime meridian information.
func (sr SpatialReference) PrimeMeridian() (string, float64) {
	var x *C.char
	offset := C.OSRGetPrimeMeridian(sr.cval, &x)
	defer C.free(unsafe.Pointer(x))
	return C.GoString(x), float64(offset)
}

// IsGeographic reports whether geographic coordinate system.
func (sr SpatialReference) IsGeographic() bool {
	val := C.OSRIsGeographic(sr.cval)
	return val != 0
}

// IsLocal reports whether local coordinate system.
func (sr SpatialReference) IsLocal() bool {
	val := C.OSRIsLocal(sr.cval)
	return val != 0
}

// IsProjected reports whether projected coordinate system.
func (sr SpatialReference) IsProjected() bool {
	val := C.OSRIsProjected(sr.cval)
	return val != 0
}

// IsCompound reports whether compound coordinate system.
func (sr SpatialReference) IsCompound() bool {
	val := C.OSRIsCompound(sr.cval)
	return val != 0
}

// IsGeocentric reports whether geocentric coordinate system.
func (sr SpatialReference) IsGeocentric() bool {
	val := C.OSRIsGeocentric(sr.cval)
	return val != 0
}

// IsVertical reports whether vertical coordinate system.
func (sr SpatialReference) IsVertical() bool {
	val := C.OSRIsVertical(sr.cval)
	return val != 0
}

// IsSameGeographicCS reports whether the geographic coordinate systems match.
func (sr SpatialReference) IsSameGeographicCS(other SpatialReference) bool {
	val := C.OSRIsSameGeogCS(sr.cval, other.cval)
	return val != 0
}

// IsSameVerticalCS reports whether the vertical coordinate systems match.
func (sr SpatialReference) IsSameVerticalCS(other SpatialReference) bool {
	val := C.OSRIsSameVertCS(sr.cval, other.cval)
	return val != 0
}

// IsSame reports whether the coordinate systems describe the same system.
func (sr SpatialReference) IsSame(other SpatialReference) bool {
	val := C.OSRIsSame(sr.cval, other.cval)
	return val != 0
}

// SetLocalCS sets the user visible local CS name.
func (sr SpatialReference) SetLocalCS(name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ErrFromOGRErr(C.OSRSetLocalCS(sr.cval, cName))
}

// SetProjectedCS sets the user visible projected CS name.
func (sr SpatialReference) SetProjectedCS(name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ErrFromOGRErr(C.OSRSetProjCS(sr.cval, cName))
}

// SetGeocentricCS sets the user visible geographic CS name.
func (sr SpatialReference) SetGeocentricCS(name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ErrFromOGRErr(C.OSRSetGeocCS(sr.cval, cName))
}

// SetWellKnownGeographicCS sets geographic CS based on well known name.
func (sr SpatialReference) SetWellKnownGeographicCS(name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ErrFromOGRErr(C.OSRSetWellKnownGeogCS(sr.cval, cName))
}

// SetFromUserInput sets spatial reference from various text formats.
func (sr SpatialReference) SetFromUserInput(name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ErrFromOGRErr(C.OSRSetFromUserInput(sr.cval, cName))
}

// CopyGeographicCSFrom wraps the corresponding GDAL/OGR operation.
func (sr SpatialReference) CopyGeographicCSFrom(other SpatialReference) error {
	return ErrFromOGRErr(C.OSRCopyGeogCSFrom(sr.cval, other.cval))
}

// SetTOWGS84 sets the Bursa-Wolf conversion to WGS84.
func (sr SpatialReference) SetTOWGS84(dx, dy, dz, ex, ey, ez, ppm float64) error {
	return ErrFromOGRErr(C.OSRSetTOWGS84(
		sr.cval,
		C.double(dx),
		C.double(dy),
		C.double(dz),
		C.double(ex),
		C.double(ey),
		C.double(ez),
		C.double(ppm),
	))
}

// TOWGS84 returns the TOWGS84 parameters if available.
func (sr SpatialReference) TOWGS84() (coeff [7]float64, err error) {
	err = ErrFromOGRErr(C.OSRGetTOWGS84(sr.cval, (*C.double)(unsafe.Pointer(&coeff[0])), 7))
	return
}

// SetCompoundCS wraps the corresponding GDAL/OGR operation.
func (sr SpatialReference) SetCompoundCS(
	name string,
	horizontal, vertical SpatialReference,
) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ErrFromOGRErr(C.OSRSetCompoundCS(sr.cval, cName, horizontal.cval, vertical.cval))
}

// SetGeographicCS sets geographic coordinate system.
func (sr SpatialReference) SetGeographicCS(
	geogName, datumName, spheroidName string,
	semiMajor, flattening float64,
	pmName string,
	offset float64,
	angularUnits string,
	toRadians float64,
) error {
	cGeogName := C.CString(geogName)
	defer C.free(unsafe.Pointer(cGeogName))
	cDatumName := C.CString(datumName)
	defer C.free(unsafe.Pointer(cDatumName))
	cSpheroidName := C.CString(spheroidName)
	defer C.free(unsafe.Pointer(cSpheroidName))
	cPMName := C.CString(pmName)
	defer C.free(unsafe.Pointer(cPMName))
	cAngularUnits := C.CString(angularUnits)
	defer C.free(unsafe.Pointer(cAngularUnits))
	return ErrFromOGRErr(C.OSRSetGeogCS(
		sr.cval,
		cGeogName,
		cDatumName,
		cSpheroidName,
		C.double(semiMajor),
		C.double(flattening),
		cPMName,
		C.double(offset),
		cAngularUnits,
		C.double(toRadians),
	))
}

// SetVerticalCS sets up the vertical coordinate system.
func (sr SpatialReference) SetVerticalCS(csName, datumName string, datumType int) error {
	cCSName := C.CString(csName)
	defer C.free(unsafe.Pointer(cCSName))
	cDatumName := C.CString(datumName)
	defer C.free(unsafe.Pointer(cDatumName))
	return ErrFromOGRErr(C.OSRSetVertCS(sr.cval, cCSName, cDatumName, C.int(datumType)))
}

// SemiMajorAxis returns spheroid semi-major axis.
func (sr SpatialReference) SemiMajorAxis() (float64, error) {
	var cErr C.OGRErr
	axis := C.OSRGetSemiMajor(sr.cval, &cErr)
	return float64(axis), ErrFromOGRErr(cErr)
}

// SemiMinorAxis returns spheroid semi-minor axis.
func (sr SpatialReference) SemiMinorAxis() (float64, error) {
	var cErr C.OGRErr
	axis := C.OSRGetSemiMinor(sr.cval, &cErr)
	return float64(axis), ErrFromOGRErr(cErr)
}

// InverseFlattening returns spheroid inverse flattening axis.
func (sr SpatialReference) InverseFlattening() (float64, error) {
	var cErr C.OGRErr
	flat := C.OSRGetInvFlattening(sr.cval, &cErr)
	return float64(flat), ErrFromOGRErr(cErr)
}

// SetAuthority wraps the corresponding GDAL/OGR operation.
func (sr SpatialReference) SetAuthority(target, authority string, code int) error {
	cTarget := C.CString(target)
	defer C.free(unsafe.Pointer(cTarget))
	cAuthority := C.CString(authority)
	defer C.free(unsafe.Pointer(cAuthority))
	return ErrFromOGRErr(C.OSRSetAuthority(sr.cval, cTarget, cAuthority, C.int(code)))
}

// AuthorityCode returns the authority code for a node.
func (sr SpatialReference) AuthorityCode(target string) string {
	cTarget := C.CString(target)
	defer C.free(unsafe.Pointer(cTarget))
	code := C.OSRGetAuthorityCode(sr.cval, cTarget)
	return C.GoString(code)
}

// AuthorityName returns the authority name for a node.
func (sr SpatialReference) AuthorityName(target string) string {
	cTarget := C.CString(target)
	defer C.free(unsafe.Pointer(cTarget))
	code := C.OSRGetAuthorityName(sr.cval, cTarget)
	return C.GoString(code)
}

// SetProjectionByName sets a projection by name.
func (sr SpatialReference) SetProjectionByName(name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ErrFromOGRErr(C.OSRSetProjection(sr.cval, cName))
}

// SetProjectionParameter sets a projection parameter value.
func (sr SpatialReference) SetProjectionParameter(name string, value float64) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ErrFromOGRErr(C.OSRSetProjParm(sr.cval, cName, C.double(value)))
}

// ProjectionParameter returns a projection parameter value.
func (sr SpatialReference) ProjectionParameter(name string, defaultValue float64) (float64, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var cErr C.OGRErr
	value := C.OSRGetProjParm(sr.cval, cName, C.double(defaultValue), &cErr)
	return float64(value), ErrFromOGRErr(cErr)
}

// SetNormalizedProjectionParameter sets a projection parameter with a normalized value.
func (sr SpatialReference) SetNormalizedProjectionParameter(name string, value float64) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ErrFromOGRErr(C.OSRSetNormProjParm(sr.cval, cName, C.double(value)))
}

// NormalizedProjectionParameter returns a normalized projection parameter value.
func (sr SpatialReference) NormalizedProjectionParameter(
	name string, defaultValue float64,
) (float64, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var cErr C.OGRErr
	value := C.OSRGetProjParm(sr.cval, cName, C.double(defaultValue), &cErr)
	return float64(value), ErrFromOGRErr(cErr)
}

// SetUTM sets UTM projection definition.
func (sr SpatialReference) SetUTM(zone int, north bool) error {
	return ErrFromOGRErr(C.OSRSetUTM(sr.cval, C.int(zone), BoolToCInt(north)))
}

// UTMZone returns UTM zone information.
func (sr SpatialReference) UTMZone() (zone int, north bool) {
	var northInt C.int
	cZone := C.OSRGetUTMZone(sr.cval, &northInt)
	return int(cZone), northInt != 0
}

// SetStatePlane sets State Plane projection definition.
func (sr SpatialReference) SetStatePlane(zone int, nad83 bool) error {
	return ErrFromOGRErr(C.OSRSetStatePlane(sr.cval, C.int(zone), BoolToCInt(nad83)))
}

// SetStatePlaneWithUnits sets State Plane projection definition.
func (sr SpatialReference) SetStatePlaneWithUnits(
	zone int,
	nad83 bool,
	unitName string,
	factor float64,
) error {
	cUnitName := C.CString(unitName)
	defer C.free(unsafe.Pointer(cUnitName))
	return ErrFromOGRErr(C.OSRSetStatePlaneWithUnits(
		sr.cval,
		C.int(zone),
		BoolToCInt(nad83),
		cUnitName,
		C.double(factor),
	))
}

// AutoIdentifyEPSG sets EPSG authority info if possible.
func (sr SpatialReference) AutoIdentifyEPSG() error {
	return ErrFromOGRErr(C.OSRAutoIdentifyEPSG(sr.cval))
}

// EPSGTreatsAsLatLong reports whether EPSG feels this coordinate system should be treated as having lat/long coordinate ordering.
func (sr SpatialReference) EPSGTreatsAsLatLong() bool {
	val := C.OSREPSGTreatsAsLatLong(sr.cval)
	return val != 0
}

// Fetch the orientation of one axis
// Unimplemented: Axis

// SetACEA sets to Albers Conic Equal Area.
func (sr SpatialReference) SetACEA(
	stdp1, stdp2, centerLat, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetACEA(
		sr.cval,
		C.double(stdp1),
		C.double(stdp2),
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetAE sets to Azimuthal Equidistant.
func (sr SpatialReference) SetAE(centerLat, centerLong, falseEasting, falseNorthing float64) error {
	return ErrFromOGRErr(C.OSRSetAE(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetBonne sets to Bonne.
func (sr SpatialReference) SetBonne(standardParallel, centralMeridian, falseEasting, falseNorthing float64) error {
	return ErrFromOGRErr(C.OSRSetBonne(
		sr.cval,
		C.double(standardParallel),
		C.double(centralMeridian),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetCEA sets to Cylindrical Equal Area.
func (sr SpatialReference) SetCEA(stdp1, centralMeridian, falseEasting, falseNorthing float64) error {
	return ErrFromOGRErr(C.OSRSetCEA(
		sr.cval,
		C.double(stdp1),
		C.double(centralMeridian),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetCS sets to Cassini-Soldner.
func (sr SpatialReference) SetCS(centerLat, centerLong, falseEasting, falseNorthing float64) error {
	return ErrFromOGRErr(C.OSRSetCS(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetEC sets to Equidistant Conic.
func (sr SpatialReference) SetEC(
	stdp1, stdp2, centerLat, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetEC(
		sr.cval,
		C.double(stdp1),
		C.double(stdp2),
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetEckert sets to Eckert I-VI.
func (sr SpatialReference) SetEckert(variation int, centralMeridian, falseEasting, falseNorthing float64) error {
	return ErrFromOGRErr(C.OSRSetEckert(
		sr.cval,
		C.int(variation),
		C.double(centralMeridian),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetEquirectangular sets to Equirectangular.
func (sr SpatialReference) SetEquirectangular(
	centerLat, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetEquirectangular(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetEquirectangularGeneralized sets to Equirectangular (generalized form).
func (sr SpatialReference) SetEquirectangularGeneralized(
	centerLat, centerLong, psuedoStdParallel, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetEquirectangular2(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(psuedoStdParallel),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetGS sets to Gall Stereographic.
func (sr SpatialReference) SetGS(centralMeridian, falseEasting, falseNorthing float64) error {
	return ErrFromOGRErr(C.OSRSetGS(
		sr.cval,
		C.double(centralMeridian),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetGH sets to Goode Homolosine.
func (sr SpatialReference) SetGH(centralMeridian, falseEasting, falseNorthing float64) error {
	return ErrFromOGRErr(C.OSRSetGH(
		sr.cval,
		C.double(centralMeridian),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetIGH sets to Interrupted Goode Homolosine.
func (sr SpatialReference) SetIGH() error {
	return ErrFromOGRErr(C.OSRSetIGH(sr.cval))
}

// SetGEOS sets to GEOS - Geostationary Satellite View.
func (sr SpatialReference) SetGEOS(
	centralMeridian, satelliteHeight, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetGEOS(
		sr.cval,
		C.double(centralMeridian),
		C.double(satelliteHeight),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetGSTM sets to Gauss Schreiber Transverse Mercator.
func (sr SpatialReference) SetGSTM(
	centerLat, centerLong, scale, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetGaussSchreiberTMercator(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(scale),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetGnomonic sets to gnomonic.
func (sr SpatialReference) SetGnomonic(
	centerLat, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetGnomonic(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetHOM sets to Hotine Oblique Mercator projection using azimuth angle.
func (sr SpatialReference) SetHOM(
	centerLat, centerLong, azimuth, rectToSkew, scale, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetHOM(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(azimuth),
		C.double(rectToSkew),
		C.double(scale),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetHOM2PNO sets to Hotine Oblique Mercator projection using two points on projection centerline.
func (sr SpatialReference) SetHOM2PNO(
	centerLat, lat1, long1, lat2, long2, scale, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetHOM2PNO(
		sr.cval,
		C.double(centerLat),
		C.double(lat1),
		C.double(long1),
		C.double(lat2),
		C.double(long2),
		C.double(scale),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetIWMPolyconic sets to International Map of the World Polyconic.
func (sr SpatialReference) SetIWMPolyconic(
	lat1, lat2, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetIWMPolyconic(
		sr.cval,
		C.double(lat1),
		C.double(lat2),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetKrovak sets to Krovak Oblique Conic Conformal.
func (sr SpatialReference) SetKrovak(
	centerLat, centerLong, azimuth, psuedoStdParallel, scale, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetKrovak(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(azimuth),
		C.double(psuedoStdParallel),
		C.double(scale),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetLAEA sets to Lambert Azimuthal Equal Area.
func (sr SpatialReference) SetLAEA(
	centerLat, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetLAEA(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetLCC sets to Lambert Conformal Conic.
func (sr SpatialReference) SetLCC(
	stdp1, stdp2, centerLat, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetLCC(
		sr.cval,
		C.double(stdp1),
		C.double(stdp2),
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetLCC1SP sets to Lambert Conformal Conic (1 standard parallel).
func (sr SpatialReference) SetLCC1SP(
	centerLat, centerLong, scale, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetLCC1SP(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(scale),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetLCCB sets to Lambert Conformal Conic (Belgium).
func (sr SpatialReference) SetLCCB(
	stdp1, stdp2, centerLat, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetLCCB(
		sr.cval,
		C.double(stdp1),
		C.double(stdp2),
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetMC sets to Miller Cylindrical.
func (sr SpatialReference) SetMC(
	centerLat, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetMC(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetMercator sets to Mercator.
func (sr SpatialReference) SetMercator(
	centerLat, centerLong, scale, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetMercator(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(scale),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetMollweide sets tp Mollweide.
func (sr SpatialReference) SetMollweide(
	centralMeridian, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetMollweide(
		sr.cval,
		C.double(centralMeridian),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetNZMG sets to New Zealand Map Grid.
func (sr SpatialReference) SetNZMG(
	centerLat, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetNZMG(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetOS sets to Oblique Stereographic.
func (sr SpatialReference) SetOS(
	originLat, meridian, scale, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetOS(
		sr.cval,
		C.double(originLat),
		C.double(meridian),
		C.double(scale),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetOrthographic sets to Orthographic.
func (sr SpatialReference) SetOrthographic(
	centerLat, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetOrthographic(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetPolyconic sets to Polyconic.
func (sr SpatialReference) SetPolyconic(
	centerLat, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetPolyconic(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetPS sets to Polar Stereographic.
func (sr SpatialReference) SetPS(
	centerLat, centerLong, scale, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetPS(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(scale),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetRobinson sets to Robinson.
func (sr SpatialReference) SetRobinson(
	centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetRobinson(
		sr.cval,
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetSinusoidal sets to Sinusoidal.
func (sr SpatialReference) SetSinusoidal(
	centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetSinusoidal(
		sr.cval,
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetStereographic sets to Stereographic.
func (sr SpatialReference) SetStereographic(
	centerLat, centerLong, scale, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetStereographic(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(scale),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetSOC sets to Swiss Oblique Cylindrical.
func (sr SpatialReference) SetSOC(
	latitudeOfOrigin, centralMeridian, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetSOC(
		sr.cval,
		C.double(latitudeOfOrigin),
		C.double(centralMeridian),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetTM sets to Transverse Mercator.
func (sr SpatialReference) SetTM(
	centerLat, centerLong, scale, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetTM(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(scale),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetTMVariant sets to Transverse Mercator variant.
func (sr SpatialReference) SetTMVariant(
	variantName string, centerLat, centerLong, scale, falseEasting, falseNorthing float64,
) error {
	cName := C.CString(variantName)
	defer C.free(unsafe.Pointer(cName))
	return ErrFromOGRErr(C.OSRSetTMVariant(
		sr.cval,
		cName,
		C.double(centerLat),
		C.double(centerLong),
		C.double(scale),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetTMG sets to Tunisia Mining Grid.
func (sr SpatialReference) SetTMG(
	centerLat, centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetTMG(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetTMSO sets to Transverse Mercator (South Oriented).
func (sr SpatialReference) SetTMSO(
	centerLat, centerLong, scale, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetTMSO(
		sr.cval,
		C.double(centerLat),
		C.double(centerLong),
		C.double(scale),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// SetVDG sets to VanDerGrinten.
func (sr SpatialReference) SetVDG(
	centerLong, falseEasting, falseNorthing float64,
) error {
	return ErrFromOGRErr(C.OSRSetVDG(
		sr.cval,
		C.double(centerLong),
		C.double(falseEasting),
		C.double(falseNorthing),
	))
}

// CleanupSR wraps the corresponding GDAL/OGR operation.
func CleanupSR() {
	C.OSRCleanup()
}

/* -------------------------------------------------------------------- */
/*      Coordinate transformation functions.                            */
/* -------------------------------------------------------------------- */

// CoordinateTransform is an exported GDAL/OGR type.
type CoordinateTransform struct {
	cval C.OGRCoordinateTransformationH
}

// CreateCoordinateTransform creates a new CoordinateTransform.
func CreateCoordinateTransform(
	source SpatialReference,
	dest SpatialReference,
) CoordinateTransform {
	ct := C.OCTNewCoordinateTransformation(source.cval, dest.cval)
	return CoordinateTransform{ct}
}

// Destroy CoordinateTransform
func (ct CoordinateTransform) Destroy() {
	C.OCTDestroyCoordinateTransformation(ct.cval)
}

// Transform wraps the corresponding GDAL/OGR operation.
func (ct CoordinateTransform) Transform(numPoints int, xPoints []float64, yPoints []float64, zPoints []float64) bool {
	if numPoints < 0 {
		return false
	}
	if numPoints == 0 {
		return true
	}
	if len(xPoints) < numPoints || len(yPoints) < numPoints {
		return false
	}

	var zPtr *C.double
	if len(zPoints) > 0 {
		if len(zPoints) < numPoints {
			return false
		}
		zPtr = float64SlicePtr(zPoints)
	}

	val := C.OCTTransform(ct.cval, C.int(numPoints), float64SlicePtr(xPoints), float64SlicePtr(yPoints), zPtr)
	return int(val) != 0
}
