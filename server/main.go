package main

import (
    "github.com/elastic/go-elasticsearch/v7"
    "github.com/elastic/go-elasticsearch/v7/esutil"
    // "github.com/elastic/go-elasticsearch/v7/esapi"
	// "os"
    "fmt"
    "time"
    "log"
)

type Document struct {
    Title string `json:"title"`
    Author string `json:"author"`
    Id int `json:"id"`
    Text string `json:"text"`
}

// func resetIndex(es){
//     putMapping()
// }

func putMapping(es *elasticsearch.Client){
    doc := Document{Title: "Test",Author: "Dhanya",Id:1,Text:"Sample"}
    res, _ := es.Index("test", esutil.NewJSONReader(&doc))
    fmt.Println(res)

    // bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
    //     Index:         indexName,        // The default index name
    //     Client:        es,               // The Elasticsearch client
    //     NumWorkers:    numWorkers,       // The number of worker goroutines
    //     FlushBytes:    int(flushBytes),  // The flush threshold in bytes
    //     FlushInterval: 30 * time.Second, // The periodic flush interval
    //   })
}

func checkConnection() *elasticsearch.Client{
    log.Println(elasticsearch.Version)
    time.Sleep(time.Second * 20)
    cfg := elasticsearch.Config{
		Addresses: []string{
			"http://elasticsearch:9200",
		},
	  }
    es, err := elasticsearch.NewClient(cfg)
    if err != nil {
   	 log.Fatalf("Error creating the client: %s", err)
    }
    res, err := es.Info()
    log.Println(res)
    if err != nil {
   	 log.Fatalf("Error getting response: %s", err)
    }
    defer res.Body.Close()
    return es
}
func main() {
    
	es:= checkConnection()
    putMapping(es)
}