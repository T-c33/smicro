package errno

import (
	"fmt"
)

type SmicroError struct {
	Code    int
	Message string
}

func (s *SmicroError) Error() string {
	return fmt.Sprintf("smicro error, code:%d message:%v", s.Code, s.Message)
}

var (
	NotHaveInstance = &SmicroError{
		Code:    1,
		Message: "not have instance",
	}
	ConnFailed = &SmicroError{
		Code:    2,
		Message: "connect failed",
	}
	InvalidNode = &SmicroError{
		Code:    3,
		Message: "invalid node",
	}
	AllNodeFailed = &SmicroError{
		Code:    4,
		Message: "all node failed",
	}
)

func IsConnectError(err error) bool {

	smicroErr, ok := err.(*SmicroError)
	if !ok {
		return false
	}
	var result bool
	if smicroErr == ConnFailed {
		result = true
	}
	return result
}
