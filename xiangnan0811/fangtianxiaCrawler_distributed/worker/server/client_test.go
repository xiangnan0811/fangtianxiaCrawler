package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"

	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/worker"

	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/rpcsupport"
)

func TestCrawlService(t *testing.T) {
	const host = ":9000"
	go rpcsupport.ServeRpc(host, worker.CrawlService{})
	time.Sleep(time.Second)

	client, err := rpcsupport.NewClient(host)
	if err != nil {
		panic(err)
	}
	req := worker.Request{
		Url: "https://chunfengshilirc028.fang.com/house/3210162372/housedetail.htm",
		Parser: worker.SerializedParser{
			FuncName: config.NewHouseParser,
			Province: "四川",
			Url:      "https://chunfengshilirc028.fang.com/house/3210162372/housedetail.htm",
		},
	}
	var result worker.ParseResult

	err = client.Call(config.CrawlServiceRpc, req, &result)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(result)
	}
}
