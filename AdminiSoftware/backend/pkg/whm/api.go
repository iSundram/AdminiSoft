
package whm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL     string
	Username    string
	AccessToken string
}

type WHMResponse struct {
	Status   int                    `json:"status"`
	Metadata map[string]interface{} `json:"metadata"`
	Data     map[string]interface{} `json:"data"`
	Errors   []string               `json:"errors"`
}

func NewClient(baseURL, username, accessToken string) *Client {
	return &Client{
		BaseURL:     baseURL,
		Username:    username,
		AccessToken: accessToken,
	}
}

func (c *Client) makeRequest(function string, params map[string]string) (*WHMResponse, error) {
	values := url.Values{}
	values.Add("api.version", "1")
	
	for key, value := range params {
		values.Add(key, value)
	}
	
	reqURL := fmt.Sprintf("%s/json-api/%s?%s", c.BaseURL, function, values.Encode())
	
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "WHM "+c.Username+":"+c.AccessToken)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var whmResp WHMResponse
	err = json.Unmarshal(body, &whmResp)
	if err != nil {
		return nil, err
	}
	
	return &whmResp, nil
}

func (c *Client) ListAccounts() ([]Account, error) {
	resp, err := c.makeRequest("listaccts", nil)
	if err != nil {
		return nil, err
	}
	
	var accounts []Account
	// Parse accounts from response
	return accounts, nil
}

func (c *Client) CreateAccount(params CreateAccountParams) error {
	reqParams := map[string]string{
		"username": params.Username,
		"domain":   params.Domain,
		"plan":     params.Plan,
		"password": params.Password,
	}
	
	_, err := c.makeRequest("createacct", reqParams)
	return err
}

func (c *Client) SuspendAccount(username string) error {
	params := map[string]string{
		"user": username,
	}
	
	_, err := c.makeRequest("suspendacct", params)
	return err
}

func (c *Client) UnsuspendAccount(username string) error {
	params := map[string]string{
		"user": username,
	}
	
	_, err := c.makeRequest("unsuspendacct", params)
	return err
}

func (c *Client) TerminateAccount(username string) error {
	params := map[string]string{
		"user": username,
	}
	
	_, err := c.makeRequest("removeacct", params)
	return err
}

type Account struct {
	Username    string `json:"username"`
	Domain      string `json:"domain"`
	Plan        string `json:"plan"`
	Suspended   bool   `json:"suspended"`
	DiskUsed    int64  `json:"disk_used"`
	DiskLimit   int64  `json:"disk_limit"`
	Email       string `json:"email"`
	Owner       string `json:"owner"`
}

type CreateAccountParams struct {
	Username string
	Domain   string
	Plan     string
	Password string
}
