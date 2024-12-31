package account

import (
	"fiber-boot/internal/dao"
	"fiber-boot/internal/domain/model"
	"fiber-boot/internal/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mdobak/go-xerrors"
)

type AccountHandler struct {
	accountDao     *dao.AccountDAO
	accessTokenDao *dao.AccessTokenDAO
}

func NewAccountHandler(accountDao *dao.AccountDAO, accessTokenDao *dao.AccessTokenDAO) *AccountHandler {
	return &AccountHandler{accountDao: accountDao, accessTokenDao: accessTokenDao}
}

func RegisterPublicRoutes(app *fiber.App, accountHandler *AccountHandler) {
	route := app.Group("/account")
	route.Post("/login", accountHandler.Login)
}

func RegisterPrivateRoutes(app *fiber.App, accountHandler *AccountHandler) {
	route := app.Group("/account")
	route.Get("/logout", accountHandler.Logout)
	route.Get("/credit", accountHandler.Credit)
}

// Logout 登出
func (h *AccountHandler) Logout(c *fiber.Ctx) error {
	token := c.Get("token")
	utils.AccessTokenCache().Delete(token)
	return utils.JsonSuccess(c, nil)
}

// Login 登录
func (h *AccountHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return &utils.MyValidError{Msg: "Invalid request body"}
	}

	var phone string = "12488887777"

	account, err := h.accountDao.GetByPhone(phone)
	if err != nil {
		return xerrors.New("GetByPhone failed", err)
	}

	if account == nil {
		account := model.Account{Phone: phone}
		err = h.accountDao.Create(&account)
		if err != nil {
			return xerrors.New("Create account failed", err)
		}
	}

	accessToken := utils.MD5(fmt.Sprintf("%s%d", req.Token, time.Now().UnixMilli()))
	confuseToken := utils.MD5(fmt.Sprintf("%s%s%d", req.Token, accessToken, time.Now().UnixMilli()))
	expires := time.Now().Add(time.Hour * 24 * 30).Unix()

	err = h.accessTokenDao.Create(&model.AccessToken{
		Stoken:  accessToken,
		Expires: expires,
		UID:     account.ID,
	})

	if err != nil {
		return xerrors.New("Create access token failed", err)
	}

	utils.AccessTokenCache().Set(accessToken, account.ID, 0)
	return utils.JsonSuccess(c, fiber.Map{
		"stoken":   accessToken,
		"mtoken":   confuseToken,
		"expires":  expires,
		"username": phone,
	})
}

// 查询积分
func (h *AccountHandler) Credit(c *fiber.Ctx) error {
	uid := c.Locals("uid").(int64)
	account, err := h.accountDao.GetById(uid)
	if err != nil {
		return xerrors.New("GetByID failed", err)
	}
	return utils.JsonSuccess(c, fiber.Map{
		"credit": account.Credit,
	})
}
