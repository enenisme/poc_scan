// pkg/utils/http.go
package utils

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/enenisme/poc_scan/pkg/model"
)

type HTTPClient struct {
	client *http.Client
}

type HTTPResponse struct {
	Status  int
	Headers http.Header
	Body    []byte
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (c *HTTPClient) DoRequest(baseURL *url.URL, req *model.Request) (*HTTPResponse, error) {
	// Construct full URL
	targetURL := *baseURL
	targetURL.Path = path.Join(targetURL.Path, req.Path)

	// Create HTTP request
	httpReq, err := http.NewRequest(req.Method, targetURL.String(), bytes.NewBufferString(req.Body))
	if err != nil {
		return nil, err
	}

	// Add headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Set default headers if not provided
	if httpReq.Header.Get("User-Agent") == "" {
		httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	}

	// Handle redirects based on configuration
	if !req.FollowRedirects {
		c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// Execute request
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HTTPResponse{
		Status:  resp.StatusCode,
		Headers: resp.Header,
		Body:    body,
	}, nil
}
