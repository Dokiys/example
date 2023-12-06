package auth

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

func Sign(values url.Values, appSecret string) string {
	text := encode(values) + "&KEY=" + appSecret

	algorithm := md5.New()
	algorithm.Write([]byte(text))
	ciphertext := algorithm.Sum(nil)

	return strings.ToUpper(hex.EncodeToString(ciphertext))
}

func encode(v url.Values) string {
	if v == nil {
		return ""
	}

	var buf bytes.Buffer
	var keys = make([]string, 0, len(v))
	for k := range v {
		if v.Get(k) == "" || k == "sign" {
			continue
		}
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		var prefix = k + "="
		for _, v := range v[k] {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}

	return buf.String()
}
