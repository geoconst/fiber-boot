package app

import (
	"errors"
	"fiber-boot/internal/utils"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mdobak/go-xerrors"
)

// 全局错误处理函数
func errorHandler(c *fiber.Ctx, err error) error {
	var verr *utils.MyValidError

	// 如果是校验错误，返回400提示
	if errors.As(err, &verr) {
		return utils.JsonFailed(c, http.StatusBadRequest, err.Error())
	}

	if _, ok := err.(xerrors.StackTracer); ok {
		slog.Error(fmt.Sprintf("%s\n%s", err.Error(), xerrors.StackTrace(err)))
	} else {
		slog.Error(err.Error())
	}
	return utils.JsonFailed(c, http.StatusInternalServerError, err.Error())
}
