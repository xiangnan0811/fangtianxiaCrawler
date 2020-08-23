package parser

import (
	"io/ioutil"
	"testing"
)

func TestParseCityList(t *testing.T) {
	contents, _ := ioutil.ReadFile("citylist_test_data.html")
	result := ParseCityList(contents)

	// verify result
	const resultSize = 1250
	expectedUrls := []string{
		"https://newhouse.fang.com/house/s",
		"https://esf.fang.com",
		"https://sh.newhouse.fang.com/house/s",
		"https://sh.esf.fang.com",
		"https://tj.newhouse.fang.com/house/s",
		"https://tj.esf.fang.com",
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
