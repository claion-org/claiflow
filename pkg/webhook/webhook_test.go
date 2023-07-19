package webhook

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestWebhook_Publish(t *testing.T) {
	var header = http.Header{}
	header.Add("X-Request-Id", "123")

	var webhook = Config{
		URL:     "http://localhost/opaque",
		Method:  "POST",
		Headers: header,
	}

	payload := `"hello"`

	if err := webhook.Publish(context.Background(), []byte(payload)); err != nil {
		t.Fatal(err)
	}
}

func TestDefer(t *testing.T) {
	ctx := context.Background()
	if true {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*1)
		defer func() {
			fmt.Println("cancel")
			cancel()
		}()
	}

	<-ctx.Done()

	fmt.Println("done")
}
