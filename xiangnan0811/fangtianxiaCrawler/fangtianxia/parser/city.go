package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/xiangnan0811/fangtianxiaCrawler/engine"
	"github.com/xiangnan0811/fangtianxiaCrawler/utils"

	jsoniter "github.com/json-iterator/go"

	"github.com/antchfx/htmlquery"
)

var erShouNextUrlRe = regexp.MustCompile(`<a href="(/house/i\d+/)">下一页</a>`)
var erShouMoreUrlRe = regexp.MustCompile(`<a href="(/house-a\d+/)">`)

func ParseCityNewHouseList(contents []byte, province string, city string, detailUrl string) engine.ParseResult {
	result := engine.ParseResult{}
	root, _ := htmlquery.Parse(strings.NewReader(string(contents[:])))
	// 下一页
	nextUrlNode := htmlquery.Find(root, "//a[@class='next']")
	if nextUrlNode != nil {
		nextUrlString := htmlquery.SelectAttr(nextUrlNode[0], "href")
		if nextUrlString != "" {
			nextUrl := detailUrl + nextUrlString
			result.Requests = append(result.Requests, engine.Request{
				Url:        nextUrl,
				ParserFunc: ParseCityNewHouseListFunc(province, city, nextUrl),
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
					Url:        moreUrl,
					ParserFunc: ParseCityNewHouseListFunc(province, city, moreUrl),
				})
			}
		}
	}

	lis := htmlquery.Find(root, "//div[@id='newhouse_loupai_list']//li[@id]")
	if len(lis) > 1 {
		for _, li := range lis {
			id := 0
			idString := strings.Split(htmlquery.SelectAttr(li, "id"), "_")[1]
			id, _ = strconv.Atoi(idString)

			var url string
			urlTextNode := htmlquery.Find(li, ".//div[@class='nlcd_name']/a")
			if urlTextNode != nil {
				urlText := htmlquery.SelectAttr(urlTextNode[0], "href")
				url = "https:" + urlText + "house/" + idString + "/housedetail.htm"
			}

			if url != "" {
				idParam := id
				urlParam := url
				provinceParam := province
				cityParam := city
				result.Requests = append(result.Requests, engine.Request{
					Url:        url,
					ParserFunc: ParseNewHouseFunc(provinceParam, cityParam, idParam, urlParam),
				})
			}
		}
	}
	return result
}

func ParseCityErShouHouseList(contents []byte, province string, city string, ershouurl string) engine.ParseResult {
	result := engine.ParseResult{}

	// 下一页
	nextUrlString := utils.ExtractString(contents, erShouNextUrlRe)
	if nextUrlString != "" && nextUrlString != "暂无资料" {
		nextUrl := ershouurl + nextUrlString
		result.Requests = append(result.Requests, engine.Request{
			Url:        nextUrl,
			ParserFunc: ParseCityErShouHouseListFunc(province, city, nextUrl),
		})
	}
	// 区县链接
	moreUrlList := utils.ExtractAll(contents, erShouMoreUrlRe)
	if len(moreUrlList) > 1 {
		for _, moreUrlString := range moreUrlList {
			if moreUrlString != "" && moreUrlString != "暂无资料" {
				moreUrl := ershouurl + moreUrlString
				result.Requests = append(result.Requests, engine.Request{
					Url:        moreUrl,
					ParserFunc: ParseCityErShouHouseListFunc(province, city, moreUrl),
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
			id, _ := strconv.Atoi(data["houseid"].(string))

			var url string
			urlTextNode := htmlquery.Find(dl, ".//h4/a")
			if urlTextNode != nil {
				urlText := htmlquery.SelectAttr(urlTextNode[0], "href")
				url = ershouurl + urlText
			}

			if url != "" {
				idParam := id
				urlParam := url
				provinceParam := province
				cityParam := city
				result.Requests = append(result.Requests, engine.Request{
					Url:        url,
					ParserFunc: ParseErShouHouseFunc(provinceParam, cityParam, idParam, urlParam),
				})
			}
		}
	}
	return result
}

func ParseErShouHouseFunc(province string, city string, id int, url string) engine.ParserFunc {
	return func(c []byte) engine.ParseResult {
		return ParseErShouHouse(c, province, city, id, url)
	}
}

func ParseNewHouseFunc(province string, city string, id int, url string) engine.ParserFunc {
	return func(c []byte) engine.ParseResult {
		return ParseNewHouse(c, province, city, id, url)
	}
}
