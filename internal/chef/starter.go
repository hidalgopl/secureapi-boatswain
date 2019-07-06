package chef

import (
	"context"
	"crypto/tls"
	"net/http"
)

type TestChan struct {
	Result   string `json:"result"`
	TestCode string `json:"test_code"`
}

type scheduleTestsReq struct {
	UserId string   `json:"user_id"`
	Tests  []string `json:"tests"`
	Url string `json:"url"`
	TestSuiteId string `json:"test_suite_id"`
}
type starter interface {
	start(ctx context.Context, req scheduleTestsReq) (chan TestChan, error)
}

type defaultStarter struct {
}

func newStarter() *defaultStarter {
	service := &defaultStarter{}
	return service
}

func (s *defaultStarter) start(ctx context.Context, req scheduleTestsReq) (chan TestChan, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resultChan := make(chan TestChan)
	r, err := client.Get(req.Url)
	defer r.Body.Close()
	if err != nil {
		panic(err)
		return resultChan, err
	}
	for  _, testCode := range req.Tests{
		go TestCodes[testCode](req.Url, r, resultChan)
	}
	return resultChan, nil
}
