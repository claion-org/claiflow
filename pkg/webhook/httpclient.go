package webhook

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

var (
	SetContentEncoding_GZIP = setHeader(http.Header{"Content-Encoding": []string{"gzip"}})
	SetContentType_Json     = setHeader(http.Header{"Content-Type": []string{"application/json"}})
	SetAccept_Json          = setHeader(http.Header{"Accept": []string{"application/json"}})
	SetAccept_Plain         = setHeader(http.Header{"Accept": []string{"text/plain"}})
)

func setHeader(header http.Header) OpRequest {
	return func(req *http.Request) error {
		// Header
		for k, v := range header {
			for i, v := range v {
				if i == 0 {
					req.Header.Set(k, v)
				} else {
					req.Header.Add(k, v)
				}
			}
		}

		return nil
	}
}

var DefaultClient = func() func() *http.Client {
	client := http.DefaultClient

	SetTransport(client, func(t *http.Transport) {
		t.Proxy = http.ProxyFromEnvironment
		t.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext
		t.MaxIdleConns = 100
		t.IdleConnTimeout = 90 * time.Second
		t.TLSHandshakeTimeout = 10 * time.Second
		t.ExpectContinueTimeout = 1 * time.Second
		t.MaxIdleConnsPerHost = runtime.GOMAXPROCS(0) + 1
		t.TLSClientConfig.InsecureSkipVerify = true
	})
	return func() *http.Client {
		return client
	}
}()

func NewClient() *http.Client {
	client := http.Client{}

	SetTransport(&client, func(t *http.Transport) {
		t.Proxy = http.ProxyFromEnvironment
		t.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext
		t.MaxIdleConns = 100
		t.IdleConnTimeout = 90 * time.Second
		t.TLSHandshakeTimeout = 10 * time.Second
		t.ExpectContinueTimeout = 1 * time.Second
		t.MaxIdleConnsPerHost = runtime.GOMAXPROCS(0) + 1
		t.TLSClientConfig.InsecureSkipVerify = true
	})

	return &client
}

type OpTransport func(*http.Transport)

func SetTransport(client *http.Client, opts ...OpTransport) *http.Client {
	if client.Transport == nil {
		client.Transport = new(http.Transport)
	}
	if client.Transport.(*http.Transport).TLSClientConfig == nil {
		client.Transport.(*http.Transport).TLSClientConfig = new(tls.Config)
	}

	for _, op := range opts {
		op(client.Transport.(*http.Transport))
	}

	return client
}

func ComposeMiddleware[A any](handler func(A) error, middleware ...func(next func(A) error) func(A) error) func(A) error {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}

type OpResponse = func(resp *http.Response) error
type OpResponseMiddleware = func(next OpResponse) OpResponse

func Do(client *http.Client, req *http.Request, ops ...OpResponseMiddleware) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// close resp.Body
	defer resp.Body.Close()

	return ComposeMiddleware(func(*http.Response) error { return nil }, ops...)(resp)
}

func CheckStatus(next OpResponse) OpResponse {
	return func(resp *http.Response) error {
		if err := next(resp); err != nil {
			return err
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, resp.Body); err != nil {
				return err
			}

			if 0 < buf.Len() {
				return fmt.Errorf("%s %s, status code : %d(%s), body : %s", resp.Request.Method, resp.Request.URL.String(), resp.StatusCode, resp.Status, strings.TrimSpace(buf.String()))
			} else {
				return fmt.Errorf("%s %s, status code : %d(%s)", resp.Request.Method, resp.Request.URL.String(), resp.StatusCode, resp.Status)
			}
		}

		return nil
	}
}

func FromJson[T any](out *T) func(next OpResponse) OpResponse {
	return func(next OpResponse) OpResponse {
		return func(resp *http.Response) error {
			if err := next(resp); err != nil {
				return err
			}

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, resp.Body); err != nil {
				return err
			}

			if err := json.Unmarshal(buf.Bytes(), out); err != nil {
				return err
			}

			return nil
		}
	}
}

func GetBody(w io.Writer) func(next OpResponse) OpResponse {
	return func(next OpResponse) OpResponse {
		return func(resp *http.Response) error {
			if err := next(resp); err != nil {
				return err
			}

			if _, err := io.Copy(w, resp.Body); err != nil {
				return err
			}

			return nil
		}
	}
}

type OpRequest = func(*http.Request) error
type OpRequestMiddleware = func(next OpRequest) OpRequest

func NewRequest(ctx context.Context, method string, url string, ops ...OpRequestMiddleware) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	err = ComposeMiddleware(func(*http.Request) error { return nil }, ops...)(req)

	return req, err
}

func ToJson[T any](in T) func(next OpRequest) OpRequest {
	b, err := json.Marshal(in)
	if err != nil {
		return func(next OpRequest) OpRequest {
			return func(req *http.Request) error {
				return err
			}
		}
	}

	return func(next OpRequest) OpRequest {
		return func(req *http.Request) error {
			// MIME
			SetContentType_Json(req)
			// BODY
			req.Body = io.NopCloser(bytes.NewReader(b))

			return next(req)
		}
	}
}

var GZipWriter = NewPool(func() *gzip.Writer {
	w, _ := gzip.NewWriterLevel(nil, gzip.BestCompression)
	return w
})

func Compress_GZip(next OpRequest) OpRequest {
	return func(req *http.Request) error {
		if err := next(req); err != nil {
			return err
		}

		// MIME
		SetContentEncoding_GZIP(req)

		// read origin body
		var body bytes.Buffer
		if _, err := io.Copy(&body, req.Body); err != nil {
			return err
		}

		// get gzip encoder from pool
		gw := GZipWriter.Get()
		defer GZipWriter.Put(gw) // put pool

		var buf bytes.Buffer
		gw.Reset(&buf)
		defer gw.Reset(nil)

		// write origin body
		if _, err := gw.Write(body.Bytes()); err != nil {
			return err
		}

		if err := gw.Close(); err != nil {
			return err
		}

		// restore compressed request body
		req.Body = io.NopCloser(&buf)

		return nil
	}
}

func SetQuery(query url.Values) func(next OpRequest) OpRequest {
	return func(next OpRequest) OpRequest {
		return func(req *http.Request) error {
			// URL.Query
			for k, v := range query {
				for i, v := range v {
					if i == 0 {
						q := req.URL.Query()
						q.Set(k, v)
						req.URL.RawQuery = q.Encode()
					} else {
						q := req.URL.Query()
						q.Add(k, v)
						req.URL.RawQuery = q.Encode()
					}
				}
			}

			return next(req)
		}
	}
}

func SetHeader(header http.Header) func(next OpRequest) OpRequest {
	return func(next OpRequest) OpRequest {
		return func(req *http.Request) error {
			// Header
			for k, v := range header {
				for i, v := range v {
					if i == 0 {
						req.Header.Set(k, v)
					} else {
						req.Header.Add(k, v)
					}
				}
			}

			return next(req)
		}
	}
}

func SetBody(r_ io.Reader) func(next OpRequest) OpRequest {
	return func(next OpRequest) OpRequest {
		return func(req *http.Request) error {

			req.Body = io.NopCloser(r_)

			return next(req)
		}
	}
}
