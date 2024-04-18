package helper

import (
	"encoding/json"
	"io"
)

type ResponseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func ReadJsonBody[T comparable](body io.ReadCloser) (*T, error) {
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var req T
	err = json.Unmarshal(bodyBytes, &req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}
