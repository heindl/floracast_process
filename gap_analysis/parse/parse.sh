#!/bin/bash

# Alabama, Arizona, Arkansas, California, Colorado, Connecticut, Delaware, Florida, Georgia, Idaho, Illinois, Indiana, Iowa, Kansas, Kentucky, Louisiana, Maine, Maryland, Massachusetts, Michigan, Minnesota, Mississippi, Missouri, Montana, Nebraska, Nevada, New Hampshire, New Jersey, New Mexico, New York, North Carolina, North Dakota, Ohio, Oklahoma, Oregon, Pennsylvania, Rhode Island, South Carolina, South Dakota, Tennessee, Texas, Utah, Vermont, Virginia, Washington, West Virginia, Wisconsin, Wyoming, District of Columbia


PA_SHAPE_COMBINED=/tmp/PADUS1_4Shapefile/PADUS1_4Combined.shp

STATES=( Georgia Alabama )

for STATE in "${STATES[@]}"
do
    STATE_PATH="/tmp/gap_analysis/$STATE"
    AREAS_PATH="$STATE_PATH/areas"

    mkdir -p "$STATE_PATH"
    mkdir -p "$AREAS_PATH"

    ogr2ogr -f 'ESRI Shapefile' \
        -t_srs 'epsg:4326' \
        -where "d_State_Nm IN ('$STATE')" \
        "$STATE_PATH/state.shp" \
        "$PA_SHAPE_COMBINED"

    shp2json "$STATE_PATH/state.shp" -o "$STATE_PATH/state.json"

    go run ./split.go --in "$STATE_PATH/state.json" --out "$AREAS_PATH/"

done



