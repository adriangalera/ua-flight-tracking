package request

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type RequestExecutor interface {
	Get(url string, headers map[string]string) string
}

type HttpRequestExecutor struct {
}

func (h HttpRequestExecutor) Get(url string, headers map[string]string) string {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Got error %v while creating the HTTP GET request", err)
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to query the provider with error %v", err)
	}
	defer response.Body.Close()
	body, errReadingBody := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body from provider %v", errReadingBody)
	}
	return string(body)
}

type MockRequestExecutor struct {
	Response       string
	RequestUrl     string
	RequestHeaders map[string]string
}

func (m MockRequestExecutor) Get(url string, headers map[string]string) string {
	m.RequestUrl = url
	m.RequestHeaders = headers
	return m.Response
}
