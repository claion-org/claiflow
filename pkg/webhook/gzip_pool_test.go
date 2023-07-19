package webhook_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/claion-org/claiflow/pkg/webhook"
)

func TestGzip(t *testing.T) {
	w := webhook.GZipWriter.Get()
	defer webhook.GZipWriter.Put(w)

	var encode bytes.Buffer
	w.Reset(&encode)

	s := "hello, world!"

	n, err := io.Copy(w, bytes.NewReader([]byte(s)))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("n", n)

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	t.Log(encode.Bytes())

	r, err := gzip.NewReader(&encode)
	if err != nil {
		t.Fatal(err)
	}

	var decode bytes.Buffer
	n, err = io.Copy(&decode, r)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("n", n)

	if err := r.Close(); err != nil {
		t.Fatal(err)
	}

	t.Log(decode.String())

}
