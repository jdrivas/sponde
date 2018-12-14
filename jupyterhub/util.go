package jupyterhub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

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
func (conn Connection) Get(cmd string, result interface{}) (resp *http.Response, err error) {
	return conn.Send(http.MethodGet, cmd, result)
}

// Post works like Get, but uses the POST verb. Post also excepts a content object
// which it will attempt to encode into JSON.
func (conn Connection) Post(cmd string, content, result interface{}) (resp *http.Response, err error) {
	return conn.sendObject(http.MethodPost, cmd, content, result)
}

// Delete works like Post but uses the DELETE verb.
func (conn Connection) Delete(cmd string, content, result interface{}) (resp *http.Response, err error) {
	return conn.sendObject(http.MethodDelete, cmd, content, result)
}

// Patch works like Post, but uses the Delete verb.
func (conn Connection) Patch(cmd string, content, result interface{}) (resp *http.Response, err error) {
	return conn.sendObject(http.MethodPatch, cmd, content, result)
}

// Send works like Get but requires a verb as its first argument.
func (conn Connection) Send(method, cmd string, result interface{}) (resp *http.Response, err error) {
	var req *http.Request
	req = conn.newRequest(method, cmd, nil)
	return sendReq(req, result)
}

// SendJSONString takes a Method, a command and content in the form of a string that is expected
// to be valid JSON. IT returns a JSON result like Get() above.
func (conn Connection) SendJSONString(method, cmd string, content string, result interface{}) (resp *http.Response, err error) {

	if verbose() {
		prettyJSON := bytes.Buffer{}
		err := json.Indent(&prettyJSON, []byte(content), "", "  ")
		if err == nil {
			fmt.Printf("%s\n%s\n", t.Title("Request Content JSON Body:"), t.Text(string(prettyJSON.Bytes())))
		} else {
			fmt.Printf("%s %s \n", t.Title("JSON Error:"), t.Fail("%v", err))
			fmt.Printf("%s\n%s\n", t.Title("Req Content Body:"), t.Text(content))
		}

	}

	buff := bytes.NewBuffer([]byte(content))
	req := conn.newRequest(method, cmd, buff)
	req.Header.Add("Content-Type", "application/json")
	resp, err = sendReq(req, result)
	return resp, err
}

// TODO: Merge the Sends into one.
// They all take an interface to content and result.
// Check type on content, if it's a string, then send it along
// if it's not then marshal
func (conn Connection) sendObject(method, cmd string, content interface{}, result interface{}) (resp *http.Response, err error) {

	//  No content, jsut send.
	if content == nil {
		resp, err = conn.Send(method, cmd, result)
	} else {

		// Otherwise, arshal the object ..,
		var b []byte
		b, err = json.Marshal(content)
		if err == nil {
			if debug() {
				prettyJSON := bytes.Buffer{}
				errI := json.Indent(&prettyJSON, b, "", "  ")

				fmt.Printf("Content to send: %#v\n", content)
				if errI == nil {
					fmt.Printf("%s %s\n", t.Title("Sending JSON:"), t.Text(string(prettyJSON.Bytes())))
				} else {
					fmt.Printf("%s %s\n", t.Title("Sending JSON:"), t.Text(string(b)))
				}
			}

			// ... and send it
			resp, err = conn.SendJSONString(method, cmd, string(b), result)
		}
	}
	return resp, err
}

//
// Private API
//

func (conn Connection) newRequest(method, cmd string, body io.Reader) *http.Request {
	req, err := conn.jhReq(method, cmd, body)
	if err != nil {
		panic(fmt.Sprintf("Coulnd't generate HTTP request - %s\n", err.Error()))
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", conn.Token))

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

		switch {
		case debug():
			fmt.Printf("%s %s\n", t.Title("Made HTTP Request:"), t.Text("%#v", req))
			fmt.Printf("%s %s\n", t.Title("Response:"), t.Text("%#v", *resp))
			fmt.Printf("HTTP: %s:%s\n", req.Method, req.URL)
			fmt.Printf("Reponse: %s\n", resp.Status)
		case verbose():
			fmt.Printf("%s %s\n", t.Title("Made HTTP Request:"), t.Text("%s %s", req.Method, req.URL))
			fmt.Printf("Reponse: %s\n", resp.Status)
		}
	}
	return resp, err
}

func (conn Connection) jhReq(method, cmd string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, conn.jhAPIURL(cmd), body)
}

func (conn Connection) jhAPIURL(cmd string) string {
	return fmt.Sprintf("%s%s", conn.HubURL, cmd)
}

// This eats the body in the response, but returns it in the
//  obj passed in.
func unmarshal(resp *http.Response, obj interface{}) (err error) {
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err == nil {

		if debug() {
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
		if debug() {
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
			err = httpErrorMesg(resp, "Check for valid argument (user, group etc).")
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

// TODO: replace this with proper logging ASAP.
func debug() bool {
	return viper.GetBool("debug")
}

func verbose() bool {
	return viper.GetBool("verbose")
}
