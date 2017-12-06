package cleaner

import (
	"fmt"
	"net/http"
	"testing"

	"golang.org/x/net/html"
)

func Test_convertor(t *testing.T) {
	url := "https://www.washingtonpost.com/world/russia-enraged-by-2018-winter-olympics-ban-over-doping/2017/12/06/5a3f2674-da08-11e7-a241-0848315642d0_story.html?hpid=hp_hp-more-top-stories_russiaoly-703am%3Ahomepage%2Fstory&utm_term=.6fe17c4a3f9b"
	resp, _ := http.Get(url)

	tr := html.NewTokenizer(resp.Body)

	b := CleanConvertor(tr)
	fmt.Printf("%s", b.String())
}
