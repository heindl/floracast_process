package main

import (
	"bitbucket.org/heindl/cxt"
	"github.com/bitly/go-nsq"
	"bitbucket.org/heindl/species/fetcher"
	"bitbucket.org/heindl/species/metafetcher"
)

const DefaultChannel = "ch"

func main() {

	// Initialize new worker instance.
	c, err := cxt.Bootstrap()
	if err != nil {
		panic(err)
	}

	for t, h := range map[string]nsq.Handler{
		cxt.NSQSpeciesFetch:   &fetcher.SpeciesFetchHandler{c},
		cxt.NSQSpeciesMetaFetch: &metafetcher.SpeciesMetaFetchHandler{c},
	}{
		go func(topic string, handler nsq.Handler, context *cxt.Context) {
			q, err := nsq.NewConsumer(topic, DefaultChannel, context.NSQConfig)
			if err != nil {
				panic(err)
			}
			q.AddConcurrentHandlers(handler, 30)
			if err := q.ConnectToNSQD(c.NSQLookup); err != nil {
				panic(err)
			}
		}(t, h, c)
	}

	<-make(chan bool)

}
