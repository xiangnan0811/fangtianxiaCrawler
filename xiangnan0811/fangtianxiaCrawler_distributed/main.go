package main

import (
	"fmt"

	worker "github.com/xiangnan0811/fangtianxiaCrawler_distributed/worker/client"

	"github.com/xiangnan0811/fangtianxiaCrawler/engine"
	"github.com/xiangnan0811/fangtianxiaCrawler/fangtianxia/parser"
	"github.com/xiangnan0811/fangtianxiaCrawler/scheduler"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"
	itemsaver "github.com/xiangnan0811/fangtianxiaCrawler_distributed/persist/client"
)

func main() {
	itemChan, err := itemsaver.ItemSaver(fmt.Sprintf(":%d", config.ItemSaverPort))
	if err != nil {
		panic(err)
	}
	processor, err := worker.CreateProcessor()
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan,
		RequestProcessor: processor,
	}
	e.Run(engine.Request{
		Url:    "https://www.fang.com/SoufunFamily.htm",
		Parser: engine.NewFuncParser(parser.ParseCityList, config.ParseCityList),
	})
}
