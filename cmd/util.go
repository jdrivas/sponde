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
	"github.com/spf13/viper"
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

func doHTTPResponse(resp *http.Response, err error) {
	if err == nil {

		fmt.Printf("HTTP %s:%s\n", resp.Request.Method, resp.Request.URL)

		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "%s\n", t.Title("Status\tLength\tEncoding\tUncompressed"))
		fmt.Fprintf(w, "%s\t%s\n",
			httpStatusFunc(resp.StatusCode)("%s", resp.Status),
			t.Text("%d\t%#v\t%t", resp.ContentLength, resp.TransferEncoding, resp.Uncompressed))
		w.Flush()

		if viper.GetBool("verbose") {
			w = ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
			fmt.Fprintf(w, "%s\n", t.Title("Header\tValue"))
			for k, v := range resp.Header {
				fmt.Fprintf(w, "%s\n", t.Text("%s\t%s", k, v))
			}
			w.Flush()

			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				// Sadly, JupyterHub doesn't seem to set the content headers very well.
				// I'm going to assume that if we get a 200 response we'll be getting
				// JSON back, otherwise it might actually be test/html like all
				// the headers from JH say.
				if resp.StatusCode < 300 {
					prettyJson := bytes.Buffer{}
					err := json.Indent(&prettyJson, body, "", "  ")
					if err == nil {
						fmt.Printf("%s\n%s\n", t.Title("Body:"), t.Text("%s", string(prettyJson.Bytes())))
					} else {
						fmt.Printf("%s\n", t.Fail("JSON Error: %s", err))
					}
				}
				if viper.GetBool("debug") {
					fmt.Printf("%s\n%s\n", t.Title("Body:"), t.Text("%s", string(body)))
				}
			}
		}
	}
	if err != nil {
		cmdError(err)
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
