package fastly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/google/go-querystring/query"
	"github.com/google/jsonapi"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/mitchellh/mapstructure"
)

// APIKeyEnvVar is the name of the environment variable where the Fastly API
// key should be read from.
const APIKeyEnvVar = "FASTLY_API_KEY"

// APIKeyHeader is the name of the header that contains the Fastly API key.
const APIKeyHeader = "Fastly-Key"

// EndpointEnvVar is the name of an environment variable that can be used
// to change the URL of API requests.
const EndpointEnvVar = "FASTLY_API_URL"

// DefaultEndpoint is the default endpoint for Fastly. Since Fastly does not
// support an on-premise solution, this is likely to always be the default.
const DefaultEndpoint = "https://api.fastly.com"

// RealtimeStatsEndpointEnvVar is the name of an environment variable that can be used
// to change the URL of realtime stats requests.
const RealtimeStatsEndpointEnvVar = "FASTLY_RTS_URL"

// DefaultRealtimeStatsEndpoint is the realtime stats endpoint for Fastly.
const DefaultRealtimeStatsEndpoint = "https://rt.fastly.com"

// ProjectURL is the url for this library.
var ProjectURL = "github.com/fastly/go-fastly"

// ProjectVersion is the version of this library.
var ProjectVersion = "6.0.0"

// UserAgent is the user agent for this particular client.
var UserAgent = fmt.Sprintf("FastlyGo/%s (+%s; %s)",
	ProjectVersion, ProjectURL, runtime.Version())

// Client is the main entrypoint to the Fastly golang API library.
type Client struct {
	// Address is the address of Fastly's API endpoint.
	Address string

	// HTTPClient is the HTTP client to use. If one is not provided, a default
	// client will be used.
	HTTPClient *http.Client

	// updateLock forces serialization of calls that modify a service.
	// Concurrent modifications have undefined semantics.
	updateLock sync.Mutex

	// apiKey is the Fastly API key to authenticate requests.
	apiKey string

	// url is the parsed URL from Address
	url *url.URL
}

// RTSClient is the entrypoint to the Fastly's Realtime Stats API.
type RTSClient struct {
	client *Client
}

// DefaultClient instantiates a new Fastly API client. This function requires
// the environment variable `FASTLY_API_KEY` is set and contains a valid API key
// to authenticate with Fastly.
func DefaultClient() *Client {
	client, err := NewClient(os.Getenv(APIKeyEnvVar))
	if err != nil {
		panic(err)
	}
	return client
}

// NewClient creates a new API client with the given key and the default API
// endpoint. Because Fastly allows some requests without an API key, this
// function will not error if the API token is not supplied. Attempts to make a
// request that requires an API key will return a 403 response.
func NewClient(key string) (*Client, error) {
	endpoint, ok := os.LookupEnv(EndpointEnvVar)

	if !ok {
		endpoint = DefaultEndpoint
	}

	return NewClientForEndpoint(key, endpoint)
}

// NewClientForEndpoint creates a new API client with the given key and API
// endpoint. Because Fastly allows some requests without an API key, this
// function will not error if the API token is not supplied. Attempts to make a
// request that requires an API key will return a 403 response.
func NewClientForEndpoint(key string, endpoint string) (*Client, error) {
	client := &Client{apiKey: key, Address: endpoint}
	return client.init()
}

// NewRealtimeStatsClient instantiates a new Fastly API client for the realtime stats.
// This function requires the environment variable `FASTLY_API_KEY` is set and contains
// a valid API key to authenticate with Fastly.
func NewRealtimeStatsClient() *RTSClient {
	endpoint, ok := os.LookupEnv(RealtimeStatsEndpointEnvVar)

	if !ok {
		endpoint = DefaultRealtimeStatsEndpoint
	}

	c, err := NewClientForEndpoint(os.Getenv(APIKeyEnvVar), endpoint)
	if err != nil {
		panic(err)
	}
	return &RTSClient{client: c}
}

// NewRealtimeStatsClientForEndpoint creates an RTSClient from a token and endpoint url.
// `token` is a Fastly API token and `endpoint` is RealtimeStatsEndpoint for the production
// realtime stats API.
func NewRealtimeStatsClientForEndpoint(token, endpoint string) (*RTSClient, error) {
	c, err := NewClientForEndpoint(token, endpoint)
	if err != nil {
		return nil, err
	}
	return &RTSClient{client: c}, nil
}

func (c *Client) init() (*Client, error) {
	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}
	c.url = u

	if c.HTTPClient == nil {
		c.HTTPClient = cleanhttp.DefaultClient()
	}

	return c, nil
}

// Get issues an HTTP GET request.
func (c *Client) Get(p string, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}
	ro.Parallel = true
	return c.Request("GET", p, ro)
}

// Head issues an HTTP HEAD request.
func (c *Client) Head(p string, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}
	ro.Parallel = true
	return c.Request("HEAD", p, ro)
}

// Patch issues an HTTP PATCH request.
func (c *Client) Patch(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("PATCH", p, ro)
}

// PatchForm issues an HTTP PUT request with the given interface form-encoded.
func (c *Client) PatchForm(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm("PATCH", p, i, ro)
}

// PatchJSON issues an HTTP PUT request with the given interface json-encoded.
func (c *Client) PatchJSON(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSON("PATCH", p, i, ro)
}

// PatchJSONAPI issues an HTTP PUT request with the given interface json-encoded.
func (c *Client) PatchJSONAPI(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPI("PATCH", p, i, ro)
}

// Post issues an HTTP POST request.
func (c *Client) Post(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("POST", p, ro)
}

// PostForm issues an HTTP POST request with the given interface form-encoded.
func (c *Client) PostForm(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm("POST", p, i, ro)
}

// PostJSON issues an HTTP POST request with the given interface json-encoded.
func (c *Client) PostJSON(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSON("POST", p, i, ro)
}

// PostJSONAPI issues an HTTP POST request with the given interface json-encoded.
func (c *Client) PostJSONAPI(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPI("POST", p, i, ro)
}

// PostJSONAPIBulk issues an HTTP POST request with the given interface json-encoded and bulk requests.
func (c *Client) PostJSONAPIBulk(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPIBulk("POST", p, i, ro)
}

// Put issues an HTTP PUT request.
func (c *Client) Put(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("PUT", p, ro)
}

// PutForm issues an HTTP PUT request with the given interface form-encoded.
func (c *Client) PutForm(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm("PUT", p, i, ro)
}

// PutFormFile issues an HTTP PUT request (multipart/form-encoded) to put a file to an endpoint.
func (c *Client) PutFormFile(urlPath string, filePath string, fieldName string, ro *RequestOptions) (*http.Response, error) {
	return c.RequestFormFile("PUT", urlPath, filePath, fieldName, ro)
}

// PutJSON issues an HTTP PUT request with the given interface json-encoded.
func (c *Client) PutJSON(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSON("PUT", p, i, ro)
}

// PutJSONAPI issues an HTTP PUT request with the given interface json-encoded.
func (c *Client) PutJSONAPI(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPI("PUT", p, i, ro)
}

// Delete issues an HTTP DELETE request.
func (c *Client) Delete(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("DELETE", p, ro)
}

// DeleteJSONAPI issues an HTTP DELETE request with the given interface json-encoded.
func (c *Client) DeleteJSONAPI(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPI("DELETE", p, i, ro)
}

// DeleteJSONAPIBulk issues an HTTP DELETE request with the given interface json-encoded and bulk requests.
func (c *Client) DeleteJSONAPIBulk(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPIBulk("DELETE", p, i, ro)
}

// Request makes an HTTP request against the HTTPClient using the given verb,
// Path, and request options.
func (c *Client) Request(verb, p string, ro *RequestOptions) (*http.Response, error) {
	req, err := c.RawRequest(verb, p, ro)
	if err != nil {
		return nil, err
	}

	if ro == nil || !ro.Parallel {
		c.updateLock.Lock()
		defer c.updateLock.Unlock()

	}
	resp, err := checkResp(c.HTTPClient.Do(req))

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RequestForm makes an HTTP request with the given interface being encoded as
// form data.
func (c *Client) RequestForm(verb, p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}

	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = "application/x-www-form-urlencoded"

	v, err := query.Values(i)
	if err != nil {
		return nil, err
	}
	body := v.Encode()

	ro.Body = strings.NewReader(body)
	ro.BodyLength = int64(len(body))

	return c.Request(verb, p, ro)
}

// RequestFormFile makes an HTTP request to upload a file to an endpoint.
func (c *Client) RequestFormFile(verb, urlPath string, filePath string, fieldName string, ro *RequestOptions) (*http.Response, error) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	defer file.Close() // #nosec G307

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("error creating multipart form: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("error copying file to multipart form: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing multipart form: %v", err)
	}

	if ro == nil {
		ro = new(RequestOptions)
	}
	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = writer.FormDataContentType()
	ro.Headers["Accept"] = "application/json"
	ro.Body = &body
	ro.BodyLength = int64(body.Len())

	return c.Request(verb, urlPath, ro)
}

func (c *Client) RequestJSON(verb, p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}

	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = "application/json"
	ro.Headers["Accept"] = "application/json"

	body, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	ro.Body = bytes.NewReader(body)
	ro.BodyLength = int64(len(body))

	return c.Request(verb, p, ro)
}

func (c *Client) RequestJSONAPI(verb, p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}

	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = jsonapi.MediaType
	ro.Headers["Accept"] = jsonapi.MediaType

	if i != nil {
		var buf bytes.Buffer
		if err := jsonapi.MarshalPayload(&buf, i); err != nil {
			return nil, err
		}

		ro.Body = &buf
		ro.BodyLength = int64(buf.Len())
	}
	return c.Request(verb, p, ro)
}

func (c *Client) RequestJSONAPIBulk(verb, p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}

	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = jsonapi.MediaType + "; ext=bulk"
	ro.Headers["Accept"] = jsonapi.MediaType + "; ext=bulk"

	var buf bytes.Buffer
	if err := jsonapi.MarshalPayload(&buf, i); err != nil {
		return nil, err
	}

	ro.Body = &buf
	ro.BodyLength = int64(buf.Len())

	return c.Request(verb, p, ro)
}

// checkResp wraps an HTTP request from the default client and verifies that the
// request was successful. A non-200 request returns an error formatted to
// included any validation problems or otherwise.
func checkResp(resp *http.Response, err error) (*http.Response, error) {
	// If the err is already there, there was an error higher up the chain, so
	// just return that.
	if err != nil {
		return resp, err
	}

	switch resp.StatusCode {
	case 200, 201, 202, 204, 205, 206:
		return resp, nil
	default:
		return resp, NewHTTPError(resp)
	}
}

// decodeBodyMap is used to decode an HTTP response body into a mapstructure struct.
// It closes `body`.
func decodeBodyMap(body io.ReadCloser, out interface{}) error {
	defer body.Close()

	var parsed interface{}
	dec := json.NewDecoder(body)
	if err := dec.Decode(&parsed); err != nil {
		return err
	}

	return decodeMap(parsed, out)
}

// decodeMap decodes an `in` struct or map to a mapstructure tagged `out`.
// It applies the decoder defaults used throughout go-fastly.
// Note that this uses opposite argument order from Go's copy().
func decodeMap(in interface{}, out interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapToHTTPHeaderHookFunc(),
			stringToTimeHookFunc(),
		),
		WeaklyTypedInput: true,
		Result:           out,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(in)
}
