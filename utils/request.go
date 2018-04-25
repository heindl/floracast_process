package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/dropbox/godropbox/errors"
	"github.com/golang/glog"
	"github.com/sethgrid/pester"
	"io"
	"io/ioutil"
	"time"
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

func PostJSON(url string, r io.Reader) (res []byte, err error) {
	client := pester.New()
	client.Concurrency = 1
	client.MaxRetries = 5
	client.Backoff = pester.ExponentialJitterBackoff
	//client.KeepLog = true
	//client.Backoff = func(retry int) time.Duration {
	//	return time.Duration(retry) * time.Second
	//}

	resp, err := client.Post(url, "application/json", r)
	if err != nil {
		return nil, errors.Wrap(err, "could not read http response body")
	}
	defer SafeClose(resp.Body, &err)

	if resp.StatusCode != 200 {
		return nil, errors.Wrapf(errors.New(resp.Status), "StatusCode: %d; URL: %s", resp.StatusCode, url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read http response body")
	}

	return body, nil
}

func request(url string) (res []byte, err error) {

	glog.Infof("Request: [%s]", url)

	client := pester.New()
	client.Concurrency = 1
	client.MaxRetries = 5
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	client.Backoff = func(retry int) time.Duration {
		return time.Duration(retry) * time.Second
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "could not get http response")
	}
	defer SafeClose(resp.Body, &err)

	if resp.StatusCode != 200 {
		return nil, errors.Wrapf(errors.New(resp.Status), "StatusCode: %d; URL: %s", resp.StatusCode, url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read http response body")
	}

	return body, nil

}
