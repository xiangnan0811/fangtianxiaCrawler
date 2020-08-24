package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/persist"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/rpcsupport"
)

var port = flag.Int("port", 0, "the port for me to listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}
	log.Fatal(serveRpc(fmt.Sprintf(":%d", *port)))

}

func serveRpc(host string) error {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return err
	}
	return rpcsupport.ServeRpc(host, &persist.ItemSaverService{Client: client})
}
