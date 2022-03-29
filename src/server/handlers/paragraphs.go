package handlers

import(
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"log"
	"context"
	"strings"
)

func Paragraphs(c *gin.Context) {
	title := c.Query("title")
	start := c.Query("start")
	end := c.Query("end") 
	var r  map[string]interface{}
	// const filter = [
	// 	map[string]interface{}{ "term": map[string]interface{}{ "title": title } },
	// 	map[string]interface{}{ "range": map[string]interface{}{ "id": map[string]interface{}{ "gte": start, "lte": end } } }
	// ]

	// const body = {
	// 	size: end - start,
	// 	sort: { location: 'asc' },
	// 	query: { bool: { filter } }
	// }
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