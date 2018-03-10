#!/usr/bin/env bash

# PAD-US is an inventory of marine and terrestrial protected areas.

# Download shape file:
# "the PAD-US database strives to be a complete “best available" inventory of areas dedicated to the
# preservation of biological diversity, and other natural, recreational or cultural uses, managed for
# these purposes through legal or other effective means."
#
# The data is updated annually.

# https://gapanalysis.usgs.gov/padus/data/metadata/
# https://gapanalysis.usgs.gov/padus/data/download/
# https://gapanalysis.usgs.gov/padus/data/statistics/

# Note that this link may not be the most recent. Found the link by using the standard download page and examining the chrome download history for a file link.

wget --directory-prefix /tmp --output-document PADUS1_4Shapefile.zip 'https://www.sciencebase.gov/catalog/file/get/56bba648e4b08d617f657960?f=__disk__bb%2F1d%2F64%2Fbb1d64d7adb8aeb0b6f75d347077ddc1406d8a53'
unzip -a /tmp/PADUS1_4Shapefile.zip -d /tmp/PADUS1_4Shapefile
rm -rf /tmp/PADUS1_4Shapefile.zip