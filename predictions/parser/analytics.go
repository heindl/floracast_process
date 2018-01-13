package parser

import (
//"fmt"
//"google.golang.org/api/iterator"
//"gopkg.in/tomb.v2"
//"cloud.google.com/go/storage"
//"context"
)

//func (Ω *predictionParser) GatherAnalytics(cxt context.Context, predictionsDirectory string) error {
//	q := &storage.Query{Prefix: fmt.Sprintf("predictions/%s/", predictionsDirectory), Delimiter: "/"}
//
//	tmb := tomb.Tomb{}
//	tmb.Go(func() error {
//		iter := Ω.Bucket.Objects(cxt, q)
//		for {
//			o, err := iter.Next()
//			if err != nil && err == iterator.Done {
//				break
//			}
//			if err != nil {
//				panic(err)
//			}
//			name := o.Name
//			tmb.Go(func() error {
//				predictions, err := Ω.parseBucketObject(cxt, name)
//				if err != nil {
//					return err
//				}
//				Ω.Lock()
//				defer Ω.Unlock()
//				for _, p := range predictions {
//					if _, ok := Ω.PredictionsOverTaxon[p.TaxonID]; !ok {
//						Ω.PredictionsOverTaxon[p.TaxonID] = stats.Float64Data{}
//					}
//					Ω.PredictionsOverTaxon[p.TaxonID] = append(Ω.PredictionsOverTaxon[p.TaxonID], p.PredictionValue)
//					if _, ok := Ω.PredictionsOverDay[p.FormattedDate]; !ok {
//						Ω.PredictionsOverDay[p.FormattedDate] = stats.Float64Data{}
//					}
//					Ω.PredictionsOverDay[p.FormattedDate] = append(Ω.PredictionsOverDay[p.FormattedDate], p.PredictionValue)
//				}
//				return nil
//			})
//		}
//		return nil
//	})
//	return tmb.Wait()
//}
