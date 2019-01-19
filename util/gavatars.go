

package util

import (
	"net/http"
	"time"
	"io/ioutil"
	"math/rand"
)

// RandAvatarData returns random avatar image byte array data from Gravatar (http://www.gravatar.com).
// Sees https://github.com/YunWisdom/BookLog/issues/131 for more details.
func RandAvatarData() (ret []byte) {
	modes := []string{"identicon", "monsterid", "wavatar", "retro", "robohash"}
	d := modes[rand.Intn(len(modes))]
	h := RandString(16)

	http.DefaultClient.Timeout = 2 * time.Second
	response, err := http.Get("http://www.gravatar.com/avatar/" + h + "?s=256&d=" + d)
	if nil != err {
		logger.Error("generate random avatar from Gavatar failed: " + err.Error())

		return nil
	}
	defer response.Body.Close()
	ret, err = ioutil.ReadAll(response.Body)
	if nil != err {
		logger.Error("generate random avatar from Gavatar failed: " + err.Error())

		return nil
	}

	return
}
