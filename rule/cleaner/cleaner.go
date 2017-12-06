package cleaner

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

var (
	title           string
	ignore          = false
	currentTag      string
	couldGetContent = false
)

var validTag = map[string]string{
	"blockquote": "",
	"h1":         "",
	"h2":         "",
	"h3":         "",
	"h4":         "",
	"p":          "",
	"strong":     "",
	"span":       "",
	"em":         "",
	"li":         "",
	"ol":         "",
	"ul":         "",
	"code":       "",
	"pre":        "",
}

func CleanConvertor(tr *html.Tokenizer) bytes.Buffer {
	for {
		tt := tr.Next()
		t := tr.Token()

		if t.Data == "title" && tt == html.StartTagToken {
			tr.Next()
			title = tr.Token().Data
			break
		}
	}

	var b bytes.Buffer

	for {
		tokenType := tr.Next()

		if tokenType == html.ErrorToken {
			break
		}

		token := tr.Token()

		if token.Data == "style" || token.Data == "script" {
			if tokenType == html.StartTagToken {
				ignore = true
			} else {
				ignore = false
			}
			continue
		}

		if token.Data == "article" {
			if tokenType == html.StartTagToken {
				b.Reset()
				continue
			} else if tokenType == html.EndTagToken {
				break
			}
		}

		if token.Data == "hr" || token.Data == "br" {
			b.WriteString(fmt.Sprintf("<%s>", token.Data))
			continue
		}

		if tokenType == html.StartTagToken {
			_, contains := validTag[token.Data]

			if contains {
				couldGetContent = true
				b.WriteString(createStartTag(token.Data))
			} else {
				couldGetContent = false
			}

		} else if tokenType == html.TextToken && couldGetContent {
			es := html.EscapeString(token.Data)

			if strings.TrimSpace(es) != "" {
				b.WriteString(es)
			}
		} else if tokenType == html.EndTagToken {
			_, has := validTag[token.Data]

			if has {
				b.WriteString(createEndTag(token.Data))
			}
		}
	}

	return b
}

func createStartTag(tagName string) string {
	return fmt.Sprintf("<%s>", tagName)
}

func createEndTag(tagName string) string {
	return fmt.Sprintf("</%s>", tagName)
}
