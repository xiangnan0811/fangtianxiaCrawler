package main

import (
	"github.com/xiangnan0811/fangtianxiaCrawler/engine"
	"github.com/xiangnan0811/fangtianxiaCrawler/model"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/rpcsupport"
	"log"
	"testing"
	"time"
)

func TestItemSaver(t *testing.T) {
	const host = ":1234"
	// start ItemSaverServer
	go log.Fatal(serveRpc(host))
	time.Sleep(time.Second)
	// start ItemSaverClient
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		panic(err)
	}
	// Call Save
	expected := engine.Item{
		OriginUrl:  "https://cd.esf.fang.com/chushou/3_213015060.htm",
		Id:         213015060,
		Province:   "四川",
		City:       "成都",
		Address:    "",
		Index:      "ershouhouse",
		GatherTime: time.Now(),
		PayLoad: model.ErShouHouseProfile{
			Title:             "景茂名都东郡清水套二，布局合理没有浪费",
			TotalPrice:        106,
			HouseType:         "2室1厅1卫",
			Room:              2,
			Hall:              1,
			Toilet:            1,
			Area:              75.45,
			UnitPrice:         14049,
			Orientation:       "北",
			Floor:             3,
			TotalFloor:        13,
			Decoration:        "简装修",
			Community:         "景茂名都东郡",
			Distinct:          "双流",
			ComArea:           "双流城区",
			BuildingAge:       2016,
			Elevator:          1,
			PropertyRight:     "商品房",
			ResidenceType:     "普通住宅",
			BuildingStructure: "平层",
			BuildingType:      "板楼",
			ListingTime:       time.Now(),
			CoreSellingPoint:  "此房为标准套二，阳台使用面积大可以做一个书房，布局合理。\n\n\n成都市景发有限公司开发，物业为景灿物业，物业费1.8，建成于2016年。\n\n\n双流万达广场直线距离500米，成都具规模的欧尚超市1500米，奥特莱斯1500米，菜市场300米，餐饮一条街200米内，双流8500亩湿地公园1800米，绕城双流收费口2000米。（数据来源百度地图）\n\n\n已建成的有516公交 816公交，双流01，02 ，06等等到达成都及周县的各路公交线路，出行您完全不用担心。地铁三号线航都大街站700米（数据来源于百度地图）",
			OwnerIntroduced:   "房东置换，诚心出售，价格可以再谈，欢迎看房，随时恭候。\n",
			TaxAnalysis:       "新房： 房税费首套房契税1至3个点。\n维修基金60元一平方\n物业费2元到3元一平方。别墅另算\n\n温馨提醒：购房按揭所需材料： 1 夫妻双方身份证原件与复印件2夫妻双方半年流水 3结婚证原件复印件 4收入证明（不同银行版本不一样）购买新房主要税费： 契税：90平以内：1% ；90-124 平收1.5%；124平以上收2%（交房后地税收取）\n房屋维修资金：建筑面积x60元/平\n燃气入户开通费：免费 房屋产权工本费300远\n注：以上收费仅供参考，均以政府相关部门和开发商终规定为准\n",
			CommunityFacility: "位置：小区地理位置，绿化，居民素质高\n配套：有健身器材，小型篮球场，可供小区居民使用\n车位：停车位充足，建有地上停车场和地下车库\n保安：小区内24保安巡逻\n",
			Services:          "本人从事房地产多年，公司有大量好房源，欢迎进入我的店铺查看，欢迎随时电话咨询，相信我的专业，为您置业安家保驾护航。",
			ImageUrlList:      []string{"https://cdnsfb.soufunimg.com/viewimage/1/2020_5/9/M17/24/21eb3f3340674416a76ab7583a39b1fb/1000x637c.jpg", "https://cdnsfb.soufunimg.com/viewimage/1/2020_5/9/M17/24/21eb3f3340674416a76ab7583a39b1fb/1000x637c.jpg"},
		},
	}
	result := ""
	err = client.Call(config.CrawlServiceRpc, expected, &result)
	if err != nil || result != "ok" {
		t.Errorf("result: %s; err: %s", result, err)
	}
}
