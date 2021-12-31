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
)

type Client struct {
	httpClient *http.Client
	username   string
	password   string
}

func NewClient(username, password string) *Client {
	return &Client{
		httpClient: &http.Client{},
		username:   username,
		password:   password,
	}
}

func (*Client) getAccessToken(instance, namespace, username, password string) (string, error) {
	params := url.Values{}
	params.Add("client_id", `node-red-admin`)
	params.Add("grant_type", `password`)
	params.Add("scope", `*`)
	params.Add("username", `nodered-operator`)
	params.Add("password", password)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s.%s.svc.cluster.local:1881/auth/token", instance+"-nodered", namespace), body)
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
		AccessToken string `json:"access_token,omitempty"`
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

func (c *Client) CreateModule(instance, namespace, packageName string) (*CreateModuleResponse, error) {
	token, err := c.getAccessToken(instance, namespace, c.username, c.password)
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
	request, err := http.NewRequest("POST", fmt.Sprintf("http://%s.%s.svc.cluster.local:1881/nodes", instance+"-nodered", namespace), bytes.NewReader(reqData))
	if err != nil {
		fmt.Println("FAIL")
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("DO")
		return nil, err
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read")
		return nil, err
	}

	var errorResp ErrorResponse
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		if unmarshalErr := json.Unmarshal(res, &errorResp); unmarshalErr != nil {
			return nil, fmt.Errorf("request failed, status code: %s", resp.Status)
		}

		fmt.Println("NRE")
		return nil, &NodeRedError{StatusCode: resp.StatusCode, Code: errorResp.Code, Message: errorResp.Message}
	}

	fmt.Printf("res: %v\n", string(res))

	var respParsed CreateModuleResponse
	if err := json.Unmarshal(res, &respParsed); err != nil {
		return nil, err
	}

	return &respParsed, nil
}
