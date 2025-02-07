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

package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	gutils "novel/app/utils"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/validation"

	"novel/app/models"
)

// 定义NovelService
type Novel struct {
}

func NewNovel() *Novel {
	return &Novel{}
}

// 判断小说是否存在
func (this *Novel) IsExists(name string) bool {
	if len(name) == 0 {
		return false
	}

	nov := models.NovelModel.GetByName(name)
	if nov != nil {
		return true
	}

	return false
}

// 获取小说采集站点列表
func (this *Novel) GetLinks(novId uint32) []*models.NovelLinks {
	args := models.ArgsNovelLinksList{}
	args.NovId = novId
	args.Limit = 100

	links, _ := models.NovelLinksModel.GetAll(args)

	return links
}

// 根据你URL 和 站点获取小说采集URL
func (this *Novel) GetByLink(link, source string) *models.NovelLinks {
	return models.NovelLinksModel.GetByLink(link, source)
}

// 添加小说采集站点
func (this *Novel) AddLink(novId uint32, link, source, chapterLink string) error {
	// 采集点已存在则直接返回
	novLink := this.GetByLink(link, source)
	if novLink.Id > 0 {
		return nil
	}

	n := &models.NovelLinks{
		NovId:       novId,
		Link:        link,
		Source:      source,
		ChapterLink: chapterLink,
	}

	return n.Insert()
}

// 判断小说是否存在
func (this *Novel) GetByName(name string) *models.Novel {
	nov := models.NovelModel.GetByName(name)

	return nov
}

// 判断小说HashKey是否存在
func (this *Novel) GetByHashKey(hashkey string) *models.Novel {
	nov := models.NovelModel.GetByHashKey(hashkey)

	return nov
}

// 获取单个管理员信息
func (this *Novel) Get(id uint32) *models.Novel {
	if id < 0 {
		return nil
	}

	var buf bytes.Buffer
	buf.WriteString("Novel:get:")
	buf.WriteString("id_" + string(id))
	res := gutils.GetRedisKeys(buf.String())
	if res != nil {
		var novs *models.Novel
		str := fmt.Sprintf("%s", res)
		if err := json.Unmarshal([]byte(str), &novs); err != nil {
			fmt.Println("err:-----", err.Error())
			return nil
		}
		return novs
	} else {
		nov := models.NewNovel()
		nov.Id = id
		err := nov.Read()
		if err != nil {
			return nil
		}
		var novs = models.Novel{
			Id:               nov.Id,
			HashKey:          nov.HashKey,
			Name:             nov.Name,
			Desc:             nov.Desc,
			Cover:            nov.Cover,
			CateId:           nov.CateId,
			CateName:         nov.CateName,
			Author:           nov.Author,
			IsOriginal:       nov.IsOriginal,
			IsHot:            nov.IsHot,
			IsRec:            nov.IsRec,
			IsTodayRec:       nov.IsTodayRec,
			IsVipRec:         nov.IsVipRec,
			IsVipReward:      nov.IsVipReward,
			IsVipUp:          nov.IsVipUp,
			IsSignNewBook:    nov.IsSignNewBook,
			IsCollect:        nov.IsCollect,
			Status:           nov.Status,
			Views:            nov.Views,
			CollectNum:       nov.Id,
			RecNum:           nov.RecNum,
			TextNum:          nov.TextNum,
			ChapterNum:       nov.ChapterNum,
			ChapterUpdatedAt: nov.ChapterUpdatedAt,
			ChapterId:        nov.ChapterId,
			ChapterTitle:     nov.ChapterTitle,
			CreatedAt:        nov.CreatedAt,
			UpdatedAt:        nov.UpdatedAt,
			DeletedAt:        nov.DeletedAt,
		}
		str, _ := json.Marshal(novs)
		gutils.SetRedisKeyValue(buf.String(), string(str))
		return nov
	}
}

// 获取今日推荐
func (this *Novel) GetTodayRecs(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.Fields = []string{"id", "hash_key", "name", "cover", "author"}
	args.FilterMaps = map[string]int{
		"is_today_rec": 1,
	}
	args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	var buf bytes.Buffer
	buf.WriteString("GetTodayRecs:")
	buf.WriteString("size_" + strconv.Itoa(size))
	buf.WriteString(":offset_" + strconv.Itoa(offset))
	res := gutils.GetRedisKeys(buf.String())
	if res != nil {
		novs := make([]*models.Novel, 0)
		str := fmt.Sprintf("%s", res)
		if err := json.Unmarshal([]byte(str), &novs); err != nil {
			fmt.Println("err:-----", err.Error())
			return nil
		}
		return novs
	} else {
		novs, _ := this.GetAll(args)
		var novlist []models.Novel
		for _, n := range novs {
			var nov = models.Novel{
				Id:               n.Id,
				HashKey:          n.HashKey,
				Name:             n.Name,
				Desc:             n.Desc,
				Cover:            n.Cover,
				CateId:           n.CateId,
				CateName:         n.CateName,
				Author:           n.Author,
				IsOriginal:       n.IsOriginal,
				IsHot:            n.IsHot,
				IsRec:            n.IsRec,
				IsTodayRec:       n.IsTodayRec,
				IsVipRec:         n.IsVipRec,
				IsVipReward:      n.IsVipReward,
				IsVipUp:          n.IsVipUp,
				IsSignNewBook:    n.IsSignNewBook,
				IsCollect:        n.IsCollect,
				Status:           n.Status,
				Views:            n.Views,
				CollectNum:       n.CollectNum,
				RecNum:           n.RecNum,
				TextNum:          n.TextNum,
				ChapterNum:       n.ChapterNum,
				ChapterUpdatedAt: n.ChapterUpdatedAt,
				ChapterId:        n.ChapterId,
				ChapterTitle:     n.ChapterTitle,
				CreatedAt:        n.CreatedAt,
				UpdatedAt:        n.UpdatedAt,
				DeletedAt:        n.DeletedAt,
			}
			novlist = append(novlist, nov)
		}
		str, _ := json.Marshal(novlist)
		gutils.SetRedisKeyValue(buf.String(), string(str))
		return novs
	}
}

// 获取推荐小说
func (this *Novel) GetRecs(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.Fields = []string{"id", "hash_key", "name", "cover", "cate_id", "cate_name", "author"}
	args.FilterMaps = map[string]int{
		"is_rec": 1,
	}
	args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons
	var buf bytes.Buffer
	buf.WriteString("GetRecs:")
	buf.WriteString("size_" + strconv.Itoa(size))
	buf.WriteString(":offset_" + strconv.Itoa(offset))
	res := gutils.GetRedisKeys(buf.String())
	if res != nil {
		novs := make([]*models.Novel, 0)
		str := fmt.Sprintf("%s", res)
		if err := json.Unmarshal([]byte(str), &novs); err != nil {
			fmt.Println("err:-----", err.Error())
			return nil
		}
		return novs
	} else {
		novs, _ := this.GetAll(args)
		var novlist []models.Novel
		for _, n := range novs {
			var nov = models.Novel{
				Id:               n.Id,
				HashKey:          n.HashKey,
				Name:             n.Name,
				Desc:             n.Desc,
				Cover:            n.Cover,
				CateId:           n.CateId,
				CateName:         n.CateName,
				Author:           n.Author,
				IsOriginal:       n.IsOriginal,
				IsHot:            n.IsHot,
				IsRec:            n.IsRec,
				IsTodayRec:       n.IsTodayRec,
				IsVipRec:         n.IsVipRec,
				IsVipReward:      n.IsVipReward,
				IsVipUp:          n.IsVipUp,
				IsSignNewBook:    n.IsSignNewBook,
				IsCollect:        n.IsCollect,
				Status:           n.Status,
				Views:            n.Views,
				CollectNum:       n.CollectNum,
				RecNum:           n.RecNum,
				TextNum:          n.TextNum,
				ChapterNum:       n.ChapterNum,
				ChapterUpdatedAt: n.ChapterUpdatedAt,
				ChapterId:        n.ChapterId,
				ChapterTitle:     n.ChapterTitle,
				CreatedAt:        n.CreatedAt,
				UpdatedAt:        n.UpdatedAt,
				DeletedAt:        n.DeletedAt,
			}
			novlist = append(novlist, nov)
		}
		str, _ := json.Marshal(novlist)
		gutils.SetRedisKeyValue(buf.String(), string(str))
		return novs
	}
}

// 获取精品推荐小说
func (this *Novel) GetVipRecs(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.Fields = []string{"id", "hash_key", "name", "cover"}
	args.FilterMaps = map[string]int{
		"is_vip_rec": 1,
	}
	args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	var buf bytes.Buffer
	buf.WriteString("GetVipRecs:")
	buf.WriteString("size_" + strconv.Itoa(size))
	buf.WriteString(":offset_" + strconv.Itoa(offset))
	res := gutils.GetRedisKeys(buf.String())
	if res != nil {
		novs := make([]*models.Novel, 0)
		str := fmt.Sprintf("%s", res)
		if err := json.Unmarshal([]byte(str), &novs); err != nil {
			fmt.Println("err:-----", err.Error())
			return nil
		}
		return novs
	} else {
		novs, _ := this.GetAll(args)
		var novlist []models.Novel
		for _, n := range novs {
			var nov = models.Novel{
				Id:               n.Id,
				HashKey:          n.HashKey,
				Name:             n.Name,
				Desc:             n.Desc,
				Cover:            n.Cover,
				CateId:           n.CateId,
				CateName:         n.CateName,
				Author:           n.Author,
				IsOriginal:       n.IsOriginal,
				IsHot:            n.IsHot,
				IsRec:            n.IsRec,
				IsTodayRec:       n.IsTodayRec,
				IsVipRec:         n.IsVipRec,
				IsVipReward:      n.IsVipReward,
				IsVipUp:          n.IsVipUp,
				IsSignNewBook:    n.IsSignNewBook,
				IsCollect:        n.IsCollect,
				Status:           n.Status,
				Views:            n.Views,
				CollectNum:       n.CollectNum,
				RecNum:           n.RecNum,
				TextNum:          n.TextNum,
				ChapterNum:       n.ChapterNum,
				ChapterUpdatedAt: n.ChapterUpdatedAt,
				ChapterId:        n.ChapterId,
				ChapterTitle:     n.ChapterTitle,
				CreatedAt:        n.CreatedAt,
				UpdatedAt:        n.UpdatedAt,
				DeletedAt:        n.DeletedAt,
			}
			novlist = append(novlist, nov)
		}
		str, _ := json.Marshal(novlist)
		gutils.SetRedisKeyValue(buf.String(), string(str))
		return novs
	}
}

// 获取原创小说
func (this *Novel) GetOriginals(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.Fields = []string{"id", "hash_key", "name", "cover", "desc", "views", "author", "cate_id", "cate_name"}
	args.FilterMaps = map[string]int{
		"is_original": 1,
	}
	args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	novs, _ := this.GetAll(args)

	return novs
}

// 获取热门小说
func (this *Novel) GetHots(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.Fields = []string{"id", "hash_key", "name", "cover", "desc", "views", "author", "cate_id", "cate_name"}
	args.FilterMaps = map[string]int{
		"is_hot": 1,
	}
	args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	var buf bytes.Buffer
	buf.WriteString("GetHots:")
	buf.WriteString("size_" + strconv.Itoa(size))
	buf.WriteString(":offset_" + strconv.Itoa(offset))
	res := gutils.GetRedisKeys(buf.String())
	if res != nil {
		novs := make([]*models.Novel, 0)
		str := fmt.Sprintf("%s", res)
		if err := json.Unmarshal([]byte(str), &novs); err != nil {
			fmt.Println("err:-----", err.Error())
			return nil
		}
		return novs
	} else {
		novs, _ := this.GetAll(args)
		var novlist []models.Novel
		for _, n := range novs {
			var nov = models.Novel{
				Id:               n.Id,
				HashKey:          n.HashKey,
				Name:             n.Name,
				Desc:             n.Desc,
				Cover:            n.Cover,
				CateId:           n.CateId,
				CateName:         n.CateName,
				Author:           n.Author,
				IsOriginal:       n.IsOriginal,
				IsHot:            n.IsHot,
				IsRec:            n.IsRec,
				IsTodayRec:       n.IsTodayRec,
				IsVipRec:         n.IsVipRec,
				IsVipReward:      n.IsVipReward,
				IsVipUp:          n.IsVipUp,
				IsSignNewBook:    n.IsSignNewBook,
				IsCollect:        n.IsCollect,
				Status:           n.Status,
				Views:            n.Views,
				CollectNum:       n.CollectNum,
				RecNum:           n.RecNum,
				TextNum:          n.TextNum,
				ChapterNum:       n.ChapterNum,
				ChapterUpdatedAt: n.ChapterUpdatedAt,
				ChapterId:        n.ChapterId,
				ChapterTitle:     n.ChapterTitle,
				CreatedAt:        n.CreatedAt,
				UpdatedAt:        n.UpdatedAt,
				DeletedAt:        n.DeletedAt,
			}
			novlist = append(novlist, nov)
		}
		str, _ := json.Marshal(novlist)
		gutils.SetRedisKeyValue(buf.String(), string(str))
		return novs
	}
}

// 获取新人签约榜列表
func (this *Novel) GetSignNewBooks(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.Fields = []string{"id", "hash_key", "name", "cover", "desc", "views", "author", "cate_id", "cate_name"}
	args.FilterMaps = map[string]int{
		"is_sign_new_book": 1,
	}
	args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	var buf bytes.Buffer
	buf.WriteString("GetSignNewBooks:")
	buf.WriteString("size_" + strconv.Itoa(size))
	buf.WriteString(":offset_" + strconv.Itoa(offset))
	res := gutils.GetRedisKeys(buf.String())
	if res != nil {
		novs := make([]*models.Novel, 0)
		str := fmt.Sprintf("%s", res)
		if err := json.Unmarshal([]byte(str), &novs); err != nil {
			fmt.Println("err:-----", err.Error())
			return nil
		}
		return novs
	} else {
		novs, _ := this.GetAll(args)
		var novlist []models.Novel
		for _, n := range novs {
			var nov = models.Novel{
				Id:               n.Id,
				HashKey:          n.HashKey,
				Name:             n.Name,
				Desc:             n.Desc,
				Cover:            n.Cover,
				CateId:           n.CateId,
				CateName:         n.CateName,
				Author:           n.Author,
				IsOriginal:       n.IsOriginal,
				IsHot:            n.IsHot,
				IsRec:            n.IsRec,
				IsTodayRec:       n.IsTodayRec,
				IsVipRec:         n.IsVipRec,
				IsVipReward:      n.IsVipReward,
				IsVipUp:          n.IsVipUp,
				IsSignNewBook:    n.IsSignNewBook,
				IsCollect:        n.IsCollect,
				Status:           n.Status,
				Views:            n.Views,
				CollectNum:       n.CollectNum,
				RecNum:           n.RecNum,
				TextNum:          n.TextNum,
				ChapterNum:       n.ChapterNum,
				ChapterUpdatedAt: n.ChapterUpdatedAt,
				ChapterId:        n.ChapterId,
				ChapterTitle:     n.ChapterTitle,
				CreatedAt:        n.CreatedAt,
				UpdatedAt:        n.UpdatedAt,
				DeletedAt:        n.DeletedAt,
			}
			novlist = append(novlist, nov)
		}
		str, _ := json.Marshal(novlist)
		gutils.SetRedisKeyValue(buf.String(), string(str))
		return novs
	}
}

// 获取收藏榜列表
func (this *Novel) GetCollects(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.Fields = []string{"id", "hash_key", "name", "cover", "desc", "views", "author", "cate_id", "cate_name"}
	args.FilterMaps = map[string]int{
		"is_collect": 1,
	}
	args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	var buf bytes.Buffer
	buf.WriteString("GetCollects:")
	buf.WriteString("size_" + strconv.Itoa(size))
	buf.WriteString(":offset_" + strconv.Itoa(offset))
	res := gutils.GetRedisKeys(buf.String())
	if res != nil {
		novs := make([]*models.Novel, 0)
		str := fmt.Sprintf("%s", res)
		if err := json.Unmarshal([]byte(str), &novs); err != nil {
			fmt.Println("err:-----", err.Error())
			return nil
		}
		return novs
	} else {
		novs, _ := this.GetAll(args)
		var novlist []models.Novel
		for _, n := range novs {
			var nov = models.Novel{
				Id:               n.Id,
				HashKey:          n.HashKey,
				Name:             n.Name,
				Desc:             n.Desc,
				Cover:            n.Cover,
				CateId:           n.CateId,
				CateName:         n.CateName,
				Author:           n.Author,
				IsOriginal:       n.IsOriginal,
				IsHot:            n.IsHot,
				IsRec:            n.IsRec,
				IsTodayRec:       n.IsTodayRec,
				IsVipRec:         n.IsVipRec,
				IsVipReward:      n.IsVipReward,
				IsVipUp:          n.IsVipUp,
				IsSignNewBook:    n.IsSignNewBook,
				IsCollect:        n.IsCollect,
				Status:           n.Status,
				Views:            n.Views,
				CollectNum:       n.CollectNum,
				RecNum:           n.RecNum,
				TextNum:          n.TextNum,
				ChapterNum:       n.ChapterNum,
				ChapterUpdatedAt: n.ChapterUpdatedAt,
				ChapterId:        n.ChapterId,
				ChapterTitle:     n.ChapterTitle,
				CreatedAt:        n.CreatedAt,
				UpdatedAt:        n.UpdatedAt,
				DeletedAt:        n.DeletedAt,
			}
			novlist = append(novlist, nov)
		}
		str, _ := json.Marshal(novlist)
		gutils.SetRedisKeyValue(buf.String(), string(str))
		return novs
	}
}

// 获取排行榜小说列表
func (this *Novel) GetRanks(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.OrderBy = "-views"
	args.Fields = []string{"id", "hash_key", "name", "cover", "desc", "views", "author", "cate_id", "cate_name", "chapter_id", "chapter_title", "updated_at", "chapter_updated_at"}
	//args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	var buf bytes.Buffer
	buf.WriteString("GetRanks:")
	buf.WriteString("size_" + strconv.Itoa(size))
	buf.WriteString(":offset_" + strconv.Itoa(offset))
	res := gutils.GetRedisKeys(buf.String())
	if res != nil {
		novs := make([]*models.Novel, 0)
		str := fmt.Sprintf("%s", res)
		if err := json.Unmarshal([]byte(str), &novs); err != nil {
			fmt.Println("err:-----", err.Error())
			return nil
		}
		return novs
	} else {
		novs, _ := this.GetAll(args)
		var novlist []models.Novel
		for _, n := range novs {
			var nov = models.Novel{
				Id:               n.Id,
				HashKey:          n.HashKey,
				Name:             n.Name,
				Desc:             n.Desc,
				Cover:            n.Cover,
				CateId:           n.CateId,
				CateName:         n.CateName,
				Author:           n.Author,
				IsOriginal:       n.IsOriginal,
				IsHot:            n.IsHot,
				IsRec:            n.IsRec,
				IsTodayRec:       n.IsTodayRec,
				IsVipRec:         n.IsVipRec,
				IsVipReward:      n.IsVipReward,
				IsVipUp:          n.IsVipUp,
				IsSignNewBook:    n.IsSignNewBook,
				IsCollect:        n.IsCollect,
				Status:           n.Status,
				Views:            n.Views,
				CollectNum:       n.CollectNum,
				RecNum:           n.RecNum,
				TextNum:          n.TextNum,
				ChapterNum:       n.ChapterNum,
				ChapterUpdatedAt: n.ChapterUpdatedAt,
				ChapterId:        n.ChapterId,
				ChapterTitle:     n.ChapterTitle,
				CreatedAt:        n.CreatedAt,
				UpdatedAt:        n.UpdatedAt,
				DeletedAt:        n.DeletedAt,
			}
			novlist = append(novlist, nov)
		}
		str, _ := json.Marshal(novlist)
		gutils.SetRedisKeyValue(buf.String(), string(str))
		return novs
	}
}

// 获取同类推荐排行榜
func (this *Novel) GetCateRanks(cateId, size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.OrderBy = "-views"
	args.Fields = []string{"id", "hash_key", "name", "cover", "desc", "views", "author", "cate_id", "cate_name"}
	args.FilterMaps = map[string]int{
		"cate_id": cateId,
	}
	//args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	novs, _ := this.GetAll(args)

	return novs
}

// 获取最新更新小说列表
func (this *Novel) GetNewUps(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.OrderBy = "-chapter_updated_at"
	args.Fields = []string{"id", "hash_key", "name", "cover", "desc", "views", "author", "cate_id", "cate_name", "chapter_id", "chapter_title", "updated_at", "chapter_updated_at"}
	//args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	var buf bytes.Buffer
	buf.WriteString("GetNewUps:")
	buf.WriteString("size_" + strconv.Itoa(size))
	buf.WriteString(":offset_" + strconv.Itoa(offset))
	res := gutils.GetRedisKeys(buf.String())
	if res != nil {
		novs := make([]*models.Novel, 0)
		str := fmt.Sprintf("%s", res)
		if err := json.Unmarshal([]byte(str), &novs); err != nil {
			fmt.Println("err:-----", err.Error())
			return nil
		}
		return novs
	} else {
		novs, _ := this.GetAll(args)
		var novlist []models.Novel
		for _, n := range novs {
			var nov = models.Novel{
				Id:               n.Id,
				HashKey:          n.HashKey,
				Name:             n.Name,
				Desc:             n.Desc,
				Cover:            n.Cover,
				CateId:           n.CateId,
				CateName:         n.CateName,
				Author:           n.Author,
				IsOriginal:       n.IsOriginal,
				IsHot:            n.IsHot,
				IsRec:            n.IsRec,
				IsTodayRec:       n.IsTodayRec,
				IsVipRec:         n.IsVipRec,
				IsVipReward:      n.IsVipReward,
				IsVipUp:          n.IsVipUp,
				IsSignNewBook:    n.IsSignNewBook,
				IsCollect:        n.IsCollect,
				Status:           n.Status,
				Views:            n.Views,
				CollectNum:       n.CollectNum,
				RecNum:           n.RecNum,
				TextNum:          n.TextNum,
				ChapterNum:       n.ChapterNum,
				ChapterUpdatedAt: n.ChapterUpdatedAt,
				ChapterId:        n.ChapterId,
				ChapterTitle:     n.ChapterTitle,
				CreatedAt:        n.CreatedAt,
				UpdatedAt:        n.UpdatedAt,
				DeletedAt:        n.DeletedAt,
			}
			novlist = append(novlist, nov)
		}
		str, _ := json.Marshal(novlist)
		gutils.SetRedisKeyValue(buf.String(), string(str))
		return novs
	}
}

// 获取新增小说榜单
func (this *Novel) GetNews(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.Fields = []string{"id", "hash_key", "name", "cover", "desc", "views", "author", "cate_id", "cate_name", "created_at"}
	args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	var buf bytes.Buffer
	buf.WriteString("GetNews:size_")
	buf.WriteString(strconv.Itoa(size))
	buf.WriteString(":offset_" + strconv.Itoa(offset))
	res := gutils.GetRedisKeys(buf.String())
	if res != nil {
		novs := make([]*models.Novel, 0)
		str := fmt.Sprintf("%s", res)
		if err := json.Unmarshal([]byte(str), &novs); err != nil {
			fmt.Println("err:-----", err.Error())
		}
		return novs
	} else {
		var novlist []models.Novel
		novs, _ := this.GetAll(args)
		for _, n := range novs {
			var nov = models.Novel{
				Id:               n.Id,
				HashKey:          n.HashKey,
				Name:             n.Name,
				Desc:             n.Desc,
				Cover:            n.Cover,
				CateId:           n.CateId,
				CateName:         n.CateName,
				Author:           n.Author,
				IsOriginal:       n.IsOriginal,
				IsHot:            n.IsHot,
				IsRec:            n.IsRec,
				IsTodayRec:       n.IsTodayRec,
				IsVipRec:         n.IsVipRec,
				IsVipReward:      n.IsVipReward,
				IsVipUp:          n.IsVipUp,
				IsSignNewBook:    n.IsSignNewBook,
				IsCollect:        n.IsCollect,
				Status:           n.Status,
				Views:            n.Views,
				CollectNum:       n.CollectNum,
				RecNum:           n.RecNum,
				TextNum:          n.TextNum,
				ChapterNum:       n.ChapterNum,
				ChapterUpdatedAt: n.ChapterUpdatedAt,
				ChapterId:        n.ChapterId,
				ChapterTitle:     n.ChapterTitle,
				CreatedAt:        n.CreatedAt,
				UpdatedAt:        n.UpdatedAt,
				DeletedAt:        n.DeletedAt,
			}
			novlist = append(novlist, nov)
		}
		str, _ := json.Marshal(novlist)
		gutils.SetRedisKeyValue(buf.String(), string(str))
		return novs
	}
}

// 获取VIP打赏
func (this *Novel) GetVipRewards(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.Fields = []string{"id", "hash_key", "name", "cover", "desc", "views", "author", "cate_id", "cate_name"}
	args.FilterMaps = map[string]int{
		"is_vip_reward": 1,
	}
	args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	novs, _ := this.GetAll(args)

	return novs
}

// 获取VIP更新
func (this *Novel) GetVipUps(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.Fields = []string{"id", "hash_key", "name", "cover", "desc", "views", "author", "cate_id", "cate_name"}
	args.FilterMaps = map[string]int{
		"is_vip_up": 1,
	}
	args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	novs, _ := this.GetAll(args)

	return novs
}

// 获取完本小说列表
func (this *Novel) GetEnds(size, offset int) []*models.Novel {
	args := models.ArgsNovelList{}
	args.Limit = size
	args.Offset = offset
	args.Fields = []string{"id", "hash_key", "name", "cover", "desc", "views", "author", "cate_id", "cate_name"}
	args.FilterMaps = map[string]int{
		"status": int(models.BOOKFINISH),
	}

	novs, _ := this.GetAll(args)

	return novs
}

// 批量获取小说列表
// 用于前台查询
func (this *Novel) GetList(size, offset int, qs map[string]interface{}) ([]*models.Novel, int64) {
	args := models.ArgsNovelList{FilterMaps: make(map[string]int, 0)}

	// 字数查询
	if textNum, ok := qs["text_num"]; ok {
		switch textNum.(int) {
		case 1:
			args.StartTextNum = 0
			args.EndTextNum = 300000
		case 2:
			args.StartTextNum = 300000
			args.EndTextNum = 500000
		case 3:
			args.StartTextNum = 500000
			args.EndTextNum = 1000000
		case 4:
			args.StartTextNum = 1000000
			args.EndTextNum = 2000000
		case 5:
			args.StartTextNum = 2000000
			args.EndTextNum = 0
		default:
		}
	}

	// 更新时间
	if upTime, ok := qs["uptime"]; ok {
		switch upTime.(int) {
		case 1:
			args.MaxChapterUpdatedAt = time.Now().AddDate(0, 0, -3).Unix()
		case 2:
			args.MaxChapterUpdatedAt = time.Now().AddDate(0, 0, -7).Unix()
		case 3:
			args.MaxChapterUpdatedAt = time.Now().AddDate(0, 0, -15).Unix()
		case 4:
			args.MaxChapterUpdatedAt = time.Now().AddDate(0, 0, -30).Unix()
		default:
		}
	}

	// 排序
	if ot, ok := qs["ot"]; ok {
		switch ot.(int) {
		case 1:
			args.ArgsBase.OrderBy = "-views"
		case 2:
			args.ArgsBase.OrderBy = "-chapter_updated_at"
		case 3:
			args.ArgsBase.OrderBy = "-text_num"
		default:
		}
	}

	if count, ok := qs["count"]; ok {
		args.Count = count.(bool)
	}

	args.Limit = size
	args.Offset = offset

	if kw, ok := qs["q"]; ok && len(kw.(string)) > 0 {
		args.Keyword = kw.(string)
	}

	if cateId, ok := qs["cate_id"]; ok && cateId.(int) > 0 {
		args.FilterMaps["cate_id"] = cateId.(int)
	}

	if status, ok := qs["status"]; ok && status.(int) > 0 {
		args.FilterMaps["status"] = status.(int)
	}

	args.Fields = []string{"id", "hash_key", "name", "cover", "desc", "author", "cate_id", "cate_name", "text_num", "status", "cate_id", "chapter_title"}
	//args.OrderBy = "un_complete_chap_num"
	cons := make(map[string]interface{})
	cons["key"] = "chapter_title"
	cons["value"] = ""
	args.Exclude = cons

	novs, count := this.GetAll(args)

	// 关键词替换
	if len(novs) > 0 {
		// 名称搜索
		if args.Keyword != "" {
			for k, v := range novs {
				novs[k].Name = strings.Replace(v.Name, args.Keyword, fmt.Sprintf(`<font color="red">%s</font>`, args.Keyword), -1)
			}
		}
	}

	return novs, count
}

// 批量获取小说列表
func (this *Novel) GetAll(args models.ArgsNovelList) ([]*models.Novel, int64) {
	novs, count := models.NovelModel.GetAll(args)

	return novs, count
}

// 批量删除小说
func (this *Novel) DeleteBatch(ids []string) error {
	if len(ids) == 0 {
		return errors.New("params error")
	}

	err := models.NovelModel.DeleteBatch(ids)
	if err != nil {
		return err
	}

	// 删除小说章节列表
	for _, id := range ids {
		u64, err := strconv.ParseUint(id, 10, 32)
		if err == nil {
			novId := uint32(u64)
			ChapterService.DelByNovId(novId)

			// 删除小说章节采集点
			models.NovelLinksModel.DelByNovId(novId)
		}
	}

	return nil
}

// 删除小说
func (this *Novel) Delete(id uint32) error {
	if id < 0 {
		return errors.New("params error")
	}

	m := models.Novel{Id: id}
	err := m.Delete(true)
	if err != nil {
		return err
	}

	// 删除小说章节列表
	ChapterService.DelByNovId(id)

	// 删除小说章节采集点
	models.NovelLinksModel.DelByNovId(id)

	return nil
}

// 清空小说章节列表
func (this *Novel) DelChapters(id uint32) error {
	if id < 0 {
		return errors.New("params error")
	}

	nov := this.Get(id)
	if nov == nil {
		return errors.New("获取小说失败")
	}

	// 删除章节列表
	err := ChapterService.DelByNovId(id)
	if err != nil {
		//return err
	}

	nov.ChapterNum = 0
	nov.ChapterId = 0
	nov.TextNum = 0
	nov.ChapterTitle = ""
	nov.ChapterUpdatedAt = 0

	nov.Update("chapter_num", "chapter_id", "text_num", "chapter_title", "chapter_updated_at")

	return nil
}

// 删除小说采集站点
func (this *Novel) DelLink(id uint32) error {
	if id < 0 {
		return errors.New("params error")
	}

	m := models.NovelLinks{Id: id}
	err := m.Delete(true)
	if err != nil {
		return err
	}

	return nil
}

// 更新浏览次数
func (this *Novel) UpViews(novId uint32) {
	nov := this.Get(novId)
	if nov == nil {
		return
	}

	nov.Views++

	nov.Update("views")
}

// 更新小说连载状态
func (this *Novel) UpStatus(novId uint32, status uint8) error {
	nov := this.Get(novId)
	if nov == nil {
		return errors.New("小说不存在")
	}

	nov.Status = status

	return nov.Update("status")
}

// 修改小说简介信息
func (this *Novel) UpNovelInfo(nov *models.Novel) error {
	// 获取分类名称
	cate := CateService.Get(nov.CateId)
	if cate == nil {
		return errors.New("分类不存在或者已被删除")
	}

	nov.CateName = cate.Name

	return nov.Update("desc", "cover", "cate_id", "cate_name", "author")
}

// 修改小说文字数
func (this *Novel) UpChapterTextNum(novId uint32, novTextNum int, isAdded bool) error {
	nov := this.Get(novId)
	if nov == nil {
		return errors.New("小说不存在")
	}

	if isAdded {
		nov.TextNum += uint32(novTextNum)
	} else {
		nov.TextNum -= uint32(novTextNum)
	}

	return nov.Update("text_num")
}

// 批量更新推荐
func (this *Novel) UpRecBatch(field string, books []string) error {
	if len(books) == 0 {
		return errors.New("参数错误")
	}

	return models.NovelModel.UpRecBatch(field, books)
}

// 修改章节信息
func (this *Novel) UpChapterInfo(novId uint32, novTextNum, chapterNum int, chapterId uint64, chapterTitle string, status uint8, unCompleteChapNum int, unCompleteChapid string) error {
	nov := this.Get(novId)
	if nov == nil {
		return errors.New("小说不存在")
	}

	nov.ChapterNum += uint32(chapterNum)
	nov.TextNum += uint32(novTextNum)
	nov.ChapterId = chapterId
	nov.ChapterTitle = chapterTitle
	nov.ChapterUpdatedAt = uint32(time.Now().Unix())
	if unCompleteChapNum > 0 {
		nov.UnCompleteChapNum = unCompleteChapNum
	}
	if unCompleteChapid != "" {
		nov.UnCompleteChapid = unCompleteChapid
	}

	if status > 0 {
		nov.Status = status
	}

	return nov.Update("text_num", "chapter_num", "chapter_id", "chapter_title", "chapter_updated_at", "status")
}

// 添加/修改
func (this *Novel) Save(novel *models.Novel) error {
	// 参数校验
	valid := validation.Validation{}
	valid.Required(novel.Name, "nameEmpty").Message("小说名称不能为空")
	valid.MaxSize(novel.Name, 100, "nameMax").Message("小说名称长度不能超过100个字符")
	valid.Required(novel.Author, "authorEmpty").Message("小说作者不能为空")
	valid.MaxSize(novel.Author, 20, "authorMax").Message("作者名称不能超过20个字符")
	valid.Required(novel.CateId, "cidEmpty").Message("小说分类不能为空")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			panic(fmt.Sprintf("Update chapter info: %s", err.Error()))
			return err
		}
	}

	// 获取分类名称
	cate := CateService.Get(novel.CateId)
	if cate == nil {
		return errors.New("分类不存在或者已被删除")
	}

	novel.CateName = cate.Name

	var err error
	if novel.Id > 0 {
		err = novel.Update("name", "desc", "cover", "cate_id", "cate_name", "author", "is_original", "is_hot", "is_rec", "is_vip_rec", "is_today_rec", "status", "is_sign_new_book", "is_collect", "is_vip_reward", "is_vip_up")
	} else {
		novel.Status = models.BOOKOPEN
		err = novel.Insert()
	}

	if err == nil && novel.IsOriginal == 0 {
		manager.AddTask(novel.Id)
	}

	return err
}
