package radosAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/smartystreets/go-aws-auth"
)

type API struct {
	host      string
	accessKey string
	secretKey string
}

// New returns client for Ceph RADOS Gateway
func New(host, accessKey, secretKey string) *API {
	return &API{host, accessKey, secretKey}
}

func (api *API) makeRequest(verb, url string) (body []byte, statusCode int, err error) {
	var apiErr apiError
	client := http.Client{}

	// fmt.Println("URL:", url)
	req, err := http.NewRequest(verb, url, nil)
	if err != nil {
		return
	}
	awsauth.SignS3(req, awsauth.Credentials{
		AccessKeyID:     api.accessKey,
		SecretAccessKey: api.secretKey,
		Expiration:      time.Now().Add(1 * time.Minute)},
	)
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	statusCode = resp.StatusCode
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if errMarshal := json.Unmarshal(body, &apiErr); errMarshal == nil && apiErr.Code != "" {
		err = errors.New(apiErr.Code)
	}
	return
}

func (api *API) get(route string, args url.Values, sub ...string) (body []byte, statusCode int, err error) {
	subreq := ""
	if len(sub) > 0 {
		subreq = sub[0] + "&"
	}
	body, statusCode, err = api.makeRequest("GET", fmt.Sprintf("%v%v?%v%s", api.host, route, subreq, args.Encode()))
	if statusCode != 200 {
		err = fmt.Errorf("[%v]: %v", statusCode, err)
	}
	return
}

func (api *API) delete(route string, args url.Values, sub ...string) (body []byte, statusCode int, err error) {
	subreq := ""
	if len(sub) > 0 {
		subreq = sub[0] + "&"
	}
	body, statusCode, err = api.makeRequest("DELETE", fmt.Sprintf("%v%v?%v%s", api.host, route, subreq, args.Encode()))
	if statusCode != 200 {
		err = fmt.Errorf("[%v]: %v", statusCode, err)
	}
	return
}

func (api *API) put(route string, args url.Values, sub ...string) (body []byte, statusCode int, err error) {
	subreq := ""
	if len(sub) > 0 {
		subreq = sub[0] + "&"
	}
	body, statusCode, err = api.makeRequest("PUT", fmt.Sprintf("%v%v?%v%s", api.host, route, subreq, args.Encode()))
	if statusCode != 200 {
		err = fmt.Errorf("[%v]: %v", statusCode, err)
	}
	return
}

func (api *API) post(route string, args url.Values, sub ...string) (body []byte, statusCode int, err error) {
	subreq := ""
	if len(sub) > 0 {
		subreq = sub[0] + "&"
	}
	body, statusCode, err = api.makeRequest("POST", fmt.Sprintf("%v%v?%v%s", api.host, route, subreq, args.Encode()))
	if statusCode != 200 {
		err = fmt.Errorf("[%v]: %v", statusCode, err)
	}
	return
}
