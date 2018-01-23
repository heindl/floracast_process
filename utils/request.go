package utils

import (
	"github.com/sethgrid/pester"
	"io/ioutil"
	"bytes"
	"github.com/dropbox/godropbox/errors"
	"encoding/json"
	"encoding/xml"
	"fmt"
)



func RequestJSON(url string, response interface{}) error {

	body, err := request(url)
	if err != nil {
		return err
	}

	if res, ok := response.(*bytes.Reader); ok {
		res.Reset(body)
		return nil
	}

	if err := json.Unmarshal(body, response); err != nil {
		return errors.Wrapf(err, "could not unmarshal http response: %s", url)
	}

	return nil
}

func RequestXML(url string, response interface{}) error {

	body, err := request(url)
	if err != nil {
		return err
	}

	//if res, ok := response.(*bytes.Reader); ok {
	//	res.Reset(body)
	//	return nil
	//}

	if err := xml.Unmarshal(body, response); err != nil {
		return errors.Wrap(err, "could not unmarshal body")
	}

	return nil
}

var totalRequests = 0


func request(url string) ([]byte, error) {
	client := pester.New()
	client.Concurrency = 1
	client.MaxRetries = 5
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true

	totalRequests += 1
	fmt.Println(totalRequests, url)

	resp, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "could not get http response")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Wrapf(errors.New(resp.Status), "StatusCode: %d; URL: %s", resp.StatusCode, url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read http response body")
	}

	return body, nil


}

