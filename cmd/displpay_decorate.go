package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jdrivas/sponde/config"
	t "github.com/jdrivas/sponde/term"
	"github.com/juju/ansiterm"
)

// We want to decorate List and Describe with some context dependent
// display of the HTTP response and any errors.
//
// List or Describe are created as methods on objects that are mirror to the
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
// The reason we don't call the  method directly o group is that
// we want to decorate the output with resp and error display depending
// on: error condition, verbose vs. debug etc. Also, we want the same display
// output in those cases where the is no object returned from the jh functions (e.g. jh.GetGroups).

//
// To do this, List(...) and Describe(...) both call render() which sets up
// a decorotor pipeline as needed.

// There are two basic display functions: List and Describe.
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

// Describe provides a detailed descript of the object.
func Describe(d Describable, resp *http.Response, err error) {
	renderer := func() {}
	if d != nil {
		renderer = d.Describe
	}
	render(renderer, resp, err)
}

// Display dispoays only the resp and error through the normal pipeline
func Display(resp *http.Response, err error) {
	render(func() {}, resp, err)
}

// This is for the HTTP direct commands which have jh ojects.
func httpDisplay(resp *http.Response, err error) {
	httpDecorate(errorDecorate(func() {}, err), resp)()
}

// private API
func render(renderer func(), resp *http.Response, err error) {

	switch {
	case config.Debug():
		httpDecorate((errorDecorate(renderer, err)), resp)()
	case config.Verbose():
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
			// HTTP Body and pretty print JSON.
			// body, err := ioutil.ReadAll(resp.Body)
			// resp.Body.Close()
			// if err == nil && resp.StatusCode != http.StatusNoContent {
			// 	// We'll just try to print JSON

			// 	prettyJSON := bytes.Buffer{}
			// 	err := json.Indent(&prettyJSON, body, "", "  ")
			// 	if err == nil {

			// 		// Yes, this is totally gratuitous.
			// 		m, err := getMessage(body)
			// 		if err == nil && m.Message != "" {
			// 			fmt.Printf("%s %s\n", t.Title("Message:"), t.Alert(m.Message))
			// 		}

			// 		// Print out the body
			// 		fmt.Printf("%s\n%s\n", t.Title("RESP JSON Body:"), t.Text("%s", string(prettyJSON.Bytes())))
			// 	} else {
			// 		fmt.Printf("%s\n%s\n", t.Title("RESP Body:"), t.Text("%s", string(body)))
			// 		fmt.Printf("%s %s \n", t.Title("JSON Error:"), t.Fail("%v", err))
			// 	}
			// } else {
			// 	if err != nil {
			// 		fmt.Printf("%s %s\n", t.Title("Body Read Error:"), t.Text("%v", err))
			// 	}
			// }

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
