package cleaner

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// TagToken type struct will be stored in the TokenStack
type TagToken struct {
	tagName string
}

// TokenStack can store tag that want to to be deferred rendering
type TokenStack []TagToken

var (
	title           string
	ignore          = false
	currentTag      string
	couldGetContent = false
	ruleTag         string
	ruleClass       string
	ruleValid       string
	ruleEnable      = false
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

// NewTagToken construct a new TagToken
func NewTagToken(name string) TagToken {
	return TagToken{
		tagName: name,
	}
}

// NewTokenStack construct a stack to actually store the TagToken
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

func DefineRules(className string) {
	ruleTag = "div"
	ruleClass = className
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
	tStack := NewTokenStack()

	for {
		tokenType := tr.Next()

		if tokenType == html.ErrorToken {
			break
		}

		token := tr.Token()

		if isInvalidTag(token.Data) {
			if tokenType == html.StartTagToken {
				tokenType = tr.Next()
				token = tr.Token()
				for !isInvalidTag(token.Data) {
					tokenType = tr.Next()
					token = tr.Token()
				}
			}
		}

		if ruleTag != "" && token.Data == ruleTag {
			if tokenType == html.StartTagToken {
				if ruleEnable {
					tStack.push(NewTagToken(ruleTag))
					continue
				}

				for _, a := range token.Attr {
					if a.Key == "class" {
						classAttrs := strings.Split(a.Val, " ")
						for _, className := range classAttrs {
							if className == ruleClass {
								tStack.push(NewTagToken(ruleClass))
								ruleEnable = true
								b.Reset()
								break
							}
						}
						continue
					}
				}
			} else if tokenType == html.EndTagToken {
				if ruleEnable {
					tStack.pop()
					if tStack.isEmpty() {
						break
					}
					continue
				}
			}
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

func isInvalidTag(name string) bool {
	return name == "script" || name == "style"
}
