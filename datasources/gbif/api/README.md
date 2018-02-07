gbif
=====

Go package to draw data from the Global Biodiversity Information Facility, http://www.gbif.org/.

```
package main

import(
    "github.com/heindl/gbif"
    "fmt"
    "time"
)

func main() {

    results, err := Occurrences(OccurrenceSearchQuery{
        TaxonKey:           7205815,
        HasCoordinate:      true,
        HasGeospatialIssue: false,
    })
    
     if err != nil {
        fmt.Errorf("ERROR: %s", err)
        return
     }
    
    fmt.Printf("The GBIF returned %d occurrences of the American white admiral butterfly", len(results))
    
}


```
