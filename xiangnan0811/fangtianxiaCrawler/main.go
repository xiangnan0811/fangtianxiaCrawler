package main

import (
	"github.com/xiangnan0811/fangtianxiaCrawler/engine"
	"github.com/xiangnan0811/fangtianxiaCrawler/fangtianxia/parser"
	"github.com/xiangnan0811/fangtianxiaCrawler/persist"
	"github.com/xiangnan0811/fangtianxiaCrawler/scheduler"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"
)

func main() {
	itemChan, err := persist.ItemSaver()
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.SimpleScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan,
		RequestProcessor: engine.Worker,
	}
	e.Run(engine.Request{
		Url:    "https://www.fang.com/SoufunFamily.htm",
		Parser: engine.NewFuncParser(parser.ParseCityList, config.ParseCityList),
	})
}
