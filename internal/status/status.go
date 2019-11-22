package status

type TestStatus string

const (
	Passed = TestStatus("passed")
	Failed = TestStatus("failed")
	Error = TestStatus("error")
)
