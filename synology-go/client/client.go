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

func (c client) Do(r api.Request, response api.Response) error {
	u := c.baseURL()

	// request can override this path by implementing APIPathProvider interface
	u.Path = "/webapi/entry.cgi"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println(u.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
	}()

	synoResponse := api.GenericResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&synoResponse); err != nil {
		return err
	}
	log.Printf("%+v\n", synoResponse)
	// -----
	// if err := mapstructure.Decode(synoResponse.Data, response.Data()); err != nil {
	// }
	if err := mapstructure.Decode(synoResponse.Data, response); err != nil {
	}
	response.SetError(c.handleErrors(response, synoResponse))
	// -----
	// response.DecodePayload(mapstructure.Decode, synoResponse.Data)
	// -----

	return nil
}

func (c client) handleErrors(errorDescriber api.ErrorDescriber, response api.GenericResponse) api.SynologyError {
	err := api.SynologyError{}
	if response.Error.Code == 0 {
		return err
	}

	err.Code = response.Error.Code

	knownErrors := append(errorDescriber.ErrorSummaries(), api.GlobalErrors)
	err.Summary = api.DescribeError(err.Code, knownErrors...)

	for _, e := range response.Error.Errors {
		item := api.ErrorItem{
			Code:    e.Code,
			Summary: api.DescribeError(e.Code, knownErrors...),
			Details: make(api.ErrorFields),
		}
		for k, v := range e.Details {
			item.Details[k] = v
		}
		// drop 'code' from map as it is represented by dedicated field
		delete(item.Details, "code")
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
