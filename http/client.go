package http

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/pushm0v/go-zoho/common"
)

type HttpClient interface {
	WithAuthorization(accessToken string)
	Request(method string, url string, params map[string]interface{}) (resp *http.Response, err error)
	BodyWriter(params map[string]interface{}) (*bytes.Buffer, *multipart.Writer)
}

type httpClient struct {
	client             *http.Client
	authorizationToken string
}

func NewHttpClient(client *http.Client) HttpClient {
	return &httpClient{
		client: client,
	}
}

func (h *httpClient) Request(method string, url string, params map[string]interface{}) (resp *http.Response, err error) {
	reqBody, bodyWriter := h.BodyWriter(params)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("Error when create new HTTP request, %v", err)
	}
	req.Header.Set("User-Agent", common.HTTP_USER_AGENT)
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	if h.authorizationToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Zoho-oauthtoken %s", h.authorizationToken))
	}

	return h.client.Do(req)
}

func (h *httpClient) BodyWriter(params map[string]interface{}) (*bytes.Buffer, *multipart.Writer) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for k, v := range params {
		bodyWriter.WriteField(k, fmt.Sprintf("%v", v))
	}
	bodyWriter.Close()
	return bodyBuf, bodyWriter
}

func (h *httpClient) WithAuthorization(accessToken string) {
	fmt.Println(accessToken)
	h.authorizationToken = accessToken
}
