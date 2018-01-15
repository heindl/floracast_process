#!/usr/bin/env bash

PA_SHAPE_COMBINED=/tmp/PADUS1_4Shapefile/PADUS1_4Combined.shp

if [[ -z "${1// }" ]]; then
    exit 0
fi

STATE_PATH="/tmp/gap_analysis/$1"
AREAS_PATH="$STATE_PATH/areas"

mkdir -p "$STATE_PATH"
mkdir -p "$AREAS_PATH"

# Input srs_def should be Albers Conical Equal Area. EPSG:4326 is expected by GeojSON.
#ogr2ogr -f 'ESRI Shapefile' \
#    -t_srs 'EPSG:4326' \
#    -where "d_State_Nm IN ('$1')" \
#    "$STATE_PATH/state.shp" \
#    "$PA_SHAPE_COMBINED"

ogr2ogr -f 'ESRI Shapefile' \
    -t_srs 'CRS:84' \
    -f GeoJSON \
    -where State_Nm="'${1}'" \
    "$STATE_PATH/state.geojson" \
    "$PA_SHAPE_COMBINED"




#shp2json "$STATE_PATH/state.shp" -o "$STATE_PATH/state.json"

#go run ./split.go --in "$STATE_PATH/state.json" --out "$AREAS_PATH/"