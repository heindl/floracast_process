#!/usr/bin/env bash

ogr2ogr -f 'ESRI Shapefile' \
    -where "d_State_Nm IN ('West Virginia', 'Virginia', 'Tennessee', 'South Carolina', 'Pennsylvania', 'Delaware', 'Georgia', 'Kentucky', 'Maryland', 'New Jersey', 'New York', 'North Carolina', 'Ohio')" \
    -t_srs 'epsg:4326' \
    ./shapefiles/parsed.shp \
    /Users/m/Downloads/PADUS1_4Shapefile/PADUS1_4Combined.shp


#ogr2ogr -f GeoJSON \
#    -where "d_State_Nm IN ('West Virginia')" \
#    -t_srs 'epsg:4326' \
#    ./shapefiles/wv.geojson \
#    /Users/m/Downloads/PADUS1_4Shapefile/PADUS1_4Combined.shp