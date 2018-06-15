package main

import (
	"testing"
	"net/http"
	"time"
	"github.com/stretchr/testify/assert"
	"bytes"
	"encoding/json"
)



func TestInitializeDatabase(t *testing.T) {
	go setDomainEnv()
	go initializeSqllite("./sample.db")
	go startServer()
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	r, _ := http.NewRequest("GET", "http://localhost:8080/", nil)

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Body)
}


func TestGetRequestWithInvalidParam(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	r1, _ := http.NewRequest("GET", "http://localhost:8080/aaaa", nil)

	resp, err1 := client.Do(r1)
	if err1 != nil {
		panic(err1)
	}

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.NotNil(t,resp.Body)
}


func TestShrinkURL(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	u := MyURL{LongURL: "http://www.google.com"}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)

	r2, _ := http.NewRequest("POST", "http://localhost:8080/shrink", b)
	r2.Header.Set("Content-Type", "application/json")

	resp2, err2 := client.Do(r2)
	if err2 != nil {
		panic(err2)
	}

	assert.Equal(t, http.StatusOK, resp2.StatusCode)

}

func TestDuplicateShrinkURLRequest(t *testing.T) {

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	u := MyURL{LongURL: "http://www.google.com"}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)

	r1, _ := http.NewRequest("POST", "http://localhost:8080/shrink", b)
	r1.Header.Set("Content-Type", "application/json")

	resp1, err1 := client.Do(r1)
	if err1 != nil {
		panic(err1)
	}

	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	var respURL MyURL
	_ = json.NewDecoder(resp1.Body).Decode(&respURL)

	assert.NotNil(t, respURL.ShortURL)


	u = MyURL{LongURL: "http://www.google.com"}
	b = new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)

	r2, _ := http.NewRequest("POST", "http://localhost:8080/shrink", b)
	r2.Header.Set("Content-Type", "application/json")

	resp2, err2 := client.Do(r2)
	if err2 != nil {
		panic(err2)
	}

	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var respURL2 MyURL
	_ = json.NewDecoder(resp2.Body).Decode(&respURL2)

	assert.NotNil(t, respURL2.ShortURL)
	assert.Equal(t, respURL.LongURL, respURL2.LongURL)
	assert.Equal(t, respURL.ShortURL, respURL2.ShortURL)
}

func TestExpandURL(t *testing.T) {

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	u := MyURL{LongURL: "http://www.google.com"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)

	r, _ := http.NewRequest("POST", "http://localhost:8080/shrink", b)
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respURL MyURL
	_ = json.NewDecoder(resp.Body).Decode(&respURL)

	assert.NotNil(t, respURL.ShortURL)


	r1, _ := http.NewRequest("GET", respURL.ShortURL, nil)

	resp, err1 := client.Do(r1)
	if err1 != nil {
		panic(err1)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t,resp.Body)

}


