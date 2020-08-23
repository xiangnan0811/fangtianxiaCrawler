package main

import (
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/persist"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/rpcsupport"
	"log"
)

func main() {
	log.Fatal(serveRpc(fmt.Sprintf(":%d", config.ItemSaverPort)))

}

func serveRpc(host string) error {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return err
	}
	return rpcsupport.ServeRpc(host, &persist.ItemSaverService{Client: client})
}
