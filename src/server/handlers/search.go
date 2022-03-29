package handlers

import(
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"log"
	"context"
	"bytes"
)

func Search(c *gin.Context) {
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