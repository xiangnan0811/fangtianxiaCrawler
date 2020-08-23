package main

import (
	"fmt"
	"github.com/xiangnan0811/fangtianxiaCrawler/engine"
	"github.com/xiangnan0811/fangtianxiaCrawler/fangtianxia/parser"
	"github.com/xiangnan0811/fangtianxiaCrawler/scheduler"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/persist/client"
)

func main() {
	itemChan, err := client.ItemSaver(fmt.Sprintf(":%d", config.ItemSaverPort))
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.SimpleScheduler{},
		WorkerCount: 100,
		ItemChan:    itemChan,
	}
	e.Run(engine.Request{
		Url:        "https://www.fang.com/SoufunFamily.htm",
		ParserFunc: parser.ParseCityList,
	})
}
