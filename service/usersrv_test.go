

package service

import "testing"

func TestGetUserByName(t *testing.T) {
	user := User.GetUserByName(testPlatformAdminName)
	if nil == user {
		t.Errorf("user is nil")

		return
	}
	user = User.GetUserByName("notfound")
	if nil != user {
		t.Errorf("user should be nil")
	}

}

func TestGetUser(t *testing.T) {
	user := User.GetUser(uint64(1))
	if nil == user {
		t.Errorf("user is nil")

		return
	}
	if 1 != user.ID {
		t.Errorf("id is not [1]")
	}
}

func TestGetBlogUsers(t *testing.T) {
	users, _ := User.GetBlogUsers(1, 1)
	if 1 > len(users) {
		t.Errorf("users is empty")

		return
	}
}

func TestGetUserBlogs(t *testing.T) {
	blogs := User.GetUserBlogs(1)
	if 1 > len(blogs) {
		t.Errorf("blogs is tempty")

		return
	}
}
