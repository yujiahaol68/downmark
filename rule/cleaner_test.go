package rule

import (
	"fmt"
	"net/http"
	"testing"

	"golang.org/x/net/html"
)

func Test_convertor(t *testing.T) {
	url := "http://www.chinadaily.com.cn/a/201712/11/WS5a2e2bbea310eefe3e9a2918.html"
	resp, _ := http.Get(url)

	tr := html.NewTokenizer(resp.Body)

	_, b := CleanConvertor(tr, "div", "id", "Content")
	fmt.Printf("%s", b.String())

	res, err := http.Get("http://www.bbc.com/news/uk-42318755")
	if err != nil {
		return
	}

	trr := html.NewTokenizer(res.Body)
	_, bf := CleanConvertor(trr, "div", "story-body__inner")
	fmt.Printf("%s", bf.String())
}
