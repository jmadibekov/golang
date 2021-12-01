// before trying `go run` this file,
// make sure that your elasticsearch is up and running

// you can do that as described here: https://developer.okta.com/blog/2021/04/23/elasticsearch-go-developers-guide
// (that is a great tutorial, by the way)
package main

import (
	"log"

	"github.com/elastic/go-elasticsearch/v7"
)

func main() {
	es, _ := elasticsearch.NewDefaultClient()
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
}
