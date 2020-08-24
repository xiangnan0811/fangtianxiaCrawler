package worker

import (
	"errors"
	"log"

	"github.com/xiangnan0811/fangtianxiaCrawler/engine"
	"github.com/xiangnan0811/fangtianxiaCrawler/fangtianxia/parser"
	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"
)

type SerializedParser struct {
	FuncName string
	Url      string
	Province string
}

type Request struct {
	Url    string
	Parser SerializedParser
}

type ParseResult struct {
	Items    []engine.Item
	Requests []Request
}

func SerializeRequest(r engine.Request) Request {
	name, province, url := r.Parser.Serialize()
	return Request{
		Url: r.Url,
		Parser: SerializedParser{
			FuncName: name,
			Url:      url,
			Province: province,
		},
	}
}

func SerializeResult(r engine.ParseResult) ParseResult {
	result := ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		result.Requests = append(result.Requests, SerializeRequest(req))
	}
	return result
}

func DeserializeRequest(r Request) (engine.Request, error) {
	p, err := DeserializeParser(r.Parser)
	if err != nil {
		return engine.Request{}, nil
	}
	return engine.Request{
		Url:    r.Url,
		Parser: p,
	}, nil
}

func DeserializeResult(r ParseResult) engine.ParseResult {
	result := engine.ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		engineReq, err := DeserializeRequest(req)
		if err != nil {
			log.Printf("error deserializing request %v: %v", req, err)
			continue
		}
		result.Requests = append(result.Requests, engineReq)
	}
	return result
}

func DeserializeParser(p SerializedParser) (engine.Parser, error) {
	switch p.FuncName {
	case config.ParseCityList:
		return engine.NewFuncParser(parser.ParseCityList, config.ParseCityList), nil
	case config.NilParser:
		return engine.NilParser{}, nil
	case config.CityNewHouseListParser:
		return parser.NewCityNewHouseListParser(p.Province, p.Url), nil
	case config.CityErShouHouseListParser:
		return parser.NewCityErShouHouseListParser(p.Province, p.Url), nil
	case config.NewHouseParser:
		return parser.NewNewHouseParser(p.Province, p.Url), nil
	case config.ErShouHouseParser:
		return parser.NewErShouHouseParser(p.Province, p.Url), nil
	default:
		return nil, errors.New("unknown parser name")
	}
}
