package extend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ExtendPlatformBrand struct {
	APIBaseURL string
	Header     http.Header
}

var (
	brandExtend = ExtendPlatformBrand{
		APIBaseURL: "https://v.paywithextend.com",
		Header: http.Header{
			"accept-language":           {"en-US,en;q=0.7"},
			"cache-control":             {"no-cache"},
			"dnt":                       {"1"},
			"origin":                    {"https://app.paywithextend.com"},
			"pragma":                    {"no-cache"},
			"priority":                  {"u=1, i"},
			"referer":                   {"https://app.paywithextend.com/"},
			"sec-ch-ua":                 {`"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`},
			"sec-ch-ua-mobile":          {"?0"},
			"sec-ch-ua-platform":        {`"macOS"`},
			"sec-fetch-dest":            {"empty"},
			"sec-fetch-mode":            {"cors"},
			"sec-fetch-site":            {"same-site"},
			"sec-gpc":                   {"1"},
			"user-agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36"},
			"x-extend-app-id":           {"app.paywithextend.com"},
			"x-extend-brand":            {"br_2F0trP1UmE59x1ZkNIAqsg"},
			"x-extend-platform":         {"web"},
			"x-extend-platform-version": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36"},
		},
	}
)

type Client struct {
	auth  Authenticator
	brand ExtendPlatformBrand
	http  *http.Client
}

func New(auth Authenticator) *Client {
	return NewWithBrand(brandExtend, auth)
}

func NewWithBrand(brand ExtendPlatformBrand, auth Authenticator) *Client {
	return &Client{
		auth:  auth,
		brand: brand,
		http:  http.DefaultClient,
	}
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.http = client
}

func (c *Client) request(ctx context.Context, method string, path string, contentType string, body io.Reader, response any) error {
	req, err := http.NewRequestWithContext(ctx, method, c.brand.APIBaseURL+path, body)
	if err != nil {
		return err
	}

	accessToken, err := c.auth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	req.Header = c.brand.Header.Clone()
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/vnd.paywithextend.v2021-03-12+json")

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("bad status: %s", res.Status)
		}

		var errorResponse apiErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return fmt.Errorf("bad status: %s", res.Status)
		}

		return errorResponse
	}

	if response != nil {
		err = json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) jsonRequest(ctx context.Context, method string, path string, body any, response any) error {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	return c.request(ctx, method, path, "application/json", bodyReader, response)
}
