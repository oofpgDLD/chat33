package model

import "fmt"

type SendResult struct {
	IsShow     int
	IsValidate int
	Data       map[string]interface{}
}

type Error struct {
	Message string
	Err     string
	Code    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code:%s,error:%s,message:%s", e.Code, e.Err, e.Message)
}
