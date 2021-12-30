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

func GetAccessToken(instance, namespace, username, password string) (string, error) {
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
		// handle err
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

func CreateModule(instance, namespace, packageName, accessToken string) error {
	type Req struct {
		Module string `json:"module"`
	}

	reqData, err := json.Marshal(&Req{Module: packageName})
	if err != nil {
		return err
	}

	fmt.Println("r", string(reqData))

	// Check via HTTP API if the desired module is installed
	request, err := http.NewRequest("POST", fmt.Sprintf("http://%s.%s.svc.cluster.local:1881/nodes", instance+"-nodered", namespace), bytes.NewReader(reqData))
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type respData struct {
		Name string `json:"name,omitempty"`
	}

	var respParsed respData
	if err := json.Unmarshal(res, &respParsed); err != nil {
		return err
	}

	fmt.Println(string(res))
	fmt.Println(respParsed.Name)

	return nil
}
