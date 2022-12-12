package westfax

import (
	"encoding/json"
	"fmt"
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
	Id        string
	Direction string
	Date      *time.Time
	FaxFiles  []FaxFile
	Format    string
	PageCount int
	Status    string
	Tag       string
}

type FaxFile struct {
	ContentLength int
	ContentType   string
	FileContents  string
}

type Response struct {
	Success bool
}

type GetFaxDocuments struct {
	Response

	Result []Fax
}

type GetInboundFaxIdentifiers struct {
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

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var securityPingResponse SecurityPingResponse
	err = json.Unmarshal(body, &securityPingResponse)

	return securityPingResponse.Result, err
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

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var faxResponse GetInboundFaxIdentifiers
	err = json.Unmarshal(body, &faxResponse)
	if err != nil {
		return nil, err
	}

	return faxResponse.Result, err
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

	defer resp.Body.Close()

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
