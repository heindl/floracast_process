#!/bin/bash

STATES=(
#AK
AL
AZ
AR
CA
CO
CT
DE
FL
GA
#HI
ID
IL
IN
IA
KS
KY
LA
ME
MD
MA
MI
MN
MS
MO
MT
NE
NV
NH
NJ
NM
NY
NC
ND
OH
OK
OR
PA
RI
SC
SD
TN
TX
UT
VT
VA
WA
WV
WI
WY
)

#STATES=(
#    GA
#    WV
#    SC
#    NC
#    VA
#    KY
#    TN
#    AL
#)

#go build

printf '%s\n' "${STATES[@]}" | parallel -j4 "./saveState.sh {.}"

echo "Parsing State Files"

#printf '%s\n' "${STATES[@]}" | parallel -j4 "./parse --in /tmp/gap_analysis/{.}/state.geojson --out /tmp/gap_analysis/{.}/areas"




