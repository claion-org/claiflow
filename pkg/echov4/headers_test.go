package echov4

import (
	"net/http"
	"testing"
)

func Test_ParseAuthorizationHeader(t *testing.T) {
	var h = http.Header{}
	h.Add(HTTP_HEAD_AUTHORIZATION, "Bearer some token")

	schema, token, ok := ParseAuthorizationHeader(h)

	t.Logf("schema: %q\n", schema)
	t.Logf("token: %q\n", token)
	t.Logf("ok: %v\n", ok)
}
