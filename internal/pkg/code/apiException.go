package code

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
	StatusRepeatError     = NewError(http.StatusInternalServerError, 200507, "问卷状态已修改，请勿重复操作！")
	SurveyNumError        = NewError(http.StatusInternalServerError, 200508, "问卷已有填写记录，无法修改！")
	TimeBeyondError       = NewError(http.StatusInternalServerError, 200509, "问卷已过截止日期，无法填写！")
	RegError              = NewError(http.StatusInternalServerError, 200510, "填写内容不符合规范！")
	UniqueError           = NewError(http.StatusInternalServerError, 200511, "填写内容重复！")
	UserExist             = NewError(http.StatusInternalServerError, 200512, "该用户已存在")
	PictureError          = NewError(http.StatusInternalServerError, 200513, "仅允许上传图片文件")
	PictureSizeError      = NewError(http.StatusInternalServerError, 200514, "图片大小超出限制")
	NotSuperAdmin         = NewError(http.StatusInternalServerError, 200513, "很抱歉，您暂无权限注册账号")
	NoPermission          = NewError(http.StatusInternalServerError, 200514, "很抱歉，您暂无权限操作")
	SurveyNotExist        = NewError(http.StatusInternalServerError, 200515, "问卷不存在")
	PermissionExist       = NewError(http.StatusInternalServerError, 200516, "该用户已有权限，请勿重复操作！")
	PermissionBelong      = NewError(http.StatusInternalServerError, 200517, "问卷为该用户所有，无需操作！")
	NotInit               = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	NotFound              = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	Unknown               = NewError(http.StatusInternalServerError, 300500, "系统异常，请稍后重试!")
)


func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}
