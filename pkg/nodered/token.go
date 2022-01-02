package nodered

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is NOT thread-safe
type Client struct {
	httpClient *http.Client
	host       string
	port       int
	username   string
	password   string

	token          string
	tokenExpiresAt time.Time
}

func NewClient(host string, port int, username, password string) *Client {
	return &Client{
		httpClient: &http.Client{},
		host:       host,
		port:       port,
		username:   username,
		password:   password,
	}
}

func (c *Client) getAccessToken() (string, error) {
	if time.Now().Before(c.tokenExpiresAt) {
		return c.token, nil
	}

	params := url.Values{}
	params.Add("client_id", `node-red-admin`)
	params.Add("grant_type", `password`)
	params.Add("scope", `*`)
	params.Add("username", c.username)
	params.Add("password", c.password)
	body := strings.NewReader(params.Encode())

	tBefore := time.Now()

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/auth/token", c.host, c.port), body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	type tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var respData tokenResp
	err = json.Unmarshal(respBodyData, &respData)
	if err != nil {
		return "", err
	}

	c.token = respData.AccessToken
	c.tokenExpiresAt = tBefore.Add(time.Second * time.Duration(respData.ExpiresIn))
	fmt.Println("Stored token")

	return respData.AccessToken, nil
}

type CreateModuleResponse struct {
	Name string `json:"name,omitempty"`
}

type ErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorCode string

var (
	ErrorCodeUnexpectedError     ErrorCode = "unexpected_error"
	ErrorCodeInvalidRequest      ErrorCode = "invalid_request"
	ErrorCodeSettingsUnavailable ErrorCode = "settings_unavailable"
	ErrorCodeModuleAlreadyLoaded ErrorCode = "module_already_loaded"
	ErrorCodeTypeInUse           ErrorCode = "type_in_use"
	ErrorCodeInvalidAPIVersion   ErrorCode = "invalid_api_version"
)

type NodeRedError struct {
	StatusCode int
	Code       ErrorCode
	Message    string
}

func (e *NodeRedError) Error() string {
	return fmt.Sprintf("http status: %d, code: %s, message: %s", e.StatusCode, e.Code, e.Message)
}

func (c *Client) GetModule() {
	// TODO
}

func (c *Client) CreateModule(packageName string) (*CreateModuleResponse, error) {
	token, err := c.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	type Req struct {
		Module string `json:"module"`
	}

	reqData, err := json.Marshal(&Req{Module: packageName})
	if err != nil {
		return nil, err
	}

	// Check via HTTP API if the desired module is installed
	request, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/nodes", c.host, c.port), bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var errorResp ErrorResponse
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		if unmarshalErr := json.Unmarshal(res, &errorResp); unmarshalErr != nil {
			return nil, fmt.Errorf("request failed, status code: %s", resp.Status)
		}
		return nil, &NodeRedError{StatusCode: resp.StatusCode, Code: errorResp.Code, Message: errorResp.Message}
	}

	var respParsed CreateModuleResponse
	if err := json.Unmarshal(res, &respParsed); err != nil {
		return nil, err
	}

	return &respParsed, nil
}
