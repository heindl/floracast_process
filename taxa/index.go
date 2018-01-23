package main

import (
	"bitbucket.org/heindl/taxa/store"
	"bitbucket.org/heindl/taxa/utils"
)

type CanonicalNameSources []CanonicalNameSource


func (Ω CanonicalNameSources) Names() []string {

	res := []string{}

	for _, ns := range Ω {
		res = utils.AddStringToSet(res, ns.Synonyms...)
		res = utils.AddStringToSet(res, ns.SynonymFor...)
		res = utils.AddStringToSet(res, ns.CanonicalName)
	}

	return res
}

type SourceMap map[store.DataSourceID]store.DataSourceTargetIDs

func (Ω SourceMap) Combine(given SourceMap) SourceMap {

	if given == nil {
		return Ω
	}
	if Ω == nil {
		return given
	}

	for sourceTypeID, targetIDs := range given {
		if _, ok := Ω[sourceTypeID]; !ok {
			Ω[sourceTypeID] = store.DataSourceTargetIDs{}
		}
		Ω[sourceTypeID] = Ω[sourceTypeID].AddToSet(targetIDs...)
	}
	return Ω
}


type CanonicalNameSource struct{
	CanonicalName string `json:",omitempty"`
	SynonymFor []string `json:",omitempty"`
	Synonyms []string `json:",omitempty"`
	Ranks []string `json:",omitempty"`
	SourceMap SourceMap `json:",omitempty"`
}

func (Ω CanonicalNameSource) Combine(src CanonicalNameSource) CanonicalNameSource {
	res := Ω
	res.Ranks = utils.AddStringToSet(Ω.Ranks, src.Ranks...)
	res.SynonymFor = utils.AddStringToSet(Ω.SynonymFor, src.SynonymFor...)
	res.Synonyms = utils.AddStringToSet(Ω.Synonyms, src.Synonyms...)
	res.SourceMap = Ω.SourceMap.Combine(src.SourceMap)
	return res
}

func (Ω CanonicalNameSource) PushSynonym(src CanonicalNameSource) CanonicalNameSource {
	res := Ω
	res.Ranks = utils.AddStringToSet(Ω.Ranks, src.Ranks...)
	//res.SynonymFor = utils.AddStringToSet(Ω.SynonymFor, src.SynonymFor...)
	res.Synonyms = utils.AddStringToSet(Ω.Synonyms, src.SynonymFor...)
	res.Synonyms = utils.AddStringToSet(Ω.Synonyms, src.Synonyms...)
	res.Synonyms = utils.AddStringToSet(Ω.Synonyms, src.CanonicalName)
	res.SourceMap = Ω.SourceMap.Combine(src.SourceMap)
	return res
}

func (Ω CanonicalNameSources) GenerateNameResults() CanonicalNameSources {


	// We assume that by joining all names, every name that is a synonym of something else will be marked as such.
	initialNameCondensed := Ω.CondenseByCanonicalName()

	return initialNameCondensed.CondenseBySynonym()

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

func (Ω CanonicalNameSources) CollectSynonyms(name string) CanonicalNameSource {

	res := CanonicalNameSource{}

	for _, ns := range Ω {
		if utils.Contains(ns.SynonymFor, name) {
			res = res.PushSynonym(ns)
		}
	}

	// Another run to ensure the synonyms don't have synonyms.
	for _, synonym := range res.Synonyms {
		res = res.PushSynonym(Ω.CollectSynonyms(synonym))
	}

	return res
}


func (Ω CanonicalNameSources) CondenseByCanonicalName() CanonicalNameSources {

	res := CanonicalNameSources{}
	// First combine all results
	for i := range Ω {
		exists := false
		for k := range res {
			if res[k].CanonicalName == Ω[i].CanonicalName {
				exists = true
				res[k] = res[k].Combine(Ω[i])
				// May be multiple matches so not exit
			}
		}
		if !exists {
			res = append(res, Ω[i])
			}
	}
	return res
}

func (Ω CanonicalNameSources) CondenseBySynonym() CanonicalNameSources {

	res := CanonicalNameSources{}
	// First combine all results
	for i := range Ω {
		exists := false
		for k := range res {
			if utils.Contains(Ω[i].SynonymFor, res[k].CanonicalName) {
				exists = true
				res[k] = res[k].PushSynonym(Ω[i])
			}
		}
		if !exists {
			res = append(res, Ω[i])
		}
	}
	return res
}