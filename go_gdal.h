// Copyright 2011 go-gdal. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#ifndef GO_GDAL_H_
#define GO_GDAL_H_

#include <gdal.h>
#include <gdal_version.h>
#include <gdal_alg.h>
#include <gdal_utils.h>
#include <gdalwarper.h>
#include <cpl_conv.h>
#include <ogr_srs_api.h>
#include <string.h>

// transform GDALProgressFunc to go func
GDALProgressFunc goGDALProgressFuncProxyB();

static inline GDALGridInverseDistanceToAPowerOptions goGDALGridInverseDistanceToAPowerOptionsInit()
{
    GDALGridInverseDistanceToAPowerOptions options;
    memset(&options, 0, sizeof(options));
#if GDAL_VERSION_NUM >= 3060000
    options.nSizeOfStructure = sizeof(options);
#endif
    return options;
}

static inline GDALGridInverseDistanceToAPowerNearestNeighborOptions goGDALGridInverseDistanceToAPowerNearestNeighborOptionsInit()
{
    GDALGridInverseDistanceToAPowerNearestNeighborOptions options;
    memset(&options, 0, sizeof(options));
#if GDAL_VERSION_NUM >= 3060000
    options.nSizeOfStructure = sizeof(options);
#endif
    return options;
}

static inline GDALGridMovingAverageOptions goGDALGridMovingAverageOptionsInit()
{
    GDALGridMovingAverageOptions options;
    memset(&options, 0, sizeof(options));
#if GDAL_VERSION_NUM >= 3060000
    options.nSizeOfStructure = sizeof(options);
#endif
    return options;
}

static inline GDALGridNearestNeighborOptions goGDALGridNearestNeighborOptionsInit()
{
    GDALGridNearestNeighborOptions options;
    memset(&options, 0, sizeof(options));
#if GDAL_VERSION_NUM >= 3060000
    options.nSizeOfStructure = sizeof(options);
#endif
    return options;
}

static inline GDALGridDataMetricsOptions goGDALGridDataMetricsOptionsInit()
{
    GDALGridDataMetricsOptions options;
    memset(&options, 0, sizeof(options));
#if GDAL_VERSION_NUM >= 3060000
    options.nSizeOfStructure = sizeof(options);
#endif
    return options;
}

static inline GDALGridLinearOptions goGDALGridLinearOptionsInit()
{
    GDALGridLinearOptions options;
    memset(&options, 0, sizeof(options));
#if GDAL_VERSION_NUM >= 3060000
    options.nSizeOfStructure = sizeof(options);
#endif
    return options;
}

#endif // GO_GDAL_H_


