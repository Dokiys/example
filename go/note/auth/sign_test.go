package auth

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	values := url.Values{}
	values.Add("action", "inquiry")
	values.Add("deviceNo", "POS01")
	values.Add("shopNo", "CN123456")
	values.Add("brand", "663")
	values.Add("body", "ewogICAgICAgICAidHJhY2VObyI6ICI5OTAwMDAwOTEwMDAxMDEwMTczMjEyMyIsCiAgICAgICAgICJvcmlnaW5hbFRyYWNlTm8iOiAiOTkwMDAwMDkxMDAwMTAxMDE3MzIxMjQiCiAgICAgfQ==")
	values.Add("mwVersion", "20161010")
	values.Add("ptlVersion", "20161010")
	values.Add("posVersion", "20161010")
	values.Add("timestamp", "1483372334")
	values.Add("sign", "F38545F4D74B5C10A9EBBC053ED9D1CF")

	sign := Sign(values, "94365019BBF9CEEAB0DF658E67754A70")
	assert.Equal(t, "F38545F4D74B5C10A9EBBC053ED9D1CF", sign)
}
