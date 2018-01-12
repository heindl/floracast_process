#!/bin/bash

# 'West Virginia', 'Virginia', 'Tennessee', 'South Carolina', 'Pennsylvania', 'Delaware', 'Georgia', 'Kentucky', 'Maryland', 'New Jersey', 'New York', 'North Carolina', 'Ohio'

STATES=( Georgia )

for STATE in "${STATES[@]}"
do
    STATE_PATH="./shapefiles/$STATE"
    AREAS_PATH="$STATE_PATH/areas"
    SVG_PATH="$STATE_PATH/svgs"

#    mkdir "$STATE_PATH"
#    mkdir "$AREAS_PATH"
#    mkdir "$SVG_PATH"

#    ogr2ogr -f 'ESRI Shapefile' \
#        -t_srs 'epsg:4326' \
#        -where "d_State_Nm IN ('$STATE')" \
#        "$STATE_PATH/state.shp" \
#        /Users/m/Downloads/PADUS1_4Shapefile/PADUS1_4Combined.shp
#
#    shp2json "$STATE_PATH/state.shp" -o "$STATE_PATH/state.json"
#
#    go run split_feature_collections --in "$STATE_PATH/state.json" --out "$AREAS_PATH/"

    for f in "$AREAS_PATH"/*.json
    do
        b=$(basename $f)
        out="$SVG_PATH/${b%.*}.svg"

        ndjson-map -r d3 'd.properties.STERADIANS = d3.geoArea(d).toFixed(10), d' < $f > $f

#        geoproject 'd3.geoAlbersUsa().fitSize([960, 960], d)' < $f | geo2svg -w 960 -h 960 --fill '#f7f7f7' --stroke '#eeeeee' > $out
    done
done



