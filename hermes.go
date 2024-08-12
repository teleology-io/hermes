package hermes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type RequestMethod string

const (
	GET     RequestMethod = "GET"
	HEAD    RequestMethod = "HEAD"
	POST    RequestMethod = "POST"
	PUT     RequestMethod = "PUT"
	DELETE  RequestMethod = "DELETE"
	CONNECT RequestMethod = "CONNECT"
	OPTIONS RequestMethod = "OPTIONS"
	TRACE   RequestMethod = "TRACE"
	PATCH   RequestMethod = "PATCH"
)

type Headers = map[string]string

type Params = map[string]string

type ClientConfiguration struct {
	BaseURL string
	Headers Headers
	Params  Params
	Timeout int
}

type Client struct {
	baseURL *url.URL
	headers Headers
	params  Params
	client  http.Client
}

type Request struct {
	Method  RequestMethod
	Headers Headers
	Params  Params
	Url     string
	Data    interface{}
}

type Response struct {
	*http.Response
	Data []byte
}

func Create(config ClientConfiguration) Client {
	var baseURL *url.URL
	var headers Headers
	var params Params
	timeout := 0

	if config.BaseURL != "" {
		baseURL, _ = url.Parse(config.BaseURL)
	}

	if config.Timeout != -1 {
		timeout = config.Timeout
	}

	if len(config.Headers) != 0 {
		headers = config.Headers
	}

	if len(config.Params) != 0 {
		params = config.Params
	}

	return Client{
		baseURL: baseURL,
		headers: headers,
		params:  params,
		client: http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func _body(data interface{}) (io.Reader, error) {
	if data == nil {
		return nil, nil
	}

	b, ok := data.(*[]byte)
	if ok {
		return bytes.NewBuffer(*b), nil
	}

	out, err := json.Marshal(data)
	if err == nil {
		return bytes.NewBuffer(out), nil
	}

	return nil, errors.New("Could not convert body to bytes.Buffer")
}

func _headers(req *http.Request, base Headers, given Headers) {
	if req == nil {
		return
	}

	for k, v := range base {
		req.Header.Set(k, v)
	}

	for k, v := range given {
		req.Header.Add(k, v)
	}
}

func _params(base Params, given Params) string {
	params := url.Values{}

	for k, v := range base {
		params.Add(k, v)
	}

	for k, v := range given {
		params.Add(k, v)
	}

	if len(params) != 0 {
		return fmt.Sprintf("?%s", params.Encode())
	}

	return ""
}

func _url(base *url.URL, given string) (string, error) {
	requestUrl := given

	if base == nil && requestUrl == "" {
		return "", errors.New("Missing baseURL or request Url")
	}

	parsedUrl, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
	}

	if base != nil {
		return base.ResolveReference(parsedUrl).String(), nil
	}

	return parsedUrl.String(), nil
}

func _method(request Request) string {
	if request.Method == "" {
		return string(GET)
	}

	return string(request.Method)
}

func (c Client) Send(request Request) (*Response, error) {
	url, err := _url(c.baseURL, request.Url)
	if err != nil {
		return nil, err
	}

	url += _params(c.params, request.Params)

	dataReq, err := _body(request.Data)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(_method(request), url, dataReq)
	if err != nil {
		return nil, err
	}

	_headers(httpReq, c.headers, request.Headers)

	httpRes, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	dataRes, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		httpRes,
		dataRes,
	}, nil
}
