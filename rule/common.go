package rule

import "fmt"

// TagToken type struct will be stored in the TokenStack
type TagToken struct {
	tagName string
}

// TokenStack can store tag that want to to be deferred rendering
type TokenStack []TagToken

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
