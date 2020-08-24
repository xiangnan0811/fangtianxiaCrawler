package main

import (
	"flag"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/rpcsupport"
	"log"
	"net/rpc"
	"strings"

	worker "github.com/xiangnan0811/fangtianxiaCrawler_distributed/worker/client"

	"github.com/xiangnan0811/fangtianxiaCrawler/engine"
	"github.com/xiangnan0811/fangtianxiaCrawler/fangtianxia/parser"
	"github.com/xiangnan0811/fangtianxiaCrawler/scheduler"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"
	itemsaver "github.com/xiangnan0811/fangtianxiaCrawler_distributed/persist/client"
)

var (
	itemSaverHost = flag.String("itemsaver_host", "", "itemsaver host")
	workerHosts   = flag.String("worker_hosts", "", "worker hosts (comma separated)")
)

func main() {
	flag.Parse()
	itemChan, err := itemsaver.ItemSaver(*itemSaverHost)
	if err != nil {
		panic(err)
	}
	pool := createClientPool(strings.Split(*workerHosts, ","))
	processor := worker.CreateProcessor(pool)

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

func createClientPool(hosts []string) chan *rpc.Client {
	var clients []*rpc.Client
	for _, h := range hosts {
		client, err := rpcsupport.NewClient(h)
		if err == nil {
			clients = append(clients, client)
			log.Printf("Creating client, connected to %s", h)
		} else {
			log.Printf("Error in creating client connect to %s: %v", h, err)
		}
	}

	out := make(chan *rpc.Client)
	go func() {
		for {
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out
}
