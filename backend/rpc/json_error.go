package rpc

import "fmt"

const (
	ErrcodeDefault                  = -42000
	ErrcodeNotificationsUnsupported = -42001
	ErrcodePanic                    = -42603
	ErrcodeMarshalError             = -42603
)

type JSONError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (err *JSONError) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("json-rpc error %d", err.Code)
	}
	return err.Message
}

func (err *JSONError) ErrorCode() int {
	return err.Code
}

func (err *JSONError) ErrorData() interface{} {
	return err.Data
}
