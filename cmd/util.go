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

func cmdError(e error) {
	fmt.Printf("Error: %s\n", t.Fail(e.Error()))
}

func checkForEmptyString(s string) (r string) {
	r = s
	if s == "" {
		r = "<empty>"
	}
	return r
}

type Listable interface {
	List()
}

type Describable interface {
	Describe()
}

func List(d Listable, resp *http.Response, err error) {
	if d != nil {
		render(d.List, resp, err)
	} else {
		standardHTTPResponse(resp, err)
	}
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
