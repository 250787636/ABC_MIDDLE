package errors

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	New    = errors.New
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)

type M map[string]interface{}

func (m M) String() string {
	var buf strings.Builder
	buf.WriteString("(")
	var index int
	for k, v := range m {
		value, _ := jsoniter.Marshal(v)
		buf.WriteString(fmt.Sprintf("%s = %s", k, value))
		if index++; index != len(m) {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	return buf.String()
}

func NewCaller(funcArg M, msg string) error {
	return errors.New(fmt.Sprintf("%s %s", Caller(2, funcArg), msg))
}
func WrapCaller(err error, funcArg M) error {
	return WithMessage(err, Caller(2, funcArg))
}
func WithCaller(err error, funcArg M, msg string) error {
	return WithMessage(err, fmt.Sprintf("%s %s", Caller(2, funcArg), msg))
}
func WithCallerf(err error, funcArg M, format string, args ...interface{}) error {
	return WithMessagef(err, fmt.Sprintf("%s %s", Caller(2, funcArg), format), args...)
}
func Caller(skip int, funcArg M) string {
	pc, filePath, line, _ := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()
	moduleName := "ky/ssp/"
	if strings.HasPrefix(funcName, moduleName) {
		funcName = funcName[len(moduleName):]
	}
	location := fmt.Sprintf("[ %s:%d ]%s%s", filepath.Base(filePath), line, funcName, funcArg)
	return location
}
