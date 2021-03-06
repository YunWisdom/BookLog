

package util

import (
	"runtime"
	"testing"
)

func TestIsWiindows(t *testing.T) {
	goos := runtime.GOOS

	if "windows" == goos && !IsWindows() {
		t.Error("runtime.GOOS returns [windows]")

		return
	}
}

func TestPwd(t *testing.T) {
	if "" == Pwd() {
		t.Error("Working directory should not be empty")

		return
	}
}

func TestHome(t *testing.T) {
	home, err := UserHome()
	if nil != err {
		t.Error("Can not get user home")

		return
	}

	t.Log(home)
}
