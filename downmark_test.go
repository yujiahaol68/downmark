package downmark

import (
	"testing"

	"golang.org/x/net/html"
)

func Test_add_link(t *testing.T) {
	d := NewDLink()
	err := d.AddLink("www.google.com")
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(d) != 1 {
		t.Errorf("Link missing")
	}
	if d[0] != "http://www.google.com" {
		t.Errorf("Did not add http prefix")
	}

	err = d.AddLink("http//google.com")
	if err == nil {
		t.Errorf("should have error because 'http//google.com' is not a valid URL")
	}
}

func exampleFunc(cv *Converted, tr *html.Tokenizer, ch chan *Converted) {
	cv.title = "passage"
	s := []string{}

	switch cv.index {
	case 0:
		s = append(s, "a")
	case 1:
		s = append(s, "b")
	default:
		s = append(s, "c")
	}

	cv.data = &s
	ch <- cv
}

func Test_handler(t *testing.T) {
	d := NewDLink()

	d.AddLink("http://www.google.com")
	d.AddLink("http://www.youtube.com")
	d.AddLink("https://www.washingtonpost.com/")

	cv, _ := d.Convert(exampleFunc)
	var testStr string

	for i := 0; i < 3; i++ {
		s := *(cv[string(i)].data)
		testStr += s[0]
	}

	if testStr != "abc" {
		t.Errorf("Expected 'abc' but got '%v' ", testStr)
	}
}
