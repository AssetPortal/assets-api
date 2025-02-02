package model

import (
	"encoding/json"
)

type Response struct {
	OK      bool            `json:"ok"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func NewResponseError(message string) Response {
	return Response{
		OK:      false,
		Message: message,
	}
}

func NewResponseData(data interface{}) Response {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return Response{
			OK:      false,
			Message: "Failed to process data",
		}
	}

	return Response{
		OK:   true,
		Data: dataJSON,
	}
}

func NewResponseEmpty() Response {
	return Response{
		OK: true,
	}
}
