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

//
// Public API
//

// Get performs an HTTP GET on the JupyterHub and returns the results
// aasumed to be JSON encoded into the  result object passed in.
// If you pass in a []map[string]interface{}, you'll get a map of
// objects back.
func Get(cmd string, result interface{}) (resp *http.Response, err error) {
	return Send(http.MethodGet, cmd, result)
}

// Post works like Get, but uses the POST verb. Post also excepts a content object
// which it will attempt to encode into JSON.
func Post(cmd string, content interface{}, result interface{}) (resp *http.Response, err error) {
	return sendObject(http.MethodPost, cmd, content, result)
}

// Delete works like Post but uses the DELETE verb.
func Delete(cmd string, content interface{}, result interface{}) (resp *http.Response, err error) {
	return sendObject(http.MethodDelete, cmd, content, result)
}

func Send(method, cmd string, result interface{}) (resp *http.Response, err error) {
	var req *http.Request
	req = newRequest(method, cmd, nil)
	return sendReq(req, result)
}

func SendJSONString(method, cmd string, content string, result interface{}) (resp *http.Response, err error) {
	buff := bytes.NewBuffer([]byte(content))
	req := newRequest(method, cmd, buff)
	resp, err = sendReq(req, result)
	return resp, err
}

//
// Private API
//

// TODO: Merge the Sends into one.
// They all take an interface to content and result.
// Check type on content, if it's a string, then send it along
// if it's not then marshal
func sendObject(method, cmd string, content interface{}, result interface{}) (resp *http.Response, err error) {
	if content == nil {
		resp, err = Send(method, cmd, result)
	} else {
		var b []byte
		b, err = json.Marshal(content)
		if err == nil {
			if viper.GetBool("debug") {
				prettyJSON := bytes.Buffer{}
				errI := json.Indent(&prettyJSON, b, "", "  ")

				fmt.Printf("Content to send: %#v\n", content)
				if errI == nil {
					fmt.Printf("%s %s\n", t.Title("Sending JSON:"), t.Text(string(prettyJSON.Bytes())))
				} else {
					fmt.Printf("%s %s\n", t.Title("Sending JSON:"), t.Text(string(b)))
				}
			}
			resp, err = SendJSONString(method, cmd, string(b), result)
		}
	}
	// if err == nil {
	// 	if viper.GetBool("debug") {
	// 		if err = checkReturnCode(*resp); err == nil {
	// 			body, err1 := ioutil.ReadAll(resp.Body)
	// 			err = err1
	// 			resp.Body.Close()
	// 			fmt.Printf("Response body: %s\n", body)
	// 		}
	// 	}
	// }
	return resp, err
}

func jhAPIURL(cmd string) string {
	return fmt.Sprintf("%s%s", config.GetHubURL(), cmd)
}

func jhReq(method, cmd string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, jhAPIURL(cmd), body)
}

func newRequest(method, cmd string, body io.Reader) *http.Request {
	req, err := jhReq(method, cmd, body)
	if err != nil {
		panic(fmt.Sprintf("Coulnd't generate HTTP request - $s\n", err.Error()))
	}

	if viper.GetBool("debug") {
		fmt.Printf("Using token authorization with token: %s\n", config.GetSafeToken(false, true))
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", config.GetToken()))

	return req
}

func sendReq(req *http.Request, result interface{}) (resp *http.Response, err error) {
	resp = nil
	resp, err = hubClient.Do(req)

	if err == nil {
		err = checkReturnCode(*resp)
		if result != nil {
			if err == nil {
				err = unmarshal(resp, result)
			}
		}

		if viper.GetBool("debug") {
			fmt.Printf("HTTP: %s:%s\n", req.Method, req.URL)
			fmt.Printf("Reponse: %s\n", resp.Status)
		}
		if viper.GetBool("debug") {
			fmt.Printf("%s %s\n", t.Title("Made HTTP Request:"), t.Text("%#v", req))
			fmt.Println("")
			fmt.Printf("%s %s\n", t.Title("Response:"), t.Text("%#v", *resp))
		}
	}
	return resp, err
}

// This eats the body in the response, but returns it in the
//  obj passed in.
func unmarshal(resp *http.Response, obj interface{}) (err error) {
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err == nil {

		if viper.GetBool("debug") {
			prettyJSON := bytes.Buffer{}
			indentErr := json.Indent(&prettyJSON, body, "", " ")
			if indentErr == nil {
				fmt.Printf("%s %s\n", t.Title("Response body is:"), t.Text("%s\n", prettyJSON))
			} else {
				fmt.Printf("%s\n", t.Fail("Error indenting JSON - %s", indentErr.Error()))
				fmt.Printf("%s %s\n", t.Title("Body:"), t.Text(string(body)))
			}
		}

		json.Unmarshal(body, &obj)
		if viper.GetBool("debug") {
			fmt.Printf("%s %s\n", t.Title("Unmarshaled object: "), t.Text("%#v", obj))
		}
	}
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
	return fmt.Errorf("HTTP Request %s:%s, HTTP Response: %s. %s",
		resp.Request.Method, resp.Request.URL, resp.Status, message)
}

func httpError(resp http.Response) error {
	return httpErrorMesg(resp, "")
}
