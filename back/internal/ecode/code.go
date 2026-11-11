// Package ecode 定义 Business 与 Order 服务共用的四位业务状态码。
package ecode

import "net/http"

type Code int32

const (
	Success Code = 1000

	InvalidArgument  Code = 1001
	NotFound         Code = 1002
	MethodNotAllowed Code = 1003
	NotImplemented   Code = 1004

	AuthRequired           Code = 1100
	InvalidToken           Code = 1101
	Forbidden              Code = 1102
	AccountDisabled        Code = 1103
	CoachNotVerified       Code = 1104
	IdempotencyKeyRequired Code = 1105
	InvalidCredentials     Code = 1106

	InvalidWechatCode    Code = 1200
	AuthNotConfigured    Code = 1201
	ChannelNotConfigured Code = 1202
	AdminAccountExists   Code = 1300
	AdminRoleInvalid     Code = 1301

	InternalError Code = 5000
)

type Definition struct {
	HTTPStatus int
	Message    string
}

var definitions = map[Code]Definition{
	Success:                {http.StatusOK, ""},
	InvalidArgument:        {http.StatusBadRequest, "请求参数无效"},
	NotFound:               {http.StatusNotFound, "资源不存在"},
	MethodNotAllowed:       {http.StatusMethodNotAllowed, "请求方法不支持"},
	NotImplemented:         {http.StatusNotImplemented, "接口业务逻辑尚未实现"},
	AuthRequired:           {http.StatusUnauthorized, "请先登录"},
	InvalidToken:           {http.StatusUnauthorized, "登录态无效或已过期"},
	Forbidden:              {http.StatusForbidden, "没有访问权限"},
	AccountDisabled:        {http.StatusForbidden, "账号当前不可登录"},
	CoachNotVerified:       {http.StatusForbidden, "当前身份不具备陪练角色"},
	IdempotencyKeyRequired: {http.StatusBadRequest, "请提供有效的Idempotency-Key"},
	InvalidWechatCode:      {http.StatusUnauthorized, "微信登录凭据无效或已使用"},
	AuthNotConfigured:      {http.StatusServiceUnavailable, "登录服务暂时不可用"},
	ChannelNotConfigured:   {http.StatusServiceUnavailable, "渠道服务暂时不可用"},
	InternalError:          {http.StatusInternalServerError, "系统内部错误"},
	InvalidCredentials:     {http.StatusUnauthorized, "账号或密码错误"},
	AdminAccountExists:     {http.StatusConflict, "管理员账号已存在"},
	AdminRoleInvalid:       {http.StatusBadRequest, "角色不存在或未启用"},
}

func Lookup(code Code) (Definition, bool) {
	definition, ok := definitions[code]
	return definition, ok
}

func Message(code Code) string {
	if definition, ok := Lookup(code); ok {
		return definition.Message
	}
	return definitions[InternalError].Message
}

func HTTPStatus(code Code) int {
	if definition, ok := Lookup(code); ok {
		return definition.HTTPStatus
	}
	return http.StatusInternalServerError
}

func FromHTTPStatus(status int) Code {
	switch status {
	case http.StatusBadRequest:
		return InvalidArgument
	case http.StatusUnauthorized:
		return InvalidToken
	case http.StatusForbidden:
		return Forbidden
	case http.StatusNotFound:
		return NotFound
	case http.StatusMethodNotAllowed:
		return MethodNotAllowed
	case http.StatusNotImplemented:
		return NotImplemented
	default:
		return InternalError
	}
}
