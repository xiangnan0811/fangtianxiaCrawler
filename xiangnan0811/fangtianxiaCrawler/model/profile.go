package model

import "time"

type NewHouseProfile struct {
	Name                      string   // 小区名
	Score                     float64  // 用户评分
	Tags                      []string // 项目特色
	AvgPrice                  string   // 均价
	HouseType                 []string // 主力户型
	PropertyType              string   // 物业类型
	PropertyRights            []string // 产权年限
	SaleAddress               string   // 售楼地址
	Open                      string   // 开盘时间
	PropertyDeveloper         string   // 开发商
	Decoration                string   // 装修情况
	LoopPosition              string   // 环线位置
	SaleStatus                string   // 销售状态
	Discount                  string   // 楼盘优惠
	Delivery                  string   // 交房时间
	ConsultingPhone           int      // 咨询电话
	Transportation            string   // 交通状况
	CommunityArea             float64  // 楼盘占地面积
	CommunityBuildingArea     float64  // 建筑面积
	PlotRatio                 float64  // 容积率
	GreeningRate              float64  // 绿化率
	PropertyManagementCompany string   // 物业公司
	PropertyFee               float64  // 物业费
	ProjectBrief              string   // 项目简介
}

type ErShouHouseProfile struct {
	Title             string    // 介绍名
	TotalPrice        float64   // 总价
	HouseType         string    // 户型字符串
	Room              int       // 室
	Hall              int       // 厅
	Toilet            int       // 卫
	Area              float64   // 建筑面积
	UnitPrice         float64   // 单价
	Orientation       string    // 朝向
	Floor             int       // 楼层
	TotalFloor        int       // 总楼层
	Decoration        string    // 装修
	Community         string    // 小区
	Distinct          string    // 区县
	ComArea           string    // 商圈
	BuildingAge       int       // 建筑年代
	Elevator          int       // 有无电梯，1为有，0为无，2为未知
	PropertyRight     string    // 产权性质
	ResidenceType     string    // 住宅类别
	BuildingStructure string    // 建筑结构
	BuildingType      string    // 建筑类别
	ListingTime       time.Time // 挂牌时间
	CoreSellingPoint  string    // 核心卖点
	OwnerIntroduced   string    // 业主心态
	TaxAnalysis       string    // 税费分析
	CommunityFacility string    // 小区配套
	Services          string    // 服务介绍
	ImageUrlList      []string  // 图片
}
