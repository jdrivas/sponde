package jupyterhub

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

// var hubAPIURL = "http://athenaeum-staging-hub:8081/hub/api"
var (
	hubAPIURL = "http://127.0.0.1:8081/hub/api"
	hubClient = &http.Client{}
)

func newRequest(method, command string) (*http.Request, error) {
	req, err := http.NewRequest(method, hubAPIURL+command, nil)
	if err == nil {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", viper.GetString("token")))
	}
	return req, err
}

func callJHGet(command string) (resp *http.Response, err error) {
	resp = nil
	req, err := newRequest(http.MethodGet, command)
	if err == nil {
		// fmt.Printf("Making HTTP Request: %#v\n", req)
		resp, err = hubClient.Do(req)
	}
	return resp, err
}
