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

	d.AddLink("://www.google.com")
	d.AddLink("http://www.youtube.com")
	d.AddLink("https://www.washingtonpost.com/")

	cv, _ := d.Convert(exampleFunc)
	var testStr string

	for i := 0; i < len(d); i++ {
		c, contains := cv[string(i)]

		if contains {
			s := *c.data
			testStr += s[0]
		}
	}

	if testStr != "abc" {
		t.Errorf("Expected 'abc' but got '%v' ", testStr)
	}
}

func Test_failure(t *testing.T) {
	d := NewDLink()

	d.AddLink("://www.nkvxxnvhh.com")
	d.AddLink("http://youtouneaube.com")
	d.AddLink("https://www.waffawshinnjbgtonpost.com/")

	cv, _ := d.Convert(exampleFunc)

	for i := 0; i < len(d); i++ {
		_, contains := cv[string(i)]
		t.Log("Expected ALL fail when request")
		if contains {
			t.Errorf("It should be no any index key in it when every request fail")
		}
	}
}

func Test_emptyConvert(t *testing.T) {
	d := NewDLink()

	_, err := d.Convert(exampleFunc)

	if err == nil {
		t.Errorf("Should have ERROR since no link inside when convert")
	}
}
