package downmark

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// DLink save all the links
type DLink []string

// DBody consist of index in the DLink
type DBody struct {
	Index      int
	ParsedBody *html.Tokenizer
}

// Conversion is a concurrent safe map with key of link index and value of pointer of the converted string slice
type Conversion map[string]*[]string

const (
	requestTimeOut = 5 * time.Second
)

var wg sync.WaitGroup

// NewDLink produce a list contains link
func NewDLink() DLink {
	d := DLink{}
	return d
}

func newDBody(i int, t *html.Tokenizer) DBody {
	return DBody{
		Index:      i,
		ParsedBody: t,
	}
}

func newConversion() Conversion {
	d := make(Conversion)
	return d
}

// AddLink append link to the list
func (d *DLink) AddLink(url string) error {
	url = completeProtocol(strings.TrimSpace(url))

	if isLink(url) {
		*d = append(*d, url)
		return nil
	}
	return fmt.Errorf("Expect URL but %v is not a url", url)
}

func isLink(l string) bool {
	_, err := url.ParseRequestURI(l)

	if err != nil {
		return false
	}
	return true
}

func completeProtocol(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	} else if strings.HasPrefix(url, "://") {
		return "https" + url
	} else {
		return "http://" + url
	}
}

// ConvertToMarkDown order go-routines to get the http body the concurrently and convert them into markdown format
func (d DLink) ConvertToMarkDown() (Conversion, error) {
	if len(d) == 0 {
		return nil, fmt.Errorf("Not URL can be used")
	}

	convertedBody := newConversion()
	tokenizerChan := make(chan *DBody)

	timeOut := time.Duration(requestTimeOut)
	client := http.Client{
		Timeout: timeOut,
	}

	for i, u := range d {
		go func(id int, url string) {
			resp, err := client.Get(url)

			if err != nil {
				fmt.Println(url + "FAIL in request")
				return
			}

			defer resp.Body.Close()

			dBody := newDBody(id, html.NewTokenizer(resp.Body))

			tokenizerChan <- &dBody
		}(i, u)
	}

	return convertedBody, nil
}
