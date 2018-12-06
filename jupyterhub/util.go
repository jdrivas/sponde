package jupyterhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jdrivas/sponde/config"
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

func Send(method, command string) (resp *http.Response, err error) {
	resp = nil
	req, err := newRequest(method, command)
	if err != nil {
		panic(fmt.Sprintf("Couldn't generate HTTP request - %s\n", err.Error()))
	}
	resp, err = hubClient.Do(req)
	if err == nil {
		if viper.GetBool("verbose") {
			fmt.Printf("HTTP: %s:%s\n", req.Method, req.URL)
			fmt.Printf("Reponse: %s\n", resp.Status)
		}
		if viper.GetBool("debug") {
			fmt.Printf("Made HTTP Request: %#v\n", req)
			fmt.Printf("Response is: %#v\n", *resp)
		}
	}
	return resp, err
}

// This eats the body in the response, but returns it in the
//  obj passed in.
func unmarshal(resp *http.Response, obj interface{}) (err error) {
	var body []byte
	if err == nil {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	if viper.GetBool("verbose") {
		fmt.Printf("Response body is: %s\n", body)
	}
	json.Unmarshal(body, &obj)
	return err
}

// getResult makes the get call with the command, and returns the
// response body in the provided object, unmarshalled from the
// JSON in the response object. The returned response will not
// have a body in it.
func getResult(cmd string, result interface{}) (*http.Response, error) {
	resp, err := Send(http.MethodGet, cmd)
	if err == nil {
		if err = checkReturnCode(*resp); err == nil {
			unmarshal(resp, result)
		}
		if viper.GetBool("verbose") {
			fmt.Printf("Unmashaled result: %#v\n", result)
		}
	}
	return resp, err
}

func Get(cmd string) (resp *http.Response, err error) {
	return Send(http.MethodGet, cmd)
}

// If verbose is on, the body is no longer in the response
func Post(cmd string) (resp *http.Response, err error) {
	resp, err = Send(http.MethodPost, cmd)
	if err == nil {
		if viper.GetBool("verbose") {
			if err = checkReturnCode(*resp); err == nil {
				body, err1 := ioutil.ReadAll(resp.Body)
				err = err1
				resp.Body.Close()
				fmt.Printf("Response body: %s\n", body)
			}
		}
	}
	return resp, err
}

func Delete(cmd string) (resp *http.Response, err error) {
	resp, err = Send(http.MethodDelete, cmd)
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
