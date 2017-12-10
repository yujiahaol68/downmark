package cleaner

import (
	"fmt"
	"net/http"
	"testing"

	"golang.org/x/net/html"
)

func Test_convertor(t *testing.T) {
	url := "http://www.bbc.com/news/world-middle-east-42297437"
	resp, _ := http.Get(url)

	tr := html.NewTokenizer(resp.Body)

	DefineRules("story-body__inner")
	b := CleanConvertor(tr)
	fmt.Printf("%s", b.String())
}
