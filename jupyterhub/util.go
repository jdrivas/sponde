package jupyterhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

// var hubAPIURL = "http://athenaeum-staging-hub:8081/hub/api"
var (
	// hubAPIURL = "http://127.0.0.1:8081/hub/api"
	hubClient = &http.Client{}
)

func newRequest(method, command string) (*http.Request, error) {
	hubURL := viper.GetString("hubURL")
	req, err := http.NewRequest(method, hubURL+command, nil)
	if err == nil {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", viper.GetString("token")))
	}
	return req, err
}

func callJHGet(command string) (resp *http.Response, err error) {
	resp = nil
	req, err := newRequest(http.MethodGet, command)
	resp, err = hubClient.Do(req)
	if err == nil {
		if viper.GetBool("debug") {
			fmt.Printf("Made HTTP Request: %#v\n", req)
			fmt.Printf("Response is: %#v\n", *resp)
		}
	}
	return resp, err
}

func unmarshal(resp *http.Response, obj interface{}) (err error) {
	body := []byte{}
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

// Get makes the get call with the command, and returns the
// results in the provided object, unmarshalled from the
// JSON in the response.
func Get(cmd string, result interface{}) error {
	resp, err := callJHGet(cmd)
	if err == nil && resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("Response Status: %d - Item Not Found", http.StatusNotFound)
	}
	if err == nil {
		unmarshal(resp, result)
	}
	return err
}
