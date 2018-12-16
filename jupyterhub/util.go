package jupyterhub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	// "strings"

	t "github.com/jdrivas/sponde/term"
	"github.com/spf13/viper"
)

var (
	// Note: I don't expect that there are any dependencies that
	// should cause this to be a problem.
	hubClient = http.DefaultClient
)

//
// Public API
//

// Send performs an HTTP request on the URL with the Token in the connection, using
// the HTTP provided by method.
// If content is non-nil, it's marshalled into the body  as a json string.
// If content is a string, it's written directly into the booy (this string is not validated as correct JSON)
// If result is non-nil Send umarshalls the response body,
// aasumed to be JSON encoded, into the result object passed in.
// If result is a []map[string]interface{}, you'll get a map of the JSON object.
func (conn Connection) Send(method, cmd string, content interface{}, result interface{}) (resp *http.Response, err error) {

	//  No content, jsut send.
	if content == nil {
		req := conn.newRequest(method, cmd, nil)
		resp, err = sendReq(req, result)
	} else {
		// Otherwise, marshal the object and send the request.
		var b []byte
		switch c := content.(type) {
		case string:
			// If we use unmarshall on the string, it escapges the quotes: "foo" => \"foo\".
			b = []byte(c)
		default:
			b, err = json.Marshal(c)
		}
		if err == nil {
			buff := bytes.NewBuffer(b)
			req := conn.newRequest(method, cmd, buff)
			req.Header.Add("Content-Type", "application/json")
			resp, err = sendReq(req, result)
		}
	}
	return resp, err
}

// Get works like Send with the GET verb,  but doesn't require a content object.
func (conn Connection) Get(cmd string, result interface{}) (resp *http.Response, err error) {
	return conn.Send(http.MethodGet, cmd, nil, result)
}

// Post works like Send using the POST verb.
func (conn Connection) Post(cmd string, content, result interface{}) (resp *http.Response, err error) {
	return conn.Send(http.MethodPost, cmd, content, result)
}

// Delete works like Send using the Delte verb.
func (conn Connection) Delete(cmd string, content, result interface{}) (resp *http.Response, err error) {
	return conn.Send(http.MethodDelete, cmd, content, result)
}

// Patch works like Send using the Patch verb.
func (conn Connection) Patch(cmd string, content, result interface{}) (resp *http.Response, err error) {
	return conn.Send(http.MethodPatch, cmd, content, result)
}

//
// Private API
//

// sendReq sends along the request with some logging along the way.
func sendReq(req *http.Request, result interface{}) (resp *http.Response, err error) {

	switch {
	// TODO: This wil dump the authorization token. Which it probably shouldn't.
	case debug():
		reqDump, dumpErr := httputil.DumpRequestOut(req, true)
		reqStr := string(reqDump)
		if dumpErr != nil {
			fmt.Printf("Error dumping request (display as generic object): %v\n", dumpErr)
			reqStr = fmt.Sprintf("%v", req)
		}
		fmt.Printf("%s %s\n", t.Title("Request"), t.Text(reqStr))
		fmt.Println()
	case verbose():
		fmt.Printf("%s %s\n", t.Title("Request:"), t.Text("%s %s", req.Method, req.URL))
		fmt.Println()
	}

	resp, err = hubClient.Do(req)
	if err == nil {

		if debug() {
			respDump, dumpErr := httputil.DumpResponse(resp, true)
			respStr := string(respDump)
			if dumpErr != nil {
				fmt.Printf("Error dumping response (display as generic object): %v\n", dumpErr)
				respStr = fmt.Sprintf("%v", resp)
			}
			fmt.Printf("%s\n%s\n", t.Title("Respose:"), t.Text(respStr))
			fmt.Println()
		}

		// Do this after the Dump, the dump reads out the response for reprting and
		// replaces the reader with anotherone that has the data.
		err = checkReturnCode(*resp)
		if result != nil {
			if err == nil {
				err = unmarshal(resp, result)
			}
		}

	}
	return resp, err
}

// newRequest creates a request as usual prepending the connections HubURL to the cmd,
// and adding the Authorization header using token.
func (conn Connection) newRequest(method, cmd string, body io.Reader) *http.Request {
	// req, err := conn.jhReq(method, cmd, body)
	req, err := http.NewRequest(method, conn.HubURL+cmd, body)
	if err != nil {
		panic(fmt.Sprintf("Coulnd't generate HTTP request - %s\n", err.Error()))
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", conn.Token))

	return req
}

// This eats the body in the response, but returns the body in
//  obj passed in. They must match of course.
func unmarshal(resp *http.Response, obj interface{}) (err error) {
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err == nil {

		if debug() {
			// TODO: Printf-ing the output of json.Indent through the bytes.Buffer.String
			// produces cruft. However writting directly to it, works o.k.
			// prettyJSON := bytes.Buffer{}
			var prettyJSON bytes.Buffer
			fmt.Fprintf(&prettyJSON, t.Title("Pretty print response body:\n"))
			indentErr := json.Indent(&prettyJSON, body, "", " ")
			if indentErr == nil {
				// fmt.Printf("%s %s\n", t.Title("Response body is:"), t.Text("%s\n", prettyJSON))
				prettyJSON.WriteTo(os.Stdout)
				fmt.Println()
				fmt.Println()
			} else {
				fmt.Printf("%s\n", t.Fail("Error indenting JSON - %s", indentErr.Error()))
				fmt.Printf("%s %s\n", t.Title("Body:"), t.Text(string(body)))
			}
		}

		json.Unmarshal(body, &obj)
		if debug() {
			fmt.Printf("%s %s\n", t.Title("Unmarshaled object: "), t.Text("%#v", obj))
			fmt.Println()
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
