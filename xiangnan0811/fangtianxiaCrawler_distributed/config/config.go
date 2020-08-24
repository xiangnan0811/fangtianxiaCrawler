package config

const (
	// Parser Name
	ParseCityList             = "ParseCityList"
	CityNewHouseListParser    = "CityNewHouseListParser"
	CityErShouHouseListParser = "CityErShouHouseListParser"
	NewHouseParser            = "NewHouseParser"
	ErShouHouseParser         = "ErShouHouseParser"
	NilParser                 = "NilParser"

	// RPC Endpoints
	ItemServiceRpc = "ItemSaverService.Save"
	// CrawlServiceRpc
	CrawlServiceRpc = "CrawlService.Process"

	//Rate limiting
	Qps = 50
)
