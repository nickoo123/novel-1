package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type ArgsSpliderList struct {
	ArgsBase
	IsRec int
}

// 蜘蛛管理模型
type Splider struct {
	Id        uint64 `orm:"auto"`
	Req       string `orm:"size(500);unique"`
	Num       uint64 `orm:"size(10)"`
	CreatedAt uint32 `orm:"size(11)"`
	UpdatedAt uint32 `orm:"size(11)"`
}

func NewSplider() *Splider {
	return &Splider{}
}

// 初始化
// 注册模型
func init() {
	orm.RegisterModelWithPrefix(TABLE_PREFIX, new(Splider))
}

func (m *Splider) query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}

func (m *Splider) InsertOrIncrement(req string) error {
	sqlStr := fmt.Sprintf("INSERT INTO %ssplider SET req=?, num=1, created_at=?, updated_at=? ON DUPLICATE KEY UPDATE `num` = `num` + 1, updated_at=?", TABLE_PREFIX)

	t := uint32(time.Now().Unix())
	fmt.Println(sqlStr)
	_, err := orm.NewOrm().Raw(sqlStr, req, t, t, t).Exec()

	return err
}

// 修改
func (m *Splider) Update(fields ...string) error {
	m.UpdatedAt = uint32(time.Now().Unix())
	if len(fields) > 0 {
		fields = append(fields, "updated_at")
	}

	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}

	return nil
}

// 删除
func (m *Splider) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

// 获取蜘蛛列表
func (m *Splider) GetAll(args ArgsSpliderList) ([]*Splider, int64) {
	list := make([]*Splider, 0)
	qs := m.query()

	var count int64 = 0
	if args.Count {
		count, _ = qs.Count()
	}

	// 分页
	if args.Limit > 0 {
		qs = qs.Limit(args.Limit, args.Offset)
	}

	orderBy := "-id"
	if args.OrderBy != "" {
		orderBy = args.OrderBy
	}

	if count > 0 || args.Count == false {
		qs.OrderBy(orderBy).All(&list, "id", "req", "num", "created_at")
	}

	return list, count
}
