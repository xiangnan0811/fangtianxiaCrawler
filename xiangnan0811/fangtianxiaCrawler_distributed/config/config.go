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
	// Service ports
	ItemSaverPort = 1234
	// Worker ports
	WorkerPort0 = 9000
	//Rate limiting
	Qps = 50
)
