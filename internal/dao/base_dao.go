package dao

import (
	"gorm.io/gorm"
)

type BaseDao[T any] struct {
	DB *gorm.DB
}

func (base *BaseDao[T]) GetById(id any) (*T, error) {
	var t T
	err := base.DB.First(&t, id).Error
	return &t, err
}

func (base *BaseDao[T]) FindAll() ([]T, error) {
	var ts []T
	err := base.DB.Find(&ts).Error
	return ts, err
}

func (base *BaseDao[T]) Create(tx *gorm.DB, t *T) error {
	if tx != nil {
		return tx.Create(t).Error
	}
	return base.DB.Create(t).Error
}

func (base *BaseDao[T]) Update(tx *gorm.DB, t *T) error {
	if tx != nil {
		return tx.Save(t).Error
	}
	return base.DB.Save(t).Error
}

func (base *BaseDao[T]) DeleteById(tx *gorm.DB, id any) error {
	var model T
	if tx != nil {
		return tx.Delete(&model, id).Error
	}
	return base.DB.Delete(&model, id).Error
}
