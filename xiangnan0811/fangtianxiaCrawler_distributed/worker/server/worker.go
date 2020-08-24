package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/worker"

	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/rpcsupport"
)

var port = flag.Int("port", 0, "the port for me to listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}
	log.Fatal(rpcsupport.ServeRpc(
		fmt.Sprintf(":%d", *port),
		worker.CrawlService{},
	))
}
