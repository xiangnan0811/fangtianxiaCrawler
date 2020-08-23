package engine

import "time"

type ParserFunc func(contents []byte) ParseResult

type Request struct {
	Url        string
	ParserFunc ParserFunc
}

type ParseResult struct {
	Requests []Request
	Items    []Item
}

func NilParser([]byte) ParseResult {
	return ParseResult{}
}

type Item struct {
	OriginUrl  string
	Id         int
	Province   string
	City       string
	Index      string
	Address    string
	GatherTime time.Time // 采集时间
	PayLoad    interface{}
}
