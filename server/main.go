package main

import(
	// "fmt"
	"context"
	"log"
	"net/http"
	"encoding/json"
	// "io/ioutil"
	// "strings"
	"bytes"
	"github.com/gin-gonic/gin"
)

const port = ":8080"

func search(c *gin.Context) {
	term := c.Query("term")
	offset := c.DefaultQuery("offset","0") 
	var buf bytes.Buffer
	var r  map[string]interface{}
	query := map[string]interface{}{
		"from" : offset,
	  	"query": map[string]interface{}{
			"match": map[string]interface{}{
				"text": map[string]interface{}{
					"query": term,
					"operator": "and",
					"fuzziness": "auto",
				},
			},
		},
		"highlight": map[string]interface{}{ 
			"fields": map[string]interface{}{ 
				"text": map[string]interface{}{}, 
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
	  log.Fatalf("Error encoding query: %s", err)
	}
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(indexName),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	  )
	  if err != nil {
		log.Fatalf("Error getting response: %s", err)
	  }
	  defer res.Body.Close()
	  if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	  }
	//   log.Printf(
	// 	"[%s] %d hits; took: %dms",
	// 	res.Status(),
	// 	int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
	// 	int(r["took"].(float64)),
	//   )
	  // Print the ID and document source for each hit.
	//   for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
	// 	log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	//   }
	// c.JSON(http.StatusOK, gin.H{"data": r})
	c.JSON(http.StatusOK, r)    
  }

  func paragraphs(c *gin.Context) {
	title := c.Query("title")
	start := c.Query("start")
	end := c.Query("end") 
	const filter = [
		map[string]interface{}{ "term": map[string]interface{}{ "title": title } },
		map[string]interface{}{ "range": map[string]interface{}{ "id": map[string]interface{}{ "gte": start, "lte": end } } }
	]

	const body = {
		size: end - start,
		sort: { location: 'asc' },
		query: { bool: { filter } }
	}
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(indexName),
		es.Search.WithBody(&strings.NewReader(body)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	  )
	  if err != nil {
		log.Fatalf("Error getting response: %s", err)
	  }
	  defer res.Body.Close()
	  if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	  }
	log.Println(title,start,end)
	c.JSON(http.StatusOK, gin.H{"start": start,"end":end})
  }

// CORS Middleware
func CORS(c *gin.Context) {

	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	// Second, we handle the OPTIONS problem
	if c.Request.Method != "OPTIONS" {
		
		c.Next()

	} else {
        
		// Everytime we receive an OPTIONS request, 
		// we just return an HTTP 200 Status Code
		// Like this, Angular can now do the real 
		// request using any other method than OPTIONS
		c.AbortWithStatus(http.StatusOK)
	}
}

func main(){
	// prepareIndex()
	es = createConnection()
	router := gin.Default()
	router.Use(CORS) 
	router.GET("/search", search)
	router.GET("/paragraphs", paragraphs)
	router.Run(port)
}

