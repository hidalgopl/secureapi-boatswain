package checks

import (
	"fmt"
	"github.com/hidalgopl/secureapi-boatswain/internal/messages"
	"github.com/hidalgopl/secureapi-boatswain/internal/status"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/textproto"
	"testing"
)

type MockPublisher struct {
	Result status.TestStatus
}

func (mp *MockPublisher) Publish(msg interface{}, subject string) error {
	m, ok := msg.(*messages.TestFinishedPub)
	if !ok {
		return errors.New("can't cast msg to TestFinishedPub")
	}
	mp.Result = m.Result
	fmt.Printf("added result: %v\n\n\n\n", mp.Result)
	return nil
}

func TestXContentTypeOptionsNoSniff(t *testing.T) {
	tt := []struct {
		testName    string
		headers     map[string][]string
		expectedErr bool
		expectedRes status.TestStatus
	}{
		{
			testName: "happy path",
			headers: http.Header{
				textproto.CanonicalMIMEHeaderKey("X-Content-Type-Options"): {"nosniff"},
			},
			expectedRes: status.Passed,
		},
	}
	mp := &MockPublisher{}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			resultChan := make(chan messages.TestFinishedPub, 1)
			err := XContentTypeOptionsNoSniff("doesnt-matter", tc.headers, resultChan, mp)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRes, mp.Result)
			close(resultChan)
		})
	}
}

func TestXFrameOptionsDeny(t *testing.T) {
	tt := []struct {
		testName    string
		headers     map[string][]string
		expectedErr bool
		expectedRes status.TestStatus
	}{
		{
			testName: "happy path",
			headers: http.Header{
				textproto.CanonicalMIMEHeaderKey("X-Frame-Options"): {"deny"},
			},
			expectedRes: status.Passed,
		},
	}
	mp := &MockPublisher{}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			resultChan := make(chan messages.TestFinishedPub, 1)
			err := XFrameOptionsDeny("doesnt-matter", tc.headers, resultChan, mp)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRes, mp.Result)
			close(resultChan)
		})
	}
}

func TestXXSSProtection(t *testing.T) {
	tt := []struct {
		testName    string
		headers     map[string][]string
		expectedErr bool
		expectedRes status.TestStatus
	}{
		{
			testName: "happy path",
			headers: http.Header{
				textproto.CanonicalMIMEHeaderKey("X-XSS-Protection"): {"1; mode=block"},
			},
			expectedRes: status.Passed,
		},
	}
	for _, tc := range tt {
		mp := &MockPublisher{}
		t.Run(tc.testName, func(t *testing.T) {
			resultChan := make(chan messages.TestFinishedPub, 1)
			err := XXSSProtection("doesnt-matter", tc.headers, resultChan, mp)
			assert.NoError(t, err)
			fmt.Printf("mp.Result: %v\n\n", mp.Result)
			assert.Equal(t, tc.expectedRes, mp.Result)
			close(resultChan)
		})
	}
}

func TestContentSecurityPolicy(t *testing.T) {
	tt := []struct {
		testName    string
		headers     map[string][]string
		expectedErr bool
		expectedRes status.TestStatus
	}{
		{
			testName: "happy path",
			headers: http.Header{
				textproto.CanonicalMIMEHeaderKey("Content-Security-Policy"): {"default-src 'none'"},
			},
			expectedRes: status.Passed,
		},
	}
	mp := &MockPublisher{}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			resultChan := make(chan messages.TestFinishedPub, 1)
			err := ContentSecurityPolicy("doesnt-matter", tc.headers, resultChan, mp)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRes, mp.Result)
			close(resultChan)
		})
	}
}

func TestDetectFingerprintHeaders(t *testing.T) {
	tt := []struct {
		testName    string
		headers     map[string][]string
		expectedErr bool
		expectedRes status.TestStatus
	}{
		{
			testName: "happy path",
			headers: http.Header{
				"x-powered-by": {"flask"},
			},
			expectedRes: status.Failed,
		},
	}
	mp := &MockPublisher{}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			resultChan := make(chan messages.TestFinishedPub, 1)
			err := DetectFingerprintHeaders("doesnt-matter", tc.headers, resultChan, mp)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRes, mp.Result)
			close(resultChan)
		})
	}
}

func TestCORSconfigured(t *testing.T) {
	tt := []struct {
		testName    string
		headers     map[string][]string
		expectedErr bool
		expectedRes status.TestStatus
	}{
		{
			testName: "happy path",
			headers: http.Header{
				"Access-Control-Allow-Origin": {"*"},
			},
			expectedRes: status.Failed,
		},
	}
	mp := &MockPublisher{}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			resultChan := make(chan messages.TestFinishedPub, 1)
			err := CORSconfigured("doesnt-matter", tc.headers, resultChan, mp)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRes, mp.Result)
			close(resultChan)
		})
	}
}
func TestStrictTransportSecurity(t *testing.T) {
	tt := []struct {
		testName    string
		headers     map[string][]string
		expectedErr bool
		expectedRes status.TestStatus
	}{
		{
			testName: "happy path",
			headers: http.Header{
				"Strict-Transport-Security": {"max-age=3600; includeSubDomains"},
			},
			expectedRes: status.Passed,
		},
	}
	mp := &MockPublisher{}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			resultChan := make(chan messages.TestFinishedPub, 1)
			err := StrictTransportSecurity("doesnt-matter", tc.headers, resultChan, mp)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRes, mp.Result)
			close(resultChan)
		})
	}
}
func TestSetCookieSecureHttpOnly(t *testing.T) {
	tt := []struct {
		testName    string
		headers     map[string][]string
		expectedErr bool
		expectedRes status.TestStatus
	}{
		{
			testName: "happy path",
			headers: http.Header{
				"Set-Cookie": {"cookie-without-secureandhttponly"},
			},
			expectedRes: status.Failed,
		},
	}
	mp := &MockPublisher{}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			resultChan := make(chan messages.TestFinishedPub, 1)
			err := SetCookieSecureHttpOnly("doesnt-matter", tc.headers, resultChan, mp)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRes, mp.Result)
			close(resultChan)
		})
	}
}
