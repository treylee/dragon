package models

type StartTestRequest struct {
	TestID           string        `json:"testID"`
	TestName         string        `json:"testName"`
	RequestPerSecond int           `json:"requestPerSecond"`
	Url              string        `json:"url"`
	Duration         int           `json:"duration"`
}
