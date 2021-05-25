package truenasapi

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func (c *client) get(url string) ([]byte, int, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 500, err
	}

	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, response.StatusCode, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}

	c.Logger.Infof("Request: GET %s %d", url, response.StatusCode)

	defer response.Body.Close()

	return body, response.StatusCode, nil
}

func (c *client) post(url string, payload *bytes.Buffer) ([]byte, int, error) {
	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, 500, err
	}

	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, 500, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}

	c.Logger.Infof("Request: POST %s %d", url, response.StatusCode)

	defer response.Body.Close()

	return body, response.StatusCode, nil
}

func (c *client) delete(url string) ([]byte, int, error) {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, 500, err
	}

	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, response.StatusCode, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}

	c.Logger.Infof("Request: DELETE %s %d", url, response.StatusCode)

	defer response.Body.Close()

	return body, response.StatusCode, nil
}
