package markdown

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var HTMLTemplate = `<html><head><title>downmark</title></head><body><div>%s</div></body></html>`

var tagP = "<p>aaaa<code>down</code>bbbb<em>cccc</em><strong>dddd</strong></p>"
var withLink = `<p>aaaa<code>down</code><a href="tecknight.xyz">bbbb</a><em>cccc</em>dddd</p>`
var codeLink = `<code><a href="http://tecknight.xyz">bbbb</a></code>`
var header = `<h1>aaaa</h1><br><h2>bbbb</h2><hr><h3>cccc</h3>`

var mTagp = "aaaa `down` bbbb*cccc***dddd**"
var mWithLink = "aaaa `down` [bbbb](http://tecknight.xyz)*cccc*dddd"
var mCodeLink = "[`bbbb`](http://tecknight.xyz)"
var mHeader = "# aaaa\n## bbbb-----### cccc"

func printAll(s *[]string) {
	ss := *s
	for _, l := range ss {
		fmt.Println(l)
	}
	fmt.Printf("\n")
}

func concatAll(s *[]string) string {
	var b bytes.Buffer
	ss := *s

	for _, l := range ss {
		b.WriteString(l)
	}
	return b.String()
}

func Test_tag_p(t *testing.T) {
	r := strings.NewReader(fmt.Sprintf(HTMLTemplate, tagP))
	tr := html.NewTokenizer(r)
	s := MdConvertor(tr)

	if s[0] != mTagp {
		t.Errorf("tag <p> TEST FAIL but get: %v", s[0])
	}
}

func Test_header(t *testing.T) {
	r := strings.NewReader(fmt.Sprintf(HTMLTemplate, header))
	tr := html.NewTokenizer(r)
	s := MdConvertor(tr)
	if data := concatAll(&s); data != mHeader {
		t.Errorf("Expected:\n%s \nBut got: \n%s", mHeader, data)
	}
}
