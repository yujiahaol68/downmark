package markdown

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type TagToken struct {
	tagName string
}

type TokenStack []TagToken

const (
	separator = "<>"
)

var (
	title  string
	ignore = true
)

var validTag = map[string]string{
	"a":          "%s",
	"blockquote": "> %s",
	"h1":         "# %s",
	"h2":         "## %s",
	"h3":         "### %s",
	"h4":         "#### %s",
	"p":          "%s",
	"strong":     "**%s**",
	"span":       "%s",
	"em":         "*%s*",
	"li":         "- %s",
	"hr":         "-----",
	"br":         "\n",
	"code":       " `%s` ",
}

var imgAttr = map[string]string{
	"img":   "src",
	"image": "href",
}

func NewTagToken(name string) TagToken {
	return TagToken{
		tagName: name,
	}
}

func NewTokenStack() TokenStack {
	return []TagToken{}
}

func (t *TokenStack) size() int {
	return len(*t)
}

func (t *TokenStack) match(newTagName string) bool {
	if t.isEmpty() {
		return false
	}
	topTag := (*t)[len(*t)-1]
	return topTag.tagName == newTagName
}

func (t *TokenStack) push(tag TagToken) {
	*t = append(*t, tag)
}

func (t *TokenStack) isEmpty() bool {
	return len(*t) == 0
}

func (t *TokenStack) pop() string {
	d := (*t)[len(*t)-1]
	(*t) = (*t)[:len(*t)-1]
	return d.tagName
}

func (t *TokenStack) peek() string {
	return (*t)[len(*t)-1].tagName
}

func (t *TokenStack) print() {
	s := *t
	for _, tag := range s {
		fmt.Printf("%s, ", tag.tagName)
	}
	fmt.Printf("\n\n")
}

func MdConvertor(tr *html.Tokenizer) []string {
	for {
		tt := tr.Next()
		t := tr.Token()

		if t.Data == "title" && tt == html.StartTagToken {
			tr.Next()
			title = tr.Token().Data
			break
		}
	}

	var converted []string
	var tempRenderStr string
	tStack := NewTokenStack()

	for {
		tokenType := tr.Next()

		if tokenType == html.ErrorToken {
			break
		}

		token := tr.Token()

		// DEBUG:
		// tStack.print()
		// fmt.Println(tempRenderStr)

		if token.Data == "hr" || token.Data == "br" {
			converted = append(converted, validTag[token.Data])
			continue
		}

		if tokenType == html.StartTagToken {
			_, contains := validTag[token.Data]

			if contains {
				tStack.push(NewTagToken(token.Data))
				ignore = false
				continue
			}
			continue
		} else if ignore == false && tokenType == html.TextToken && !tStack.isEmpty() {

			if tempRenderStr == "" || tStack.size() == 1 {
				tempRenderStr += token.Data
			} else {
				tempRenderStr = tempRenderStr + separator + token.Data
			}

		} else if tokenType == html.EndTagToken {
			_, contains := validTag[token.Data]

			if contains && tStack.match(token.Data) {
				if tStack.size() == 1 {
					converted = append(converted, renderNormal(tStack.pop(), tempRenderStr))
					tempRenderStr = ""
					ignore = true
				} else {
					ss := strings.Split(tempRenderStr, separator)
					if len(ss) == 2 {
						tempRenderStr = concat(ss[0], renderNormal(tStack.pop(), ss[1]))
					}
				}
			}
		}
	}
	return converted
}

func renderNormal(tagName string, s string) string {
	return fmt.Sprintf(validTag[tagName], s)
}

func concat(s1 string, s2 string) string {
	var b bytes.Buffer
	b.WriteString(s1)
	b.WriteString(s2)
	return b.String()
}

func getMemory(s string) string {
	ss := strings.Split(s, separator)
	return ss[len(ss)-1]
}
