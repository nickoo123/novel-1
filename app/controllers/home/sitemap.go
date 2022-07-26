package home

import (
	"fmt"
	"novel/app/services"
	"novel/app/utils/sitemap"
	"time"
)

type SitemapController struct {
	BaseController
}

// 首页
func (this *SitemapController) Index() {
	st := sitemap.NewSiteMap()
	st.SetPretty(true)

	services.NovelService.GetNewUps(5000, 0)
	url := sitemap.NewUrl()
	url.SetLoc("https://www.biqugesk.cc/book/detail?id=28999&novid=2918")
	url.SetLastmod(time.Now())
	url.SetChangefreq("daily")
	url.SetPriority(1)
	st.AppendUrl(url)
	bt, err := st.ToXml()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", bt)
}
