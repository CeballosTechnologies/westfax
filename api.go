package westfax

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const baseUrl = "https://api2.westfax.com/REST"

type Client struct {
	url       url.URL
	username  string
	password  string
	productId string
}

type Fax struct {
	CreatedBy       string
	CreatedVia      string
	DocPageCount    int
	Date            *time.Time
	Direction       string
	FaxCallInfoList []FaxCallInfo
	FaxFiles        []FaxFile
	FaxQuality      string
	FilterValue     string
	Format          string
	Id              string
	JobName         string
	LoginId         string
	PageCount       int
	Status          string
	Tag             string
}

type FaxCallInfo struct {
	CallId        string
	CallPageCount int
	CompletedUtc  time.Time
	FilterFlag    int
	OrigNumber    string
	OrigCSID      string
	Result        string
	TermCSID      string
	TermNumber    string
}

type FaxFile struct {
	ContentLength int
	ContentType   string
	FileContents  string
}

type Response struct {
	Success bool
}

type GetFaxDescriptionResponse struct {
	Response
	Result []Fax
}

type GetFaxDocuments struct {
	Response

	Result []Fax
}

type GetInboundFaxIdentifiersResponse struct {
	Response
	Result []Fax
}

type SecurityPingResponse struct {
	Response
	Result string
}

func New(username string, password string, productId string) *Client {
	client := new(Client)
	client.username = username
	client.password = password
	client.productId = productId

	u, _ := url.Parse(baseUrl)
	client.url = *u

	return client
}

func (c *Client) SecurityPing(ping string) (string, error) {
	form := url.Values{}
	form.Add("StringParams1", ping)

	resp, err := http.PostForm(fmt.Sprintf("%s/Security_Ping/json", &c.url), form)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var securityPingResponse SecurityPingResponse
	err = json.Unmarshal(body, &securityPingResponse)

	return securityPingResponse.Result, err
}

func (c *Client) GetFaxDescription(faxId string) (Fax, error) {
	var fax Fax

	form := url.Values{}
	form.Add("Username", c.username)
	form.Add("Password", c.password)
	form.Add("Cookies", "false")
	form.Add("ProductId", c.productId)
	form.Add("FaxIds1", fmt.Sprintf(`{"Id":"%s"}`, faxId))

	resp, err := http.PostForm(fmt.Sprintf("%s/Fax_GetFaxDescriptionsUsingIds/json", &c.url), form)
	if err != nil {
		return fax, err
	}

	if err != nil {
		return fax, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fax, err
	}

	var faxResponse GetFaxDescriptionResponse
	err = json.Unmarshal(body, &faxResponse)
	if err != nil {
		return fax, err
	}

	return faxResponse.Result[0], err
}

func (c *Client) GetFaxDocument(faxId string) (Fax, error) {
	var fax Fax

	form := url.Values{}
	form.Add("Username", c.username)
	form.Add("Password", c.password)
	form.Add("Cookies", "false")
	form.Add("ProductId", c.productId)
	form.Add("FaxIds1", fmt.Sprintf(`{"Id":"%s"}`, faxId))
	form.Add("Format", "pdf")

	resp, err := http.PostForm(fmt.Sprintf("%s/Fax_GetFaxDocuments/json", &c.url), form)
	if err != nil {
		return fax, err
	}

	if err != nil {
		return fax, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fax, err
	}

	var faxResponse GetFaxDocuments
	err = json.Unmarshal(body, &faxResponse)
	if err != nil {
		return fax, err
	}

	return faxResponse.Result[0], err
}

func (c *Client) GetInboundFaxIdentifiers(startDate string) ([]Fax, error) {
	form := url.Values{}
	form.Add("Username", c.username)
	form.Add("Password", c.password)
	form.Add("Cookies", "false")
	form.Add("ProductId", c.productId)
	form.Add("StartDate", startDate)
	form.Add("FaxDirection", "Inbound")

	resp, err := http.PostForm(fmt.Sprintf("%s/Fax_GetFaxIdentifiers/json", &c.url), form)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var faxResponse GetInboundFaxIdentifiersResponse
	err = json.Unmarshal(body, &faxResponse)
	if err != nil {
		return nil, err
	}

	return faxResponse.Result, err
}
