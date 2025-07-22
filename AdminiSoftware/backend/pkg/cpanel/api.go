
package cpanel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL  string
	Username string
	Password string
	Token    string
}

type APIResponse struct {
	Result struct {
		Status int                    `json:"status"`
		Data   map[string]interface{} `json:"data"`
		Errors []string               `json:"errors"`
	} `json:"result"`
}

func NewClient(baseURL, username, password string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
	}
}

func (c *Client) NewTokenClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
	}
}

func (c *Client) makeRequest(module, function string, params map[string]string) (*APIResponse, error) {
	values := url.Values{}
	values.Add("cpanel_jsonapi_module", module)
	values.Add("cpanel_jsonapi_func", function)
	values.Add("cpanel_jsonapi_version", "2")
	
	for key, value := range params {
		values.Add(key, value)
	}
	
	req, err := http.NewRequest("POST", c.BaseURL+"/json-api/cpanel", strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	if c.Token != "" {
		req.Header.Set("Authorization", "cpanel "+c.Username+":"+c.Token)
	} else {
		req.SetBasicAuth(c.Username, c.Password)
	}
	
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
	
	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}
	
	return &apiResp, nil
}

func (c *Client) ListDomains() ([]string, error) {
	resp, err := c.makeRequest("DomainLookup", "getbasedomains", nil)
	if err != nil {
		return nil, err
	}
	
	domains := []string{}
	if data, ok := resp.Result.Data["domains"].([]interface{}); ok {
		for _, domain := range data {
			if domainStr, ok := domain.(string); ok {
				domains = append(domains, domainStr)
			}
		}
	}
	
	return domains, nil
}

func (c *Client) CreateEmailAccount(email, password, quota string) error {
	params := map[string]string{
		"email":    email,
		"password": password,
		"quota":    quota,
	}
	
	_, err := c.makeRequest("Email", "addpop", params)
	return err
}

func (c *Client) CreateDatabase(dbname, dbuser, dbpass string) error {
	params := map[string]string{
		"db":   dbname,
		"user": dbuser,
		"pass": dbpass,
	}
	
	_, err := c.makeRequest("Mysql", "adddb", params)
	return err
}
package cpanel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL  string
	Username string
	Password string
	Token    string
}

type APIResponse struct {
	Status   string      `json:"status"`
	Messages []string    `json:"messages"`
	Data     interface{} `json:"data"`
}

func NewClient(baseURL, username, password string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
	}
}

func (c *Client) MakeRequest(endpoint string, params map[string]interface{}) (*APIResponse, error) {
	url := fmt.Sprintf("%s/execute/%s", c.BaseURL, endpoint)
	
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}

func (c *Client) CreateAccount(domain, username, password, packageName string) error {
	params := map[string]interface{}{
		"domain":   domain,
		"username": username,
		"password": password,
		"pkg":      packageName,
	}

	_, err := c.MakeRequest("Accounts/create_user_session", params)
	return err
}

func (c *Client) ListAccounts() ([]interface{}, error) {
	resp, err := c.MakeRequest("Accounts/list_users", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	if accounts, ok := resp.Data.([]interface{}); ok {
		return accounts, nil
	}

	return nil, fmt.Errorf("unexpected response format")
}

func (c *Client) SuspendAccount(username string) error {
	params := map[string]interface{}{
		"user": username,
	}

	_, err := c.MakeRequest("Accounts/suspend_user", params)
	return err
}

func (c *Client) UnsuspendAccount(username string) error {
	params := map[string]interface{}{
		"user": username,
	}

	_, err := c.MakeRequest("Accounts/unsuspend_user", params)
	return err
}

func (c *Client) DeleteAccount(username string) error {
	params := map[string]interface{}{
		"user": username,
	}

	_, err := c.MakeRequest("Accounts/remove_user", params)
	return err
}
