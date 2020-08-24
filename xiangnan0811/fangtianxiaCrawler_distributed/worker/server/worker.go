package main

import (
	"fmt"
	"log"

	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/worker"

	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"

	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/rpcsupport"
)

func main() {
	log.Fatal(rpcsupport.ServeRpc(
		fmt.Sprintf(":%d", config.WorkerPort0),
		worker.CrawlService{},
	))
}
