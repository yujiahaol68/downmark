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
	title        string
	ignore       = true
	insideLink   = false
	imgTemplate  = "![](%s)"
	linkTemplate = "[%s](%s)"
	aText        string
	href         string
)

var validTag = map[string]string{
	"a":          "%s",
	"blockquote": "> %s",
	"h1":         "# %s",
	"h2":         "## %s",
	"h3":         "### %s",
	"h4":         "#### %s",
	"p":          "%s",
	"strong":     "__%s__",
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

		if token.Data == "a" && tokenType == html.StartTagToken {
			for _, a := range token.Attr {
				if a.Key == "href" {
					href = a.Val
					insideLink = true
					break
				}
			}
			continue
		}

		if isImgTag(token.Data) {
			attrName := imgAttr[token.Data]
			for _, a := range token.Attr {
				if a.Key == attrName {
					converted = append(converted, fmt.Sprintf(imgTemplate, a.Val))
					break
				}
			}
			continue
		}

		if tokenType == html.StartTagToken {
			_, contains := validTag[token.Data]

			if contains && !insideLink {
				tStack.push(NewTagToken(token.Data))
				ignore = false
				continue
			}
			continue
		} else if ignore == false && tokenType == html.TextToken && !tStack.isEmpty() {

			if insideLink {
				// Has separator or not to judge whether render outter tag first
				if strings.Contains(tempRenderStr, separator) {
					tempRenderStr += fmt.Sprintf(linkTemplate, token.Data, href)
				} else {
					topTagName := tStack.peek()
					s, needPop := renderLink(token.Data, href, topTagName)
					if needPop {
						tStack.pop()
					}
					tempRenderStr += s
				}

				insideLink = false
				href = ""
				continue
			}

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

func renderLink(text string, href string, outterT string) (string, bool) {
	if text == "" {
		text = href
	}
	if outterT == "p" || outterT == "span" {
		return fmt.Sprintf(linkTemplate, text, href), false
	}
	return fmt.Sprintf(linkTemplate, renderNormal(outterT, text), href), true
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

func isImgTag(tagName string) bool {
	_, contain := imgAttr[tagName]
	return contain
}
