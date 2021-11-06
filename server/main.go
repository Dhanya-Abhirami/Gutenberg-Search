package main

import (
    "github.com/elastic/go-elasticsearch/v7"
	// "os"
    "time"
    "log"
)

type Document struct {
    Title string `json:"title"`
    Author string `json:"author"`
    Id int `json:"id"`
    Text string `json:"text"`
}

func resetIndex(es){
    es
    putMapping()
}

func putMapping(es){

}

func checkConnection(){
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

}