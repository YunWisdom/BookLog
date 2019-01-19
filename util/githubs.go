

package util

import (
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

// GitHubUserInfo returns GitHub user info specified by the given access token.
func GitHubUserInfo(accessToken string) (ret map[string]interface{}) {
	response, data, errors := gorequest.New().Get("https://api.github.com/user?access_token=" + accessToken).Timeout(7 * time.Second).
		Set("User-Agent", "Pipe; +https://github.com/YunWisdom/BookLog").EndStruct(&ret)
	if nil != errors || http.StatusOK != response.StatusCode {
		logger.Errorf("get github user info failed: %+v, %s", errors, data)

		return nil
	}

	return
}
