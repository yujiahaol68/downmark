package downmark

import "testing"

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
