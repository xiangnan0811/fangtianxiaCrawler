package engine

import (
	"time"

	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"
)

type ParserFunc func(contents []byte) ParseResult

type Parser interface {
	Parse(contents []byte) ParseResult
	Serialize() (name string, province string, url string)
}

type Request struct {
	Url    string
	Parser Parser
}

type ParseResult struct {
	Requests []Request
	Items    []Item
}

type NilParser struct {
}

func (n NilParser) Parse(_ []byte) ParseResult {
	return ParseResult{}
}

func (n NilParser) Serialize() (name string, province string, url string) {
	return config.NilParser, "", ""
}

type Item struct {
	OriginUrl  string
	Id         string
	Province   string
	City       string
	Index      string
	Address    string
	GatherTime time.Time // 采集时间
	PayLoad    interface{}
}

type FuncParser struct {
	parser ParserFunc
	name   string
}

func (f *FuncParser) Parse(contents []byte) ParseResult {
	return f.parser(contents)
}

func (f *FuncParser) Serialize() (name string, province string, url string) {
	return f.name, province, url
}

func NewFuncParser(p ParserFunc, name string) *FuncParser {
	return &FuncParser{parser: p, name: name}
}
