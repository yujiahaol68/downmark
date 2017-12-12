package rule

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var HTMLTemplate = `<html><head><title>downmark</title></head><body><div>%s</div></body></html>`

var tagP = "<p>aaaa<code>down</code>bbbb<em>cccc</em><strong>dddd</strong></p>"
var withLink = `<p>aaaa<code>down</code><a href="http://tecknight.xyz">bbbb</a><em>cccc</em>dddd</p>`
var codeLink = `<p><code><a href="http://tecknight.xyz"><em>bbbb</em></a></code></p>`
var linkCase = `<p><code><a href="https://golang.org/pkg/builtin/#recover" data-href="https://golang.org/pkg/builtin/#recover"><em class="markup--em markup--p-em">recover()</em></a></code><em class="markup--em markup--p-em"> returns the value provided to </em><code class="markup--code markup--p-code"><a href="https://golang.org/pkg/builtin/#panic" data-href="https://golang.org/pkg/builtin/#panic" class="markup--anchor markup--p-anchor" rel="noopener" target="_blank"><em class="markup--em markup--p-em">panic()</em></a></code><em class="markup--em markup--p-em"> which let’s you decide what you’d do with it. You can also pass an error or other types of values to panic, then you can check whether the panic was caused by the value you’re looking for. More </em><a href="https://blog.golang.org/defer-panic-and-recover" data-href="https://blog.golang.org/defer-panic-and-recover" class="markup--anchor markup--p-anchor" rel="noopener" target="_blank"><em class="markup--em markup--p-em">here</em></a><em class="markup--em markup--p-em">.</em></p>`
var header = `<h1>aaaa</h1><br><h2>bbbb</h2><hr><h3>cccc</h3>`
var img = `<img class="progressiveMedia-image" src="https://cdn-images-1.medium.com/max/1000/1*H0luK0YxgVlSkXqFsyhSnw.png">`
var pre = `<pre><strong>type</strong> Car <strong>struct</strong> {<br>  model <strong>string</strong><br>}</pre>`

var mTagp = "aaaa `down` bbbb*cccc*__dddd__"
var mWithLink = "aaaa `down` [bbbb](http://tecknight.xyz)*cccc*dddd"
var mCodeLink = "[ `bbbb` ](http://tecknight.xyz)"
var mHeader = "# aaaa\n## bbbb-----### cccc"
var mImg = "![](https://cdn-images-1.medium.com/max/1000/1*H0luK0YxgVlSkXqFsyhSnw.png)"

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

func injectTestCaseAndExec(tCase string) []string {
	r := strings.NewReader(fmt.Sprintf(HTMLTemplate, tCase))
	tr := html.NewTokenizer(r)
	_, s := MdConvertor(tr)
	return s
}

func Test_tag_p(t *testing.T) {
	s := injectTestCaseAndExec(tagP)

	if s[0] != mTagp {
		t.Errorf("tag <p> TEST FAIL but get: %v", s[0])
	}
}

func Test_header(t *testing.T) {
	s := injectTestCaseAndExec(header)
	if data := concatAll(&s); data != mHeader {
		t.Errorf("Expected:\n%s \nBut got: \n%s", mHeader, data)
	}
}

func Test_img(t *testing.T) {
	s := injectTestCaseAndExec(img)
	if s[0] != mImg {
		t.Errorf("Expected:\n%s \nBut got: \n%s", mImg, s[0])
	}
}

func Test_link(t *testing.T) {
	s := injectTestCaseAndExec(withLink)
	cs := injectTestCaseAndExec(codeLink)

	if s[0] != mWithLink {
		t.Errorf("\nExpected:\n%s\nBut got\n%s", mWithLink, s[0])
	}

	if cs[0] != mCodeLink {
		t.Errorf("\nExpected:\n%s\nBut got\n%s", mCodeLink, cs[0])
	}
}
