package services

import (
	"errors"

	"novel/app/models"
)

// 定义SpliderService
type Splider struct {
}

func NewSplider() *Splider {
	return &Splider{}
}

// 获取Splider关键字
func (this *Splider) GetRes(size int) []*models.Splider {
	args := models.ArgsSpliderList{}
	args.Limit = size

	list, _ := models.SpliderModel.GetAll(args)

	return list
}

// 获取搜索记录列表
func (this *Splider) GetAll(args models.ArgsSpliderList) ([]*models.Splider, int64) {
	return models.SpliderModel.GetAll(args)
}

// 添加搜索记录
func (this *Splider) InsertOrIncrement(req string) error {
	if len(req) == 0 {
		return errors.New("params error")
	}

	err := models.SpliderModel.InsertOrIncrement(req)

	return err
}

// 删除搜索记录
func (this *Splider) Delete(id uint64) error {
	if id < 0 {
		return errors.New("params error")
	}

	s := models.Splider{Id: id}
	err := s.Delete()
	if err != nil {
		return err
	}

	return nil
}
