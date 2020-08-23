package parser

import (
	"strings"

	"github.com/xiangnan0811/fangtianxiaCrawler/engine"

	"github.com/antchfx/htmlquery"
)

func ParseCityList(contents []byte) engine.ParseResult {
	result := engine.ParseResult{}
	root, _ := htmlquery.Parse(strings.NewReader(string(contents[:])))
	trs := htmlquery.Find(root, "//div[@id='c02']//tr")
	trs = trs[:len(trs)-1] // 过滤海外城市
	province := ""

	for _, tr := range trs {
		tds := htmlquery.Find(tr, ".//td[not(@class)]")
		// 获取省份信息
		provinceNode := htmlquery.Find(tds[0], "./strong")
		if len(provinceNode) != 0 {
			// 某省城市多的话第二行没有省份信息，故而使用上一次的省份信息
			provinceText := strings.TrimSpace(htmlquery.InnerText(provinceNode[0]))
			if provinceText != "" {
				province = provinceText
			}
		}

		// 获取城市信息
		cityLinks := htmlquery.Find(tds[1], ".//a")
		for _, cityLink := range cityLinks {
			city := htmlquery.InnerText(cityLink)
			// 设置直辖市的省份信息
			if city == "北京" || city == "上海" || city == "天津" || city == "重庆" {
				province = city
			}
			cityUrl := htmlquery.SelectAttr(cityLink, "href")

			// 构建新房、二手房的url链接
			urlModule := strings.Split(cityUrl, "//")
			var newHouseUrl, erShouHouseUrl string
			// 北京新房二手房链接特殊些，其余均为特定格式
			if urlModule[1] == "bj.fang.com/" {
				newHouseUrl = "https://newhouse.fang.com/house/s"
				erShouHouseUrl = "https://esf.fang.com"
			} else if urlModule[1] == "liupanshui.fang.com/" {
				newHouseUrl = "https://lps.newhouse.fang.com/house/s"
				erShouHouseUrl = "https://lps.esf.fang.com"
			} else if urlModule[1] == "macau.fang.com/" || urlModule[1] == "hk.esf.fang.com/" {
				newHouseUrl = ""
				erShouHouseUrl = ""
			} else if urlModule[1] == "wenchang.fang.com/" || urlModule[1] == "gaoyang.fang.com/" || urlModule[1] == "baoying.fang.com/" {
				newHouseUrl = ""
				erShouHouseUrl = "https://wenchang.esf.fang.com"
			} else if urlModule[1] == "qionghai.fang.com/" || urlModule[1] == "zhaozhou.fang.com/" || urlModule[1] == "kangping.fang.com/" {
				newHouseUrl = ""
				erShouHouseUrl = "https://wenchang.esf.fang.com"
			} else {
				urlCity := strings.Split(urlModule[1], ".")[0]
				newHouseUrl = "https://" + urlCity + ".newhouse.fang.com/house/s"
				erShouHouseUrl = "https://" + urlCity + ".esf.fang.com"
			}

			provinceParam := province
			cityParam := city
			newHouseUrlParam := newHouseUrl
			if newHouseUrl != "" {
				result.Requests = append(result.Requests, engine.Request{
					Url:        newHouseUrl,
					ParserFunc: ParseCityNewHouseListFunc(provinceParam, cityParam, newHouseUrlParam),
				})
			}

			erShouUrl := erShouHouseUrl
			if erShouHouseUrl != "" {
				result.Requests = append(result.Requests, engine.Request{
					Url:        erShouHouseUrl,
					ParserFunc: ParseCityErShouHouseListFunc(provinceParam, cityParam, erShouUrl),
				})
			}
		}
	}
	return result
}

func ParseCityNewHouseListFunc(province string, city string, url string) engine.ParserFunc {
	return func(c []byte) engine.ParseResult {
		return ParseCityNewHouseList(c, province, city, url)
	}
}

func ParseCityErShouHouseListFunc(province string, city string, url string) engine.ParserFunc {
	return func(c []byte) engine.ParseResult {
		return ParseCityErShouHouseList(c, province, city, url)
	}
}
