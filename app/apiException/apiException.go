package apiException

import "net/http"

type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

var (
	ServerError           = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	ParamError            = NewError(http.StatusInternalServerError, 200501, "参数错误")
	UserNotFind           = NewError(http.StatusInternalServerError, 200502, "该用户不存在")
	NotLogin              = NewError(http.StatusInternalServerError, 200503, "未登录")
	NoThatPasswordOrWrong = NewError(http.StatusInternalServerError, 200504, "密码错误")
	HttpTimeout           = NewError(http.StatusInternalServerError, 200505, "系统异常，请稍后重试!")
	RequestError          = NewError(http.StatusInternalServerError, 200506, "系统异常，请稍后重试!")
	NotInit               = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	NotFound              = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	Unknown               = NewError(http.StatusInternalServerError, 300500, "系统异常，请稍后重试!")
)

func OtherError(message string) *Error {
	return NewError(http.StatusForbidden, 100403, message)
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}
