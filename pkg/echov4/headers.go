package echov4

import (
	"net/http"
	"strings"
)

const (
	HTTP_HEAD_AUTHORIZATION     = "Authorization"
	AUTHORIZATION_SCHEMA_BEARER = "Bearer"
)

func GetAuthorizationHeader(header http.Header) string {
	return header.Get(HTTP_HEAD_AUTHORIZATION)
}

func ParseAuthorizationHeader(header http.Header) (schema string, token string, ok bool) {
	auth := header.Get(HTTP_HEAD_AUTHORIZATION)

	substr := " "
	split := strings.Index(auth, substr)
	if split < 0 {
		return
	}

	tokens := []string{
		auth[:split],
		auth[split+len(substr):],
	}

	ss := make([]string, 2)
	return ss[0], ss[1], copy(ss, tokens) == 2
}

func SeHttpHeader(header http.Header, key string, value ...string) {
	for i, value := range value {
		switch i {
		case 0:
			header.Set(key, value)
		default:
			header.Add(key, value)
		}
	}
}
