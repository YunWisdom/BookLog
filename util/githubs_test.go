

package util

import (
	"testing"
)

func TestGitHubUserInfo(t *testing.T) {
	user := GitHubUserInfo("error tk")
	if nil != user {
		t.Error("get a user")

		return
	}
}
