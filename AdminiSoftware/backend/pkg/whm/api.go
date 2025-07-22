
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
package whm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL   string
	Username  string
	AccessKey string
}

type WHMResponse struct {
	Status   int                    `json:"status"`
	Messages map[string]interface{} `json:"messages"`
	Data     map[string]interface{} `json:"data"`
	Metadata map[string]interface{} `json:"metadata"`
}

func NewClient(baseURL, username, accessKey string) *Client {
	return &Client{
		BaseURL:   baseURL,
		Username:  username,
		AccessKey: accessKey,
	}
}

func (c *Client) MakeRequest(endpoint string, params map[string]string) (*WHMResponse, error) {
	reqURL := fmt.Sprintf("%s/json-api/%s", c.BaseURL, endpoint)
	
	// Add authentication parameters
	values := url.Values{}
	values.Set("user", c.Username)
	values.Set("accesskey", c.AccessKey)
	
	for key, value := range params {
		values.Set(key, value)
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var whmResp WHMResponse
	if err := json.NewDecoder(resp.Body).Decode(&whmResp); err != nil {
		return nil, err
	}

	return &whmResp, nil
}

func (c *Client) CreateAccount(domain, username, password, packageName string) error {
	params := map[string]string{
		"domain":   domain,
		"username": username,
		"password": password,
		"pkg":      packageName,
	}

	resp, err := c.MakeRequest("createacct", params)
	if err != nil {
		return err
	}

	if resp.Status != 1 {
		return fmt.Errorf("failed to create account: %v", resp.Messages)
	}

	return nil
}

func (c *Client) ListAccounts() (*WHMResponse, error) {
	return c.MakeRequest("listaccts", map[string]string{})
}

func (c *Client) SuspendAccount(username, reason string) error {
	params := map[string]string{
		"user":   username,
		"reason": reason,
	}

	resp, err := c.MakeRequest("suspendacct", params)
	if err != nil {
		return err
	}

	if resp.Status != 1 {
		return fmt.Errorf("failed to suspend account: %v", resp.Messages)
	}

	return nil
}

func (c *Client) UnsuspendAccount(username string) error {
	params := map[string]string{
		"user": username,
	}

	resp, err := c.MakeRequest("unsuspendacct", params)
	if err != nil {
		return err
	}

	if resp.Status != 1 {
		return fmt.Errorf("failed to unsuspend account: %v", resp.Messages)
	}

	return nil
}

func (c *Client) DeleteAccount(username string, keepDNS bool) error {
	params := map[string]string{
		"user": username,
	}
	
	if keepDNS {
		params["keepdns"] = "1"
	}

	resp, err := c.MakeRequest("removeacct", params)
	if err != nil {
		return err
	}

	if resp.Status != 1 {
		return fmt.Errorf("failed to delete account: %v", resp.Messages)
	}

	return nil
}

func (c *Client) GetAccountInfo(username string) (*WHMResponse, error) {
	params := map[string]string{
		"user": username,
	}

	return c.MakeRequest("accountsummary", params)
}

func (c *Client) CreatePackage(name string, limits map[string]interface{}) error {
	params := map[string]string{
		"name": name,
	}

	for key, value := range limits {
		params[key] = fmt.Sprintf("%v", value)
	}

	resp, err := c.MakeRequest("addpkg", params)
	if err != nil {
		return err
	}

	if resp.Status != 1 {
		return fmt.Errorf("failed to create package: %v", resp.Messages)
	}

	return nil
}
