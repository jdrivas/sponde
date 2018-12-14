package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	t "github.com/jdrivas/sponde/term"
	"github.com/juju/ansiterm"
)

// We want to decorate List and Describe with some context dependent
// display of the HTTP response and any errors.
//
// List, Describe are created as methods on objects that are mirror to the
// jupyterhub objects. We want to use method displatch to deal with the
// different kinds of objects and sicne we can't add functions outside of
// a package, we'll create mirror types.
// e.g.   type Group jh.Group
//
// To use these you cast an object from jh to the mirror type, then
// call the function List() or Describe() with the result from the JH function:
//
//   groups, resp, err := jh.GetGroups()
//   List(Groups(groups), resp, err)
//
// The method's are not called directly to alllow decorattion of the output with resp and
// error display dependent  on: error condition, verbose vs. debug etc.
// Also, we want the same display output in those cases where the is no object returned
// from the jh functions (e.g. jh.GetGroups).
//
// Display() is a method that is called when there is only an http.Resoionse and error returned.
// This is typical, for example, on Delete calls.

//
// To do this, List(...) and Describe(...) , Display(...) all call render() which sets up
// a decorotor pipeline as needed.

// There are two basic object display functions: List and Describe.
// Not every data object supports both. Generally, they
// all spport List, and some also support Describe.

// Listable suppots List()
type Listable interface {
	List()
}

// Describable supports Describe()
type Describable interface {
	Describe()
}

// List and Describe display their objects by calling render, but
// first checking that an object is there. If not they send along
// an empty function for the descoroator to call.
func List(d Listable, resp *http.Response, err error) {
	renderer := func() {}
	if d != nil {
		renderer = d.List
	}
	render(renderer, resp, err)
}

// Describe provides detailed output on the object.
func Describe(d Describable, resp *http.Response, err error) {
	renderer := func() {}
	if d != nil {
		renderer = d.Describe
	}
	render(renderer, resp, err)
}

// Display dispolays only the resp and error through the normal pipeline
func Display(resp *http.Response, err error) {
	render(func() {}, resp, err)
}

// DisplayF calls the displayRenderer function as part of the standard render pipeline.
// This is useful for printing out status information bracketed by the usual
// verbose/debugt etc. influenced response and error output from the normal pipeline.
func DisplayF(displayRenderer func(), resp *http.Response, err error) {
	render(displayRenderer, resp, err)

}

// This is for the HTTP direct commands which have jh ojects.
func httpDisplay(resp *http.Response, err error) {
	httpDecorate(errorDecorate(func() {}, err), resp)()
}

func displayServerStartedF(started bool, resp *http.Response, err error) func() {
	return (func() {
		result := t.Success("started")
		if started == false {
			result = t.Success("requested")
			if err != nil {
				result = t.Fail("probably not started")
			}
		}
		fmt.Printf("%s %s\n", t.Title("Server"), result)
	})
}

func displpayServerStopedF(stopped bool, resp *http.Response, err error) func() {
	return (func() {
		result := t.Success("stopped")
		if stopped == false {
			result = t.Success("requested")
			if err != nil {
				result = t.Fail("probably not stopped")
			}
		}
		fmt.Printf("%s %s\n", t.Title("Server"), result)
	})
}

// private API
func render(renderer func(), resp *http.Response, err error) {

	switch {
	case Debug():
		httpDecorate((errorDecorate(renderer, err)), resp)()
	case Verbose():
		shortHTTPDecorate((errorDecorate(renderer, err)), resp)()
	default:

		if err == nil {
			errorDecorate(renderer, err)()
		} else {
			errorHTTPDecorate((errorDecorate(renderer, err)), resp)()
		}
	}
}

// The Decorators are built as pre function call. That is print, the
// call youre argument. So this goes frist to last, with the list.
// Thus httpDecorate(errorDecorate(d.List)) will first print
// the http Response, then the error message, then the List().
func errorDecorate(f func(), err error) func() {
	return (func() {
		if err != nil {
			fmt.Printf(fmt.Sprintf("%s\n", t.Error(err)))
		}

		f()
	})
}

// What to say if there is no response.
func nilResp() {
	fmt.Printf("Nil HTTP Response.\n")
}

// One linear update on the response.
func shortHTTPDecorate(f func(), resp *http.Response) func() {
	return (func() {
		if resp == nil {
			nilResp()
		} else {
			fmt.Printf("%s %s\n", t.Title("HTTP Response: "), httpStatusFunc(resp.StatusCode)("%s", resp.Status))
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil && resp.StatusCode != http.StatusNoContent {
				// m, err := getMessage(body)
				// if err == nil && m.Message != "" {
				// 	fmt.Printf("%s %s\n", t.Title("Message:"), t.Alert(m.Message))
				// }
				prettyPrintBody(body)
			}
		}

		f()
	})
}

func errorHTTPDecorate(f func(), resp *http.Response) func() {
	return (func() {
		if resp == nil {
			nilResp()
		} else {
			fmt.Printf("%s %s\n", t.Title("HTTP Response: "), httpStatusFunc(resp.StatusCode)("%s", resp.Status))
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil && resp.StatusCode != http.StatusNoContent {
				m, err := getMessage(body)
				if err == nil && m.Message != "" {
					fmt.Printf("%s %s\n", t.Title("Message:"), t.Alert(m.Message))
				}
			}
		}

		f()
	})
}

// Message is a simple struct to pull out JSON that is often
// embedded in error returns.
type Message struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Tabled based HTTP reseponse with headers.
func httpDecorate(f func(), resp *http.Response) func() {
	return (func() {
		if resp != nil {
			w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
			fmt.Fprintf(w, "%s\n", t.Title("Status\tLength\tEncoding\tUncompressed"))
			fmt.Fprintf(w, "%s\t%s\n",
				httpStatusFunc(resp.StatusCode)("%s", resp.Status),
				t.Text("%d\t%#v\t%t", resp.ContentLength, resp.TransferEncoding, resp.Uncompressed))
			w.Flush()

			// Headers
			w = ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
			fmt.Fprintf(w, "%s\n", t.Title("Header\tValue"))
			for k, v := range resp.Header {
				fmt.Fprintf(w, "%s\n", t.Text("%s\t%s", k, v))
			}
			w.Flush()

			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil && resp.StatusCode != http.StatusNoContent {
				prettyPrintBody(body)
			} else {
				fmt.Printf("%s %s\n", t.Title("Body Read Error:"), t.Text("%v", err))
			}
		} else {
			nilResp()
		}

		f()

	})
}

func getMessage(jsonString []byte) (m Message, err error) {
	err = json.Unmarshal(jsonString, &m)
	return m, err
}

// Assume it's json and try to pretty print.
func prettyPrintBody(body []byte) {
	prettyJSON := bytes.Buffer{}
	err := json.Indent(&prettyJSON, body, "", "  ")
	if err == nil {

		// Yes, this is totally gratuitous.
		m, err := getMessage(body)
		if err == nil && m.Message != "" {
			fmt.Printf("%s %s\n", t.Title("Message:"), t.Alert(m.Message))
		}

		// Print out the body
		fmt.Printf("%s\n%s\n", t.Title("RESP JSON Body:"), t.Text("%s", string(prettyJSON.Bytes())))
	} else {
		fmt.Printf("%s\n%s\n", t.Title("RESP Body:"), t.Text("%s", string(body)))
		fmt.Printf("%s %s \n", t.Title("JSON Error:"), t.Fail("%v", err))
	}

}

func httpStatusFunc(httpStatus int) (f t.ColorSprintfFunc) {
	switch {
	case httpStatus < 300:
		f = t.Success
	case httpStatus < 400:
		f = t.Warn
	default:
		f = t.Fail
	}
	return f
}

func cmdError(e error) {
	fmt.Printf("Error: %s\n", t.Fail(e.Error()))
}
