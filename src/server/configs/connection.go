package configs

import (
    "github.com/elastic/go-elasticsearch/v7"
    "log"
    "time"
)

func createConnection() *elasticsearch.Client{
    // log.Println(elasticsearch.Version)
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
    if err != nil {
   	 log.Fatalf("Error getting response: %s", err)
    }
    defer res.Body.Close()
    return es
}

