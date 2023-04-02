package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"

	"github.com/maksym-nazarenko/terraform-provider-synology/synology-go/client/api"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/publicsuffix"
)

type client struct {
	httpClient *http.Client
	host       string
}

func New(host string, skipCertificateVerification bool) (*client, error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipCertificateVerification,
		},
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{
		Transport: transport,
		Jar:       jar,
	}

	return &client{
		httpClient: httpClient,
		host:       host,
	}, nil
}

func (c *client) Login(user, password, sessionName string) error {
	u := c.baseURL()

	u.Path = "/webapi/entry.cgi"
	q := u.Query()
	q.Add("api", "SYNO.API.Auth")
	q.Add("version", "7")
	q.Add("method", "login")
	q.Add("account", user)
	q.Add("passwd", password)
	q.Add("session", sessionName)
	q.Add("format", "cookie")
	u.RawQuery = q.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
	}()

	return nil
}

func (c client) FileStationInfo() {
	u := c.baseURL()

	u.Path = "/webapi/entry.cgi"
	q := u.Query()
	q.Add("api", "SYNO.FileStation.Info")
	q.Add("version", "2")
	q.Add("method", "get")
	u.RawQuery = q.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		log.Print(err)
		return
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	defer func() {
		_, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
	}()
}

func (c client) Do(r api.Request) (api.Response, error) {
	u := c.baseURL()

	u.Path = r.APIPath()
	q := u.Query()
	q.Add("api", r.APIName())
	q.Add("version", strconv.Itoa(r.APIVersion()))
	q.Add("method", r.APIMethod())
	for k, v := range r.RequestParams() {
		q.Add(k, v)
	}

	u.RawQuery = q.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println(u.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
	}()

	synoResponse := api.GenericResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&synoResponse); err != nil {
		return nil, err
	}
	if !synoResponse.Success {
		synoResponse.Error = c.handleErrors(r, synoResponse)
	}
	log.Printf("%+v\n", synoResponse)

	response := r.NewResponseInstance()
	if err := mapstructure.Decode(synoResponse.Data, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c client) handleErrors(req api.ErrorDescriber, response api.GenericResponse) api.SynologyError {
	err := api.SynologyError{}
	err.Code = response.Error.Code

	knownErrors := append(req.ErrorSummaries(), api.GlobalErrors)
	err.Summary = api.DescribeError(err.Code, knownErrors...)

	for _, e := range response.Error.Errors {
		item := api.ErrorItem{
			Code:    e.Code,
			Summary: api.DescribeError(e.Code, knownErrors...),
			Details: make(api.ErrorFields),
		}
		for k, v := range e.Details {
			if k == "code" {
				continue
			}
			item.Details[k] = v
		}
		err.Errors = append(err.Errors, item)
	}

	return err
}

func (c client) baseURL() url.URL {
	return url.URL{
		Scheme: "https",
		Host:   c.host,
	}
}
