package parser

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xiangnan0811/fangtianxiaCrawler/engine"
	"github.com/xiangnan0811/fangtianxiaCrawler/model"
	"github.com/xiangnan0811/fangtianxiaCrawler/utils"
)

var (
	nameRe                      = regexp.MustCompile(`<h1><a.*?target="_blank">(.*?)</a></h1>`)
	scoreRe                     = regexp.MustCompile(`<span style="margin-right: 5px;">([\d.]+)</span>`)
	tagsRe                      = regexp.MustCompile(`<span class="tag">(.*?)</span>`)
	propertyRightsRe            = regexp.MustCompile(`<p style="width: 130px;float: left;">(.*?)</p>`)
	propertyRightsRe2           = regexp.MustCompile(`产权年限：</div>\s*.*?">(.*?)\s*</div>`)
	propertyTypeRe              = regexp.MustCompile(`物业类别：</div>[^>]+title="(.*?)">`)
	avgPriceRe                  = regexp.MustCompile(`价.*?格：(.|\n)*?[em]?>\s*(.*?)\s*</em>`)
	decorationRe                = regexp.MustCompile(`装修状况：</div>[^>]+right">((.|\n)*?)</div>`)
	propertyDeveloperRe         = regexp.MustCompile(`开.*?发<i\s*.*?商：\s*</div>\s*.*?target="_blank">(.*?)</a></div>`)
	addressRe                   = regexp.MustCompile(`楼盘地址：</div>\s*.*?">\s*(.*?)\s*?<`)
	loopPositionRe              = regexp.MustCompile(`环线位置：</div>\s*.*?">\s*(.*?)\s*</div>`)
	houseTypeRe                 = regexp.MustCompile(`<a href="//[a-z0-9A-Z]+.fang.com//photo/list_\d+_\d+.htm">(.*?)</a>`)
	saleAddressRe               = regexp.MustCompile(`售楼地址：</div>[^>]+title="(.*?)">`)
	openTimeRe                  = regexp.MustCompile(`开盘时间：</div>\s*.*?">(.*?)<`)
	saleStatusRe                = regexp.MustCompile(`销售状态：</div>\s*.*?">\s*(.*?)</div>`)
	discountRe                  = regexp.MustCompile(`楼盘优惠：</div>\s*.*?">\s*(.*?)\s*</div>`)
	deliveryTimeRe              = regexp.MustCompile(`交房时间：</div>\s*.*?">(.*?)<`)
	consultingPhoneRe           = regexp.MustCompile(`咨询电话：</div>\s*.*?">\s*(\d+)\s*</div>`)
	transportationRe            = regexp.MustCompile(`<h3>交通状况</h3>\s*(.*?)\s*</div>`)
	transportationRe2           = regexp.MustCompile(`<span>交通</span>\s*(.*?)\s*</li>`)
	communityAreaRe             = regexp.MustCompile(`占地面积：</div>\s*.*?">([.\d]+)平方米\s*</div>`)
	communityBuildingAreaRe     = regexp.MustCompile(`建筑面积：</div>\s*.*?">([.\d]+)平方米\s*</div>`)
	plotRatioRe                 = regexp.MustCompile(`容.*?积<i\s*.*?率：\s*</div>\s*.*?">([.\d]+).*?\s*</div>`)
	greeningRate                = regexp.MustCompile(`绿.*?化<i\s*.*?率：\s*</div>\s*.*?">([.\d]+).*?\s*</div>`)
	propertyManagementCompanyRe = regexp.MustCompile(`物业公司：</div>\s*.*?blank">\s*(.*?)\s*</a>`)
	propertyFeeRe               = regexp.MustCompile(`物.*?业.*?费：</div>\s*.*?">([.\d]+).*?\s*</div>`)
	projectBriefRe              = regexp.MustCompile(`项目简介</h3>\s*.*?">\s*((.|\n)*?)\s*</p>`)
)

func ParseNewHouse(contents []byte, province string, city string, id int, detailUrl string) engine.ParseResult {
	result := engine.ParseResult{}

	s := strings.Split(detailUrl, "?rfss=")
	if len(s) > 1 {
		return result
	}

	if redirect := redirectRe.FindAllSubmatch(contents, -1); redirect != nil {
		rfss := utils.ExtractAll(contents, rfssRe)
		if lenRfss := len(rfss); lenRfss >= 2 {
			newUrl := detailUrl + "?rfss=" + rfss[lenRfss-2]
			result.Requests = append(result.Requests, engine.Request{
				Url:        newUrl,
				ParserFunc: ParseNewHouseFunc(province, city, id, newUrl),
			})
			return result
		}
	} else if verify := VerifyRe.FindAllSubmatch(contents, -1); verify != nil {
		newUrl := detailUrl + "?rfss=2-0-1"
		result.Requests = append(result.Requests, engine.Request{
			Url:        newUrl,
			ParserFunc: ParseNewHouseFunc(province, city, id, newUrl),
		})
		return result
	} else {
		newHouseProfile := model.NewHouseProfile{}
		newHouseProfile.Name = utils.ExtractString(contents, nameRe)
		newHouseProfile.Score = utils.ExtractFloat64(contents, scoreRe)
		newHouseProfile.Tags = utils.ExtractAll(contents, tagsRe)

		avgPriceMatches := avgPriceRe.FindSubmatch(contents)
		if len(avgPriceMatches) >= 3 {
			avgPrice := strings.Split(string(avgPriceMatches[2]), ">")
			if len(avgPrice) >= 1 {
				newHouseProfile.AvgPrice = avgPrice[len(avgPrice)-1]
			}
		}

		newHouseProfile.HouseType = utils.ExtractAll(contents, houseTypeRe)
		newHouseProfile.PropertyType = utils.ExtractString(contents, propertyTypeRe)
		propertyRights := utils.ExtractAll(contents, propertyRightsRe)
		if propertyRights == nil {
			propertyRights = utils.ExtractAll(contents, propertyRightsRe2)
		}
		newHouseProfile.PropertyRights = propertyRights
		newHouseProfile.SaleAddress = utils.ExtractString(contents, saleAddressRe)
		newHouseProfile.OpenTime = utils.ExtractString(contents, openTimeRe)
		newHouseProfile.PropertyDeveloper = utils.ExtractString(contents, propertyDeveloperRe)
		newHouseProfile.Decoration = utils.ExtractString(contents, decorationRe)

		newHouseProfile.LoopPosition = utils.ExtractString(contents, loopPositionRe)

		newHouseProfile.SaleStatus = utils.ExtractString(contents, saleStatusRe)
		newHouseProfile.Discount = utils.ExtractString(contents, discountRe)
		newHouseProfile.DeliveryTime = utils.ExtractString(contents, deliveryTimeRe)
		if consultingPhone, err := strconv.Atoi(utils.ExtractString(contents, consultingPhoneRe)); err == nil {
			newHouseProfile.ConsultingPhone = consultingPhone
		}
		transportation := utils.ExtractString(contents, transportationRe)
		if transportation == "" {
			transportation = utils.ExtractString(contents, transportationRe2)
		}
		newHouseProfile.Transportation = transportation
		newHouseProfile.CommunityArea = utils.ExtractFloat64(contents, communityAreaRe)
		newHouseProfile.CommunityBuildingArea = utils.ExtractFloat64(contents, communityBuildingAreaRe)
		newHouseProfile.PlotRatio = utils.ExtractFloat64(contents, plotRatioRe)
		newHouseProfile.GreeningRate = utils.ExtractFloat64(contents, greeningRate) / 100
		newHouseProfile.PropertyManagementCompany = utils.ExtractString(contents, propertyManagementCompanyRe)
		newHouseProfile.PropertyFee = utils.ExtractFloat64(contents, propertyFeeRe)
		newHouseProfile.ProjectBrief = utils.ExtractString(contents, projectBriefRe)

		item := engine.Item{
			OriginUrl:  detailUrl,
			Id:         id,
			Province:   province,
			City:       city,
			Index:      "newhouse",
			Address:    utils.ExtractString(contents, addressRe),
			GatherTime: time.Now(),
			PayLoad:    newHouseProfile,
		}
		result.Items = append(result.Items, item)
	}
	return result
}
