package handler

type StartTestRequest struct {
	TestName         string        `json:"testName"`
	RequestPerSecond int           `json:"requestPerSecond"`
	Url              string        `json:"url"`
	Duration         int           `json:"duration"`
}