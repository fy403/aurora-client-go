package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type AuroraConnector struct {
	client      *http.Client
	authRequest *AuthRequest
	cookies     []*http.Cookie
	interval    time.Duration
	loginUrl    string
	tasksUrl    string
	touchUrl    string
	connError   chan error
}

// NewAuroraConnector create a client instance
func NewAuroraConnector(loginUrl, tasksUrl, touchUrl string) *AuroraConnector {
	defaultTransport := &http.Transport{
		DialContext: (&net.Dialer{
			// Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost: 5,
		MaxIdleConns:        0,
		// IdleConnTimeout:     10 * time.Second,
	}
	return &AuroraConnector{
		client: &http.Client{
			Timeout:   time.Second * time.Duration(10),
			Transport: defaultTransport,
		},
		authRequest: &AuthRequest{},
		loginUrl:    loginUrl,
		tasksUrl:    tasksUrl,
		touchUrl:    touchUrl,
		interval:    10 * time.Millisecond,
	}
}

// Init set some personal information for login
func (conn *AuroraConnector) Init(userName, password string) error {
	conn.authRequest.Name = userName
	conn.authRequest.Password = password
	return nil
}

// Login send a logit post http request to Aurora
func (conn *AuroraConnector) login() error {
	var err error
	requestOBJ, err := json.Marshal(conn.authRequest)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(requestOBJ)
	request, err := http.NewRequest("POST", conn.loginUrl, bodyReader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := conn.client.Do(request)
	if err != nil {
		return err
	}
	if response == nil {
		return errors.New("Response is nil")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("Login failed error: %s", string(body))
	}
	conn.cookies = response.Cookies()
	return nil
}

func (conn *AuroraConnector) SendAsync(requestOBJ *CenterRequest) <-chan *CenterResponse {
	channel := make(chan *CenterResponse, 10)
	go func() {
		responseOBJ, err := conn.SendSync(requestOBJ)
		if err != nil {
			close(channel)
			return
		}
		channel <- responseOBJ
	}()
	return channel
}

// SendSync send a sync http request to Aurora
func (conn *AuroraConnector) SendSync(requestOBJ *CenterRequest) (*CenterResponse, error) {
	var err error
	req, err := json.Marshal(requestOBJ)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(req)
	request, err := http.NewRequest("POST", conn.tasksUrl, bodyReader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cache-control", " no-cache")
	for _, c := range conn.cookies {
		request.AddCookie(c)
	}
	response, err := conn.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("Response is nil")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	switch response.StatusCode {
	case http.StatusOK:
		responseOBJ := &CenterResponse{}
		err = json.Unmarshal(body, responseOBJ)
		if err != nil {
			return nil, err
		}
		return responseOBJ, nil
	case http.StatusPartialContent:
		responseOBJ := &CenterResponse{}
		err = json.Unmarshal(body, responseOBJ)
		if err != nil {
			return nil, err
		}
		var signatuires []*Signature
		var callback *Signature
		if responseOBJ.TaskType == "group" {
			for _, taskResponse := range responseOBJ.TaskResponses {
				signatuires = append(signatuires, taskResponse.Signatures...)
			}
		} else {
			signatuires = responseOBJ.TaskResponses[0].Signatures
		}

		if responseOBJ.TaskType == "chord" {
			callback = responseOBJ.TaskResponses[0].CallBack
		} else {
			callback = nil
		}
		req := &CenterRequest{
			BatchID:         responseOBJ.BatchID,
			Timestamp:       time.Now().Local().Unix(),
			TaskType:        responseOBJ.TaskType,
			Signatures:      signatuires,
			TimeoutDuration: 800, // invalid, just for validate
			SleepDuration:   50,  // invalid, just for validate
			SendConcurrency: 5,   // invalid, just for validate
			CallBack:        callback,
		}
		return conn.sendTouchWithTimeout(req, conn.interval)
	case http.StatusForbidden:
		err := conn.login()
		if err != nil {
			return nil, err
		}
		return conn.SendSync(requestOBJ)
	default:
		return nil, fmt.Errorf("Send fail: %s", string(body))
	}
}

func (conn *AuroraConnector) sendTouchWithTimeout(requestOBJ *CenterRequest, sleepDuration time.Duration) (*CenterResponse, error) {
	for {
		results, err := conn.sendTouch(requestOBJ)
		if results == nil && err == nil {
			time.Sleep(sleepDuration)
		} else {
			return results, err
		}
	}
}

// sendTouch will send immediately without time duration, if return nil, nil, means the task is not finished yet
func (conn *AuroraConnector) sendTouch(requestOBJ *CenterRequest) (*CenterResponse, error) {
	var err error
	req, err := json.Marshal(requestOBJ)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(req)
	request, err := http.NewRequest("POST", conn.touchUrl, bodyReader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cache-control", " no-cache")
	for _, c := range conn.cookies {
		request.AddCookie(c)
	}
	response, err := conn.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("Response is nil")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	switch v := response.StatusCode; v {
	case http.StatusBadRequest:
		return nil, fmt.Errorf("Bad request: %s", string(body))
	case http.StatusBadGateway:
		return nil, fmt.Errorf("Task has failed: %s", string(body))
	case http.StatusOK:
		break
	default:
		return nil, fmt.Errorf("Unknow StatusCode: %d, %s", v, string(body))
	}

	responseOBJ := &CenterResponse{}
	err = json.Unmarshal(body, responseOBJ)
	if err != nil {
		return nil, err
	}
	// if results is empty, Your choice
	for _, taskResponse := range responseOBJ.TaskResponses {
		if len(taskResponse.Results) == 0 {
			return nil, nil
		}
	}
	return responseOBJ, nil
}

func (conn *AuroraConnector) Close() error {
	conn.client.CloseIdleConnections()
	return nil
}
