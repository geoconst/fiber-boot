package dao

import (
	"fiber-boot/internal/domain/model"
	"time"

	"gorm.io/gorm"
)

type AccountDAO struct {
	BaseDao[model.Account]
}

func NewAccountDAO(db *gorm.DB) *AccountDAO {
	return &AccountDAO{
		BaseDao: BaseDao[model.Account]{db},
	}
}

func (dao *AccountDAO) GetByPhone(phone string) (*model.Account, error) {
	var account model.Account
	err := dao.DB.Where("phone =?", phone).First(&account).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &account, nil
}

func (dao *AccountDAO) Create(account *model.Account) error {
	unix := time.Now().Unix()
	account.Ctime = &unix
	return dao.DB.Create(account).Error
}
