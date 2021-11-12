package main

import (
    "github.com/elastic/go-elasticsearch/v7"
    "github.com/elastic/go-elasticsearch/v7/esutil"
    // "github.com/elastic/go-elasticsearch/v7/esapi"
	"os"
    "context"
    "log"
    "time"
    "sync/atomic"
    "regexp"
    "strings"
    "strconv"
    "bytes"
    "encoding/json"
)

const indexName = "gutenberg"
const bookPath = "./books/"

type Document struct {
    Title string `json:"title"`
    Author string `json:"author"`
    Paragraphs []string `json:"paragraphs"`
    // Id int `json:"id"`
    // Text string `json:"text"`
}

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

func resetIndex(es *elasticsearch.Client){
    res, _ := es.Indices.Exists([]string{indexName})
    if(res.StatusCode == 200){
        res, _ = es.Indices.Delete([]string{indexName})
    }
}

func parseBookFile (filePath string) (string, string, []string) {
    b, err := os.ReadFile(filePath)
    if err!=nil {
        log.Fatal(err)
    }
	book := string(b)
	title := regexp.MustCompile(`/^Title:\s(.+)$/m`).FindString(book)
	authorMatch := regexp.MustCompile(`/^Author:\s(.+)$/m`).FindString(book)
	var author string
    var paragraphs []string
    if len(authorMatch)==0 || strings.TrimSpace(authorMatch) == "" {
		author = "Unknown Author"
	} else {
		author = authorMatch
	}
	
	log.Printf("Reading Book - %s By %s",title,author)
  
    // const startOfBookMatch = book.match(/^\*{3}\s*START OF (THIS|THE) PROJECT GUTENBERG EBOOK.+\*{3}$/m)
    // const startOfBookIndex = startOfBookMatch.index + startOfBookMatch[0].length
    // const endOfBookIndex = book.match(/^\*{3}\s*END OF (THIS|THE) PROJECT GUTENBERG EBOOK.+\*{3}$/m).index
  
    // const paragraphs = book
    //   .slice(startOfBookIndex, endOfBookIndex) // Remove Guttenberg header and footer
    //   .split(/\n\s+\n/g) // Split each paragraph into it's own array entry
    //   .map(line => line.replace(/\r\n/g, ' ').trim()) // Remove paragraph line breaks and whitespace
    //   .map(line => line.replace(/_/g, '')) // Guttenberg uses "_" to signify italics.  We'll remove it, since it makes the raw text look messy.
    //   .filter((line) => (line && line !== '')) // Remove empty lines
  
    // log.Println(`Parsed ${paragraphs.length} Paragraphs\n`)
    return title, author, paragraphs
  }

func loadBooks(es *elasticsearch.Client){
    resetIndex(es)
    bi, _ := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
        Index:  indexName,        
        Client: es,               
    })

    f, err := os.Open(bookPath)
    if err!=nil{
        log.Fatal(err)
        return
    }
    files, err := f.Readdir(0)
    f.Close()
    if err!=nil{
        log.Fatal(err)
    }
    var countSuccessful uint64
    for i, file := range files {
        match, _ := regexp.MatchString("[0-9A-Za-z-]+(\\.txt)", file.Name())
        if(match){   
            title, author, paragraphs := parseBookFile(bookPath+file.Name())
            d := &Document{
                Title: title,
                Author: author,
                Paragraphs: paragraphs,
            }
            data, err := json.Marshal(d)
            if err != nil {
                log.Fatalf("Cannot encode article %d: %s", d.Title, err)
            }
            err = bi.Add(context.Background(), esutil.BulkIndexerItem{
                        Action:     "index",
                        DocumentID: strconv.Itoa(i),
                        Body: bytes.NewReader(data),
                        OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
                            atomic.AddUint64(&countSuccessful, 1)
                          },
                        OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
                        if err != nil {
                            log.Printf("ERROR: %s", err)
                        } else {
                            log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
                        }
                        },
                    })
            if err != nil {
                log.Fatalf("Unexpected error: %s", err)
            }
        }
    }

    if err := bi.Close(context.Background()); err != nil {
        log.Fatalf("Unexpected error: %s", err)
    }
    stats := bi.Stats()
    log.Println(stats)
}

func main() {
    
	es:= createConnection()
    loadBooks(es)
}