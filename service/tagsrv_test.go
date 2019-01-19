

package service

import "testing"

func TestGetTags(t *testing.T) {
	tags := Tag.GetTags(2, 1)
	if nil == tags {
		t.Errorf("tags is nil")
	}
}
