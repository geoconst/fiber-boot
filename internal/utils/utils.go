package utils

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

// 判断是否是生产环境
func IsProd() bool {
	return os.Getenv("ENV") == "prod" || os.Getenv("ENV") == "pro" ||
		os.Getenv("env") == "prod" || os.Getenv("env") == "pro"
}

func MD5(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

type WebResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

func JsonSuccess(ctx *fiber.Ctx, data any) error {
	return ctx.JSON(WebResp{
		Code: http.StatusOK,
		Msg:  "success",
		Data: data,
	})
}

func JsonFailed(ctx *fiber.Ctx, code int, msg string) error {
	ctx.Status(code)
	return ctx.JSON(WebResp{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// MyValidError 自定义验证错误
type MyValidError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func (ve *MyValidError) Error() string {
	return ve.Msg
}
