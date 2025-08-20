package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RequestOptions captures all inputs required to make an HTTP request.
type RequestOptions struct {
	// BaseURL is the API base like "https://api.example.com" (without trailing slash).
	BaseURL string
	// Path is the API path like "/v1/orders". Leading slash optional.
	Path string
	// Method is the HTTP method: GET, POST, PUT, PATCH, DELETE, etc.
	Method string
	// QueryParams are appended to the URL.
	QueryParams map[string]string
	// Headers to include in the request. Content-Type and Accept can be set here.
	Headers map[string]string
	// Body to send. If Body is an io.Reader, it is sent as-is. Otherwise it is JSON-encoded.
	Body any
	// Timeout is optional per-request timeout. If zero, a default client without custom timeout is used.
	Timeout time.Duration
}

// MakeRequest performs an HTTP call based on the provided RequestOptions and returns
// the HTTP status code, response body bytes, and response headers.
//
// Behavior:
// - If opts.Body is an io.Reader, it is streamed as the request body.
// - If opts.Body is non-nil and not an io.Reader, it is JSON-encoded and Content-Type is set to application/json unless already set.
// - For GET/HEAD methods, non-nil bodies are ignored.
func MakeRequest(ctx context.Context, opts RequestOptions) (int, []byte, http.Header, error) {
	if strings.TrimSpace(opts.Method) == "" {
		return 0, nil, nil, fmt.Errorf("http method is required")
	}

	fullURL, err := buildURL(opts.BaseURL, opts.Path, opts.QueryParams)
	if err != nil {
		return 0, nil, nil, err
	}

	var bodyReader io.Reader
	// Only include body for methods that allow a body
	method := strings.ToUpper(opts.Method)
	if opts.Body != nil && method != http.MethodGet && method != http.MethodHead {
		if r, ok := opts.Body.(io.Reader); ok {
			bodyReader = r
		} else {
			encoded, encErr := json.Marshal(opts.Body)
			if encErr != nil {
				return 0, nil, nil, fmt.Errorf("failed to encode request body: %w", encErr)
			}
			bodyReader = bytes.NewReader(encoded)
			// Set Content-Type if not explicitly provided
			if opts.Headers == nil {
				opts.Headers = map[string]string{}
			}
			if _, exists := opts.Headers["Content-Type"]; !exists {
				opts.Headers["Content-Type"] = "application/json"
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return 0, nil, nil, err
	}

	// Default Accept header for APIs
	if opts.Headers == nil {
		opts.Headers = map[string]string{}
	}
	if _, ok := opts.Headers["Accept"]; !ok {
		opts.Headers["Accept"] = "application/json"
	}
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	if opts.Timeout > 0 {
		client.Timeout = opts.Timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, resp.Header, err
	}

	return resp.StatusCode, respBytes, resp.Header, nil
}

func buildURL(base, path string, query map[string]string) (string, error) {
	base = strings.TrimSpace(base)
	path = strings.TrimSpace(path)
	if base == "" && strings.HasPrefix(path, "http") {
		// Absolute URL provided in path
		u, err := url.Parse(path)
		if err != nil {
			return "", err
		}
		q := u.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		return u.String(), nil
	}

	if base == "" {
		return "", fmt.Errorf("base url is required when path is not absolute")
	}

	// Normalize slashes
	if strings.HasSuffix(base, "/") {
		base = strings.TrimSuffix(base, "/")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	u, err := url.Parse(base + path)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
