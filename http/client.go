package http

import (
	"bytes"
	"fmt"
	"github.com/pushm0v/go-zoho/common"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type HttpClient interface {
	WithAuthorization(accessToken string)
	Get(url string, params map[string]interface{}) (resp *http.Response, err error)
	Post(url string, params map[string]interface{}) (resp *http.Response, err error)
	UploadZIP(url string, params map[string]interface{}, headers map[string]interface{}, handler io.Reader) (resp *http.Response, err error)
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
	req, err := http.NewRequest(common.HTTP_METHOD_POST, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("Error when create new HTTP POST request, %v", err)
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	return h.request(req)
}

func (h *httpClient) UploadZIP(url string, params map[string]interface{}, headers map[string]interface{}, handler io.Reader) (resp *http.Response, err error) {
	reqBody, bodyWriter := h.BodyWriter(params)

	mimeHeader := textproto.MIMEHeader{}
	mimeHeader.Set("Content-Disposition", `form-data; name="file"; filename="file.zip"`)
	mimeHeader.Set("Content-Type", mime.TypeByExtension(".zip"))
	part, err := bodyWriter.CreatePart(mimeHeader)
	if err != nil {
		return nil, fmt.Errorf("Error when create new HTTP UploadZIP request, %v", err)
	}
	_, err = io.Copy(part, handler)
	if err != nil {
		return nil, fmt.Errorf("Error when create new HTTP UploadZIP request, %v", err)
	}

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
	h.authorizationToken = accessToken
}
