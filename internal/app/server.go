package app

import (
	"errors"
	"fiber-boot/internal/dao"
	"fiber-boot/internal/module/account"
	"fiber-boot/internal/utils"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mdobak/go-xerrors"
	slogfiber "github.com/samber/slog-fiber"
	"gorm.io/gorm"
)

type Server struct {
	config *Config
	fiber  *fiber.App
}

func (s *Server) Start() error {
	return s.fiber.Listen(fmt.Sprintf(":%d", s.config.Port))
}

var tokenDao *dao.AccessTokenDAO

func NewServer(
	config *Config,
	accountHandler *account.AccountHandler,
	accessTokenDao *dao.AccessTokenDAO) *Server {

	tokenDao = accessTokenDao
	logConfig := slogfiber.Config{
		DefaultLevel:      slog.LevelDebug,
		ClientErrorLevel:  slog.LevelWarn,
		ServerErrorLevel:  slog.LevelError,
		WithRequestHeader: true,
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 30,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	app.Use(slogfiber.NewWithConfig(slog.Default(), logConfig))
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	// app.Use(signMiddleware)

	account.RegisterPublicRoutes(app, accountHandler)

	app.Use(authMiddleware)
	account.RegisterPrivateRoutes(app, accountHandler)

	return &Server{
		config: config,
		fiber:  app,
	}
}

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

// 签名中间件
func SignMiddleware(c *fiber.Ctx) error {
	//todo 待正式环境下完善
	if !utils.IsProd() {
		return c.Next()
	}

	sts := c.Get("ts")
	ts, err := strconv.ParseInt(sts, 10, 64)
	if err != nil {
		slog.Info("ts convert error", "ts", sts)
		return utils.JsonFailed(c, http.StatusBadRequest, "sign error")
	}
	if time.Now().UnixMilli()-ts > 300000 {
		slog.Info("sign expired", "ts", sts)
		return utils.JsonFailed(c, http.StatusBadRequest, "sign error")
	}

	stoken := c.Get("stoken")
	nonce := c.Get("nonce")
	path := c.Path()
	salt := "fiber-boot"

	signStr := fmt.Sprintf("%s%s%s%s%s", stoken, nonce, path, sts, salt)
	sign2 := utils.MD5(signStr)

	sign := c.Get("sign")
	if sign != sign2 {
		slog.Info("sign error", "sign", sign, "sign2", sign2)
		return utils.JsonFailed(c, http.StatusBadRequest, "sign error")
	}
	return c.Next()
}

// 鉴权中间件
func authMiddleware(c *fiber.Ctx) error {
	stoken := c.Get("token")
	if stoken == "" {
		return utils.JsonFailed(c, http.StatusUnauthorized, "unauthorized")
	}

	value, found := utils.AccessTokenCache().Get(stoken)
	if !found {
		accessToken, err := tokenDao.GetByToken(stoken)
		if err != nil && err != gorm.ErrRecordNotFound {
			return xerrors.New("token not found", err)
		}
		if accessToken == nil {
			return utils.JsonFailed(c, http.StatusUnauthorized, "unauthorized")
		}
		if accessToken.Expires < time.Now().Unix() {
			return utils.JsonFailed(c, http.StatusUnauthorized, "unauthorized")
		}
		value = accessToken.UID
	}

	uid := value.(int64)
	c.Locals("uid", uid)
	return c.Next()
}
