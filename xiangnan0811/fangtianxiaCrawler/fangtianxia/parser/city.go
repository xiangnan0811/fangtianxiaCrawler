package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"

	"github.com/xiangnan0811/fangtianxiaCrawler/engine"
	"github.com/xiangnan0811/fangtianxiaCrawler/utils"

	jsoniter "github.com/json-iterator/go"

	"github.com/antchfx/htmlquery"
)

var erShouNextUrlRe = regexp.MustCompile(`<a href="(/house/i\d+/)">下一页</a>`)
var erShouMoreUrlRe = regexp.MustCompile(`<a href="(/house-a\d+/)">`)

func ParseCityNewHouseList(contents []byte, province string, detailUrl string) engine.ParseResult {
	result := engine.ParseResult{}
	root, _ := htmlquery.Parse(strings.NewReader(string(contents[:])))
	// 下一页
	nextUrlNode := htmlquery.Find(root, "//a[@class='next']")
	if nextUrlNode != nil {
		nextUrlString := htmlquery.SelectAttr(nextUrlNode[0], "href")
		if nextUrlString != "" {
			nextUrl := detailUrl + nextUrlString
			result.Requests = append(result.Requests, engine.Request{
				Url:    nextUrl,
				Parser: NewCityNewHouseListParser(province, nextUrl),
			})
		}
	}

	// 区县链接
	moreUrlNodes := htmlquery.Find(root, "//li[@id='quyu_name']/a[not(@id)]")
	if len(moreUrlNodes) > 1 {
		for _, moreUrlNode := range moreUrlNodes {
			moreUrlString := htmlquery.SelectAttr(moreUrlNode, "href")
			if moreUrlString != "" {
				moreUrl := detailUrl + moreUrlString
				result.Requests = append(result.Requests, engine.Request{
					Url:    moreUrl,
					Parser: NewCityNewHouseListParser(province, moreUrl),
				})
			}
		}
	}

	lis := htmlquery.Find(root, "//div[@id='newhouse_loupai_list']//li[@id]")
	if len(lis) > 1 {
		for _, li := range lis {
			idString := strings.Split(htmlquery.SelectAttr(li, "id"), "_")[1]

			var url string
			urlTextNode := htmlquery.Find(li, ".//div[@class='nlcd_name']/a")
			if urlTextNode != nil {
				urlText := htmlquery.SelectAttr(urlTextNode[0], "href")
				url = "https:" + urlText + "house/" + idString + "/housedetail.htm"
			}

			if url != "" {
				urlParam := url
				provinceParam := province
				result.Requests = append(result.Requests, engine.Request{
					Url:    url,
					Parser: NewNewHouseParser(provinceParam, urlParam),
				})
			}
		}
	}
	return result
}

func ParseCityErShouHouseList(contents []byte, province string, ershouurl string) engine.ParseResult {
	result := engine.ParseResult{}

	// 下一页
	nextUrlString := utils.ExtractString(contents, erShouNextUrlRe)
	if nextUrlString != "" && nextUrlString != "暂无资料" {
		nextUrl := ershouurl + nextUrlString
		result.Requests = append(result.Requests, engine.Request{
			Url:    nextUrl,
			Parser: NewCityErShouHouseListParser(province, nextUrl),
		})
	}
	// 区县链接
	moreUrlList := utils.ExtractAll(contents, erShouMoreUrlRe)
	if len(moreUrlList) > 1 {
		for _, moreUrlString := range moreUrlList {
			if moreUrlString != "" && moreUrlString != "暂无资料" {
				moreUrl := ershouurl + moreUrlString
				result.Requests = append(result.Requests, engine.Request{
					Url:    moreUrl,
					Parser: NewCityErShouHouseListParser(province, moreUrl),
				})
			}
		}
	}

	root, _ := htmlquery.Parse(strings.NewReader(string(contents[:])))
	dls := htmlquery.Find(root, "//dl[@dataflag='bg']")
	if len(dls) > 1 {
		for _, dl := range dls {
			dataBg := htmlquery.SelectAttr(dl, "data-bg")
			data := make(map[string]interface{})
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			err := json.Unmarshal([]byte(dataBg), &data)
			if err != nil {
				fmt.Println("json unmarshal error ", err)
			}

			var url string
			urlTextNode := htmlquery.Find(dl, ".//h4/a")
			if urlTextNode != nil {
				urlText := htmlquery.SelectAttr(urlTextNode[0], "href")
				url = ershouurl + urlText
			}

			if url != "" {
				urlParam := url
				provinceParam := province
				result.Requests = append(result.Requests, engine.Request{
					Url:    url,
					Parser: NewErShouHouseParser(provinceParam, urlParam),
				})
			}
		}
	}
	return result
}

type NewHouseParser struct {
	province string
	url      string
}

func (n *NewHouseParser) Parse(contents []byte) engine.ParseResult {
	return ParseNewHouse(contents, n.province, n.url)
}

func (n *NewHouseParser) Serialize() (name string, province string, url string) {
	return config.NewHouseParser, n.province, n.url
}

func NewNewHouseParser(province string, url string) *NewHouseParser {
	return &NewHouseParser{
		province: province,
		url:      url,
	}
}

type ErShouHouseParser struct {
	province string
	url      string
}

func (e *ErShouHouseParser) Parse(contents []byte) engine.ParseResult {
	return ParseErShouHouse(contents, e.province, e.url)
}

func (e *ErShouHouseParser) Serialize() (name string, province string, url string) {
	return config.ErShouHouseParser, e.province, e.url
}

func NewErShouHouseParser(province string, url string) *ErShouHouseParser {
	return &ErShouHouseParser{
		province: province,
		url:      url,
	}
}
