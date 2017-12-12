package rule

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

var (
	mutex sync.Mutex
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

// CleanConvertor return the title and the content of the html in cleaner way
func CleanConvertor(tr *html.Tokenizer, rules ...string) (string, bytes.Buffer) {
	var title string
	couldGetContent := false
	titleFound := false
	ruleEnable := false
	var ruleTag string
	var ruleValue string
	ruleAttr := "class"

	if len(rules) == 2 {
		ruleTag = rules[0]
		ruleValue = rules[1]
	} else if len(rules) == 3 {
		ruleTag = rules[0]
		ruleAttr = rules[1]
		ruleValue = rules[2]
	}

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

		if !titleFound && isTitleTag(token.Data) && tokenType == html.EndTagToken {
			b.Reset()
			titleFound = true
			continue
		}

		if ruleTag != "" && token.Data == ruleTag {
			if tokenType == html.StartTagToken {
				if ruleEnable {
					tStack.push(NewTagToken(ruleTag))
					continue
				}

				for _, a := range token.Attr {
					if a.Key == ruleAttr {
						vals := strings.Split(a.Val, " ")
						for _, val := range vals {
							if val == ruleValue {
								tStack.push(NewTagToken(ruleTag))
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

		if token.Data == "img" {
			for _, a := range token.Attr {
				if a.Key == "src" {
					b.WriteString(fmt.Sprintf("<img src=\"%s\">", a.Val))
					continue
				}
			}
		}

		if token.Data == "hr" || token.Data == "br" {
			b.WriteString(fmt.Sprintf("<%s>", token.Data))
			continue
		}

		if tokenType == html.StartTagToken {
			mutex.Lock()
			_, contains := validTag[token.Data]
			mutex.Unlock()

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
			mutex.Lock()
			_, has := validTag[token.Data]
			mutex.Unlock()

			if has {
				b.WriteString(createEndTag(token.Data))
			}
		}
	}

	return title, b
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

func isTitleTag(tagName string) bool {
	return tagName == "h1" || tagName == "h2"
}
