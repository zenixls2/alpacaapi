package alpacaapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/c9s/requestgen"
)

const defaultHTTPTimeout = time.Second * 10
const RestBaseURL = "https://api.alpaca.markets"
const PaperBaseURL = "https://paper-api.alpaca.markets"

var DebugRequestResponse = false

type AuthMethod int

const (
	AuthMethodOAuth AuthMethod = iota
	AuthMethodAPIKey
	AuthMethodBrokerKey
)

type RestClient struct {
	requestgen.BaseAPIClient
	authMethod   AuthMethod
	APIKey       string
	APISecret    string
	BrokerKey    string
	BrokerSecret string
	OAuthToken   string
	RetryLimit   int
	RetryDelay   time.Duration

	AccountService *AccountService
	OrderService   *OrderService
}

func NewClient() *RestClient {
	u, err := url.Parse(RestBaseURL)
	if err != nil {
		panic(err)
	}
	client := &RestClient{
		BaseAPIClient: requestgen.BaseAPIClient{
			BaseURL: u,
			HttpClient: &http.Client{
				Timeout: defaultHTTPTimeout,
			},
		},
		RetryLimit: 3,
		RetryDelay: time.Second,
	}
	client.AccountService = &AccountService{client: client}
	client.OrderService = &OrderService{client: client}
	return client
}

func NewPaperClient() *RestClient {
	u, err := url.Parse(PaperBaseURL)
	if err != nil {
		panic(err)
	}
	client := &RestClient{
		BaseAPIClient: requestgen.BaseAPIClient{
			BaseURL: u,
			HttpClient: &http.Client{
				Timeout: defaultHTTPTimeout,
			},
		},
		RetryLimit: 3,
		RetryDelay: time.Second,
	}
	client.AccountService = &AccountService{client: client}
	client.OrderService = &OrderService{client: client}
	return client
}

func (c *RestClient) SetAuthByOAuth(token string) {
	c.OAuthToken = token
	c.authMethod = AuthMethodOAuth
}

func (c *RestClient) SetAuthByAPIKey(key, secret string) {
	c.APIKey = key
	c.APISecret = secret
	c.authMethod = AuthMethodAPIKey
}

func (c *RestClient) SetAuthByBrokerKey(key, secret string) {
	c.BrokerKey = key
	c.BrokerSecret = secret
	c.authMethod = AuthMethodBrokerKey
}

func (c *RestClient) NewAuthenticatedRequest(ctx context.Context, method, refURL string, params url.Values, payload interface{}) (*http.Request, error) {
	rel, err := url.Parse(refURL)
	if err != nil {
		return nil, err
	}

	if params != nil {
		rel.RawQuery = params.Encode()
	}

	pathURL := c.BaseURL.ResolveReference(rel)
	path := pathURL.Path
	if rel.RawQuery != "" {
		path += "?" + rel.RawQuery
	}

	body, err := castPayload(payload)
	if err != nil {
		return nil, err
	}
	if DebugRequestResponse {
		log.Printf("Request URL: %s, Method: %s, Body: %s\n", path, method, string(body))
	}
	req, err := http.NewRequestWithContext(ctx, method, pathURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	switch c.authMethod {
	case AuthMethodOAuth:
		if c.OAuthToken == "" {
			return nil, errors.New("OAuth token is not set")
		}
		req.Header.Set("Authorization", "Bearer "+c.OAuthToken)
	case AuthMethodAPIKey:
		if c.APIKey == "" || c.APISecret == "" {
			return nil, errors.New("API key or secret is not set")
		}
		req.Header.Set("APCA-API-KEY-ID", c.APIKey)
		req.Header.Set("APCA-API-SECRET-KEY", c.APISecret)
	case AuthMethodBrokerKey:
		if c.BrokerKey == "" || c.BrokerSecret == "" {
			return nil, errors.New("Broker key or secret is not set")
		}
		req.SetBasicAuth(c.BrokerKey, c.BrokerSecret)
	default:
		return nil, errors.New("authentication method is not set")
	}
	return req, nil
}

func castPayload(payload interface{}) ([]byte, error) {
	if payload == nil {
		return nil, nil
	}

	switch v := payload.(type) {
	case string:
		return []byte(v), nil

	case []byte:
		return v, nil

	}
	return json.Marshal(payload)
}

var defaultHttpClient = &http.Client{
	Timeout: defaultHTTPTimeout,
}

func (c *RestClient) NewRequest(ctx context.Context, method, refURL string, params url.Values, payload interface{},
) (*http.Request, error) {
	return c.BaseAPIClient.NewRequest(ctx, method, refURL, params, payload)
}

func (c *RestClient) SendRequest(req *http.Request) (*requestgen.Response, error) {
	if c.HttpClient == nil {
		c.HttpClient = defaultHttpClient
	}

	var resp *http.Response
	var err error

	for i := 0; ; i++ {
		resp, err = c.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusTooManyRequests {
			break
		}
		if i >= c.RetryLimit {
			break
		}
		time.Sleep(c.RetryDelay)
	}

	if err = verify(resp); err != nil {
		return nil, err
	}

	response, err := requestgen.NewResponse(resp)
	if DebugRequestResponse {
		log.Printf("Response Status: %s, err: %v, Body: %s\n", resp.Status, err, string(response.Body))
	}
	if err != nil {
		return response, err
	}

	if response.IsError() {
		return response, &ErrResponse{Response: response, Body: response.Body, Request: req}
	}

	return response, nil
}

type ErrResponse struct {
	Response *requestgen.Response
	Body     []byte
	Request  *http.Request
}

func (e *ErrResponse) Error() string {
	return fmt.Sprintf("request failed with status code: %d, body: %q", e.Response.StatusCode, string(e.Body))
}

func verify(resp *http.Response) error {
	if resp == nil {
		return errors.New("response is nil")
	}
	if resp.StatusCode >= http.StatusMultipleChoices {
		defer resp.Body.Close()
		return APIErrorFromResponse(resp)
	}
	return nil
}

type APIError struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Body       string `json:"-"`
}

func APIErrorFromResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		// If the error is not in our JSON format, we simply return the HTTP response
		return fmt.Errorf("%s (HTTP %d)", body, resp.StatusCode)
	}
	apiErr.StatusCode = resp.StatusCode
	apiErr.Body = strings.TrimSpace(string(body))
	return &apiErr
}

func (e *APIError) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("%s (HTTP %d, Code %d)", e.Message, e.StatusCode, e.Code)
	}
	return fmt.Sprintf("%s (HTTP %d)", e.Message, e.StatusCode)
}

type APIResponse json.RawMessage
