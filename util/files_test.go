

package util

import (
	"strconv"
	"testing"
)

func TestGetFileSize(t *testing.T) {
	size := File.GetFileSize(".")

	t.Log("size of file [.] is [" + strconv.FormatInt(size, 10) + "]")
}

func TestIsExist(t *testing.T) {
	if !File.IsExist(".") {
		t.Error(". must exist")

		return
	}
}

func TestIdBinary(t *testing.T) {
	if File.IsBinary("not binary content") {
		t.Error("The content should not be binary")

		return
	}
}

func TestIsImg(t *testing.T) {
	if !File.IsImg(".jpg") {
		t.Error(".jpg should be a valid extension of a image file")

		return
	}
}

func TestIsDir(t *testing.T) {
	if !File.IsDir(".") {
		t.Error(". should be a directory")

		return
	}
}
