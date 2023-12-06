package auth

import (
	"crypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRsaEnCrypt(t *testing.T) {
	var text = "123456"
	encrypt, err := RSAEncrypt([]byte(text), []byte(PublicKey))
	assert.NoError(t, err)
	decrypt, err := RSADecrypt(encrypt, []byte(PrivateKey))
	assert.NoError(t, err)

	assert.Equal(t, text, decrypt)
}

func TestRSASign(t *testing.T) {
	var text = "123456"
	signed, err := RSASign([]byte(text), []byte(PrivateKey), crypto.MD5)
	assert.NoError(t, err)
	err = RSAVerify([]byte(text), signed, []byte(PublicKey), crypto.MD5)
	assert.NoError(t, err)
}

const PublicKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuhSuHlW9k8hDU/bz1XSr
eoVHeUo1KKTdBGJdeDrdN09rc4A+xLStYtUUb1l5OBAo/S8tkLMSznlCh1hpRN8X
aqHLeY6zQUFtJIOCyWmKAXZ+1+gzNRLx9oDl2z2RWEDgvdsgH9k9IeyED++889PO
i73gzXzo4MlTObNBc2ms5TGWG3E3BoAiwCFDh5CaaUH2A7RJqn2lJchkk68COoNt
mjhBot81+Ayzo9FejRzAArgMC4ggGCnthmVLppJZL5nkdrcb3hnzOV+t+Rn265sl
iUTsU7x7npMT1BB3VAAi8JCuIfXkq4v1hVpIrnh/AoA/83TvSl8Nyzhd8xrIX9Rw
hQIDAQAB
-----END PUBLIC KEY-----
`
const PrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAuhSuHlW9k8hDU/bz1XSreoVHeUo1KKTdBGJdeDrdN09rc4A+
xLStYtUUb1l5OBAo/S8tkLMSznlCh1hpRN8XaqHLeY6zQUFtJIOCyWmKAXZ+1+gz
NRLx9oDl2z2RWEDgvdsgH9k9IeyED++889POi73gzXzo4MlTObNBc2ms5TGWG3E3
BoAiwCFDh5CaaUH2A7RJqn2lJchkk68COoNtmjhBot81+Ayzo9FejRzAArgMC4gg
GCnthmVLppJZL5nkdrcb3hnzOV+t+Rn265sliUTsU7x7npMT1BB3VAAi8JCuIfXk
q4v1hVpIrnh/AoA/83TvSl8Nyzhd8xrIX9RwhQIDAQABAoIBAD0pcPWtjutVJrQh
dpHRkE9sIUr/lituzOqU/k33YyC77QCAxaDYFilnChlzWkGQJjjZ2es2wNa2yqQ7
7hXiEgtNdNjWi5rGS35173YOR3slnJmQy8wFFxzDz0lZmj1czcH5lTyCpfJRkDLA
xOWL19bcqVhBVzYTVlSf08KAmICYVP9IdhMdaL3yGH5JeDVq5IyLpHU/OxgAdBxA
+jKV2Ux8qULzAvqKaTdexOzUICmSPGQUDdtGI9ejC+dCmejDcNHm5mAzTasRvl4G
zK9w2KAKn7KG6kyqzS6R1N1x+dXoXA2ok6RDyA1uElWfwtvTpiUUclQpHAl6WOEL
cR57rwECgYEA6guOkqmI+jg+JbJuB5SpjnjN6ePCJyW77lA1Gpbfja7QA5yhXk6w
SV6sMvj6pzkMKr6bK49O85Wp4SggSOv0NABKHhHVxNeAhFdZUHB4eB4sObj6LB6b
v/MEJABRtl/lkij77mZZb+brWly/Fdv+jxMDeHA/7hBkzfoWWm0ibMECgYEAy4lK
Psf0Raq74ZYf9wpStxpUU+e6GQEox+yWxDefgKy3yXe4vS0/nEUfISYkrQM15J/W
UP2bXr7qFEsZbaWX9LvhiZ2cwawMeOoyhHED/cgKbX7V4tWssioyKPS9TMalBmPz
lqOc+lUL1JfjW0xXGQ7PDGKYiCBinBD45lZGwMUCgYEAr/IIOJFi/FiTv9snhGNq
JEUE57PlHXDsmveJNHf/j4+/qTdyyGb3d/DIG3m5VUU5tFieZlzRyaTVlQKJYsif
SQh+r6RQxC4N22+fIS2sIwDr1mkNCWXpSJ/0mOv0gdoNx5cv7cTbr5g1jjTzIgfE
kuKEVWJtbRa98Wr0qv4oRwECgYBmSpx2yvVYIgOEz3dHJ/gEMAZbmOVtdaiyOZRY
DwBpoeRIK5Q790a12gNYHJxoG2n1eeMzFxID3v0zr76a3ZNuGxKxn/XNXBN0nXdA
GrB/1g5vk0QZWXwOmqhU7xNIR7leadNdTOMy6JUmhiNsmgRYAppKRi8UkvocJ2eA
E7JBZQKBgHTL1HTsUVWGyPkIN4VBFaHwYK0xK+ElF/SanwpA304lY0wwNoqB9lio
bBD3fuXOEsR4/LhZjvfKEFa3VbT4UmrP/pVX3wksC6/J2VtNCgk05jmgHj0hLkzM
OSrb9zbsouFTQZUDyEw5obZiFlZclPXYFD1e84qMd14ss/fzdGhA
-----END RSA PRIVATE KEY-----
`
