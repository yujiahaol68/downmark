package downmark

import (
	"fmt"
	"io"
	"strings"
)

// DLink save all the links
type DLink []string

// DBody is a concurrent safe map with key of link index and the res.Body of the link
type DBody map[string]io.ReadCloser

// NewDLink produce a list contains link
func NewDLink() DLink {
	d := DLink{}
	return d
}

func newDBody() DBody {
	d := make(DBody)
	return d
}

// AddLink append link to the list
func (d *DLink) AddLink(url string) error {
	if isLink(url) {
		*d = append(*d, completeProtocol(url))
		return nil
	}
	return fmt.Errorf("Expect URL but %v is not a url", url)
}

func isLink(l string) bool {
	// use regex to deal it
	return true
}

func completeProtocol(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}
	return "http://" + url
}

// GetBodys order go-routines to get the http body the concurrent
func (d DLink) GetBodys() DBody {
	dBo := newDBody()

	// go routine get the body and save in DBody map concurrent-safe way
	return dBo
}
