package jupyterhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jdrivas/jhmon/config"
	"github.com/spf13/viper"
)

var (
	hubClient = &http.Client{}
)

func newRequest(method, command string) (*http.Request, error) {
	hubURL := config.GetHubURL()
	token := config.GetToken()
	req, err := http.NewRequest(method, hubURL+command, nil)

	if err == nil {
		if viper.GetBool("debug") {
			fmt.Printf("Using token authorization with token: %s\n", token)
		}
		req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	}
	return req, err
}

func callJHGet(command string) (resp *http.Response, err error) {
	resp = nil
	req, err := newRequest(http.MethodGet, command)
	resp, err = hubClient.Do(req)
	if err == nil {
		if viper.GetBool("verbose") {
			fmt.Printf("Made HTTP Request: %#v\n", req)
			fmt.Printf("Response is: %#v\n", *resp)
		}
	}
	return resp, err
}

func unmarshal(resp *http.Response, obj interface{}) (err error) {
	body := []byte{}
	if err == nil {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	if viper.GetBool("debug") {
		fmt.Printf("Response body is: %s\n", body)
	}
	json.Unmarshal(body, &obj)
	return err
}

// Get makes the get call with the command, and returns the
// results in the provided object, unmarshalled from the
// JSON in the response.
func get(cmd string, result interface{}) (*http.Response, error) {
	resp, err := callJHGet(cmd)
	if err == nil {
		if err = checkReturnCode(*resp); err == nil {
			unmarshal(resp, result)
		}
	}
	return resp, err
}

// Returns an "informative" error if not 200
func checkReturnCode(resp http.Response) (err error) {
	err = nil
	if resp.StatusCode >= 300 {
		switch resp.StatusCode {
		case http.StatusNotFound:
			err = httpErrorMesg(resp, "Check for misbehaving connection, or missing token.")
		case http.StatusUnauthorized:
			err = httpErrorMesg(resp, "Check for valid token.")
		case http.StatusForbidden:
			err = httpErrorMesg(resp, "Check for valid token and token user must be an admin")
		default:
			err = httpError(resp)
		}
	}
	return err
}

func httpErrorMesg(resp http.Response, message string) error {
	return fmt.Errorf("HTTP Request %s:%s, HTTP Response -> %s. %s",
		resp.Request.Method, resp.Request.URL, resp.Status, message)
}

func httpError(resp http.Response) error {
	return httpErrorMesg(resp, "")
}
