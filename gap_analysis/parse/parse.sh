#!/usr/bin/env bash

PA_SHAPE_COMBINED=/tmp/PADUS1_4Shapefile/PADUS1_4Combined.shp

STATE_PATH="/tmp/gap_analysis/$1"
AREAS_PATH="$STATE_PATH/areas"

mkdir -p "$STATE_PATH"
mkdir -p "$AREAS_PATH"

ogr2ogr -f 'ESRI Shapefile' \
    -t_srs 'epsg:4326' \
    -where "d_State_Nm IN ('$1')" \
    "$STATE_PATH/state.shp" \
    "$PA_SHAPE_COMBINED"

shp2json "$STATE_PATH/state.shp" -o "$STATE_PATH/state.json"

go run ./split.go --in "$STATE_PATH/state.json" --out "$AREAS_PATH/"