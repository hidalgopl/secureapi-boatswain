package chef

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/hidalgopl/secureapi-boatswain/internal/status"
)

func XContentTypeOptionsNoSniff(url string, r *http.Response, resultChan chan TestChan) error {
	var Status string
	header := r.Header.Get("X-Content-Type-Options")
	if header == "nosniff" {
		Status = "passed"
	} else {
		Status = "failed"
	}
	result := TestChan{
		Result:   Status,
		TestCode: "SEC#0001",
	}
	resultChan <- result
	return nil
}

func XFrameOptionsDeny(url string, r *http.Response, resultChan chan TestChan) error {
	var Status string
	header := r.Header.Get("X-Frame-Options")
	if header == "deny" {
		Status = "passed"
	} else {
		Status = "failed"
	}
	result := TestChan{
		Result:   Status,
		TestCode: "SEC#0002",
	}
	resultChan <- result
	return nil
}

func XXSSProtection(url string, r *http.Response, resultChan chan TestChan) error {
	var Status string
	header := r.Header.Get("X-XSS-Protection")
	if header == "1" || header == "1; mode=block" {
		Status = status.Passed
	} else {
		Status = status.Failed
	}
	result := TestChan{
		Result:   Status,
		TestCode: "SEC#0003",
	}
	resultChan <- result
	return nil
}

func ContentSecurityPolicy(url string, r *http.Response, resultChan chan TestChan) error {
	var Status string
	header := r.Header.Get("Content-Security-Policy")
	if header == "default-src 'none'" {
		Status = status.Passed
	} else {
		Status = status.Failed
	}
	result := TestChan{
		Result:   Status,
		TestCode: "SEC#0004",
	}
	resultChan <- result
	return nil
}

func OptionsRequestNotAllowed(url string, r *http.Response, resultChan chan TestChan) error {
	var Status string
	requestBody, _ := json.Marshal(map[string]string{})
	body := bytes.NewBuffer(requestBody)
	req, err := http.NewRequest(http.MethodOptions, url, body)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Status = status.Error
	} else {
		if resp.StatusCode == http.StatusMethodNotAllowed {
			Status = status.Passed
		} else {
			Status = status.Failed
		}
	}
	result := TestChan{
		Result:   Status,
		TestCode: "SEC#0005",
	}
	resultChan <- result
	return nil

}

var (
	TestCodes = map[string]func(string, *http.Response, chan TestChan) error{
		"SEC#0001": XContentTypeOptionsNoSniff,
		"SEC#0002": XFrameOptionsDeny,
		"SEC#0003": XXSSProtection,
		"SEC#0004": ContentSecurityPolicy,
		"SEC#0005": OptionsRequestNotAllowed,
	}
)
