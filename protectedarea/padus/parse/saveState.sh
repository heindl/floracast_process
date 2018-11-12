#!/usr/bin/env bash

PA_SHAPE_COMBINED=/tmp/PADUS1_4Shapefile/PADUS1_4Combined.shp

if [[ -z "${1// }" ]]; then
    exit 0
fi

STATE_PATH="/tmp/gap_analysis/$1"
AREAS_PATH="$STATE_PATH/areas"

# Clear if exists
rm -rf "$AREAS_PATH"

mkdir -p "$STATE_PATH"
mkdir -p "$AREAS_PATH"


ogr2ogr -f 'ESRI Shapefile' \
    -f KML \
    -where State_Nm="'${1}'" \
    "$STATE_PATH/state.kml" \
    "$PA_SHAPE_COMBINED"

#ogr2ogr -f 'ESRI Shapefile' \
#    -t_srs 'CRS:84' \
#    -f GeoJSON \
#    -where State_Nm="'${1}'" \
#    "$STATE_PATH/state.geojson" \
#    "$PA_SHAPE_COMBINED"


#shp2json "$STATE_PATH/state.shp" -o "$STATE_PATH/state.json"