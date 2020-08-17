package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/pushm0v/go-zoho/common"
)

type HttpClient interface {
	WithAuthorizationFunc(f func() string)
	Get(url string, params map[string]interface{}) (resp *http.Response, err error)
	Post(url string, params map[string]interface{}) (resp *http.Response, err error)
	PostJson(url string, params interface{}) (resp *http.Response, err error)
	UploadZIP(url string, params map[string]interface{}, headers map[string]interface{}, handler io.Reader) (resp *http.Response, err error)
	BodyWriter(params map[string]interface{}) (*bytes.Buffer, *multipart.Writer)
}

type httpClient struct {
	client               *http.Client
	useAuthorization     bool
	onAuthorizationToken func() string
}

func NewHttpClient(client *http.Client) HttpClient {
	return &httpClient{
		client:           client,
		useAuthorization: false,
	}
}

func (h *httpClient) Get(url string, params map[string]interface{}) (resp *http.Response, err error) {
	req, err := http.NewRequest(common.HTTP_METHOD_GET, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error when create new HTTP GET request, %v", err)
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, fmt.Sprintf("%v", v))
	}

	req.URL.RawQuery = q.Encode()
	return h.request(req)
}

func (h *httpClient) Post(url string, params map[string]interface{}) (resp *http.Response, err error) {
	reqBody, bodyWriter := h.BodyWriter(params)

	//Dont forget to close writer before NewRequest, see : https://stackoverflow.com/questions/47452046/how-to-compute-content-length-of-multipart-file-request-with-formdata-in-go
	bodyWriter.Close()
	req, err := http.NewRequest(common.HTTP_METHOD_POST, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("Error when create new HTTP POST request, %v", err)
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	return h.request(req)
}

func (h *httpClient) PostJson(url string, params interface{}) (resp *http.Response, err error) {
	reqBody, err := json.Marshal(params)
	if err != nil {
		return
	}

	req, err := http.NewRequest(common.HTTP_METHOD_POST, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("Error when create new HTTP POST Json request, %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return h.request(req)
}

func (h *httpClient) UploadZIP(url string, params map[string]interface{}, headers map[string]interface{}, handler io.Reader) (resp *http.Response, err error) {
	reqBody, bodyWriter := h.BodyWriter(params)
	part, err := bodyWriter.CreateFormFile("file", "file.zip")
	if err != nil {
		return nil, fmt.Errorf("Error when create new HTTP UploadZIP request, %v", err)
	}
	_, err = io.Copy(part, handler)
	if err != nil {
		return nil, fmt.Errorf("Error when create new HTTP UploadZIP request, %v", err)
	}

	//Dont forget to close writer before NewRequest, see : https://stackoverflow.com/questions/47452046/how-to-compute-content-length-of-multipart-file-request-with-formdata-in-go
	bodyWriter.Close()
	req, err := http.NewRequest(common.HTTP_METHOD_POST, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("Error when create new HTTP UploadZIP request, %v", err)
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	return h.request(req)
}

func (h *httpClient) request(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("User-Agent", common.HTTP_USER_AGENT)
	if h.useAuthorization {
		accessToken := h.onAuthorizationToken()
		req.Header.Set("Authorization", fmt.Sprintf("Zoho-oauthtoken %s", accessToken))
	}

	return h.client.Do(req)
}

func (h *httpClient) BodyWriter(params map[string]interface{}) (*bytes.Buffer, *multipart.Writer) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for k, v := range params {
		bodyWriter.WriteField(k, fmt.Sprintf("%v", v))
	}

	return bodyBuf, bodyWriter
}

func (h *httpClient) WithAuthorizationFunc(f func() string) {
	h.useAuthorization = true
	h.onAuthorizationToken = f
}
