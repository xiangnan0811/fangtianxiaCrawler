package client2

import (
	"net/rpc"

	"github.com/xiangnan0811/fangtianxiaCrawler/engine"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/worker"
)

func CreateProcessor(clientChan chan *rpc.Client) engine.Processor {

	return func(req engine.Request) (engine.ParseResult, error) {
		sReq := worker.SerializeRequest(req)
		var sResult worker.ParseResult
		c := <-clientChan
		err := c.Call(config.CrawlServiceRpc, sReq, &sResult)
		if err != nil {
			return engine.ParseResult{}, err
		}
		return worker.DeserializeResult(sResult), nil
	}
}
