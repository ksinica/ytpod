package youtube

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
)

var (
	errNotFound = errors.New(http.StatusText(http.StatusNotFound))
)

func drainAndClose(resp *http.Response) error {
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		resp.Body.Close()
		return err
	}
	return resp.Body.Close()
}

func httpFetch(ctx context.Context, url string, f func(r io.Reader) error) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient(ctx).Do(req)
	if err != nil {
		return err
	}
	defer drainAndClose(resp)

	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		return f(resp.Body)
	case http.StatusNotFound:
		return errNotFound
	default:
		return errors.New(http.StatusText(resp.StatusCode))
	}
}

type ctxhttpClient struct{}

func WithHttpClient(ctx context.Context, client *http.Client) context.Context {
	return context.WithValue(ctx, ctxhttpClient{}, client)
}

func httpClient(ctx context.Context) *http.Client {
	if client, ok := ctx.Value(ctxhttpClient{}).(*http.Client); ok {
		return client
	}
	return http.DefaultClient
}

func UseHTTPClient(client *http.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := WithHttpClient(r.Context(), client)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func httpRequestScheme(r *http.Request) string {
	scheme := r.Header.Get("X-Forwarded-Proto")
	if len(scheme) > 0 {
		return strings.ToLower(scheme)
	}
	return "http"
}
