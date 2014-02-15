package xenforo_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func New(apiURL *url.URL, apiKey string, client *http.Client) *API {
	return &API{
		APIURL:     apiURL,
		HttpClient: client,
		LoginHash:  apiKey,
	}
}

type APIError struct {
	ErrorNumber  int    `json:"error,omitempty"`
	ErrorMessage string `json:"message,omitempty"`
}

func (x *APIError) Error() string {
	return fmt.Sprintf("%d: %s", x.ErrorNumber, x.ErrorMessage)
}

type API struct {
	APIURL     *url.URL
	HttpClient *http.Client
	LoginHash  string
}

func (x *API) GetCallURL(action string) *url.URL {
	newAPIURL, err := url.Parse(x.APIURL.String()) // TODO: Find a better way to clone the URL
	if err != nil {                                // This really shouldn't error
		panic(err)
	}

	q := newAPIURL.Query()
	q.Set("action", action)
	if len(x.LoginHash) > 0 {
		q.Set("hash", x.LoginHash)
	} else {
		q.Del("hash")
	}

	newAPIURL.RawQuery = q.Encode()

	return newAPIURL
}

func (x *API) MakeCall(callUrl *url.URL, dst interface{}) error {
	res, err := x.HttpClient.Get(callUrl.String())
	if err != nil {
		return err
	}

	fullBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	baseResponse := new(APIError)
	if err := json.Unmarshal(fullBody, baseResponse); err != nil {
		return err
	}

	if baseResponse.ErrorNumber != 0 {
		return baseResponse
	}

	if err := json.Unmarshal(fullBody, dst); err != nil {
		return err
	}

	return nil
}
