package jupyterhub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/jdrivas/sponde/config"
	t "github.com/jdrivas/sponde/term"
	"github.com/spf13/viper"
)

var (
	hubClient = &http.Client{}
)

func jhAPIURL(cmd string) string {
	return fmt.Sprintf("%s%s", config.GetHubURL(), cmd)
}

func jhReq(method, cmd string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, jhAPIURL(cmd), body)
}

func newRequest(method, cmd string, body io.Reader) *http.Request {
	req, err := jhReq(method, cmd, body)
	if err == nil {
		if viper.GetBool("debug") {
			fmt.Printf("Using token authorization with token: %s\n", config.GetSafeToken(false, true))
		}
		req.Header.Add("Authorization", fmt.Sprintf("token %s", config.GetToken()))
	} else {
		panic(fmt.Sprintf("Coulnd't generate HTTP request - $s\n", err.Error()))
	}

	return req
}

func sendReq(req *http.Request) (resp *http.Response, err error) {
	resp = nil
	resp, err = hubClient.Do(req)
	if err == nil {
		if viper.GetBool("debug") {
			fmt.Printf("HTTP: %s:%s\n", req.Method, req.URL)
			fmt.Printf("Reponse: %s\n", resp.Status)
		}
		if viper.GetBool("debug") {
			fmt.Printf("%s %s\n", t.Title("Made HTTP Request:"), t.Text("%#v\n", req))
			fmt.Println("")
			fmt.Printf("%s %s\n", t.Title("Response:"), t.Text("%#v\n", *resp))
		}
	}
	return resp, err
}

// getResult makes the get call with the command, and returns the
// response body in the provided object, unmarshalled from the
// JSON in the response object. The returned response will not
// have a body in it.
func getResult(cmd string, result interface{}) (*http.Response, error) {
	resp, err := Get(cmd)
	if err == nil {
		if err = checkReturnCode(*resp); err == nil {
			unmarshal(resp, result)
		}
		if viper.GetBool("debug") {
			fmt.Printf("Unmashaled result: %#v\n", result)
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
	if viper.GetBool("debug") {
		fmt.Printf("Response body is: %s\n", body)
	}
	json.Unmarshal(body, &obj)
	return err
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

//
// PUBLIC API
//

func Send(method, cmd string) (resp *http.Response, err error) {
	var req *http.Request
	req = newRequest(method, cmd, nil)
	resp, err = hubClient.Do(req)
	return resp, err
}

func Get(cmd string) (resp *http.Response, err error) {
	req := newRequest(http.MethodGet, cmd, nil)
	return sendReq(req)
}

// func PostContent(cmd string, content interface{}) (resp *http.Response, err error) {
func PostContent(cmd string, content string) (resp *http.Response, err error) {
	jsonBytes := []byte(content)
	// c := content.(string)
	// ca := []string{c}
	// jsonBytes, err := json.Marshal(ca)
	if viper.GetBool("debug") {
		fmt.Printf("JSON string is: %s\n", jsonBytes)
	}
	// if err == nil {
	// 	fmt.Printf("POST error marshaling JSON %s\n", err)
	// 	fmt.Print("Tried to Marshal object: %#v\n", content)
	// }
	req := newRequest(http.MethodPost, cmd, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	// bodyReader, readErr := req.GetBody()
	// if readErr == nil {
	// 	b, readErr := ioutil.ReadAll(bodyReader)
	// 	if readErr == nil {
	// 		fmt.Printf("Body of REQ looks like: %s\n", b)
	// 	} else {
	// 		fmt.Printf("Some error in reading request body\n")
	// 	}
	// }
	resp, err = sendReq(req)
	return resp, err
}

// If verbose is on, the body is no longer in the response
func Post(cmd string, content interface{}) (resp *http.Response, err error) {
	req := newRequest(http.MethodPost, cmd, nil)
	resp, err = sendReq(req)
	if err == nil {
		if viper.GetBool("debug") {
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

func Delete(cmd string, content interface{}) (resp *http.Response, err error) {
	req := newRequest(http.MethodDelete, cmd, nil)
	resp, err = sendReq(req)
	return resp, err
}
