package main

import (
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/utils"
	"fmt"
)

type CanonicalNameSources []CanonicalNameSource


func (Ω CanonicalNameSources) Names() []string {

	res := []string{}

	for _, ns := range Ω {
		res = utils.AddStringToSet(res, ns.Synonyms...)
		res = utils.AddStringToSet(res, ns.SynonymFor)
		res = utils.AddStringToSet(res, ns.CanonicalName)
	}

	return res
}

type SourceOccurrenceCount map[store.DataSourceID]map[store.DataSourceTargetID]int

func (Ω SourceOccurrenceCount) Combine(given SourceOccurrenceCount) SourceOccurrenceCount {

	if given == nil {
		return Ω
	}
	if Ω == nil {
		return given
	}

	for sourceTypeID, occurrences := range given {
		if _, ok := Ω[sourceTypeID]; !ok {
			Ω[sourceTypeID] = map[store.DataSourceTargetID]int{}
		}
		for targetID, count := range occurrences {
			if _, ok := Ω[sourceTypeID][targetID]; ok {
				Ω[sourceTypeID][targetID] += count
			} else {
				Ω[sourceTypeID][targetID] = count
			}
		}
	}
	return Ω
}

type CanonicalNameSource struct{
	CanonicalName     string                `json:",omitempty"`
	SynonymFor        string                `json:",omitempty"`
	Synonyms          []string              `json:",omitempty"`
	Ranks             []string              `json:",omitempty"`
	SourceOccurrences SourceOccurrenceCount `json:",omitempty"`
}

func (Ω CanonicalNameSource) Combine(src CanonicalNameSource) CanonicalNameSource {
	res := Ω
	res.Ranks = utils.AddStringToSet(Ω.Ranks, src.Ranks...)
	if len(Ω.SynonymFor) > 0 || len(src.SynonymFor) > 0 {
		panic("combine should already have synonym fors removed")
	}
	res.Synonyms = utils.AddStringToSet(Ω.Synonyms, src.Synonyms...)
	res.SourceOccurrences = Ω.SourceOccurrences.Combine(src.SourceOccurrences)
	return res
}

func (Ω CanonicalNameSource) PushSynonym(src CanonicalNameSource) CanonicalNameSource {
	res := Ω
	res.Ranks = utils.AddStringToSet(Ω.Ranks, src.Ranks...)
	//res.SynonymFor = utils.AddStringToSet(Ω.SynonymFor, src.SynonymFor...)
	res.Synonyms = utils.AddStringToSet(Ω.Synonyms, src.SynonymFor)
	res.Synonyms = utils.AddStringToSet(Ω.Synonyms, src.Synonyms...)
	res.Synonyms = utils.AddStringToSet(Ω.Synonyms, src.CanonicalName)
	res.SourceOccurrences = Ω.SourceOccurrences.Combine(src.SourceOccurrences)
	return res
}

func (Ω CanonicalNameSources) GenerateNameResults() CanonicalNameSources {


	// We assume that by joining all names, every name that is a synonym of something else will be marked as such.
	srcs := Ω.DemoteSynonymFor().CondenseByName().AbsorbNamesThatAreSynonyms().CondenseBySynonyms()


	for i := range srcs {
		for srcID, occurrences := range srcs[i].SourceOccurrences {
			for targetID, count := range occurrences {
				if count == 0 {
					delete(srcs[i].SourceOccurrences[srcID], targetID)
				}
			}
		}
	}

	fmt.Println(utils.JsonOrSpew(srcs))

	return nil

	//return initialNameCondensed.CondenseBySynonym()

	//res := CanonicalNameSources{}

	//// Keep only those without parent synonyms.
	//for _, c := range initialNameCondensed {
	//	if len(c.SynonymFor) == 0 {
	//		res = append(res, c)
	//	}
	//}

	// Gather Synonyms
	//for i := range res {
	//	initialNameCondensed.CollectSynonyms(res[i].CanonicalName)
	//	res[i] = res[i].PushSynonym()
	//}
	//
	//return res
}

//func (Ω CanonicalNameSources) CollectSynonyms(name string) CanonicalNameSource {
//
//	res := CanonicalNameSource{}
//
//	for _, ns := range Ω {
//		if utils.ContainsString(ns.SynonymFor, name) {
//			res = res.PushSynonym(ns)
//		}
//	}
//
//	// Another run to ensure the synonyms don't have synonyms.
//	for _, synonym := range res.Synonyms {
//		res = res.PushSynonym(Ω.CollectSynonyms(synonym))
//	}
//
//	return res
//}

func (Ω CanonicalNameSources) DemoteSynonymFor() CanonicalNameSources {
	res := CanonicalNameSources{}
	for _, a := range Ω {
		if len(a.SynonymFor) == 0 {
			res = append(res, a)
		} else {
			c := CanonicalNameSource{CanonicalName: a.SynonymFor}.PushSynonym(a)
			res = append(res, c)
		}
	}
	return res
}

func (Ω CanonicalNameSources) CondenseByName() CanonicalNameSources {

	res := CanonicalNameSources{}
	// First combine all results
	for i, a := range Ω {
		exists := false
		for k, b := range res {

			sameName := a.CanonicalName == b.CanonicalName

			if sameName {
				exists = true
				res[k] = a.Combine(b)
			}
		}
		if !exists {
			res = append(res, Ω[i])
		}
	}
	return res
}

func (Ω CanonicalNameSources) AbsorbNamesThatAreSynonyms() CanonicalNameSources {

	res := CanonicalNameSources{}
	// First combine all results
	for i, a := range Ω {
		exists := false
		for k, b := range res {

			if utils.ContainsString(a.Synonyms, b.CanonicalName) {
				exists = true
				res[k] = a.PushSynonym(b)
			} else if utils.ContainsString(b.Synonyms, a.CanonicalName) {
				exists = true
				res[k] = b.PushSynonym(a)
			}
		}
		if !exists {
			res = append(res, Ω[i])
		}
	}
	return res
}




func (Ω CanonicalNameSources) CondenseBySynonyms() CanonicalNameSources {

	res := CanonicalNameSources{}
	// First combine all results
	for i, a := range Ω {
		exists := false
		for k, b := range res {

			if utils.IntersectsStrings(a.Synonyms, b.Synonyms) {
				exists = true
				// Keep the name with the largest number of synonyms
				if len(a.Synonyms) > len(b.Synonyms) {
					res[k] = a.Combine(b)
				} else {
					res[k] = b.Combine(a)
				}

			}
		}
		if !exists {
			res = append(res, Ω[i])
			}
	}
	return res
}

//func (Ω CanonicalNameSources) CondenseBySynonymFor() CanonicalNameSources {
//
//	res := CanonicalNameSources{}
//	// First combine all results
//	for i := range Ω {
//		exists := false
//		for k := range res {
//			if utils.ContainsString(Ω[i].SynonymFor, res[k].CanonicalName) {
//				exists = true
//				res[k] = res[k].PushSynonym(Ω[i])
//			} else if utils.ContainsString(res[k].SynonymFor, Ω[i].CanonicalName) {
//				exists = true
//				res[k] = Ω[i].PushSynonym(res[k])
//			} else if utils.IntersectsStrings(res[k].Synonyms, Ω[i].Synonyms) {
//				exists = true
//				res[k] = Ω[i].Combine(res[k])
//			}
//		}
//		if !exists {
//			res = append(res, Ω[i])
//		}
//	}
//	return res
//}



//func (Ω CanonicalNameSources) CondenseBySynonyms() CanonicalNameSources {
//
//	res := CanonicalNameSources{}
//	// First combine all results
//	for i := range Ω {
//		exists := false
//		for k := range res {
//			if utils.ContainsString(Ω[i].SynonymFor, res[k].CanonicalName) {
//				exists = true
//				res[k] = res[k].PushSynonym(Ω[i])
//			} else if utils.ContainsString(res[k].SynonymFor, Ω[i].CanonicalName) {
//				exists = true
//				res[k] = Ω[i].PushSynonym(res[k])
//			} else if utils.IntersectsStrings(res[k].Synonyms, Ω[i].Synonyms) {
//				res[k] = Ω[i].Combine(res[k])
//			}
//		}
//		if !exists {
//			res = append(res, Ω[i])
//		}
//	}
//	return res
//}