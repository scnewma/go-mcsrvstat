package mcsrvstat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "0.0.1"
	defaultBaseURL = "https://api.mcsrvstat.us/"
	userAgent      = "go-mcsrvstat/" + libraryVersion
	mediaType      = "application/json"
	apiVersion     = "2"
)

type Client struct {
	client *http.Client

	BaseURL   *url.URL
	UserAgent string
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}

	return c
}

type ClientOpt interface {
	Apply(*Client) error
}

type ClientOptFunc func(*Client) error

func (cof ClientOptFunc) Apply(c *Client) error {
	return cof(c)
}

type SetBaseURL string

func (bu SetBaseURL) Apply(c *Client) error {
	u, err := url.Parse(string(bu))
	if err != nil {
		return err
	}
	c.BaseURL = u
	return nil
}

type SetUserAgent string

func (ua SetUserAgent) Apply(c *Client) error {
	c.UserAgent = fmt.Sprintf("%s %s", ua, c.UserAgent)
	return nil
}

func (c *Client) newRequest(ctx context.Context, method, urlStr string) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return nil, err
		}
	}

	return resp, err
}

type Status struct {
	Online bool   `json:"online"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	Debug  struct {
		Ping          bool `json:"ping"`
		Query         bool `json:"query"`
		SRV           bool `json:"srv"`
		QueryMismatch bool `json:"querymismatch"`
		IPInSRV       bool `json:"ipinsrv"`
		AnimatedMOTD  bool `json:"animatedmotd"`
		ProxyPipe     bool `json:"proxypipe"`
		CacheTime     int  `json:"cachetime"`
	} `json:"debug"`
	MOTD struct {
		Raw   []string `json:"raw"`
		Clean []string `json:"clean"`
		HTML  []string `json:"html"`
	} `json:"motd"`
	Players struct {
		Online int      `json:"online"`
		Max    int      `json:"max"`
		List   []string `json:"list"`
	} `json:"players"`
	Version  string `json:"version"`
	Protocol int    `json:"protocol"`
	Hostname string `json:"hostname"`
	Icon     string `json:"icon"`
	Software string `json:"software"`
	Map      string `json:"map"`
	Plugins  struct {
		Names []string `json:"names"`
		Raw   []string `json:"raw"`
	} `json:"plugins"`
	Mods struct {
		Names []string `json:"names"`
		Raw   []string `json:"raw"`
	} `json:"mods"`
	Info struct {
		Raw   []string `json:"raw"`
		Clean []string `json:"clean"`
		HTML  []string `json:"html"`
	} `json:"info"`
}

func (c *Client) Status(ctx context.Context, addr string) (Status, *http.Response, error) {
	path := fmt.Sprintf("%s/%s", apiVersion, addr)
	req, err := c.newRequest(ctx, http.MethodGet, path)
	if err != nil {
		return Status{}, nil, err
	}

	var status Status
	resp, err := c.do(ctx, req, &status)
	if err != nil {
		return Status{}, resp, err
	}

	return status, resp, nil
}
