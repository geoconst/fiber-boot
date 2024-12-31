package dao

import (
	"fiber-boot/internal/domain/model"
	"time"

	"gorm.io/gorm"
)

type AccessTokenDAO struct {
	BaseDao[model.AccessToken]
}

func NewAccessTokenDAO(db *gorm.DB) *AccessTokenDAO {
	return &AccessTokenDAO{
		BaseDao: BaseDao[model.AccessToken]{db},
	}
}

func (dao *AccessTokenDAO) GetByToken(token string) (*model.AccessToken, error) {
	var accessToken model.AccessToken
	err := dao.DB.Where("stoken =?", token).First(&accessToken).Error
	if err != nil {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &accessToken, err
}

func (dao *AccessTokenDAO) Create(accessToken *model.AccessToken) error {
	unix := time.Now().Unix()
	accessToken.Ctime = &unix
	return dao.DB.Create(accessToken).Error
}
