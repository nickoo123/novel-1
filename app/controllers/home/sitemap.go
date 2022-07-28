package home

import (
	"bytes"
	"fmt"
	"novel/app/services"
	"novel/app/utils/sitemap"
	"strconv"
	"time"
)

type SitemapController struct {
	BaseController
}

// 首页
func (this *SitemapController) Index() {
	page := 1
	st := sitemap.NewSiteMap()
	st.SetPretty(true)

	novs := services.NovelService.GetNews(5000, page)
	for _, nov := range novs {
		var buf bytes.Buffer
		buf.WriteString("https://www.biqugesk.cc/book/index?id=")
		buf.WriteString(fmt.Sprintf("%v", nov.Id))
		url := sitemap.NewUrl()
		url.SetLoc(buf.String())
		ins, _ := strconv.ParseInt(fmt.Sprintf("%v", nov.CreatedAt), 10, 64)
		te := time.Unix(ins, 0)
		url.SetLastmod(te)
		url.SetChangefreq("daily")
		url.SetPriority(1)
		st.AppendUrl(url)
	}
	this.Data["xml"] = st
	this.ServeXML()
	this.TplName = "xml.tpl"
}
