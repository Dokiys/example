package myerrors

import (
	sql "database/sql"
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
	runtime "runtime"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsInvalidParameter(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_InvalidParameter.String() && e.Code == 400
}

func ErrorInvalidParameter(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_InvalidParameter.String(), fmt.Sprintf(format, args...))
}

func ErrorInvalidParameterWithMeta(err error) *errors.Error {
	m := map[string]string{"location": location(), "error": err.Error()}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorInvalidParameter("数据不存在！").WithMetadata(m)
	}
	return ErrorSystemError("系统繁忙！").WithMetadata(m)
}

func IsAccessForbidden(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_AccessForbidden.String() && e.Code == 403
}

func ErrorAccessForbidden(format string, args ...interface{}) *errors.Error {
	return errors.New(403, ErrorReason_AccessForbidden.String(), fmt.Sprintf(format, args...))
}

func ErrorAccessForbiddenWithMeta(err error) *errors.Error {
	m := map[string]string{"location": location(), "error": err.Error()}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorAccessForbidden("数据不存在！").WithMetadata(m)
	}
	return ErrorSystemError("系统繁忙！").WithMetadata(m)
}

func IsUnauthenticated(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_Unauthenticated.String() && e.Code == 401
}

func ErrorUnauthenticated(format string, args ...interface{}) *errors.Error {
	return errors.New(401, ErrorReason_Unauthenticated.String(), fmt.Sprintf(format, args...))
}

func ErrorUnauthenticatedWithMeta(err error) *errors.Error {
	m := map[string]string{"location": location(), "error": err.Error()}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorUnauthenticated("数据不存在！").WithMetadata(m)
	}
	return ErrorSystemError("系统繁忙！").WithMetadata(m)
}

func IsBusinessError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_BusinessError.String() && e.Code == 400
}

func ErrorBusinessError(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_BusinessError.String(), fmt.Sprintf(format, args...))
}

func ErrorBusinessErrorWithMeta(err error) *errors.Error {
	m := map[string]string{"location": location(), "error": err.Error()}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorBusinessError("数据不存在！").WithMetadata(m)
	}
	return ErrorSystemError("系统繁忙！").WithMetadata(m)
}

func IsSystemError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_SystemError.String() && e.Code == 500
}

func ErrorSystemError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, ErrorReason_SystemError.String(), fmt.Sprintf(format, args...))
}

func ErrorSystemErrorWithMeta(err error) *errors.Error {
	m := map[string]string{"location": location(), "error": err.Error()}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorSystemError("数据不存在！").WithMetadata(m)
	}
	return ErrorSystemError("系统繁忙！").WithMetadata(m)
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_NotFound.String() && e.Code == 400
}

func ErrorNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_NotFound.String(), fmt.Sprintf(format, args...))
}

func ErrorNotFoundWithMeta(err error) *errors.Error {
	m := map[string]string{"location": location(), "error": err.Error()}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorNotFound("数据不存在！").WithMetadata(m)
	}
	return ErrorSystemError("系统繁忙！").WithMetadata(m)
}

func IsOrderNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_OrderNotFound.String() && e.Code == 400
}

func ErrorOrderNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_OrderNotFound.String(), fmt.Sprintf(format, args...))
}

func ErrorOrderNotFoundWithMeta(err error) *errors.Error {
	m := map[string]string{"location": location(), "error": err.Error()}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorOrderNotFound("数据不存在！").WithMetadata(m)
	}
	return ErrorSystemError("系统繁忙！").WithMetadata(m)
}

func IsItemNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_ItemNotFound.String() && e.Code == 400
}

func ErrorItemNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(400, ErrorReason_ItemNotFound.String(), fmt.Sprintf(format, args...))
}

func ErrorItemNotFoundWithMeta(err error) *errors.Error {
	m := map[string]string{"location": location(), "error": err.Error()}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrorItemNotFound("数据不存在！").WithMetadata(m)
	}
	return ErrorSystemError("系统繁忙！").WithMetadata(m)
}

func SystemErrorWithMeta(err error) *errors.Error {
	m := map[string]string{"location": location(), "error": err.Error()}
	return errors.New(500, ErrorReason_SystemError.String(), "系统繁忙！").WithMetadata(m)
}

func location() string {
	_, lcErr, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	_, lcCaller, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	var split [2]int
	for i := 0; i < len(lcCaller); i++ {
		if lcCaller[i] == '/' {
			split[0], split[1] = i, split[0]
		}
		if lcCaller[i] != lcErr[i] {
			break
		}
	}

	if split[1]+1 > len(lcCaller) {
		return ""
	}
	return fmt.Sprintf("%s:%d\n", lcCaller[split[1]+1:], line)
}
