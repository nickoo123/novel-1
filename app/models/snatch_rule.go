// Copyright 2017 Vckai Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

// 采集规则配置模型
type SnatchRule struct {
	Id        uint32 `orm:"auto;size(11)"`
	Name      string `orm:"size(50)"`
	Code      string `orm:"size(20)"`
	Url       string `orm:"size(50)"`
	State     uint8  `orm:"size(1)"`
	IsUpdate  uint8  `orm:"size(1);default(0)"`
	Charset   string `orm:"size(10)"`
	CateMap   string `orm:"size(2555)"`
	Rule      string `orm:"size(2555)"`
	TestData  string `orm:"size(2555)"`
	CreatedAt uint32 `orm:"size(11)"`
	UpdatedAt uint32 `orm:"size(11)"`
	DeletedAt uint32 `orm:"size(11)"`

	Test     *SnatchTestData `orm:"-"`
	Rules    *Rule           `orm:"-"`
	CateMaps []*CateMap      `orm:"-"`
}

// 分类映射
type CateMap struct {
	Id   uint32 `json:"id"`
	Name string `json:"name"`
}

// 采集测试对比数据
type SnatchTestData struct {
	// 小说菜单
	BookMenuList string `json:"book_menu_list"` // 小说菜单列表
	BookMenuNode string `json:"book_menu_node"` // 小说菜单节点

	// 小说列表
	BookList    string `json:"book_list"`     // 小说列表
	BookListSub string `json:"book_list_sub"` // 小说列表子级

	// 小说简介页面
	BookURL        string `json:"book_url"`         // 小说简介URL
	BookTitle      string `json:"book_title"`       // 小说标题
	BookAuthor     string `json:"book_author"`      // 小说作者
	BookCate       string `json:"book_cate"`        // 小说分类名
	BookDesc       string `json:"book_desc"`        // 小说简介
	BookCover      string `json:"book_cover"`       // 小说封面LOGO
	BookChapterURL string `json:"book_chapter_url"` // 小说章节页面URL

	// 章节目录
	ChapterURL          string `json:"chapter_url"`           // 章节目录URL
	ChapterCatalogURL   string `json:"chapter_catalog_url"`   // 章节目录中第一个章节的URL地址
	ChapterCatalogTitle string `json:"chapter_catalog_title"` // 章节目录中第一个章节的标题

	// 详情页面
	InfoURL         string `json:"info_url"`           // 小说章节页面URL
	InfoTitle       string `json:"info_title"`         // 小说章节标题
	InfoDesc        string `json:"info_desc"`          // 小说章节内容
	InfoPrePageURL  string `json:"info_pre_page_url"`  // 上一章节URL
	InfoNextPageURL string `json:"info_next_page_url"` // 下一章节URL

	// 搜索页面获得URL
	FindBookKw  string `json:"find_book_kw"`  // 搜索关键词
	FindBookURL string `json:"find_book_url"` // 搜索结果小说简介页面的URL
}

// 采集规则
type Rule struct {
	// 是否小说简介页面
	IsBookURL string `json:"is_book_url"`

	// 是否可爬虫页面
	IsCrawlerURL string `json:"is_crawler_url"`

	// 小说简介页面采集规则
	BookMenuListSelector         string `json:"book_menu_list_selector"`          // 小说菜单列表选择器
	BookMenuListAttr             string `json:"book_menu_list_attr"`              // 小说菜单列表获取属性
	BookMenuListFilter           string `json:"book_menu_list_filter"`            // 小说菜单列表过滤规则
	BookMenuNodeSelector         string `json:"book_menu_node_selector"`          // 小说菜单节点选择器
	BookMenuNodeAttr             string `json:"book_menu_node_attr"`              // 小说菜单节点获取属性
	BookMenuNodeFilter           string `json:"book_menu_node_filter"`            // 小说菜单节点过滤规则
	BookListSelector             string `json:"book_list_selector"`               // 小说列表选择器
	BookListAttr                 string `json:"book_list_attr"`                   // 小说列表获取属性
	BookListFilter               string `json:"book_list_filter"`                 // 小说列表过滤规则
	BookListSubSelector          string `json:"book_list_sub_selector"`           // 小说列表子级选择器
	BookListSubAttr              string `json:"book_list_sub_attr"`               // 小说列表子级选择属性
	BookTitleSelector            string `json:"book_title_selector"`              // 小说标题选择器
	BookTitleAttr                string `json:"book_title_attr"`                  // 小说标题获取属性
	BookTitleFilter              string `json:"book_title_filter"`                // 小说标题过滤规则
	BookAuthorSelector           string `json:"book_author_selector"`             // 作者选择器
	BookAuthorAttr               string `json:"book_author_attr"`                 // 作者选择属性
	BookAuthorFilter             string `json:"book_author_filter"`               // 作者过滤规则
	BookCateSelector             string `json:"book_cate_selector"`               // 分类选择器
	BookCateAttr                 string `json:"book_cate_attr"`                   // 分类选择属性
	BookCateFilter               string `json:"book_cate_filter"`                 // 分类过滤规则
	BookDescSelector             string `json:"book_desc_selector"`               // 简介选择器
	BookDescAttr                 string `json:"book_desc_attr"`                   // 简介选择属性
	BookDescFilter               string `json:"book_desc_filter"`                 // 简介过滤规则
	BookCoverSelector            string `json:"book_cover_selector"`              // 缩略图选择器
	BookCoverAttr                string `json:"book_cover_attr"`                  // 缩略图选择属性
	BookNoCover                  string `json:"book_no_cover"`                    // 无缩略图标志
	BookChapterURLSelector       string `json:"book_chapter_url_selector"`        // 章节目录页链接CSS选择器
	BookChapterURLAttr           string `json:"book_chapter_url_attr"`            // 章节目录页链接选择属性
	BookLastChapterTitleSelector string `json:"book_last_chapter_title_selector"` // 最新章节标题CSS选择器
	BookLastChapterTitleAttr     string `json:"book_last_chapter_title_attr"`     // 最新章节选择属性

	// 章节目录采集规则
	ChapterCatalogSelector  string `json:"chapter_catalog_selector"`   // 章节目录选择器
	ChapterNextPageSelector string `json:"chapter_next_page_selector"` // 章节目录下一页选择器
	ChapterAbandonNum       int    `json:"chapter_abandon_num"`        // 丢弃章节数（最新章节）

	// 详情页面采集规则
	InfoTitleSelector    string `json:"info_title_selector"`     // 章节标题选择器
	InfoTitleFilter      string `json:"info_title_filter"`       // 章节标题过滤规则
	InfoDescSelector     string `json:"info_desc_selector"`      // 章节内容选择器
	InfoDescFilter       string `json:"info_desc_filter"`        // 章节内容过滤规则
	InfoPrePageSelector  string `json:"info_pre_page_selector"`  // 上一页详情页面选择器
	InfoNextPageSelector string `json:"info_next_page_selector"` // 下一页详情页面选择器

	// 搜索页面采集规则
	FindURL             string `json:"find_url"`               // 搜索小说URL
	FindCharset         string `json:"find_charset"`           // 搜索页面编码
	FindBookURLSelector string `json:"find_book_url_selector"` // 小说链接CSS选择器
}

type ArgsListRule struct {
	Ids []interface{}
}

func NewSnatchRule() *SnatchRule {
	return &SnatchRule{
		Rules:    &Rule{},
		Test:     &SnatchTestData{},
		CateMaps: make([]*CateMap, 0),
	}
}

// 初始化
// 注册模型
func init() {
	orm.RegisterModelWithPrefix("nov_", new(SnatchRule))
}

func (m *SnatchRule) query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}

// 添加
func (m *SnatchRule) Insert() error {
	m.CreatedAt = uint32(time.Now().Unix())
	m.UpdatedAt = uint32(time.Now().Unix())

	m.Encode()
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

// 读取
func (m *SnatchRule) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}

	m.Decode()

	return nil
}

// 根据接口获取权限id
func (m *SnatchRule) GetByCode(code string) *SnatchRule {
	var r SnatchRule

	m.query().Filter("code", code).Filter("deleted_at", 0).One(&r, "id", "name", "code", "is_update", "url", "state", "charset", "cate_map", "rule")

	r.Decode()

	return &r
}

// 修改
func (m *SnatchRule) Update(fields ...string) error {
	m.UpdatedAt = uint32(time.Now().Unix())
	if len(fields) > 0 {
		fields = append(fields, "updated_at")
	}
	m.Encode()
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

// 批量删除
func (m *SnatchRule) DeleteBatch(ids []string) error {
	marks := make([]string, len(ids))
	for i := range marks {
		marks[i] = "?"
	}

	query := fmt.Sprintf("IN (%s)", strings.Join(marks, ", "))
	sqlStr := fmt.Sprintf("UPDATE nov_snatch_rule SET deleted_at=? WHERE `id` %s", query)

	_, err := orm.NewOrm().Raw(sqlStr, uint32(time.Now().Unix()), ids).Exec()

	return err
}

// 删除
func (m *SnatchRule) Delete(forceDelete ...bool) error {
	// 软删除
	if len(forceDelete) == 0 {
		m.DeletedAt = uint32(time.Now().Unix())
		return m.Update("deleted_at")
	}

	// 强制删除
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *SnatchRule) Encode() error {
	if m.Rules != nil {
		data, err := json.Marshal(m.Rules)
		if err != nil {
			return errors.New("Json marshaling failed：" + err.Error())
		}

		m.Rule = string(data)
	}

	if m.CateMaps != nil {
		data, err := json.Marshal(m.CateMaps)
		if err != nil {
			return errors.New("Json marshaling failed：" + err.Error())
		}

		m.CateMap = string(data)
	}

	if m.Test != nil {
		data, err := json.Marshal(m.Test)
		if err != nil {
			return errors.New("Json marshaling failed：" + err.Error())
		}

		m.TestData = string(data)
	}

	return nil
}

func (m *SnatchRule) Decode() error {
	m.Rules = &Rule{}
	m.Test = &SnatchTestData{}
	m.CateMaps = make([]*CateMap, 0)

	if len(m.Rule) > 0 {
		err := json.Unmarshal([]byte(m.Rule), &m.Rules)
		if err != nil {
			return errors.New("Json marshaling failed：" + err.Error())
		}
	}

	if len(m.CateMap) > 0 {
		err := json.Unmarshal([]byte(m.CateMap), &m.CateMaps)
		if err != nil {
			return errors.New("Json marshaling failed：" + err.Error())
		}
	}

	if len(m.TestData) > 0 {
		err := json.Unmarshal([]byte(m.TestData), &m.Test)
		if err != nil {
			return errors.New("Json marshaling failed：" + err.Error())
		}
	}

	return nil
}

// 获取配置列表
func (m *SnatchRule) GetAll(args *ArgsListRule) []*SnatchRule {
	list := make([]*SnatchRule, 0)

	q := m.query().Filter("deleted_at", 0)

	if len(args.Ids) > 0 {
		q = q.Filter("id__in", args.Ids...)
	}

	q.OrderBy("-id").All(&list, "id", "name", "code", "url", "is_update", "state", "charset", "cate_map", "rule", "test_data", "created_at")

	return list
}

// 获取是否启用
func (m SnatchRule) StateName() string {
	if m.State == 0 {
		return `<span class="layui-btn layui-btn-primary layui-btn-xs btn-set-state">启用</span>`
	}

	return `<span class="layui-btn layui-btn-xs btn-set-un-state">禁用</span>`
}
