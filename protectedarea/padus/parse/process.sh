#!/bin/bash

FULL_STATES=(
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

STATES=(OR)

printf '%s\n' "${STATES[@]}" | parallel -j4 "./getState.sh {.}"

print "Parsing State Files"

printf '%s\n' "${STATES[@]}" | go run ./parse.go --in "/tmp/gap_analysis/{.}/state.geojson" --out "/tmp/gap_analysis/{.}/areas"




