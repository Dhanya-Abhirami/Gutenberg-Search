package main

import (
    "github.com/elastic/go-elasticsearch/v7"
    "github.com/elastic/go-elasticsearch/v7/esutil"
    // "github.com/elastic/go-elasticsearch/v7/esapi"
	"os"
    "context"
    "log"
    "sync/atomic"
    "regexp"
    "strings"
    "strconv"
    "bytes"
    "encoding/json"
)

const indexName = "gutenberg"
const bookPath = "./books/"
var es *elasticsearch.Client

type Document struct {
    Title string `json:"title"`
    Author string `json:"author"`
    Paragraphs []string `json:"paragraphs"`
    // Id int `json:"id"`
    // Text string `json:"text"`
}

func resetIndex(){
    res, _ := es.Indices.Exists([]string{indexName})
    if(res.StatusCode == 200){
        log.Println("Resetting")
        res, _ = es.Indices.Delete([]string{indexName})
    }
}

func parseBookFile (filePath string) (string, string, []string) {
    b, err := os.ReadFile(filePath)
    if err!=nil {
        log.Fatal(err)
    }
	book := string(b)
	title := regexp.MustCompile(`(?m)(^Title:\s)(.+)(\s)$`).FindStringSubmatch(book)[2]
	authorMatch := regexp.MustCompile(`(?m)(^Author:\s)(.+)(\s)$`).FindStringSubmatch(book)
	var author string
    var paragraphs []string
    if len(authorMatch)==0 || strings.TrimSpace(authorMatch[2]) == "" {
		author = "Unknown Author"
	} else {
		author = authorMatch[2]
	}
	
	log.Printf("Reading Book - %s By - %s",title,author)
    startOfBookMatch := regexp.MustCompile(`(?m)(^\*{3}\s*START OF (THIS|THE) PROJECT GUTENBERG EBOOK.+\*{3})$`)
    endOfBookMatch := regexp.MustCompile(`(?m)(^\*{3}\s*END OF (THIS|THE) PROJECT GUTENBERG EBOOK.+\*{3})$`)
    log.Println(len(startOfBookMatch.FindStringIndex(book)),len(endOfBookMatch.FindStringIndex(book)))
    if len(startOfBookMatch.FindStringIndex(book)) > 0 && len(endOfBookMatch.FindStringIndex(book)) > 0 {
        startOfBookIndex := startOfBookMatch.FindStringIndex(book)[0] + len(startOfBookMatch.FindString(book))
        endOfBookIndex := endOfBookMatch.FindStringIndex(book)[0]
        
        // paragraphs = 
        //   .slice() // Remove Guttenberg header and footer
        //   .split(/\n\s+\n/g) // Split each paragraph into it's own array entry
        //   .map(line => line.replace(/\r\n/g, ' ').trim()) // Remove paragraph line breaks and whitespace
        //   .map(line => line.replace(/_/g, '')) // Guttenberg uses "_" to signify italics.  We'll remove it, since it makes the raw text look messy.
        //   .filter((line) => (line && line !== '')) // Remove empty lines
        book = book[startOfBookIndex:endOfBookIndex]
        paragraphBreakmatch := regexp.MustCompile(`\n\s+\n`)
        raw_paragraphs := paragraphBreakmatch.Split(book, -1) 
        r := strings.NewReplacer("\r\n", " ", 
        "_", " ")
        for _, line := range raw_paragraphs {
            
            line = r.Replace(line)
            if line!=""{
                paragraphs = append(paragraphs,line)
            }
        }
    
    }

    log.Printf("Parsed %d Paragraphs\n",len(paragraphs))
    return title, author, paragraphs
  }

func loadBooks(){
    resetIndex()
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

func prepareIndex() {
	es = createConnection()
    loadBooks()
}