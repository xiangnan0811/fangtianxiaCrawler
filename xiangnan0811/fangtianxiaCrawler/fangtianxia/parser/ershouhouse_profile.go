package parser

import (
	"regexp"
	"strings"
	"time"

	"github.com/xiangnan0811/fangtianxiaCrawler/engine"
	"github.com/xiangnan0811/fangtianxiaCrawler/model"
	"github.com/xiangnan0811/fangtianxiaCrawler/utils"

	"github.com/antchfx/htmlquery"
)

var (
	erShouHouseCityRe    = regexp.MustCompile(`pageConfig.city='(.*?)';`)
	redirectRe           = regexp.MustCompile(`(<title>跳转...</title>)`)
	locationRe           = regexp.MustCompile(`//location.href="(.*?)";`)
	VerifyRe             = regexp.MustCompile(`<title>(访问验证.*?)</title>`)
	totalPriceRe         = regexp.MustCompile(`pageConfig.price ?= ?'?(\d+.?\d+)'?;`)
	roomRe               = regexp.MustCompile(`pageConfig.room ?= ?'?(\d+)'?;`)
	hallRe               = regexp.MustCompile(`pageConfig.hall ?= ?'?(\d+)'?;`)
	toiletRe             = regexp.MustCompile(`pageConfig.toilet ?= ?'?(\d+)'?;`)
	areaRe               = regexp.MustCompile(`pageConfig.area ?= ?'?(\d+\.?\d+)'?;`)
	unitPriceRe          = regexp.MustCompile(`pageConfig.unitprice ?= ?'?(\d+.?\d+)'?;`)
	orientationRe        = regexp.MustCompile(`pageConfig.forward ?= ?'?(.*?)'?;`)
	floorRe              = regexp.MustCompile(`pageConfig.floor ?= ?'?(\d+)'?;`)
	totalFloorRe         = regexp.MustCompile(`pageConfig.totalfloor ?= ?'?(\d+)'?;`)
	communityRe          = regexp.MustCompile(`pageConfig.projname ?= ?'(.*?)';`)
	distinctRe           = regexp.MustCompile(`pageConfig.district ?= ?'(.*?)',`)
	distinctRe2          = regexp.MustCompile(`pageConfig.district ?= ?"(.*?)";`)
	comAreaRe            = regexp.MustCompile(`pageConfig.comarea ?= ?'(.*?)';`)
	erShouHouseAddressRe = regexp.MustCompile(`pageConfig.address ?= ?'(.*?)';`)

	buildingAgeRe       = regexp.MustCompile(`建筑年代</span>\s*.*?">(\d+).*?年</span>`)
	elevatorRe          = regexp.MustCompile(`有无电梯</span>\s*.*?">(.*?)\s*</span>`)
	propertyRightRe     = regexp.MustCompile(`产权性质</span>\s*.*?">(.*?)\s*</span>`)
	residenceType       = regexp.MustCompile(`住宅类别</span>\s*.*?">(.*?)\s*</span>`)
	buildingStructureRe = regexp.MustCompile(`建筑结构</span>\s*.*?">(.*?)\s*</span>`)
	buildingTypeRe      = regexp.MustCompile(`建筑类别</span>\s*.*?">(.*?)\s*</span>`)
	listingTimeRe2      = regexp.MustCompile(`挂牌时间</span>\s*.*?">\s*([-\d]+)\s*</span>`)
	listingTimeRe       = regexp.MustCompile(`pageConfig.inserttime ?= ?'([- :\d]+) \d{3}';`)
	coreSellingPointRe  = regexp.MustCompile(`核心卖点</div><div class="fyms_con floatl gray3"><div>((.|\n)*?)</div></div>`)
	ownerIntroducedRe   = regexp.MustCompile(`业主心态</div><div class="fyms_con floatl gray3">((.|\n)*?)</div>`)
	taxAnalysisRe       = regexp.MustCompile(`税费分析</div><div class="fyms_con floatl gray3">((.|\n)*?)</div>`)
	communityFacilityRe = regexp.MustCompile(`小区配套</div><div class="fyms_con floatl gray3">((.|\n)*?)</div>`)
	servicesRe          = regexp.MustCompile(`服务介绍</div><div class="fyms_con floatl gray3">((.|\n)*?)</div>`)
)

func ParseErShouHouse(contents []byte, province string, url string) engine.ParseResult {
	result := engine.ParseResult{}

	if redirect := redirectRe.FindAllSubmatch(contents, -1); redirect != nil {
		locationUrl := utils.ExtractString(contents, locationRe)

		if locationUrl != "暂无资料" && locationUrl != "" {
			result.Requests = append(result.Requests, engine.Request{
				Url:    locationUrl,
				Parser: NewErShouHouseParser(province, locationUrl),
			})
			return result
		}
	} else if verify := VerifyRe.FindAllSubmatch(contents, -1); verify != nil {
		// TODO: 验证码未处理，只是将请求发回调度器，稍后重新请求处理
		url = strings.Split(url, "?")[0]
		result.Requests = append(result.Requests, engine.Request{
			Url:    url,
			Parser: NewNewHouseParser(province, url),
		})
		return result
	} else {
		erShouHouseProfile := model.ErShouHouseProfile{}

		root, _ := htmlquery.Parse(strings.NewReader(string(contents[:])))
		titleNode := htmlquery.Find(root, "//meta[@name='keywords']")
		if titleNode != nil {
			erShouHouseProfile.Title = strings.TrimSpace(htmlquery.SelectAttr(titleNode[0], "content"))
		} else {
			titleNode = htmlquery.Find(root, "//h1")
			if titleNode != nil {
				erShouHouseProfile.Title = strings.TrimSpace(htmlquery.InnerText(titleNode[0]))
			}
		}

		erShouHouseProfile.TotalPrice = utils.ExtractFloat64(contents, totalPriceRe)
		erShouHouseProfile.Room = utils.ExtractInt(contents, roomRe)
		erShouHouseProfile.Hall = utils.ExtractInt(contents, hallRe)
		erShouHouseProfile.Toilet = utils.ExtractInt(contents, toiletRe)
		erShouHouseProfile.Area = utils.ExtractFloat64(contents, areaRe)
		erShouHouseProfile.UnitPrice = utils.ExtractFloat64(contents, unitPriceRe)
		erShouHouseProfile.Orientation = utils.ExtractString(contents, orientationRe)
		if erShouHouseProfile.Orientation == "" || erShouHouseProfile.Orientation == "暂无资料" {
			orientationNode := htmlquery.Find(root, "//div[contains(@class, 'w146')]/div[@class='tt']")
			if len(orientationNode) >= 2 {
				erShouHouseProfile.Orientation = strings.TrimSpace(htmlquery.InnerText(orientationNode[1]))
			}
		}

		erShouHouseProfile.Floor = utils.ExtractInt(contents, floorRe)
		erShouHouseProfile.TotalFloor = utils.ExtractInt(contents, totalFloorRe)
		houseTypeNode := htmlquery.Find(root, "//div[contains(@class, 'w146')]/div[@class='tt']")
		if houseTypeNode != nil {
			erShouHouseProfile.HouseType = strings.TrimSpace(htmlquery.InnerText(houseTypeNode[0]))
		}
		decorationNode := htmlquery.Find(root, "//div[contains(@class, 'w132')]/div[@class='tt']")
		if len(decorationNode) >= 2 {
			erShouHouseProfile.Decoration = strings.TrimSpace(htmlquery.InnerText(decorationNode[1]))
		}

		erShouHouseProfile.Community = utils.ExtractString(contents, communityRe)
		distinct := utils.ExtractString(contents, distinctRe2)
		if distinct == "暂无资料" || distinct == "" {
			distinct = utils.ExtractString(contents, distinctRe)
		}
		erShouHouseProfile.Distinct = distinct
		erShouHouseProfile.ComArea = utils.ExtractString(contents, comAreaRe)

		erShouHouseProfile.BuildingAge = utils.ExtractInt(contents, buildingAgeRe)
		elevator := utils.ExtractString(contents, elevatorRe)
		if elevator == "有" {
			erShouHouseProfile.Elevator = 1
		} else if elevator == "无" || elevator == "没有" {
			erShouHouseProfile.Elevator = 0
		} else {
			erShouHouseProfile.Elevator = 2
		}
		erShouHouseProfile.PropertyRight = utils.ExtractString(contents, propertyRightRe)
		erShouHouseProfile.ResidenceType = utils.ExtractString(contents, residenceType)
		erShouHouseProfile.BuildingStructure = utils.ExtractString(contents, buildingStructureRe)

		erShouHouseProfile.BuildingType = utils.ExtractString(contents, buildingTypeRe)
		listTimeString := utils.ExtractString(contents, listingTimeRe)
		if listTimeString != "" && listTimeString != "暂无资料" {
			listTime, err := time.ParseInLocation("2006-01-02 15:04:05", listTimeString, time.Local)
			if err == nil {
				erShouHouseProfile.ListingTime = listTime
			}
		}
		if listTimeString == "" || listTimeString == "暂无资料" {
			listTimeString = utils.ExtractString(contents, listingTimeRe2)
			if listTimeString != "" && listTimeString != "暂无资料" {
				listTime, err := time.ParseInLocation("2006-01-02", listTimeString, time.Local)
				if err == nil {
					erShouHouseProfile.ListingTime = listTime
				}
			}
		}

		erShouHouseProfile.CoreSellingPoint = utils.ExtractString(contents, coreSellingPointRe)
		erShouHouseProfile.OwnerIntroduced = utils.ExtractString(contents, ownerIntroducedRe)
		erShouHouseProfile.TaxAnalysis = utils.ExtractString(contents, taxAnalysisRe)
		erShouHouseProfile.CommunityFacility = utils.ExtractString(contents, communityFacilityRe)
		erShouHouseProfile.Services = utils.ExtractString(contents, servicesRe)

		imageUrlList := htmlquery.Find(root, "//img[@class='loadimg']")
		if len(imageUrlList) >= 2 {
			for _, node := range imageUrlList {
				erShouHouseProfile.ImageUrlList = append(erShouHouseProfile.ImageUrlList, "https:"+htmlquery.SelectAttr(node, "data-src"))
			}
		}

		idStringList := strings.Split(strings.Split(url, ".htm")[0], "_")
		cityString := utils.ExtractString(contents, erShouHouseCityRe)
		if len(idStringList) == 2 {
			idString := idStringList[1]

			item := engine.Item{
				OriginUrl:  url,
				Id:         idString,
				Province:   province,
				City:       cityString,
				Index:      "ershouhouse",
				Address:    utils.ExtractString(contents, erShouHouseAddressRe),
				GatherTime: time.Now(),
				PayLoad:    erShouHouseProfile,
			}
			result.Items = append(result.Items, item)
		}
	}
	return result
}
