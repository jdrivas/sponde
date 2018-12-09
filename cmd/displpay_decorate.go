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

// We want to decorate the display of the Hub data we get.
// So for each struct in sponde/jupyterhub/*
// we create a 'shadow' type that we can add display methods to.
// e.g. type  Indo jh.Info
//  Then we can decorate with the functions below.

// There are two basic display functions: List and Describe.
// Not every data object supports both. Generally, they
// all spport List, and some also support Describe.
type Listable interface {
	List()
}

type Describable interface {
	Describe()
}

// We want to decorate listing with some context dpendent
// additional display of the response and any errors.

func List(d Listable, resp *http.Response, err error) {

	lister := func() {}
	if d != nil {
		lister = d.List
	}

	switch {
	case config.Debug():
		httpDecorate((errorDecorate(lister, err)), resp)()
	case config.Verbose():
		shortHTTPDecorate((errorDecorate(lister, err)), resp)()
	default:
		errorDecorate(lister, err)()
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
		}

		f()
	})
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

			// HTTP Body and pretty print JSON.
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				if resp.Header.Get("Content-Type") == "application/json" {
					// if resp.StatusCode < 300 {
					prettyJSON := bytes.Buffer{}
					err := json.Indent(&prettyJSON, body, "", "  ")
					if err == nil {
						fmt.Printf("%s\n%s\n", t.Title("Body:"), t.Text("%s", string(prettyJSON.Bytes())))
					} else {
						fmt.Printf("%s\n", t.Fail("JSON Error: %s", err))
					}
				} else {
					fmt.Printf("%s\n%s\n", t.Title("Body:"), t.Text("%s", string(body)))
				}
			} else {
				fmt.Printf("%s %s\n", t.Title("Body:"), t.Text("Already Read!"))
			}

		} else {
			nilResp()
		}

		f()

	})
}

func Describe(d Describable, resp *http.Response, err error) {
	if d != nil {
		render(d.Describe, resp, err)
	} else {
		standardHTTPResponse(resp, err)
	}
}

// type responseFlags
func render(renderFunc func(), resp *http.Response, err error) {
	if err == nil {
		standardHTTPResponse(resp, err)
		if renderFunc != nil {
			renderFunc()
		}
	} else {
		displayHTTPResponse(resp, true, true, true)
		cmdError(err)
	}
}

func standardHTTPResponse(resp *http.Response, err error) {
	if err == nil {
		if config.Debug() {
			displayHTTPResponse(resp, true, true, true)
		} else if config.Verbose() {
			displayHTTPResponse(resp, true, false, false)
		}
	} else {
		displayHTTPResponse(resp, true, true, true)
		cmdError(err)
	}
}

// This is for the HTTP direct commands.
func doHTTPResponse(resp *http.Response, err error) {
	if err == nil {
		if config.Verbose() {
			displayHTTPResponse(resp, true, true, true)
		} else {
			displayHTTPResponse(resp, true, false, false)
		}
	} else {
		cmdError(err)
	}
}

func displayHTTPResponse(resp *http.Response, status, headers, body bool) {

	fmt.Printf("HTTP %s:%s\n", resp.Request.Method, resp.Request.URL)

	if status {
		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "%s\n", t.Title("Status\tLength\tEncoding\tUncompressed"))
		fmt.Fprintf(w, "%s\t%s\n",
			httpStatusFunc(resp.StatusCode)("%s", resp.Status),
			t.Text("%d\t%#v\t%t", resp.ContentLength, resp.TransferEncoding, resp.Uncompressed))
		w.Flush()
	}

	// HTTP response headers
	if headers {
		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "%s\n", t.Title("Header\tValue"))
		for k, v := range resp.Header {
			fmt.Fprintf(w, "%s\n", t.Text("%s\t%s", k, v))
		}
		w.Flush()
	}

	// HTTP Body and pretty print JSON.
	if body {
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err == nil {
			if resp.Header.Get("Content-Type") == "application/json" {
				// if resp.StatusCode < 300 {
				prettyJson := bytes.Buffer{}
				err := json.Indent(&prettyJson, body, "", "  ")
				if err == nil {
					fmt.Printf("%s\n%s\n", t.Title("Body:"), t.Text("%s", string(prettyJson.Bytes())))
				} else {
					fmt.Printf("%s\n", t.Fail("JSON Error: %s", err))
				}
			} else {
				fmt.Printf("%s\n%s\n", t.Title("Body:"), t.Text("%s", string(body)))
			}
		}
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
