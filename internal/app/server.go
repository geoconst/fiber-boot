package app

import (
	"fiber-boot/internal/core"
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
	config         *core.Config
	fiber          *fiber.App
	tokenDao       *dao.AccessTokenDAO
	accountHandler *account.AccountHandler
}

func NewServer(
	config *core.Config,
	accountHandler *account.AccountHandler,
	accessTokenDao *dao.AccessTokenDAO) *Server {

	f := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 30,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	return &Server{
		config:         config,
		fiber:          f,
		tokenDao:       accessTokenDao,
		accountHandler: accountHandler,
	}
}

func (s *Server) Start() error {
	s.RegisterRoutes()
	return s.fiber.Listen(fmt.Sprintf(":%d", s.config.Port))
}

// 注册路由
func (s *Server) RegisterRoutes() {
	logConfig := slogfiber.Config{
		DefaultLevel:      slog.LevelDebug,
		ClientErrorLevel:  slog.LevelWarn,
		ServerErrorLevel:  slog.LevelError,
		WithRequestHeader: true,
	}

	s.fiber.Use(slogfiber.NewWithConfig(slog.Default(), logConfig))
	s.fiber.Use(recover.New(recover.Config{EnableStackTrace: true}))
	// s.fiber.Use(s.SignMiddleware)

	account.RegisterPublicRoutes(s.fiber, s.accountHandler)

	s.fiber.Use(s.authMiddleware)
	account.RegisterPrivateRoutes(s.fiber, s.accountHandler)
}

// 签名中间件
func (s *Server) SignMiddleware(c *fiber.Ctx) error {
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
func (s *Server) authMiddleware(c *fiber.Ctx) error {
	stoken := c.Get("token")
	if stoken == "" {
		return utils.JsonFailed(c, http.StatusUnauthorized, "unauthorized")
	}

	value, found := utils.AccessTokenCache().Get(stoken)
	if !found {
		accessToken, err := s.tokenDao.GetByToken(stoken)
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
