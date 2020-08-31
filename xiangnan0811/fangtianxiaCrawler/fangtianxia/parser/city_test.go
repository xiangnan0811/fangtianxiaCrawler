package parser

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/axgle/mahonia"
)

func TestParseCityNewHouseList(t *testing.T) {
	c, err := ioutil.ReadFile("city_newhouse_list_test_data.html")
	if err != nil {
		fmt.Println(err)
	}
	bodyString := mahonia.NewDecoder("gbk").ConvertString(string(c))
	result := ParseCityNewHouseList([]byte(bodyString), "上海", "https://sh.newhouse.fang.com/house/s/?rfss=2-0-1")

	// verify result
	const resultSize = 43
	expectedUrls := []string{
		"https://sh.newhouse.fang.com/house/s/b92/?ctm=1.sh.xf_search.page.0",
		"https://sh.newhouse.fang.com/house/s/pudong/?ctm=1.sh.xf_search.lpsearch_area.0",
		"https://sh.newhouse.fang.com/house/s/jiading/?ctm=1.sh.xf_search.lpsearch_area.1",
		"https://sh.newhouse.fang.com/house/s/baoshan/?ctm=1.sh.xf_search.lpsearch_area.2",
		"https://sh.newhouse.fang.com/house/s/minhang/?ctm=1.sh.xf_search.lpsearch_area.3",
		"https://sh.newhouse.fang.com/house/s/songjiang/?ctm=1.sh.xf_search.lpsearch_area.4",
	}

	if realReqSize := len(result.Requests); realReqSize != resultSize {
		t.Errorf("result should have %d requests; but had %d", resultSize, realReqSize)
	}

	for i, url := range expectedUrls {
		if result.Requests[i].Url != url {
			t.Errorf("expected url #%d: %s; but was %s", i, url, result.Requests[i].Url)
		}
	}
}

func TestParseCityErShouHouseList(t *testing.T) {
	c, err := ioutil.ReadFile("city_ershouhouse_list_test_data.html")
	if err != nil {
		fmt.Println(err)
	}
	bodyString := mahonia.NewDecoder("gbk").ConvertString(string(c))
	result := ParseCityErShouHouseList([]byte(bodyString), "上海", "https://sh.esf.fang.com")

	// verify result
	const resultSize = 117
	expectedUrls := []string{
		"https://sh.esf.fang.com/house/i32/",
		"https://sh.esf.fang.com/house-a025/",
		"https://sh.esf.fang.com/house-a029/",
		"https://sh.esf.fang.com/house-a030/",
		"https://sh.esf.fang.com/house-a018/",
		"https://sh.esf.fang.com/house-a0586/",
		"https://sh.esf.fang.com/house-a028/",
	}

	if realReqSize := len(result.Requests); realReqSize != resultSize {
		t.Errorf("result should have %d requests; but had %d", resultSize, realReqSize)
	}

	for i, url := range expectedUrls {
		if result.Requests[i].Url != url {
			t.Errorf("expected url #%d: %s; but was %s", i, url, result.Requests[i].Url)
		}
	}
}
