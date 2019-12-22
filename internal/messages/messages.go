package messages

import (
	"fmt"
	"github.com/hidalgopl/secureapi-boatswain/internal/status"
	"strings"
	"time"
)

type StartTestSuitePub struct {
	TestSuiteID string    `json:"test_suite_id"`
	Url         string    `json:"url"`
	Tests       []string  `json:"tests"`
	Timestamp   time.Time `json:"timestamp"`
	UserID      string    `json:"user_id"`
}

func (msg *StartTestSuitePub) Print() string {
	return fmt.Sprintf("[%s] URL: %s {%s}", msg.Timestamp, msg.Url, strings.Join(msg.Tests[:], ","))
}

type TestFinishedPub struct {
	TestSuiteID string            `json:"test_suite_id"`
	Result      status.TestStatus `json:"result"`
	TestCode    string            `json:"test_code"`
	Timestamp   time.Time         `json:"timestamp"`
}

type TestSuiteFinishedPub struct {
	TestSuiteID string            `json:"test_suite_id"`
	Url         string            `json:"url"`
	Tests       []TestFinishedPub `json:"tests"`
	Timestamp   time.Time         `json:"timestamp"`
	UserID      string            `json:"user_id"`
}
