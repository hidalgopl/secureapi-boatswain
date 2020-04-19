package checks

import (
	"fmt"
	"github.com/hidalgopl/secureapi-boatswain/internal/messages"
	"github.com/hidalgopl/secureapi-boatswain/internal/publisher"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"

	"github.com/hidalgopl/secureapi-boatswain/internal/status"
)

var (
	FingerPrintHeaders = []string{
		"x-powered-by",
		"x-generator",
		"server",
		"x-aspnet-version",
		"x-aspnetmvc-version",
	}
)

type TestChan struct {
	Result   status.TestStatus `json:"result"`
	TestCode string            `json:"test_code"`
}

func NotifyCheckFinished(testSuiteID string, testCode string, status status.TestStatus, resultChan chan messages.TestFinishedPub, publisher publisher.Publisher) error {
	msg := &messages.TestFinishedPub{
		TestSuiteID: testSuiteID,
		Result:      status,
		TestCode:    testCode,
		Timestamp:   time.Now(),
	}
	logrus.Debugf("[%s] %s finished with %v, publishing results...", testSuiteID, testCode, status)
	resultChan <- *msg
	logrus.Debug("Pushed info to channel")
	finishedSubject := fmt.Sprintf("test_suite.%s.test.%s.finished", testSuiteID, testCode)
	err := publisher.Publish(msg, finishedSubject)
	if err != nil {
		logrus.Errorf("error during publishing: %v", err)
		return err
	}
	logrus.Infof("published %s to NATS", finishedSubject)
	return nil
}

func XContentTypeOptionsNoSniff(testSuiteID string, headers http.Header, resultChan chan messages.TestFinishedPub, publisher publisher.Publisher) error {
	var Status status.TestStatus
	testCode := "SEC0001"
	header := headers.Get("X-Content-Type-Options")
	if header == "nosniff" {
		Status = status.Passed
	} else {
		Status = status.Failed
	}
	err := NotifyCheckFinished(testSuiteID, testCode, Status, resultChan, publisher)
	if err != nil {
		return err
	}
	return nil
}

func XFrameOptionsDeny(testSuiteID string, headers http.Header, resultChan chan messages.TestFinishedPub, publisher publisher.Publisher) error {
	var Status status.TestStatus
	testCode := "SEC0002"
	header := headers.Get("X-Frame-Options")
	header = strings.ToLower(header)
	if header == "deny" || header == "sameorigin"{
		Status = status.Passed
	} else {
		Status = status.Failed
	}
	err := NotifyCheckFinished(testSuiteID, testCode, Status, resultChan, publisher)
	if err != nil {
		return err
	}
	return nil
}

func XXSSProtection(testSuiteID string, headers http.Header, resultChan chan messages.TestFinishedPub, publisher publisher.Publisher) error {
	var Status status.TestStatus
	testCode := "SEC0003"
	header := headers.Get("X-XSS-Protection")
	if header == "1" || header == "1; mode=block" {
		Status = status.Passed
	} else {
		Status = status.Failed
	}
	err := NotifyCheckFinished(testSuiteID, testCode, Status, resultChan, publisher)
	if err != nil {
		return err
	}
	return nil
}

func ContentSecurityPolicy(testSuiteID string, headers http.Header, resultChan chan messages.TestFinishedPub, publisher publisher.Publisher) error {
	var Status status.TestStatus
	testCode := "SEC0004"
	header := headers.Get("Content-Security-Policy")
	if header != "" {
		Status = status.Passed
	} else {
		Status = status.Failed
	}
	err := NotifyCheckFinished(testSuiteID, testCode, Status, resultChan, publisher)
	if err != nil {
		return err
	}
	return nil
}

func DetectFingerprintHeaders(testSuiteID string, headers http.Header, resultChan chan messages.TestFinishedPub, publisher publisher.Publisher) error {
	var Status status.TestStatus
	testCode := "SEC0005"
	Status = status.Passed
	for _, key := range FingerPrintHeaders {
		if _, ok := headers[key]; ok {
			Status = status.Failed
		}
	}
	err := NotifyCheckFinished(testSuiteID, testCode, Status, resultChan, publisher)
	if err != nil {
		return err
	}
	return nil

}

func CORSconfigured(testSuiteID string, headers http.Header, resultChan chan messages.TestFinishedPub, publisher publisher.Publisher) error {
	var Status status.TestStatus
	testCode := "SEC0006"
	Status = status.Passed
	header := headers.Get("Access-Control-Allow-Origin")
	if header == "*" {
		Status = status.Failed
	}
	err := NotifyCheckFinished(testSuiteID, testCode, Status, resultChan, publisher)
	if err != nil {
		return err
	}
	return nil
}

func StrictTransportSecurity(testSuiteID string, headers http.Header, resultChan chan messages.TestFinishedPub, publisher publisher.Publisher) error {
	var Status status.TestStatus
	testCode := "SEC0007"
	Status = status.Passed
	header := headers.Get("Strict-Transport-Security")
	properlySet := strings.Contains(header, "max-age=")
	if !properlySet {
		Status = status.Failed
	}
	err := NotifyCheckFinished(testSuiteID, testCode, Status, resultChan, publisher)
	if err != nil {
		return err
	}
	return nil
}

func SetCookieSecureHttpOnly(testSuiteID string, headers http.Header, resultChan chan messages.TestFinishedPub, publisher publisher.Publisher) error {
	var Status status.TestStatus
	testCode := "SEC0008"
	Status = status.Passed
	header := headers.Get("Set-Cookie")
	if !strings.Contains(header, "Secure") {
		Status = status.Failed
	}
	if !strings.Contains(header, "HttpOnly") {
		Status = status.Failed
	}
	err := NotifyCheckFinished(testSuiteID, testCode, Status, resultChan, publisher)
	if err != nil {
		return err
	}
	return nil
}

func CacheControlOrExpires(testSuiteID string, headers http.Header, resultChan chan messages.TestFinishedPub, publisher publisher.Publisher) error {
	var Status status.TestStatus
	testCode := "SEC0009"
	Status = status.Passed
	headerCacheControl := headers.Get("Cache-Control")
	headerExpires := headers.Get("Expires")
	if headerCacheControl+headerExpires == "" {
		Status = status.Failed
	}
	err := NotifyCheckFinished(testSuiteID, testCode, Status, resultChan, publisher)
	if err != nil {
		return err
	}
	return nil
}

var (
	TestCodes = map[string]func(string, http.Header, chan messages.TestFinishedPub, publisher.Publisher) error{
		"SEC0001": XContentTypeOptionsNoSniff,
		"SEC0002": XFrameOptionsDeny,
		"SEC0003": XXSSProtection,
		"SEC0004": ContentSecurityPolicy,
		"SEC0005": DetectFingerprintHeaders,
		"SEC0006": CORSconfigured,
		"SEC0007": StrictTransportSecurity,
		"SEC0008": SetCookieSecureHttpOnly,
		"SEC0009": CacheControlOrExpires,
	}
)
